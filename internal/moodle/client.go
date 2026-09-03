package moodle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Client is a Moodle Web Services HTTP client.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client

	siteInfoOnce  sync.Once
	siteInfoCache *SiteInfo
	siteInfoErr   error
}

// NewClient constructs a new Moodle client.
func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// request is the generic POST dispatcher for the Moodle Web Services REST API.
func (c *Client) request(ctx context.Context, wsFunction string, params url.Values, result any) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("wstoken", c.token)
	params.Set("moodlewsrestformat", "json")
	params.Set("wsfunction", wsFunction)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/webservice/rest/server.php", strings.NewReader(params.Encode()))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	// Moodle returns HTTP 200 even on errors; detect the error envelope.
	if bytes.HasPrefix(bytes.TrimSpace(body), []byte("{")) {
		var moodleError struct {
			Exception string `json:"exception"`
			ErrorCode string `json:"errorcode"`
			Message   string `json:"message"`
		}
		if err := json.Unmarshal(body, &moodleError); err == nil {
			if moodleError.Exception != "" || moodleError.ErrorCode != "" {
				return fmt.Errorf("moodle error %s: %s", moodleError.ErrorCode, moodleError.Message)
			}
		}
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// PreSeedUserID populates the GetSiteInfo cache with a known user ID from
// config.json, so the server never makes a GetSiteInfo network call during
// normal operation. Only the UserID field is populated; other fields remain zero.
func (c *Client) PreSeedUserID(userID int) {
	c.siteInfoOnce.Do(func() {
		c.siteInfoCache = &SiteInfo{UserID: userID}
	})
}

// GetSiteInfo returns site and user info. The result is cached after the first
// successful call so subsequent calls do not make a network request.
func (c *Client) GetSiteInfo(ctx context.Context) (*SiteInfo, error) {
	c.siteInfoOnce.Do(func() {
		var info SiteInfo
		if err := c.request(ctx, "core_webservice_get_site_info", nil, &info); err != nil {
			c.siteInfoErr = err
			return
		}
		c.siteInfoCache = &info
	})
	return c.siteInfoCache, c.siteInfoErr
}

// GetEnrolledCourses returns all courses the given user is enrolled in.
func (c *Client) GetEnrolledCourses(ctx context.Context, userID int) ([]Course, error) {
	params := url.Values{"userid": {fmt.Sprintf("%d", userID)}}
	var courses []Course
	if err := c.request(ctx, "core_enrol_get_users_courses", params, &courses); err != nil {
		return nil, err
	}
	return courses, nil
}

// GetCourseContents returns all sections and modules for the given course.
func (c *Client) GetCourseContents(ctx context.Context, courseID int) ([]Section, error) {
	params := url.Values{"courseid": {fmt.Sprintf("%d", courseID)}}
	var sections []Section
	if err := c.request(ctx, "core_course_get_contents", params, &sections); err != nil {
		return nil, err
	}
	return sections, nil
}

// GetCalendarEvents returns upcoming action events sorted by due date.
func (c *Client) GetCalendarEvents(ctx context.Context, timesortfrom int64, limitnum int) ([]CalendarEvent, error) {
	params := url.Values{
		"timesortfrom": {fmt.Sprintf("%d", timesortfrom)},
		"limitnum":     {fmt.Sprintf("%d", limitnum)},
	}
	var result struct {
		Events []CalendarEvent `json:"events"`
	}
	if err := c.request(ctx, "core_calendar_get_action_events_by_timesort", params, &result); err != nil {
		return nil, err
	}
	return result.Events, nil
}

// GetCourseGrades returns grade items for the given user in the given course.
func (c *Client) GetCourseGrades(ctx context.Context, courseID, userID int) ([]UserGrade, error) {
	params := url.Values{
		"courseid": {fmt.Sprintf("%d", courseID)},
		"userid":   {fmt.Sprintf("%d", userID)},
	}
	var result struct {
		UserGrades []UserGrade `json:"usergrades"`
	}
	if err := c.request(ctx, "gradereport_user_get_grade_items", params, &result); err != nil {
		return nil, err
	}
	return result.UserGrades, nil
}

// GetForumsByCourse returns all forums in a course.
func (c *Client) GetForumsByCourse(ctx context.Context, courseID int) ([]Forum, error) {
	params := url.Values{"courseids[0]": {fmt.Sprintf("%d", courseID)}}
	var forums []Forum
	if err := c.request(ctx, "mod_forum_get_forums_by_courses", params, &forums); err != nil {
		return nil, err
	}
	return forums, nil
}

// GetForumDiscussions returns discussions for a given forum, newest first.
func (c *Client) GetForumDiscussions(ctx context.Context, forumID int) ([]Discussion, error) {
	params := url.Values{
		"forumid":   {fmt.Sprintf("%d", forumID)},
		"sortorder": {"1"},
	}
	var result struct {
		Discussions []Discussion `json:"discussions"`
	}
	if err := c.request(ctx, "mod_forum_get_forum_discussions", params, &result); err != nil {
		return nil, err
	}
	return result.Discussions, nil
}

// BaseURL returns the Moodle base URL configured for this client.
// Used by tool handlers to construct direct links (e.g. assignment review URLs).
func (c *Client) BaseURL() string { return c.baseURL }

// StageAssignmentDraft saves a file and/or online text to a Moodle assignment
// as a Draft (not submitted for grading). The student must manually click
// "Submit assignment" in the Moodle UI to finalize.
//
// At least one of draftItemID (>0) or text (non-empty) must be provided.
func (c *Client) StageAssignmentDraft(ctx context.Context, assignmentID, draftItemID int, text string) error {
	params := url.Values{
		"assignmentid": {fmt.Sprintf("%d", assignmentID)},
	}
	if draftItemID != 0 {
		params.Set("plugindata[files_filemanager]", fmt.Sprintf("%d", draftItemID))
	}
	if text != "" {
		params.Set("plugindata[onlinetext_editor][text]", text)
		params.Set("plugindata[onlinetext_editor][format]", "1")
		params.Set("plugindata[onlinetext_editor][itemid]", "0")
	}
	var warnings []any
	return c.request(ctx, "mod_assign_save_submission", params, &warnings)
}

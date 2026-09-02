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

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client

	siteInfoOnce sync.Once
	siteInfoCache *SiteInfo
	siteInfoErr   error
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Send the HTTP request to get payload - Generic for reuse
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
	if bytes.HasPrefix(bytes.TrimSpace(body), []byte("{")) {
		var moodleError struct {
			Exception string `json:"exception"`
			ErrorCode string `json:"errorcode"`
			Message   string `json:"message"`
		}
		if err := json.Unmarshal(body, &moodleError); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		if moodleError.Exception != "" || moodleError.ErrorCode != "" {
			return fmt.Errorf("moodle error %s: %s", moodleError.ErrorCode, moodleError.Message)
		}
	}
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

// Specific get functions to get different datas

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

func (c *Client) GetEnrolledCourses(ctx context.Context, userID int) ([]Course, error) {
	params := url.Values{"userid": {fmt.Sprintf("%d", userID)}}
	var courses []Course
	if err := c.request(ctx, "core_enrol_get_users_courses", params, &courses); err != nil {
		return nil, err
	}
	return courses, nil
}

func (c *Client) GetCourseContents(ctx context.Context, courseID int) ([]Section, error) {
	params := url.Values{"courseid": {fmt.Sprintf("%d", courseID)}}
	var sections []Section
	if err := c.request(ctx, "core_course_get_contents", params, &sections); err != nil {
		return nil, err
	}
	return sections, nil
}

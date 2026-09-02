package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tuanha/vgu-mcp/internal/moodle"
)

// RegisterAnnouncementTools registers MCP tools for reading course announcements.
func RegisterAnnouncementTools(s *server.MCPServer, client *moodle.Client) {
	s.AddTool(mcp.NewTool(
		"read_course_announcements",
		mcp.WithDescription("Read the latest announcements from the course announcement forum"),
		mcp.WithNumber("courseid", mcp.Required(), mcp.Description("The Moodle course ID")),
	), readCourseAnnouncements(client))
}

func readCourseAnnouncements(client *moodle.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		courseID, err := request.RequireInt("courseid")
		if err != nil {
			return mcp.NewToolResultError("missing or invalid required parameter 'courseid'"), nil
		}

		forums, err := client.GetForumsByCourse(ctx, courseID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("fetch forums: %v", err)), nil
		}

		// Locate the announcement forum: type == "news" or name contains "Announcement".
		var announcementForumID int
		for _, f := range forums {
			if f.Type == "news" || strings.Contains(strings.ToLower(f.Name), "announcement") {
				announcementForumID = f.ID
				break
			}
		}
		if announcementForumID == 0 {
			return mcp.NewToolResultText("No announcement forum found for this course."), nil
		}

		discussions, err := client.GetForumDiscussions(ctx, announcementForumID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("fetch discussions: %v", err)), nil
		}
		return jsonResult(discussions)
	}
}

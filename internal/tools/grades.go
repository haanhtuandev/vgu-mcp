package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tuanha/vgu-mcp/internal/moodle"
)

// RegisterGradeTools registers MCP tools for course grades.
func RegisterGradeTools(s *server.MCPServer, client *moodle.Client) {
	s.AddTool(mcp.NewTool(
		"get_course_grades",
		mcp.WithDescription("Fetch grade items and scores for the authenticated user in a given course"),
		mcp.WithNumber("courseid", mcp.Required(), mcp.Description("The Moodle course ID")),
	), getCourseGrades(client))
}

func getCourseGrades(client *moodle.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		courseID, err := request.RequireInt("courseid")
		if err != nil {
			return mcp.NewToolResultError("missing or invalid required parameter 'courseid'"), nil
		}
		siteInfo, err := client.GetSiteInfo(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("fetch site info: %v", err)), nil
		}
		grades, err := client.GetCourseGrades(ctx, courseID, siteInfo.UserID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("fetch grades: %v", err)), nil
		}
		return jsonResult(grades)
	}
}

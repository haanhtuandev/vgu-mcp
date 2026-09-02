package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tuanha/vgu-mcp/internal/moodle"
)

func RegisterCourseTools(server *server.MCPServer, client *moodle.Client) {
	server.AddTool(mcp.NewTool(
		"get_enrolled_courses",
		mcp.WithDescription("Fetch all Moodle courses enrolled by the authenticated user"),
	), getEnrolledCourses(client))

	server.AddTool(mcp.NewTool(
		"get_course_contents",
		mcp.WithDescription("Fetch sections, modules, and learning materials for a given course ID"),
		mcp.WithNumber("courseid", mcp.Required(), mcp.Description("The Moodle course ID")),
	), getCourseContents(client))
}

func getEnrolledCourses(client *moodle.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		siteInfo, err := client.GetSiteInfo(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("fetch site info: %v", err)), nil
		}
		courses, err := client.GetEnrolledCourses(ctx, siteInfo.UserID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("fetch courses: %v", err)), nil
		}
		return jsonResult(courses)
	}
}

func getCourseContents(client *moodle.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		courseID, err := request.RequireInt("courseid")
		if err != nil {
			return mcp.NewToolResultError("missing or invalid required parameter 'courseid'"), nil
		}
		sections, err := client.GetCourseContents(ctx, courseID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("fetch course contents: %v", err)), nil
		}
		return jsonResult(sections)
	}
}

func jsonResult(value any) (*mcp.CallToolResult, error) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return mcp.NewToolResultError("format JSON output"), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

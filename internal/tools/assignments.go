package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tuanha/vgu-mcp/internal/moodle"
)

// RegisterAssignmentTools registers MCP tools for assignment submission.
func RegisterAssignmentTools(s *server.MCPServer, client *moodle.Client) {
	s.AddTool(mcp.NewTool(
		"submit_assignment_draft",
		mcp.WithDescription("Submit an online-text draft for a Moodle assignment"),
		mcp.WithNumber("assignmentid", mcp.Required(), mcp.Description("The Moodle assignment ID")),
		mcp.WithString("text", mcp.Required(), mcp.Description("The submission text content (plain text or HTML)")),
	), submitAssignmentDraft(client))
}

func submitAssignmentDraft(client *moodle.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		assignmentID, err := request.RequireInt("assignmentid")
		if err != nil {
			return mcp.NewToolResultError("missing or invalid required parameter 'assignmentid'"), nil
		}
		text, err := request.RequireString("text")
		if err != nil {
			return mcp.NewToolResultError("missing or invalid required parameter 'text'"), nil
		}

		if err := client.SaveAssignmentSubmission(ctx, assignmentID, text); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("submit assignment: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Assignment %d submitted successfully.", assignmentID)), nil
	}
}

package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tuanha/vgu-mcp/internal/moodle"
)

// RegisterDeadlineTools registers MCP tools for upcoming deadlines.
func RegisterDeadlineTools(s *server.MCPServer, client *moodle.Client) {
	s.AddTool(mcp.NewTool(
		"get_upcoming_deadlines",
		mcp.WithDescription("Fetch upcoming assignment and activity deadlines for the authenticated user"),
		mcp.WithNumber("limitnum", mcp.Description("Maximum number of events to return (default: 10)")),
	), getUpcomingDeadlines(client))
}

func getUpcomingDeadlines(client *moodle.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		limitnum := 10
		if n, err := request.RequireInt("limitnum"); err == nil {
			limitnum = n
		}

		events, err := client.GetCalendarEvents(ctx, time.Now().Unix(), limitnum)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("fetch calendar events: %v", err)), nil
		}
		return jsonResult(events)
	}
}

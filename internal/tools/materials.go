package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tuanha/vgu-mcp/internal/moodle"
)

// RegisterMaterialTools registers MCP tools for downloading course materials.
func RegisterMaterialTools(s *server.MCPServer, client *moodle.Client) {
	s.AddTool(mcp.NewTool(
		"download_course_material",
		mcp.WithDescription("Stream-download a course file (PDF, ZIP, etc.) from Moodle to a local directory"),
		mcp.WithString("file_url", mcp.Required(), mcp.Description("The fileurl from get_course_contents (pluginfile.php URL)")),
		mcp.WithString("destination_dir", mcp.Description("Local directory to save the file (default: ./downloads)")),
	), downloadCourseMaterial(client))
}

func downloadCourseMaterial(client *moodle.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fileURL, err := request.RequireString("file_url")
		if err != nil {
			return mcp.NewToolResultError("missing or invalid required parameter 'file_url'"), nil
		}
		destDir := "./downloads"
		if d, err := request.RequireString("destination_dir"); err == nil && d != "" {
			destDir = d
		}

		path, err := client.Download(ctx, fileURL, destDir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("download file: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("File saved to: %s", path)), nil
	}
}

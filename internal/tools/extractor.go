package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tuanha/vgu-mcp/internal/extractor"
	"github.com/tuanha/vgu-mcp/internal/moodle"
)

// RegisterExtractorTools registers MCP tools for in-place text extraction
// from course material files (currently PDF only).
func RegisterExtractorTools(s *server.MCPServer, client *moodle.Client) {
	s.AddTool(mcp.NewTool(
		"extract_course_material_text",
		mcp.WithDescription(
			"Extract plain text from a course material PDF and return it as Markdown. "+
				"Provide file_url (a pluginfile URL from get_course_contents) to download and extract on-the-fly, "+
				"or local_filepath to extract from an already-downloaded file. "+
				"The AI can immediately use the returned text to summarise, quiz, or explain the material.",
		),
		mcp.WithString("file_url",
			mcp.Description("Moodle pluginfile URL from get_course_contents (must end in .pdf)")),
		mcp.WithString("local_filepath",
			mcp.Description("Absolute or relative path to a local PDF file")),
	), extractCourseMaterialText(client))
}

func extractCourseMaterialText(client *moodle.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fileURL, _ := req.RequireString("file_url")
		localPath, _ := req.RequireString("local_filepath")

		if fileURL == "" && localPath == "" {
			return mcp.NewToolResultError(
				"provide either 'file_url' (Moodle URL) or 'local_filepath'"), nil
		}

		var (
			data     []byte
			filename string
			err      error
		)

		if fileURL != "" {
			if !isPDFPath(fileURL) {
				return mcp.NewToolResultError(
					"only PDF files are supported — use download_course_material for other formats"), nil
			}
			filename = filepath.Base(strings.SplitN(fileURL, "?", 2)[0])
			data, err = client.Fetch(ctx, fileURL)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("fetch file: %v", err)), nil
			}
		} else {
			if !isPDFPath(localPath) {
				return mcp.NewToolResultError(
					"only PDF files are supported — use download_course_material for other formats"), nil
			}
			filename = filepath.Base(localPath)
			data, err = os.ReadFile(localPath)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("read file: %v", err)), nil
			}
		}

		result, err := extractor.ExtractPDF(data, filename)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("extract PDF: %v", err)), nil
		}
		if result.Content == "" {
			return mcp.NewToolResultError(
				"no text could be extracted — the PDF may be a scanned image with no text layer"), nil
		}

		out, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(out)), nil
	}
}

// isPDFPath returns true if the path or URL (before any query string) ends in .pdf.
func isPDFPath(path string) bool {
	base := strings.SplitN(path, "?", 2)[0]
	return strings.EqualFold(filepath.Ext(base), ".pdf")
}

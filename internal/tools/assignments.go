package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tuanha/vgu-mcp/internal/moodle"
)

// stageDraftResult is the structured JSON response returned to the AI agent.
type stageDraftResult struct {
	Status          string `json:"status"`
	AssignmentID    int    `json:"assignment_id"`
	Filename        string `json:"filename,omitempty"`
	DraftItemID     int    `json:"draft_item_id,omitempty"`
	MoodleReviewURL string `json:"moodle_review_url,omitempty"`
	Message         string `json:"message"`
}

// RegisterAssignmentTools registers MCP tools for assignment staging.
func RegisterAssignmentTools(s *server.MCPServer, client *moodle.Client) {
	s.AddTool(mcp.NewTool(
		"stage_assignment_draft",
		mcp.WithDescription(
			"Uploads a local file (ZIP, PDF, archive) and/or a text note to a Moodle assignment. "+
				"Saves strictly as a DRAFT — the student must manually click 'Submit assignment' "+
				"in Moodle to finalize for grading. Returns a direct review URL.",
		),
		mcp.WithNumber("assignment_id", mcp.Required(),
			mcp.Description(
				"The assignment's INSTANCE ID — use the 'instance' field (NOT the 'id' field) "+
					"from the module object returned by get_course_contents. "+
					"These are two different numbers: 'id' is the course module ID used in URLs, "+
					"'instance' is the assignment record ID required by the submission API. "+
					"Example: if get_course_contents returns {\"id\":34852, \"instance\":2432}, use 2432.",
			)),
		mcp.WithString("review_url",
			mcp.Description(
				"The assignment's Moodle page URL — use the 'url' field from the same module object "+
					"returned by get_course_contents. "+
					"Example: \"https://moodle.vgu.edu.vn/mod/assign/view.php?id=34852\". "+
					"Pass this directly; do not construct it from the assignment_id.",
			)),
		mcp.WithString("file_path",
			mcp.Description("Local path to the file to upload (e.g. ./dist/lab2.zip)")),
		mcp.WithString("text_content",
			mcp.Description("Online text or notes to include with the draft submission")),
	), stageAssignmentDraft(client))
}

func stageAssignmentDraft(client *moodle.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		assignmentID, err := request.RequireInt("assignment_id")
		if err != nil {
			return mcp.NewToolResultError("missing or invalid required parameter 'assignment_id'"), nil
		}

		reviewURL, _ := request.RequireString("review_url")
		filePath, _ := request.RequireString("file_path")
		textContent, _ := request.RequireString("text_content")

		if filePath == "" && textContent == "" {
			return mcp.NewToolResultError(
				"provide at least one of 'file_path' or 'text_content'"), nil
		}

		// Step 1: Upload file to draft area (if provided).
		var upload moodle.UploadResponse
		if filePath != "" {
			// Validate the file exists before attempting an upload.
			if _, err := os.Stat(filePath); err != nil {
				return mcp.NewToolResultError(
					fmt.Sprintf("file not found: %s", filePath)), nil
			}
			upload, err = client.UploadDraftFile(ctx, filePath)
			if err != nil {
				return mcp.NewToolResultError(
					fmt.Sprintf("upload file to Moodle draft area: %v", err)), nil
			}
		}

		// Step 2: Bind draft file and/or text to the assignment.
		if err := client.StageAssignmentDraft(ctx, assignmentID, upload.ItemID, textContent); err != nil {
			return mcp.NewToolResultError(
				fmt.Sprintf("stage assignment draft: %v", err)), nil
		}

		// Step 3: Return structured confirmation.
		// Use the review_url passed directly from module.url — never construct it
		// from assignment_id, which is the instance ID (not the course module ID).
		result := stageDraftResult{
			Status:          "draft_staged",
			AssignmentID:    assignmentID,
			Filename:        upload.Filename,
			DraftItemID:     upload.ItemID,
			MoodleReviewURL: reviewURL,
			Message: "File successfully uploaded and saved as a DRAFT. " +
				"It has NOT been submitted for grading. " +
				"Visit the review URL and click 'Submit assignment' to finalize.",
		}
		data, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	}
}

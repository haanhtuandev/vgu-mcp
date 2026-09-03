package moodle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// UploadDraftFile uploads a local file to Moodle's user draft area via a
// multipart/form-data POST to /webservice/upload.php.
//
// The returned UploadResponse contains the draft ItemID, which must be passed
// to StageAssignmentDraft to bind the file to an assignment submission.
func (c *Client) UploadDraftFile(ctx context.Context, filePath string) (UploadResponse, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return UploadResponse{}, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	// Build multipart body — Moodle expects the file under the field "file_1".
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	part, err := w.CreateFormFile("file_1", filepath.Base(filePath))
	if err != nil {
		return UploadResponse{}, fmt.Errorf("create multipart field: %w", err)
	}
	if _, err := io.Copy(part, f); err != nil {
		return UploadResponse{}, fmt.Errorf("write multipart data: %w", err)
	}
	w.Close()

	uploadURL := c.baseURL + "/webservice/upload.php?token=" + c.token
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, body)
	if err != nil {
		return UploadResponse{}, fmt.Errorf("create upload request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return UploadResponse{}, fmt.Errorf("upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UploadResponse{}, fmt.Errorf("upload http error: %s", resp.Status)
	}

	var uploads []UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploads); err != nil {
		return UploadResponse{}, fmt.Errorf("decode upload response: %w", err)
	}
	if len(uploads) == 0 || uploads[0].ItemID == 0 {
		return UploadResponse{}, fmt.Errorf("upload succeeded but no draft itemid returned")
	}
	return uploads[0], nil
}

package moodle

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Download streams a Moodle pluginfile to destinationDir using O(1) memory.
//
// If fileURL contains /pluginfile.php (but not /webservice/pluginfile.php), the
// path is rewritten so the token can be appended as a query parameter.
func (c *Client) Download(ctx context.Context, fileURL, destinationDir string) (string, error) {
	// Rewrite URL to webservice path if needed.
	if strings.Contains(fileURL, "/pluginfile.php") &&
		!strings.Contains(fileURL, "/webservice/pluginfile.php") {
		fileURL = strings.Replace(fileURL, "/pluginfile.php", "/webservice/pluginfile.php", 1)
	}

	// Append the auth token.
	sep := "?"
	if strings.Contains(fileURL, "?") {
		sep = "&"
	}
	authURL := fileURL + sep + "token=" + c.token

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, authURL, nil)
	if err != nil {
		return "", fmt.Errorf("create download request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("download request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download http error: %s", resp.Status)
	}

	if err := os.MkdirAll(destinationDir, 0o755); err != nil {
		return "", fmt.Errorf("create destination dir: %w", err)
	}

	// Derive filename from the URL path (before the query string).
	urlPath := strings.SplitN(fileURL, "?", 2)[0]
	filename := filepath.Base(urlPath)
	destPath := filepath.Join(destinationDir, filename)

	f, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("create file %s: %w", destPath, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", fmt.Errorf("stream file: %w", err)
	}
	return destPath, nil
}

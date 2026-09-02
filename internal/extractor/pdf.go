package extractor

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

// Result is the structured output of a PDF text extraction.
type Result struct {
	Filename   string `json:"filename"`
	TotalPages int    `json:"total_pages"`
	Content    string `json:"content"`
}

// ExtractPDF extracts plain text from raw PDF bytes and formats it as Markdown
// with a "## Page N" header before each page's content.
//
// It returns an error if the bytes are not a valid PDF. If a page has no
// extractable text (e.g. scanned image), that page is silently skipped.
// If no pages yield text, Content will be an empty string — callers should
// check for this and return a user-friendly message.
func ExtractPDF(data []byte, filename string) (*Result, error) {
	ra := bytes.NewReader(data)
	r, err := pdf.NewReader(ra, int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("open PDF: %w", err)
	}

	totalPages := r.NumPage()
	var sb strings.Builder

	for i := 1; i <= totalPages; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		trimmed := strings.TrimSpace(text)
		if trimmed == "" {
			continue
		}
		fmt.Fprintf(&sb, "## Page %d\n\n%s\n\n", i, trimmed)
	}

	return &Result{
		Filename:   filepath.Base(filename),
		TotalPages: totalPages,
		Content:    strings.TrimSpace(sb.String()),
	}, nil
}

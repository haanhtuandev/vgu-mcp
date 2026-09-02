package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/tuanha/vgu-mcp/internal/moodle"
	"github.com/tuanha/vgu-mcp/internal/tools"
)

func main() {

	// set the logging output as stderr because stdio is used as the communication pipe between MCP server and MCP client (opencode)
	log.SetOutput(os.Stderr)

	moodleURL, moodleToken := os.Getenv("MOODLE_URL"), os.Getenv("MOODLE_TOKEN")
	if moodleURL == "" || moodleToken == "" {
		log.Fatal("MOODLE_URL and MOODLE_TOKEN must be set")
	}

	mcpServer := server.NewMCPServer("VGU Moodle Server", "1.0.0", server.WithToolCapabilities(true))
	tools.RegisterCourseTools(mcpServer, moodle.NewClient(moodleURL, moodleToken))

	fmt.Fprintln(os.Stderr, "Moodle MCP Server running on stdio")
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

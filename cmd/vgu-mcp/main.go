package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/tuanha/vgu-mcp/internal/config"
	"github.com/tuanha/vgu-mcp/internal/moodle"
	"github.com/tuanha/vgu-mcp/internal/tools"
)

// version is injected at build time via -ldflags "-X main.version=<tag>".
var version = "dev"

func main() {
	log.SetOutput(os.Stderr)

	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(version)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "setup" {
		runSetup()
		return
	}

	runServer()
}

func runServer() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	if cfg.MoodleURL == "" || cfg.MoodleToken == "" {
		log.Fatal("no credentials configured — run `vgu-mcp setup` to get started")
	}

	client := moodle.NewClient(cfg.MoodleURL, cfg.MoodleToken)
	// Skip the GetSiteInfo network call on first tool invocation if we already
	// know the user ID from the saved config.
	if cfg.UserID != 0 {
		client.PreSeedUserID(cfg.UserID)
	}

	mcpServer := server.NewMCPServer("VGU MCP Server", version, server.WithToolCapabilities(true))
	tools.RegisterCourseTools(mcpServer, client)
	tools.RegisterDeadlineTools(mcpServer, client)
	tools.RegisterGradeTools(mcpServer, client)
	tools.RegisterMaterialTools(mcpServer, client)
	tools.RegisterExtractorTools(mcpServer, client)
	tools.RegisterAnnouncementTools(mcpServer, client)
	tools.RegisterAssignmentTools(mcpServer, client)

	fmt.Fprintln(os.Stderr, "VGU MCP Server running on stdio")
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

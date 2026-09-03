# Contributing

`vgu-mcp` is built by VGU students, for VGU students. Contributions are welcome — whether that's a new tool, a bug fix, or better documentation.

---

## Project structure

```
vgu-mcp/
├── cmd/vgu-mcp/
│   ├── main.go              # Entry point & subcommand routing (setup vs server)
│   └── setup.go             # Interactive one-time credential setup
├── internal/
│   ├── config/
│   │   └── config.go        # Load() and Save() config from env or ~/.config/vgu-mcp/
│   ├── extractor/
│   │   └── pdf.go           # Pure-Go PDF text extraction via ledongthuc/pdf
│   ├── moodle/
│   │   ├── auth.go          # Login() via /login/token.php
│   │   ├── client.go        # Moodle Web Services HTTP client + all API methods
│   │   ├── downloader.go    # Download() to disk, Fetch() to memory, pluginfileAuthURL()
│   │   ├── types.go         # All Moodle API request/response types
│   │   └── uploader.go      # UploadDraftFile() multipart upload to draft area
│   └── tools/
│       ├── announcements.go # read_course_announcements
│       ├── assignments.go   # stage_assignment_draft
│       ├── courses.go       # get_enrolled_courses, get_course_contents
│       ├── deadlines.go     # get_upcoming_deadlines
│       ├── extractor.go     # extract_course_material_text
│       ├── grades.go        # get_course_grades
│       └── materials.go     # download_course_material
├── docs/                    # User-facing documentation
├── go.mod
├── go.sum
└── Makefile
```

---

## Adding a new tool

Each tool is ~40–80 lines of Go. The pattern is consistent across all existing tools — follow it and yours will fit right in.

### Step 1 — Add a Moodle API method (if needed)

If your tool calls a Moodle Web Service function that doesn't exist yet, add it to `internal/moodle/client.go`:

```go
// GetSomething returns ... for the given parameters.
func (c *Client) GetSomething(ctx context.Context, id int) ([]Something, error) {
    params := url.Values{"id": {fmt.Sprintf("%d", id)}}
    var result []Something
    if err := c.request(ctx, "mod_something_get_data", params, &result); err != nil {
        return nil, err
    }
    return result, nil
}
```

Add any new response types to `internal/moodle/types.go`.

### Step 2 — Create the tool file

Create `internal/tools/something.go`:

```go
package tools

import (
    "context"
    "fmt"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
    "github.com/tuanha/vgu-mcp/internal/moodle"
)

func RegisterSomethingTools(s *server.MCPServer, client *moodle.Client) {
    s.AddTool(mcp.NewTool(
        "get_something",
        mcp.WithDescription("What this tool does, in one sentence."),
        mcp.WithNumber("id", mcp.Required(), mcp.Description("The something ID")),
    ), getSomething(client))
}

func getSomething(client *moodle.Client) server.ToolHandlerFunc {
    return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        id, err := req.RequireInt("id")
        if err != nil {
            return mcp.NewToolResultError("missing required parameter 'id'"), nil
        }

        data, err := client.GetSomething(ctx, id)
        if err != nil {
            return mcp.NewToolResultError(fmt.Sprintf("get something: %v", err)), nil
        }
        return jsonResult(data) // defined in courses.go
    }
}
```

> **Error handling rule:** always return `(mcp.NewToolResultError(...), nil)` for errors — never `(nil, err)`. The MCP protocol expects the error to be in the result, not the Go error return.

### Step 3 — Register the tool

In `cmd/vgu-mcp/main.go`, add one line inside `runServer()`:

```diff
 tools.RegisterMaterialTools(mcpServer, client)
+tools.RegisterSomethingTools(mcpServer, client)
 tools.RegisterAnnouncementTools(mcpServer, client)
```

### Step 4 — Test it

```bash
make build
make test
./vgu-mcp   # start the server and test via your AI client
```

---

## Makefile targets

```bash
make build    # compile binary for the current platform
make run      # build and run (for local testing)
make test     # run all tests
make tidy     # tidy go.mod / go.sum
make install  # install binary to /usr/local/bin/vgu-mcp
make dist     # cross-compile for all 5 platforms (output: dist/)
make clean    # remove dist/ and binary
```

### Distribution targets

`make dist` produces 5 binaries:

| File | Platform |
|---|---|
| `vgu-mcp-*-linux-amd64/vgu-mcp` | Linux x86_64 |
| `vgu-mcp-*-linux-arm64/vgu-mcp` | Linux ARM (Raspberry Pi, cloud) |
| `vgu-mcp-*-darwin-amd64/vgu-mcp` | macOS Intel |
| `vgu-mcp-*-darwin-arm64/vgu-mcp` | macOS Apple Silicon |
| `vgu-mcp-*-windows-amd64/vgu-mcp.exe` | Windows x86_64 |

---

## Ideas for new tools

Things that would make `vgu-mcp` significantly more useful:

- **`get_assignment_details`** — fetch submission instructions, due date, and grading rubric for a specific assignment
- **`get_quiz_attempts`** — list past quiz attempts and scores
- **`send_message`** — send a Moodle message to a course participant or instructor
- **`get_forum_post`** — read a specific discussion thread in full
- **DOCX/PPTX extraction** — extend `extract_course_material_text` using a pure-Go DOCX reader

---

## Pull requests

1. Fork the repo
2. Create a branch: `git checkout -b feat/my-new-tool`
3. Make your changes and ensure `make test` passes
4. Open a pull request with a clear description of what the tool does and why it's useful

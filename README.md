# VGU MCP Server

> **Let your AI assistant talk to Moodle.** Ask OpenCode, Claude, or Cursor "what assignments do I have?" and get real answers from VGU's Moodle — no copy-pasting, no browser digging.

---

## What is this?

`vgu-mcp` is a [Model Context Protocol](https://modelcontextprotocol.io) server that connects AI coding assistants to VGU's Moodle platform. Once set up, your AI can:

- 📚 List your enrolled courses
- 📋 Fetch course content, sections, and materials
- ⏰ Check upcoming assignment deadlines
- 📊 View your grades
- 📢 Read course announcements
- 📥 Download lecture files directly
- 📝 Submit assignment drafts

---

## Quick Start

### 1. Download the binary

Grab the latest release for your platform from [Releases](https://github.com/tuanha/vgu-mcp/releases):

| Platform | File |
|---|---|
| macOS (Apple Silicon M1/M2/M3) | `vgu-mcp-*-darwin-arm64.zip` |
| macOS (Intel) | `vgu-mcp-*-darwin-amd64.zip` |
| Linux (x86_64) | `vgu-mcp-*-linux-amd64.zip` |
| Linux (ARM) | `vgu-mcp-*-linux-arm64.zip` |
| Windows | `vgu-mcp-*-windows-amd64.zip` |

Or build from source (requires [Go 1.21+](https://go.dev/dl/)):

```bash
git clone https://github.com/tuanha/vgu-mcp
cd vgu-mcp
make build
```

### 2. Run setup (one-time)

```bash
./vgu-mcp setup
```

You'll be prompted for your VGU student credentials:

```
Moodle URL [https://moodle.vgu.edu.vn]: ↵
Username: 10423118
Password: ••••••••

✓ Login successful. Welcome, Ha Anh Tuan!
✓ Credentials saved to ~/.config/vgu-mcp/config.json
```

> **Your password is never saved.** Only the Moodle web service token is stored locally at `~/.config/vgu-mcp/config.json`.

If automatic login fails (e.g. your account uses SSO), you'll be prompted to paste your token manually:
> Moodle → Profile → Preferences → **Security keys** → copy your Web Service token

### 3. Install (optional but recommended)

```bash
make install
# or manually:
sudo cp vgu-mcp /usr/local/bin/vgu-mcp
```

This lets AI clients find the binary by name without a full path.

---

## Connect to Your AI Client

### OpenCode

Add this to `~/.config/opencode/opencode.json`:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "vgu-moodle": {
      "type": "local",
      "command": ["vgu-mcp"]
    }
  }
}
```

### Claude Desktop

Add this to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "vgu-moodle": {
      "command": "vgu-mcp"
    }
  }
}
```

### Cursor / Other MCP clients

Use `vgu-mcp` (or the full path if not installed) as the command. No arguments, no env vars needed.

---

## Available Tools

Once connected, your AI can use these tools automatically:

| Tool | What it does |
|---|---|
| `get_enrolled_courses` | Lists all your enrolled Moodle courses |
| `get_course_contents` | Fetches sections, modules, and files for a course |
| `get_upcoming_deadlines` | Shows upcoming assignment deadlines |
| `get_course_grades` | Retrieves your grades for a course |
| `read_course_announcements` | Reads the announcement forum for a course |
| `download_course_material` | Downloads a lecture file to your local machine |
| `submit_assignment_draft` | Submits online-text content for an assignment |

**Example prompts to try:**
- *"What courses am I enrolled in this semester?"*
- *"Show me all the materials from my Distributed Systems course"*
- *"Do I have any assignments due this week?"*
- *"What grades do I have in Operating Systems?"*
- *"Download the lecture slides from week 3 of my OS course"*

---

## Configuration

Credentials are stored in `~/.config/vgu-mcp/config.json` after running `vgu-mcp setup`:

```json
{
  "moodle_url": "https://moodle.vgu.edu.vn",
  "moodle_token": "your_token_here",
  "userid": 2197
}
```

You can also configure via environment variables (useful for CI or advanced setups):

```bash
export MOODLE_URL="https://moodle.vgu.edu.vn"
export MOODLE_TOKEN="your_token_here"
vgu-mcp
```

Environment variables take precedence over the config file.

---

## Re-authenticate

If your token expires or you get auth errors, just re-run setup:

```bash
vgu-mcp setup
```

---

## Build from Source

```bash
git clone https://github.com/tuanha/vgu-mcp
cd vgu-mcp

make build    # build binary
make test     # run tests
make install  # install to /usr/local/bin
make dist     # cross-compile for all platforms (output: dist/)
make tidy     # tidy go.mod
```

---

## Project Structure

```
vgu-mcp/
├── cmd/vgu-mcp/
│   ├── main.go          # Entry point & subcommand routing
│   └── setup.go         # Interactive credential setup
├── internal/
│   ├── config/
│   │   └── config.go    # Load/Save config from env or ~/.config/vgu-mcp/
│   ├── moodle/
│   │   ├── auth.go      # Login via login/token.php
│   │   ├── client.go    # Moodle Web Services HTTP client
│   │   ├── downloader.go # Streaming file downloader
│   │   └── types.go     # Moodle API types
│   └── tools/
│       ├── courses.go       # get_enrolled_courses, get_course_contents
│       ├── deadlines.go     # get_upcoming_deadlines
│       ├── grades.go        # get_course_grades
│       ├── materials.go     # download_course_material
│       ├── announcements.go # read_course_announcements
│       └── assignments.go   # submit_assignment_draft
└── Makefile
```

---

## Contributing

This project was made by VGU students, for VGU students. If you want to add more tools (SIS integration, timetable, etc.), contributions are very welcome.

1. Fork the repo
2. Add your tool in `internal/tools/`
3. Register it in `cmd/vgu-mcp/main.go`
4. Open a pull request

See the existing tools as reference — each one is just ~40 lines of Go.

---

## License

MIT

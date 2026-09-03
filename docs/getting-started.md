# Getting Started

This guide walks you through installing `vgu-mcp`, authenticating with your VGU Moodle account, and connecting it to your AI client.

---

## Requirements

- A VGU student account with access to [moodle.vgu.edu.vn](https://moodle.vgu.edu.vn)
- One of: [OpenCode](https://opencode.ai), [Claude Desktop](https://claude.ai/download), or [Cursor](https://cursor.com)
- macOS, Linux, or Windows (x86_64)

---

## 1. Download the binary

Grab the latest release for your platform from [GitHub Releases](https://github.com/haanhtuandev/vgu-mcp/releases):

| Platform | File |
|---|---|
| macOS (Apple Silicon M1/M2/M3) | `vgu-mcp-*-darwin-arm64.zip` |
| macOS (Intel) | `vgu-mcp-*-darwin-amd64.zip` |
| Linux (x86_64) | `vgu-mcp-*-linux-amd64.zip` |
| Linux (ARM) | `vgu-mcp-*-linux-arm64.zip` |
| Windows | `vgu-mcp-*-windows-amd64.zip` |

Unzip and place the binary somewhere convenient (e.g. your home directory or project folder).

### Build from source

If you prefer to build yourself, you need [Go 1.21+](https://go.dev/dl/):

```bash
git clone https://github.com/haanhtuandev/vgu-mcp
cd vgu-mcp
make build
```

---

## 2. Authenticate (one-time)

Run the setup command:

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

> **Your password is never saved.** Only a Moodle web service token is stored at `~/.config/vgu-mcp/config.json`.

### If automatic login fails (SSO accounts)

Some VGU accounts use SSO and can't authenticate via the standard login endpoint. In that case, setup will show a fallback prompt:

```
Automatic login failed: ...
Please paste your Moodle web service token instead.
(Moodle → Profile → Preferences → Security keys)
Token: _
```

Go to Moodle in your browser → click your profile photo → **Preferences** → **Security keys** → copy your Web Service token and paste it.

---

## 3. Install globally (optional but recommended)

Installing to `/usr/local/bin` lets AI clients find the binary by name without a full path:

```bash
make install
# or manually:
sudo cp vgu-mcp /usr/local/bin/vgu-mcp
```

On Windows, add the folder containing `vgu-mcp.exe` to your `PATH`.

---

## 4. Connect to your AI client

### OpenCode

Add to `~/.config/opencode/opencode.json`:

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

If you haven't installed globally, use the full path:

```json
"command": ["/Users/yourname/vgu-mcp"]
```

### Claude Desktop

**macOS** — edit `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "vgu-moodle": {
      "command": "vgu-mcp"
    }
  }
}
```

**Windows** — edit `%APPDATA%\Claude\claude_desktop_config.json` with the same content (use `vgu-mcp.exe` if needed).

Restart Claude Desktop after saving.

### Cursor

Open **Settings → Features → MCP → Add new MCP server**:

- **Name:** `vgu-moodle`
- **Type:** `command`
- **Command:** `vgu-mcp`

### Other MCP clients

Any client that supports the Model Context Protocol can use `vgu-mcp` as a local stdio server. Set the command to `vgu-mcp` (or full path). No arguments, no environment variables required.

---

## 5. Verify it's working

In your AI chat, type:

> *"What courses am I enrolled in on Moodle?"*

The AI should call `get_enrolled_courses` and return your course list. If it doesn't, check:

1. The binary is in your `PATH` or the full path is correct in the config
2. `~/.config/vgu-mcp/config.json` exists and contains a valid token (re-run `vgu-mcp setup` if needed)
3. Your AI client was restarted after adding the MCP config

---

## Re-authenticating

If your token expires or you see authentication errors, just run setup again:

```bash
vgu-mcp setup
```

This overwrites the existing token without touching any other settings.

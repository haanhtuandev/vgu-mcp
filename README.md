# VGU MCP

[![Go](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![MCP](https://img.shields.io/badge/MCP-compatible-6f42c1?style=flat)](https://modelcontextprotocol.io)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat)](LICENSE)
[![Release](https://img.shields.io/github/v/release/haanhtuandev/vgu-mcp?style=flat)](https://github.com/haanhtuandev/vgu-mcp/releases)

> **Your AI assistant, connected to VGU Moodle.**
> Ask it about deadlines, read lecture PDFs, check grades, stage assignment files —
> all without opening a browser.

---

## What you can do

- Ask about upcoming deadlines and grades in plain English
- Browse and download course materials
- Read PDF lecture notes — extracted and summarised instantly, no extra tools needed
- Stage assignment files as drafts on Moodle for your review before submitting

---

## Get started in 3 steps

**1. Download** the binary for your platform from
[Releases →](https://github.com/haanhtuandev/vgu-mcp/releases)

**2. Authenticate** once with your VGU credentials:
```bash
./vgu-mcp setup
```
Your password is never saved — only a web service token is stored locally.

**3. Connect** to your AI client:

<details>
<summary><strong>OpenCode</strong></summary>

Add to `~/.config/opencode/opencode.json`:
```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "vgu-moodle": { "type": "local", "command": ["vgu-mcp"] }
  }
}
```
</details>

<details>
<summary><strong>Claude Desktop</strong></summary>

Add to `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or
`%APPDATA%\Claude\claude_desktop_config.json` (Windows):
```json
{
  "mcpServers": {
    "vgu-moodle": { "command": "vgu-mcp" }
  }
}
```
</details>

<details>
<summary><strong>Cursor / other MCP clients</strong></summary>

Set the command to `vgu-mcp`. No arguments or environment variables required.
</details>

→ See the [Getting Started guide](docs/getting-started.md) for full setup instructions and troubleshooting.

---

## Try these prompts

Once connected, just type naturally in your AI chat:

- *"What assignments do I have due this week?"*
- *"Summarise the week 3 OS lecture slides for me"*
- *"Show me my grades for Distributed Systems"*
- *"Stage my lab2.zip to the component design assignment as a draft"*
- *"Quiz me on the Lab 2 specification PDF"*

---

## Docs

| | |
|---|---|
| [Getting Started](docs/getting-started.md) | Download, setup, AI client configuration |
| [Tool Reference](docs/tools.md) | All tools, parameters, and example prompts |
| [Configuration](docs/configuration.md) | Config file, environment variables, re-authentication |
| [Contributing](docs/contributing.md) | Project structure, adding tools, building from source |

---

MIT License — made by VGU students, for VGU students.

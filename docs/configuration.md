# Configuration

`vgu-mcp` reads credentials from a local config file or environment variables. No tokens or passwords are ever required inside your AI client's config.

---

## Config file

After running `vgu-mcp setup`, credentials are saved to:

```
~/.config/vgu-mcp/config.json
```

**Contents:**

```json
{
  "moodle_url": "https://moodle.vgu.edu.vn",
  "moodle_token": "abc123...",
  "userid": 2197
}
```

| Field | Description |
|---|---|
| `moodle_url` | Base URL of the Moodle instance |
| `moodle_token` | Web service token obtained during setup — this is what authenticates API calls |
| `userid` | Your Moodle user ID, cached at setup to avoid a network call on each server start |

**File permissions:** `0600` (owner read/write only). The token is sensitive — treat it like a password.

---

## Environment variables

You can override the config file with environment variables. This is useful for CI pipelines or shared machines where you don't want a config file on disk.

| Variable | Maps to |
|---|---|
| `MOODLE_URL` | `moodle_url` |
| `MOODLE_TOKEN` | `moodle_token` |

Environment variables take **precedence** over the config file. If both `MOODLE_URL` and `MOODLE_TOKEN` are set, the config file is not read at all.

```bash
export MOODLE_URL="https://moodle.vgu.edu.vn"
export MOODLE_TOKEN="your_token_here"
vgu-mcp
```

---

## Re-authentication

Moodle web service tokens do not have a fixed expiry, but if you see authentication errors, re-run setup:

```bash
vgu-mcp setup
```

This overwrites the existing token in `~/.config/vgu-mcp/config.json` without touching any other settings.

---

## Security

- **Your password is never stored.** The setup command exchanges your password for a token via `POST /login/token.php` and immediately discards the password.
- The token stored in `config.json` has the same access level as your student Moodle account — it can read course data and submit to assignments you are enrolled in.
- If you suspect your token is compromised, go to **Moodle → Profile → Preferences → Security keys** and reset it, then re-run `vgu-mcp setup`.

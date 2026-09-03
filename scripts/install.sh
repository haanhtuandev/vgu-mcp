#!/bin/sh
#
# vgu-mcp installer for macOS and Linux.
#
# Downloads the release binary for your platform, verifies its SHA-256
# checksum, and installs it to ~/.local/bin (or $VGU_MCP_INSTALL_DIR).
# No sudo required.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/haanhtuandev/vgu-mcp/main/scripts/install.sh | sh
#
# Optional environment variables:
#   VGU_MCP_VERSION      Install a specific release tag, e.g. v0.1.2 (default: latest)
#   VGU_MCP_INSTALL_DIR  Directory to install into (default: ~/.local/bin)
#   VGU_MCP_REPO         GitHub repo, for testing (default: haanhtuandev/vgu-mcp)
#
# It is safe to re-run: the script overwrites the binary with the newest
# release and is idempotent.

set -eu

REPO="${VGU_MCP_REPO:-haanhtuandev/vgu-mcp}"
INSTALL_DIR="${VGU_MCP_INSTALL_DIR:-"$HOME/.local/bin"}"
API_URL="${VGU_MCP_API_URL:-https://api.github.com/repos/$REPO/releases/latest}"
BASE_URL="${VGU_MCP_BASE_URL:-https://github.com/$REPO/releases/download}"

log() { printf '\033[1;32m%s\033[0m\n' "$*"; }
die() { printf '\033[1;31m%s\033[0m\n' "$*" >&2; exit 1; }

# --- resolve version ---------------------------------------------------------

if [ -n "${VGU_MCP_VERSION:-}" ]; then
  TAG="${VGU_MCP_VERSION#v}"
  TAG="v$TAG"
else
  log "Resolving latest release for $REPO ..."
  TAG="$(curl -fsSL "$API_URL" | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1)"
  [ -n "$TAG" ] || die "Could not determine the latest release from $API_URL"
fi
VERSION="${TAG#v}"
log "Installing vgu-mcp $VERSION"

# --- detect platform ---------------------------------------------------------

OS="$(uname -s)"
case "$OS" in
  Darwin) os=darwin ;;
  Linux)  os=linux ;;
  *) die "Unsupported operating system: $OS (only macOS and Linux are supported by this script)" ;;
esac

ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64) arch=amd64 ;;
  arm64|aarch64) arch=arm64 ;;
  *) die "Unsupported architecture: $ARCH" ;;
esac

ASSET="vgu-mcp_${VERSION}_${os}_${arch}.tar.gz"

# --- download + verify -------------------------------------------------------

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT INT TERM

log "Downloading $ASSET ..."
curl -fL --retry 3 --progress-bar -o "$TMP/$ASSET" "$BASE_URL/$TAG/$ASSET"

log "Verifying SHA-256 checksum ..."
curl -fsSL -o "$TMP/checksums.txt" "$BASE_URL/$TAG/checksums.txt"
EXPECTED="$(awk -v f="$ASSET" '$2 == f { print $1 }' "$TMP/checksums.txt")"
[ -n "$EXPECTED" ] && [ "$EXPECTED" != "null" ] \
  || die "No checksum found for $ASSET in checksums.txt"

if command -v shasum >/dev/null 2>&1; then
  ACTUAL="$(shasum -a 256 "$TMP/$ASSET" | awk '{ print $1 }')"
elif command -v sha256sum >/dev/null 2>&1; then
  ACTUAL="$(sha256sum "$TMP/$ASSET" | awk '{ print $1 }')"
else
  die "Missing checksum tool: install shasum or coreutils (sha256sum) and re-run"
fi

[ "$ACTUAL" = "$EXPECTED" ] || die "Checksum mismatch for $ASSET (expected $EXPECTED, got $ACTUAL)"

# --- install ---------------------------------------------------------------

mkdir -p "$INSTALL_DIR"
tar -xzf "$TMP/$ASSET" -C "$TMP" "vgu-mcp"
DEST="$INSTALL_DIR/vgu-mcp"
install -m 0755 "$TMP/vgu-mcp" "$DEST"

if [ "$os" = "darwin" ] && command -v xattr >/dev/null 2>&1; then
  # Binary is not notarized; clear the quarantine flag so Gatekeeper lets it run.
  xattr -dr com.apple.quarantine "$DEST" 2>/dev/null || true
fi

log "Installed $DEST"

# --- PATH -------------------------------------------------------------------

case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    rc="${RC_FILE:-}"
    if [ -z "$rc" ]; then
      for f in "$HOME/.zshrc" "$HOME/.bashrc" "$HOME/.profile"; do
        [ -f "$f" ] && rc="$f" && break
      done
    fi
    if [ -n "$rc" ]; then
      if ! grep -q "vgu-mcp installer" "$rc" 2>/dev/null; then
        {
          printf '\n# >>> vgu-mcp installer >>>\nexport PATH="%s:$PATH"\n# <<< vgu-mcp installer <<<\n' "$INSTALL_DIR"
        } >>"$rc"
      fi
      log "Added $INSTALL_DIR to PATH in $rc (restart your shell or run: export PATH=\"$INSTALL_DIR:\$PATH\")"
    else
      printf '\033[1;33mAdd %s to your PATH, then run:\n  export PATH="%s:$PATH"\n\033[0m\n' "$INSTALL_DIR" "$INSTALL_DIR"
    fi
    ;;
esac

# --- done ------------------------------------------------------------------

log "vgu-mcp $VERSION installed."
printf '\nNext steps:\n'
printf '  1. Authenticate once:  vgu-mcp setup\n'
printf '  2. Add it to your AI client (see https://github.com/%s#readme)\n\n' "$REPO"
printf 'To upgrade later, just re-run this command.\n'
printf 'To uninstall: rm "%s"\n' "$DEST"
if [ -n "${rc:-}" ]; then
  printf '  and remove the "vgu-mcp installer" block from %s\n' "$rc"
fi
printf '\n'

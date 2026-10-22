#!/bin/sh
# Lipi installer for macOS and Linux.
#
#   curl -fsSL https://raw.githubusercontent.com/RusilKoirala/lipi/master/install.sh | sh
#
# Downloads the right lipi binary from GitHub Releases, checks it, puts it in
# ~/.lipi/bin and adds that folder to your PATH. Run it again to update.
#
# Options (environment variables):
#   LIPI_VERSION=v0.1.0   install a specific version instead of the latest
#   LIPI_INSTALL=~/.lipi  install somewhere else

set -eu

REPO="RusilKoirala/lipi"
VERSION="${LIPI_VERSION:-latest}"
LIPI_HOME="${LIPI_INSTALL:-$HOME/.lipi}"
BIN_DIR="$LIPI_HOME/bin"

# ---------- output ----------

if [ -t 1 ]; then
  BOLD="$(printf '\033[1m')"; DIM="$(printf '\033[2m')"; RED="$(printf '\033[31m')"
  GREEN="$(printf '\033[32m')"; RESET="$(printf '\033[0m')"
else
  BOLD=""; DIM=""; RED=""; GREEN=""; RESET=""
fi

info() { printf '%s\n' "$*"; }
# Show paths under your home folder as ~/...
tilde() { case "$1" in "$HOME"/*) printf '~%s' "${1#"$HOME"}" ;; *) printf '%s' "$1" ;; esac; }
step() { printf '%s›%s %s\n' "$DIM" "$RESET" "$*"; }
fail() {
  printf '%sGalti:%s %s\n' "$RED" "$RESET" "$*" >&2
  exit 1
}

# ---------- which binary? ----------

case "$(uname -s)" in
  Darwin) os="darwin" ;;
  Linux) os="linux" ;;
  MINGW* | MSYS* | CYGWIN*) fail "Windows ko lagi lipi-setup.exe chalaunus: https://github.com/$REPO/releases/latest" ;;
  *) fail "yo system ($(uname -s)) ko lagi lipi banaieko chhaina" ;;
esac

case "$(uname -m)" in
  x86_64 | amd64) arch="amd64" ;;
  arm64 | aarch64) arch="arm64" ;;
  *) fail "yo CPU ($(uname -m)) ko lagi lipi banaieko chhaina" ;;
esac

# An Intel terminal running under Rosetta on an Apple Silicon Mac: use the native build.
if [ "$os" = "darwin" ] && [ "$arch" = "amd64" ]; then
  if [ "$(sysctl -n sysctl.proc_translated 2>/dev/null || echo 0)" = "1" ]; then
    arch="arm64"
  fi
fi

asset="lipi-$os-$arch"

if [ -n "${LIPI_DOWNLOAD_BASE:-}" ]; then
  base="$LIPI_DOWNLOAD_BASE" # used for testing against a local server
elif [ "$VERSION" = "latest" ]; then
  base="https://github.com/$REPO/releases/latest/download"
else
  base="https://github.com/$REPO/releases/download/$VERSION"
fi

# ---------- download ----------

download() {
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL --retry 2 -o "$2" "$1"
  elif command -v wget >/dev/null 2>&1; then
    wget -q -O "$2" "$1"
  else
    fail "curl wa wget chahincha"
  fi
}

sha256() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | cut -d ' ' -f 1
  elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | cut -d ' ' -f 1
  else
    fail "sha256sum wa shasum chahincha"
  fi
}

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT INT TERM

info ""
info "${BOLD}Lipi${RESET} install hudai chha ($os/$arch)"

step "Downloading $asset"
download "$base/$asset" "$tmp/lipi" || fail "download bhayena: $base/$asset"
download "$base/checksums.txt" "$tmp/checksums.txt" || fail "checksums.txt download bhayena"

step "Checking the download"
expected="$(grep " $asset\$" "$tmp/checksums.txt" | cut -d ' ' -f 1 || true)"
[ -n "$expected" ] || fail "checksums.txt ma $asset bhetiyena"
actual="$(sha256 "$tmp/lipi")"
[ "$expected" = "$actual" ] || fail "download bigriyo (checksum milena). Feri try garnus."

# ---------- install ----------

step "Installing to $(tilde "$BIN_DIR")"
mkdir -p "$BIN_DIR"
chmod +x "$tmp/lipi"
mv "$tmp/lipi" "$BIN_DIR/lipi"

installed="$("$BIN_DIR/lipi" --version 2>/dev/null)" || fail "lipi install bhayo tara chalena"

# ---------- PATH ----------

case ":$PATH:" in
  *":$BIN_DIR:"*) on_path=1 ;;
  *) on_path=0 ;;
esac

rc=""
if [ "$on_path" = "0" ]; then
  shell_name="$(basename "${SHELL:-sh}")"
  case "$shell_name" in
    zsh)
      rc="${ZDOTDIR:-$HOME}/.zshrc"
      line="export PATH=\"$BIN_DIR:\$PATH\""
      ;;
    bash)
      if [ "$os" = "darwin" ]; then rc="$HOME/.bash_profile"; else rc="$HOME/.bashrc"; fi
      line="export PATH=\"$BIN_DIR:\$PATH\""
      ;;
    fish)
      rc="$HOME/.config/fish/config.fish"
      line="fish_add_path \"$BIN_DIR\""
      ;;
    *)
      rc="$HOME/.profile"
      line="export PATH=\"$BIN_DIR:\$PATH\""
      ;;
  esac

  mkdir -p "$(dirname "$rc")"
  touch "$rc"
  if grep -qF "$BIN_DIR" "$rc"; then
    step "PATH is already set in $(tilde "$rc")"
  else
    step "Adding lipi to your PATH in $(tilde "$rc")"
    printf '\n# lipi\n%s\n' "$line" >>"$rc"
  fi
fi

# ---------- done ----------

info ""
info "${GREEN}✓${RESET} ${BOLD}$installed${RESET} install bhayo."
info ""
if [ -n "$rc" ]; then
  info "  Open a new terminal, or run:"
  info "    ${BOLD}source $(tilde "$rc")${RESET}"
  info ""
fi
info "  Then try:"
info "    ${BOLD}lipi --version${RESET}"
info "    ${BOLD}lipi namaste.lip${RESET}"
info ""
info "  ${DIM}Docs: https://github.com/$REPO${RESET}"
info ""

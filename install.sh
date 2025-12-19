#!/usr/bin/env sh
set -euf

GITHUB_REPO="rickyseezy/gclone"
VERSION="${GCLONE_VERSION:-}"
INSTALL_DIR="${GCLONE_INSTALL_DIR:-/usr/local/bin}"

if [ -z "$VERSION" ]; then
  VERSION="$(curl -fsSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | awk -F '"' '/tag_name/{print $4; exit}')"
fi

VERSION="${VERSION#v}"

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Darwin) OS="darwin" ;;
  Linux) OS="linux" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
 esac

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
 esac

ASSET="gclone_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${GITHUB_REPO}/releases/download/v${VERSION}/${ASSET}"

TMP_DIR="$(mktemp -d)"
cleanup() { rm -rf "$TMP_DIR"; }
trap cleanup EXIT

curl -fsSL "$URL" -o "$TMP_DIR/$ASSET"

tar -xzf "$TMP_DIR/$ASSET" -C "$TMP_DIR"

if [ ! -w "$INSTALL_DIR" ]; then
  echo "No write access to $INSTALL_DIR."
  echo "Re-run with: sudo GCLONE_INSTALL_DIR=$INSTALL_DIR sh install.sh"
  echo "Or choose a writable dir: GCLONE_INSTALL_DIR=~/.local/bin sh install.sh"
  exit 1
fi

install -m 0755 "$TMP_DIR/gclone" "$INSTALL_DIR/gclone"

echo "Installed gclone to $INSTALL_DIR/gclone"

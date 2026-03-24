#!/bin/sh
set -e

REPO="chrisfentiman/dbx"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="dbx"

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

case "$OS" in
  darwin|linux) ;;
  *) echo "Unsupported OS: $OS" && exit 1 ;;
esac

TARGET="dbx-${OS}-${ARCH}"
URL="https://github.com/${REPO}/releases/latest/download/${TARGET}.tar.gz"

echo "Downloading ${TARGET}..."
TMPDIR=$(mktemp -d)
curl -sL "$URL" | tar xz -C "$TMPDIR"

echo "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
sudo mv "$TMPDIR/$TARGET" "$INSTALL_DIR/$BINARY_NAME"
rm -rf "$TMPDIR"

echo "Done. Run 'dbx version' to verify."

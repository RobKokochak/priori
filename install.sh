#!/bin/bash

set -e

# Determine OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Convert architecture names
case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
esac

VERSION=$(curl -s https://api.github.com/repos/robkokochak/priori/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | sed 's/v//')
BINARY_NAME="priori_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/robkokochak/priori/releases/download/v${VERSION}/${BINARY_NAME}"

printf "Retrieving priori v%s for %s-%s...\n" "$VERSION" "$OS" "$ARCH"

TMP_DIR=$(mktemp -d)
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

curl -L -s "$DOWNLOAD_URL" -o "$TMP_DIR/${BINARY_NAME}"
tar -xzf "$TMP_DIR/${BINARY_NAME}" -C "$TMP_DIR"

sudo mv "$TMP_DIR/priori" /usr/local/bin/priori

echo "priori installed successfully"

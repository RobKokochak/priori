#!/bin/bash

set -e

OS=$(uname -s)
ARCH=$(uname -m)

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

if [ ! -f "$TMP_DIR/${BINARY_NAME}" ]; then
    echo "Error: Download failed"
    exit 1
fi

tar -xzf "$TMP_DIR/${BINARY_NAME}" -C "$TMP_DIR"

if [ ! -f "$TMP_DIR/priori" ]; then
    echo "Error: Extraction failed or binary not found"
    exit 1
fi

sudo mv "$TMP_DIR/priori" /usr/local/bin/priori
sudo chmod +x /usr/local/bin/priori

echo "priori installed successfully"

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

# Set version and download URL
VERSION=$(curl -s https://api.github.com/repos/robkokochak/priori/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | sed 's/v//')
BINARY_NAME="priori-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/robkokochak/priori/releases/download/v${VERSION}/${BINARY_NAME}"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

# Download and verify
echo "Downloading priori..."
curl -L "$DOWNLOAD_URL" -o "$TMP_DIR/$BINARY_NAME"
chmod +x "$TMP_DIR/$BINARY_NAME"

# Install
echo "Installing priori..."
sudo mv "$TMP_DIR/$BINARY_NAME" /usr/local/bin/priori

echo "priori installed successfully"

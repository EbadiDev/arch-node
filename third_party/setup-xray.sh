#!/bin/bash

# Script to setup xray binaries from latest releases
# This script downloads the latest xray binaries and extracts them

set -e

XRAY_VERSION=${1:-"latest"}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
THIRD_PARTY_DIR="$PROJECT_ROOT/third_party"

echo "Setting up Xray binaries..."

# Create third_party directory if it doesn't exist
mkdir -p "$THIRD_PARTY_DIR"

# Function to get latest release version
get_latest_version() {
    curl -s https://api.github.com/repos/XTLS/Xray-core/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

# Function to download and extract xray binary
download_xray() {
    local platform=$1
    local arch=$2
    local version=$3
    
    local filename="Xray-linux-64.zip"
    local extract_dir="xray-linux-64"
    
    if [[ "$platform" == "darwin" ]]; then
        filename="Xray-macos-arm64-v8a.zip"
        extract_dir="xray-macos-arm64"
    fi
    
    local download_url="https://github.com/XTLS/Xray-core/releases/download/${version}/${filename}"
    local zip_path="$THIRD_PARTY_DIR/$filename"
    local extract_path="$THIRD_PARTY_DIR/$extract_dir"
    
    echo "Downloading $filename version $version..."
    curl -L -f --retry 3 --retry-delay 1 -o "$zip_path" "$download_url" || {
        echo "Error: Failed to download $filename"
        echo "Please download manually from: $download_url"
        exit 1
    }
    
    echo "Extracting to $extract_path..."
    rm -rf "$extract_path"
    mkdir -p "$extract_path"
    unzip -q "$zip_path" -d "$extract_path"
    
    # Make xray executable
    chmod +x "$extract_path/xray"
    
    echo "Cleaning up zip file..."
    rm "$zip_path"
    
    # Create version file
    echo "$version" > "$extract_path/VERSION"
    
    echo "✓ $extract_dir setup complete"
}

# Get version to download
if [[ "$XRAY_VERSION" == "latest" ]]; then
    XRAY_VERSION=$(get_latest_version)
    echo "Latest version: $XRAY_VERSION"
fi

# Download Linux x64 binary (always needed)
download_xray "linux" "amd64" "$XRAY_VERSION"

# Download macOS ARM64 binary if requested
if [[ "${2:-}" == "all" ]] || [[ "${2:-}" == "macos" ]]; then
    download_xray "darwin" "arm64" "$XRAY_VERSION"
fi

echo "Xray binaries setup complete!"
echo "Version: $XRAY_VERSION"

# Update the version in config.go
CONFIG_FILE="$PROJECT_ROOT/internal/config/config.go"
if [[ -f "$CONFIG_FILE" ]]; then
    sed -i "s/const XrayCoreVersion = \"Xray v[^\"]*\"/const XrayCoreVersion = \"Xray $XRAY_VERSION\"/" "$CONFIG_FILE"
    echo "Updated XrayCoreVersion in config.go to: Xray $XRAY_VERSION"
fi

echo ""
echo "Available binaries:"
ls -la "$THIRD_PARTY_DIR"/xray-*/xray 2>/dev/null || echo "No xray binaries found"

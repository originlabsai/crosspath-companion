#!/bin/bash

# Cross-platform build script for CrossPath Companion
# Usage:
#   ./build.sh           Build production (all platforms)
#   ./build.sh dev       Build dev (local backend, current platform only)
#   ./build.sh prod      Build production (all platforms)
#   ./build.sh oss       Build OSS (localhost:3000 default, all platforms)

set -e

MODE="${1:-prod}"
PKG="github.com/crosspath/mcp-client/internal/config"
VERSION="0.1.5"

# Find Go binary — snap Go can fail with apparmor, use direct path as fallback
GO_BIN="go"
if ! command -v go &>/dev/null || ! go version &>/dev/null 2>&1; then
    if [ -x /snap/go/current/bin/go ]; then
        GO_BIN="/snap/go/current/bin/go"
    elif [ -x /usr/local/go/bin/go ]; then
        GO_BIN="/usr/local/go/bin/go"
    else
        echo "ERROR: Go not found"
        exit 1
    fi
fi

# Backend URLs
PROD_URL="wss://crosspath.im/mcp/connect"
DEV_URL="ws://localhost:3001/mcp/connect"
OSS_URL="ws://localhost:3000/mcp/connect"

if [ "$MODE" = "dev" ]; then
    BACKEND_URL="$DEV_URL"
    LDFLAGS="-X ${PKG}.DefaultBackendURL=${BACKEND_URL} -X main.Version=${VERSION}-dev"
    echo "Building CrossPath Companion (DEV → ${DEV_URL})"
    echo ""

    # Dev: build only for current platform
    echo "Building for current platform..."
    $GO_BIN build -ldflags "$LDFLAGS" -o cross_companion ./cmd/cross_companion
    echo "✓ Built ./cross_companion (dev)"
    echo ""
    echo "Run with: ./cross_companion"
    exit 0
fi

if [ "$MODE" = "oss" ]; then
    BACKEND_URL="$OSS_URL"
    LDFLAGS="-X ${PKG}.DefaultBackendURL=${BACKEND_URL} -X main.Version=${VERSION}-oss -s -w"
    echo "Building CrossPath Companion (OSS → ${OSS_URL})"
    echo ""
else
    # Production build
    BACKEND_URL="$PROD_URL"
    LDFLAGS="-X ${PKG}.DefaultBackendURL=${BACKEND_URL} -X main.Version=${VERSION} -s -w"
    echo "Building CrossPath Companion (PROD → ${PROD_URL})"
    echo ""
fi

# Create bin directory
mkdir -p bin

# Windows
echo "[1/5] Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 $GO_BIN build -ldflags "$LDFLAGS" -o bin/cross_companion-windows-amd64.exe ./cmd/cross_companion
echo "✓ Windows build complete"

# Linux
echo "[2/5] Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 $GO_BIN build -ldflags "$LDFLAGS" -o bin/cross_companion-linux-amd64 ./cmd/cross_companion
echo "✓ Linux build complete"

# macOS Intel
echo "[3/5] Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 $GO_BIN build -ldflags "$LDFLAGS" -o bin/cross_companion-darwin-amd64 ./cmd/cross_companion
echo "✓ macOS build complete"

# macOS Apple Silicon
echo "[4/5] Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 $GO_BIN build -ldflags "$LDFLAGS" -o bin/cross_companion-darwin-arm64 ./cmd/cross_companion
echo "✓ macOS ARM build complete"

# Linux ARM64
echo "[5/5] Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 $GO_BIN build -ldflags "$LDFLAGS" -o bin/cross_companion-linux-arm64 ./cmd/cross_companion
echo "✓ Linux ARM build complete"

echo ""
echo "========================================"
echo "✓ All builds completed successfully!"
echo "========================================"

# Generate checksums
echo ""
echo "Generating checksums..."
cd bin
if command -v sha256sum &>/dev/null; then
    sha256sum cross_companion-* > checksums.txt
elif command -v shasum &>/dev/null; then
    shasum -a 256 cross_companion-* > checksums.txt
else
    echo "⚠️  Neither sha256sum nor shasum found, skipping checksum generation"
fi
cd ..

if [ -f bin/checksums.txt ]; then
    echo "✓ Checksums generated"
    echo ""
    cat bin/checksums.txt
fi

# Copy installer scripts
echo ""
echo "Copying installer scripts..."
if [ -d scripts ]; then
    cp scripts/install.sh bin/ 2>/dev/null && echo "✓ Copied install.sh"
    cp scripts/install.ps1 bin/ 2>/dev/null && echo "✓ Copied install.ps1"
fi

echo ""
echo "========================================"
echo "Release artifacts ready in bin/:"
echo "========================================"
ls -lh bin/
echo ""
echo "To run:"
echo "  Windows: bin/cross_companion-windows-amd64.exe"
echo "  Linux:   bin/cross_companion-linux-amd64"
echo "  macOS:   bin/cross_companion-darwin-amd64"
echo "  macOS (M1/M2): bin/cross_companion-darwin-arm64"
echo ""

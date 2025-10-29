#!/bin/bash

# Script to verify cross-compiled binaries
# Tests basic executability and help command functionality

set -e

BINARY_DIR="bin"
BINARY_NAME="shantilly"
PLATFORMS=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")

echo "Verifying cross-compiled binaries..."
echo "===================================="

# Check if bin directory exists
if [ ! -d "$BINARY_DIR" ]; then
    echo "ERROR: Binary directory '$BINARY_DIR' does not exist"
    exit 1
fi

SUCCESS_COUNT=0
TOTAL_COUNT=${#PLATFORMS[@]}

for platform in "${PLATFORMS[@]}"; do
    echo ""
    echo "Checking $platform..."

    platform_dir="$BINARY_DIR/$platform"

    if [ ! -d "$platform_dir" ]; then
        echo "  ❌ Directory $platform_dir does not exist"
        continue
    fi

    # Determine binary name (add .exe for Windows)
    if [[ "$platform" == windows/* ]]; then
        binary_path="$platform_dir/$BINARY_NAME.exe"
    else
        binary_path="$platform_dir/$BINARY_NAME"
    fi

    if [ ! -f "$binary_path" ]; then
        echo "  ❌ Binary $binary_path does not exist"
        continue
    fi

    echo "  ✅ Binary exists: $binary_path"

    # Check file size
    size=$(stat -c%s "$binary_path" 2>/dev/null || stat -f%z "$binary_path" 2>/dev/null || echo "unknown")
    echo "  📏 Size: $size bytes"

    # Test if binary is executable (on current platform only)
    if [[ "$platform" == "linux/amd64" ]]; then
        echo "  🧪 Testing execution on current platform..."

        # Test help command
        if "$binary_path" --help >/dev/null 2>&1; then
            echo "  ✅ Help command works"
        else
            echo "  ❌ Help command failed"
            continue
        fi

        # Test basic execution (should show usage)
        if "$binary_path" >/dev/null 2>&1; then
            echo "  ✅ Basic execution works (shows usage)"
        else
            echo "  ❌ Basic execution failed"
            continue
        fi
    else
        echo "  ⏭️  Skipping execution test (different platform)"
    fi

    echo "  ✅ $platform verification passed"
    ((SUCCESS_COUNT++))
done

echo ""
echo "===================================="
echo "Verification Summary: $SUCCESS_COUNT/$TOTAL_COUNT platforms verified"
echo ""

if [ "$SUCCESS_COUNT" -eq "$TOTAL_COUNT" ]; then
    echo "🎉 All binaries verified successfully!"
    exit 0
else
    echo "⚠️  Some binaries failed verification"
    exit 1
fi

#!/bin/bash

# ComfyUI Go SDK - Progress Tracking Demo Runner
# This script demonstrates the progress tracking example

set -e

echo "╔════════════════════════════════════════════════════════════╗"
echo "║     ComfyUI Go SDK - Progress Tracking Demo Runner        ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Check if ComfyUI is running
echo "🔍 Checking if ComfyUI is running..."
if ! curl -s http://127.0.0.1:8188/system_stats > /dev/null 2>&1; then
    echo "❌ Error: ComfyUI is not running on http://127.0.0.1:8188"
    echo ""
    echo "Please start ComfyUI first:"
    echo "  cd /path/to/ComfyUI"
    echo "  python main.py"
    echo ""
    exit 1
fi

echo "✅ ComfyUI is running"
echo ""

# Build the example if needed
if [ ! -f "./bin/progress" ]; then
    echo "🔨 Building progress example..."
    make build
    echo ""
fi

# Run the example
echo "🚀 Running progress tracking demo..."
echo ""
./bin/progress

echo ""
echo "✨ Demo completed! Check the current directory for generated images."

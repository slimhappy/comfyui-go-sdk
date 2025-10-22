#!/bin/bash

# Quick start script for execute_from_json example

set -e

echo "╔════════════════════════════════════════════════════════════════════╗"
echo "║     ComfyUI Go SDK - Execute from JSON Quick Start                ║"
echo "╚════════════════════════════════════════════════════════════════════╝"
echo ""

# Check if ComfyUI is running
echo "🔍 Checking if ComfyUI is running..."
if ! curl -s http://127.0.0.1:8188/system_stats > /dev/null 2>&1; then
    echo "❌ ComfyUI is not running!"
    echo ""
    echo "Please start ComfyUI first:"
    echo "  cd /data/ComfyUI"
    echo "  python main.py"
    echo ""
    exit 1
fi
echo "✅ ComfyUI is running"
echo ""

# Build the example
echo "🔨 Building execute_from_json example..."
cd /data/comfyui-go-sdk
make build-execute-json
echo ""

# Check if workflow file exists
WORKFLOW_FILE="examples/execute_from_json/workflow_example.json"
if [ ! -f "$WORKFLOW_FILE" ]; then
    echo "❌ Example workflow file not found: $WORKFLOW_FILE"
    exit 1
fi

echo "📋 Available commands:"
echo ""
echo "1. Basic execution:"
echo "   ./bin/execute_from_json $WORKFLOW_FILE"
echo ""
echo "2. With custom seed:"
echo "   ./bin/execute_from_json $WORKFLOW_FILE seed=12345"
echo ""
echo "3. With custom parameters:"
echo "   ./bin/execute_from_json $WORKFLOW_FILE seed=42 steps=25 cfg=7.5"
echo ""
echo "4. With custom prompts:"
echo "   ./bin/execute_from_json $WORKFLOW_FILE prompt=\"beautiful landscape\" negative=\"blurry\""
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Ask user if they want to run the example
read -p "Would you like to run the basic example now? (y/n) " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "🚀 Running example..."
    echo ""
    ./bin/execute_from_json "$WORKFLOW_FILE"
else
    echo ""
    echo "👍 You can run the example manually using the commands above."
fi

echo ""
echo "📚 For more information, see: examples/execute_from_json/README.md"

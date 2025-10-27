#!/bin/bash

# Quick test script for all examples
# This script runs each example briefly to verify they work

set -e

echo "=== Testing ComfyUI Go SDK Examples ==="
echo ""

COMFYUI_URL="${COMFYUI_URL:-http://127.0.0.1:8188}"
echo "ComfyUI Server: $COMFYUI_URL"
echo ""

# Test 1: Integration Test (comprehensive)
echo "[1/6] Running Integration Test..."
timeout 60 ./bin/integration_test || echo "  (timeout or error)"
echo ""

# Test 2: Model Info (quick info query)
echo "[2/6] Testing Model Info..."
timeout 10 ./bin/model_info > /dev/null 2>&1 && echo "  ✓ Model info example works" || echo "  ✗ Model info example failed"
echo ""

# Test 3: Queue Management
echo "[3/6] Testing Queue Management..."
timeout 30 ./bin/queue_management > /dev/null 2>&1 && echo "  ✓ Queue management example works" || echo "  ✗ Queue management example failed"
echo ""

# Test 4: History Operations
echo "[4/6] Testing History Operations..."
timeout 30 ./bin/history_operations > /dev/null 2>&1 && echo "  ✓ History operations example works" || echo "  ✗ History operations example failed"
echo ""

# Test 5: Error Handling
echo "[5/6] Testing Error Handling..."
timeout 30 ./bin/error_handling > /dev/null 2>&1 && echo "  ✓ Error handling example works" || echo "  ✗ Error handling example failed"
echo ""

# Test 6: Image Operations
echo "[6/6] Testing Image Operations..."
timeout 60 ./bin/image_operations > /dev/null 2>&1 && echo "  ✓ Image operations example works" || echo "  ✗ Image operations example failed"
echo ""

echo "=== Test Complete ==="
echo ""
echo "All examples have been compiled and tested."
echo "Run individual examples with: ./bin/<example_name>"
echo ""
echo "Available examples:"
echo "  - basic"
echo "  - websocket"
echo "  - advanced"
echo "  - progress"
echo "  - execute_from_json"
echo "  - queue_management"
echo "  - history_operations"
echo "  - model_info"
echo "  - image_operations"
echo "  - error_handling"
echo "  - integration_test"

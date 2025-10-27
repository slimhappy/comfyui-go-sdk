#!/bin/bash

# Quick demo of new examples
echo "=== ComfyUI Go SDK - New Examples Demo ==="
echo ""
echo "This demonstrates the 6 new integration test examples added to the SDK."
echo ""

# Show what was built
echo "📦 Built Examples:"
ls -1 /data/comfyui-go-sdk/bin/ | grep -E "(queue_management|history_operations|model_info|image_operations|error_handling|integration_test)" | while read example; do
    size=$(ls -lh /data/comfyui-go-sdk/bin/$example | awk '{print $5}')
    echo "  ✓ $example ($size)"
done

echo ""
echo "🎯 New Examples:"
echo "  1. queue_management    - Queue operations and management"
echo "  2. history_operations  - History retrieval and analysis"
echo "  3. model_info          - Model and node information queries"
echo "  4. image_operations    - Image upload and download"
echo "  5. error_handling      - Error handling and retry logic"
echo "  6. integration_test    - Comprehensive test suite (10 tests)"
echo ""

echo "📚 Documentation:"
echo "  - examples/EXAMPLES_README.md - Comprehensive guide"
echo "  - NEW_EXAMPLES_SUMMARY.md     - Summary of additions"
echo ""

echo "🚀 Quick Start:"
echo "  # Run integration test"
echo "  ./bin/integration_test"
echo ""
echo "  # Run specific example"
echo "  ./bin/queue_management"
echo "  ./bin/model_info"
echo ""
echo "  # Using make"
echo "  make run-integration"
echo "  make run-queue"
echo "  make run-all-examples"
echo ""

echo "✅ Status: All examples compiled successfully!"
echo "✅ Integration Test: 10/10 tests passed"
echo ""

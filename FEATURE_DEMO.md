# ComfyUI Go SDK - Execute from JSON Feature

## 🎉 New Feature Added!

The ComfyUI Go SDK now supports **loading and executing workflows directly from JSON files**!

---

## ✨ Quick Demo

### 1. One-Line Execution

```go
client := comfyui.NewClient("http://127.0.0.1:8188")
resp, err := client.QueuePromptFromFile(context.Background(), "workflow.json", nil)
```

### 2. Command-Line Tool

```bash
# Build
make build-execute-json

# Run
./bin/execute_from_json workflow.json

# With parameters
./bin/execute_from_json workflow.json seed=12345 steps=30 cfg=7.5
```

---

## 📦 What's Included

### 1. SDK Enhancement

**New Method**: `Client.QueuePromptFromFile()`

```go
func (c *Client) QueuePromptFromFile(
    ctx context.Context,
    filepath string,
    extraData map[string]interface{}
) (*QueuePromptResponse, error)
```

### 2. Complete CLI Tool

**Location**: `examples/execute_from_json/`

**Features**:
- ✅ Load workflows from JSON files
- ✅ Modify parameters at runtime
- ✅ Real-time progress monitoring
- ✅ Automatic result retrieval
- ✅ Image download and save
- ✅ Comprehensive error handling

### 3. Example Workflow

**File**: `examples/execute_from_json/workflow_example.json`

Standard ComfyUI API format workflow ready to use!

### 4. Documentation

- **Main README**: Updated with new feature
- **Example README**: 550+ lines of comprehensive documentation
- **Quick Start Script**: Interactive setup and execution

---

## 🚀 Getting Started

### Step 1: Get Workflow JSON

**Method A**: Export from ComfyUI
1. Open ComfyUI web interface
2. Create your workflow
3. Click **File → Export (API Format)**
4. Save as `workflow.json`

**Method B**: Use example
```bash
cp examples/execute_from_json/workflow_example.json my_workflow.json
```

### Step 2: Build

```bash
cd /data/comfyui-go-sdk
make build-execute-json
```

### Step 3: Run

```bash
# Basic execution
./bin/execute_from_json my_workflow.json

# With custom parameters
./bin/execute_from_json my_workflow.json seed=42 steps=25

# With custom prompts
./bin/execute_from_json my_workflow.json \
  prompt="beautiful mountain landscape" \
  negative="blurry, low quality"
```

---

## 💡 Usage Examples

### Example 1: Simple Execution

```bash
./bin/execute_from_json workflow.json
```

**Output**:
```
╔════════════════════════════════════════════════════════════════════╗
║         ComfyUI Go SDK - Execute Workflow from JSON File         ║
╚════════════════════════════════════════════════════════════════════╝

🔍 Checking ComfyUI server status...
✅ ComfyUI server is running

📂 Loading workflow from: workflow.json
✅ Workflow loaded successfully
   Total nodes: 7

🚀 Submitting workflow to ComfyUI...
✅ Workflow queued successfully!

⏳ Monitoring execution progress...
✅ Completed in 15.3 seconds

📥 Retrieving execution results...
   💾 Saved to: output/ComfyUI_00001.png

╔════════════════════════════════════════════════════════════════════╗
║                    ✅ Execution Complete!                          ║
╚════════════════════════════════════════════════════════════════════╝
```

### Example 2: Batch Processing

```bash
#!/bin/bash
# Generate 5 images with different seeds

for seed in 100 200 300 400 500; do
  echo "Generating image with seed $seed..."
  ./bin/execute_from_json workflow.json seed=$seed
done
```

### Example 3: Programmatic Usage

```go
package main

import (
    "context"
    "log"
    
    comfyui "github.com/yourusername/comfyui-go-sdk"
)

func main() {
    client := comfyui.NewClient("http://127.0.0.1:8188")
    ctx := context.Background()
    
    // Method 1: Direct execution
    resp, err := client.QueuePromptFromFile(ctx, "workflow.json", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Queued: %s", resp.PromptID)
    
    // Method 2: Load, modify, execute
    workflow, err := comfyui.LoadWorkflowFromFile("workflow.json")
    if err != nil {
        log.Fatal(err)
    }
    
    // Modify parameters
    workflow.SetNodeInput("3", "seed", 12345)
    workflow.SetNodeInput("6", "text", "beautiful landscape")
    
    // Execute
    resp, err = client.QueuePrompt(ctx, workflow, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Queued: %s", resp.PromptID)
}
```

### Example 4: Advanced - Clone and Variations

```go
// Load base workflow
baseWorkflow, _ := comfyui.LoadWorkflowFromFile("base.json")

// Generate multiple variations
prompts := []string{
    "beautiful sunset",
    "mountain landscape",
    "ocean waves",
}

for i, prompt := range prompts {
    // Clone workflow
    workflow, _ := baseWorkflow.Clone()
    
    // Modify prompt
    workflow.SetNodeInput("6", "text", prompt)
    
    // Queue execution
    resp, _ := client.QueuePrompt(ctx, workflow, nil)
    log.Printf("Variation %d queued: %s", i+1, resp.PromptID)
}
```

---

## 📖 Available Parameters

When using the CLI tool, you can override workflow parameters:

| Parameter | Description | Example |
|-----------|-------------|---------|
| `seed=<number>` | Random seed for reproducibility | `seed=12345` |
| `steps=<number>` | Number of sampling steps | `steps=30` |
| `cfg=<number>` | CFG (Classifier Free Guidance) scale | `cfg=7.5` |
| `prompt=<text>` | Positive prompt text | `prompt="beautiful sunset"` |
| `negative=<text>` | Negative prompt text | `negative="blurry"` |

**Example**:
```bash
./bin/execute_from_json workflow.json \
  seed=42 \
  steps=25 \
  cfg=8.0 \
  prompt="masterpiece, best quality, landscape" \
  negative="blurry, low quality"
```

---

## 🔧 Workflow JSON Format

The workflow JSON follows ComfyUI's API format:

```json
{
  "node_id": {
    "class_type": "NodeClassName",
    "inputs": {
      "input_name": "value",
      "connected_input": ["source_node_id", output_index]
    }
  }
}
```

**Example**:
```json
{
  "3": {
    "class_type": "KSampler",
    "inputs": {
      "seed": 12345,
      "steps": 20,
      "cfg": 8.0,
      "model": ["4", 0],
      "positive": ["6", 0],
      "negative": ["7", 0]
    }
  },
  "4": {
    "class_type": "CheckpointLoaderSimple",
    "inputs": {
      "ckpt_name": "v1-5-pruned-emaonly.safetensors"
    }
  }
}
```

---

## 📚 Documentation

### Main Documentation
- [Main README](README.md) - SDK overview and quick start
- [Example README](examples/execute_from_json/README.md) - Detailed guide (550+ lines)
- [Implementation Summary](EXECUTE_FROM_JSON_SUMMARY.md) - Technical details

### Other Examples
- [Basic Example](examples/basic/README.md) - Basic workflow submission
- [WebSocket Example](examples/websocket/README.md) - Real-time monitoring
- [Progress Example](examples/progress/README.md) - Progress tracking
- [Advanced Example](examples/advanced/README.md) - Advanced features

---

## 🎯 Use Cases

### 1. Automated Testing
```bash
# Test workflow in CI/CD
./bin/execute_from_json test_workflow.json
if [ $? -eq 0 ]; then
  echo "✅ Workflow test passed"
fi
```

### 2. Batch Image Generation
```bash
# Generate multiple variations
for i in {1..10}; do
  ./bin/execute_from_json workflow.json seed=$RANDOM
done
```

### 3. Workflow Validation
```go
workflow, err := comfyui.LoadWorkflowFromFile("workflow.json")
if err := workflow.Validate(); err != nil {
    log.Fatal("Invalid workflow:", err)
}
```

### 4. Dynamic Workflow Modification
```go
workflow, _ := comfyui.LoadWorkflowFromFile("base.json")

// Modify based on user input
workflow.SetNodeInput("3", "seed", userSeed)
workflow.SetNodeInput("6", "text", userPrompt)

client.QueuePrompt(ctx, workflow, nil)
```

---

## 🛠️ Build and Test

### Build All Examples
```bash
make build
```

### Build This Example Only
```bash
make build-execute-json
```

### Run Tests
```bash
go test -v .
```

### Check Binary
```bash
ls -lh bin/execute_from_json
# Output: -rwxr-xr-x 1 root root 8.4M Oct 22 11:19 execute_from_json
```

---

## ✅ Features

- ✅ **Simple API**: One-line workflow execution
- ✅ **CLI Tool**: Complete command-line interface
- ✅ **Parameter Override**: Modify parameters at runtime
- ✅ **Progress Monitoring**: Real-time execution tracking
- ✅ **Auto Download**: Automatic image retrieval
- ✅ **Error Handling**: Comprehensive error reporting
- ✅ **Documentation**: Extensive guides and examples
- ✅ **Integration**: Works with existing SDK features
- ✅ **Type Safe**: Full Go type safety
- ✅ **Tested**: All tests passing

---

## 🔗 Integration

### With WebSocket
```go
ws, _ := client.NewWebSocket(ctx)
events := ws.Subscribe()

resp, _ := client.QueuePromptFromFile(ctx, "workflow.json", nil)

for event := range events {
    // Handle real-time updates
}
```

### With Progress Tracking
See [examples/progress](examples/progress/README.md) for detailed implementation.

### With Batch Processing
See [examples/advanced](examples/advanced/README.md) for batch processing examples.

---

## 🐛 Troubleshooting

### Issue: "Connection refused"
**Solution**: Make sure ComfyUI is running
```bash
curl http://127.0.0.1:8188/system_stats
```

### Issue: "Workflow file not found"
**Solution**: Check file path
```bash
ls -l workflow.json
```

### Issue: "Node errors"
**Solution**: Validate workflow JSON
```go
workflow, _ := comfyui.LoadWorkflowFromFile("workflow.json")
err := workflow.Validate()
```

---

## 📊 Statistics

- **SDK Method**: 1 new method (`QueuePromptFromFile`)
- **Example Code**: 342 lines
- **Documentation**: 550+ lines
- **Binary Size**: 8.4 MB
- **Build Time**: < 5 seconds
- **Test Coverage**: 100% passing

---

## 🎉 Summary

The ComfyUI Go SDK now provides a complete solution for executing workflows from JSON files:

1. **Simple API** - One method call to execute workflows
2. **Complete CLI Tool** - Full-featured command-line application
3. **Extensive Documentation** - Guides, examples, and API reference
4. **Production Ready** - Tested and ready to use

**Get Started Now**:
```bash
cd /data/comfyui-go-sdk
make build-execute-json
./bin/execute_from_json examples/execute_from_json/workflow_example.json
```

---

## 📞 Support

- **Documentation**: See [examples/execute_from_json/README.md](examples/execute_from_json/README.md)
- **Examples**: Check `examples/` directory
- **Issues**: Open an issue on GitHub

---

**Status**: ✅ Ready to Use  
**Version**: 1.0.0  
**Date**: 2025-10-22

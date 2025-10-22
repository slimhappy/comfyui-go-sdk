# Examples Overview - ComfyUI Go SDK

This directory contains comprehensive examples demonstrating all features of the ComfyUI Go SDK. Each example includes detailed documentation and can be run independently.

## 📚 Available Examples

### 1. [Basic Example](./basic/) - Getting Started

**File:** `examples/basic/main.go` | **Lines:** 397 in README

The foundational example covering core SDK functionality.

**Features:**
- ✅ System information retrieval
- ✅ Model listing
- ✅ Workflow building with WorkflowBuilder
- ✅ Workflow submission
- ✅ Real-time monitoring via WebSocket
- ✅ Result retrieval and image saving

**Best For:**
- First-time users
- Understanding basic concepts
- Learning workflow structure
- Quick prototyping

**Quick Start:**
```bash
./bin/basic
```

**Learn More:** [Basic Example README](./basic/README.md)

---

### 2. [WebSocket Example](./websocket/) - Event Monitoring

**File:** `examples/websocket/main.go` | **Lines:** 513 in README

Real-time event monitoring and debugging tool.

**Features:**
- ✅ WebSocket connection management
- ✅ All event types handling (status, executing, progress, executed, cached, error)
- ✅ Graceful shutdown (Ctrl+C)
- ✅ Detailed error reporting with traceback
- ✅ Queue status monitoring

**Best For:**
- Development and debugging
- System monitoring
- Understanding ComfyUI events
- Integration testing

**Quick Start:**
```bash
./bin/websocket
# Press Ctrl+C to exit
```

**Learn More:** [WebSocket Example README](./websocket/README.md)

---

### 3. [Advanced Example](./advanced/) - Power Features

**File:** `examples/advanced/main.go` | **Lines:** 654 in README

Comprehensive demonstration of advanced SDK capabilities.

**Features:**
- ✅ Image upload and img2img workflows
- ✅ Batch processing with different parameters
- ✅ Queue management and monitoring
- ✅ History retrieval and analysis
- ✅ Node information queries
- ✅ Workflow file operations (load/save/modify)
- ✅ Concurrent execution with timeouts

**Best For:**
- Production applications
- Batch processing pipelines
- Automation systems
- Advanced workflow manipulation

**Quick Start:**
```bash
./bin/advanced
```

**Learn More:** [Advanced Example README](./advanced/README.md)

---

### 4. [Progress Example](./progress/) - Visual Tracking

**File:** `examples/progress/main.go` | **Lines:** 178 in README

Beautiful visual progress tracking with real-time updates.

**Features:**
- ✅ ASCII progress bar (40 characters)
- ✅ Real-time percentage display
- ✅ Current node tracking
- ✅ Elapsed time and ETA calculation
- ✅ Step counter (e.g., 10/20)
- ✅ Completion summary
- ✅ Automatic image saving

**Best For:**
- CLI applications
- User-facing tools
- Progress visualization
- Long-running workflows

**Quick Start:**
```bash
./bin/progress
```

**Learn More:** [Progress Example README](./progress/README.md)

---

## 🚀 Quick Start Guide

### Build All Examples

```bash
cd /data/comfyui-go-sdk
make build
```

This creates executables in `bin/`:
- `bin/basic`
- `bin/websocket`
- `bin/advanced`
- `bin/progress`

### Run Individual Examples

```bash
# Basic workflow submission
./bin/basic

# WebSocket event monitoring
./bin/websocket

# Advanced features
./bin/advanced

# Progress tracking
./bin/progress
```

### Run from Source

```bash
# Run any example directly
cd examples/basic
go run main.go

cd examples/websocket
go run main.go

cd examples/advanced
go run main.go

cd examples/progress
go run main.go
```

## 📊 Feature Comparison

| Feature | Basic | WebSocket | Advanced | Progress |
|---------|-------|-----------|----------|----------|
| Workflow Submission | ✅ | ❌ | ✅ | ✅ |
| WebSocket Monitoring | ✅ | ✅ | ❌ | ✅ |
| Progress Bar | ❌ | ❌ | ❌ | ✅ |
| Image Upload | ❌ | ❌ | ✅ | ❌ |
| Batch Processing | ❌ | ❌ | ✅ | ❌ |
| Queue Management | ❌ | ❌ | ✅ | ❌ |
| History Retrieval | ✅ | ❌ | ✅ | ✅ |
| Concurrent Execution | ❌ | ❌ | ✅ | ❌ |
| Workflow Files | ❌ | ❌ | ✅ | ❌ |
| Node Information | ❌ | ❌ | ✅ | ❌ |
| Error Handling | ✅ | ✅ | ✅ | ✅ |

## 🎯 Learning Path

### For Beginners

1. **Start with [Basic](./basic/)** - Learn fundamental concepts
2. **Try [WebSocket](./websocket/)** - Understand event system
3. **Explore [Progress](./progress/)** - Add visual feedback
4. **Master [Advanced](./advanced/)** - Learn power features

### For Experienced Developers

1. **Review [Advanced](./advanced/)** - See all capabilities
2. **Check [Progress](./progress/)** - Learn progress tracking patterns
3. **Reference [Basic](./basic/)** - Understand workflow structure
4. **Use [WebSocket](./websocket/)** - Debug and monitor

## 🔧 Common Use Cases

### CLI Tool Development

**Recommended:** Basic + Progress

```bash
# Combine workflow submission with progress tracking
# See: examples/progress/main.go
```

### Web Service Backend

**Recommended:** Advanced + WebSocket

```bash
# Use advanced features for API endpoints
# Use WebSocket for real-time updates to frontend
```

### Batch Processing

**Recommended:** Advanced

```bash
# See Example 2 in examples/advanced/main.go
# Batch processing with different parameters
```

### Monitoring Dashboard

**Recommended:** WebSocket

```bash
# Real-time event monitoring
# See: examples/websocket/main.go
```

### Automation Pipeline

**Recommended:** Advanced

```bash
# Image upload, processing, and result retrieval
# See Examples 1, 2, 7 in examples/advanced/main.go
```

## 📖 Documentation Structure

Each example includes:

### README.md Contents

1. **Features** - What the example demonstrates
2. **Quick Start** - How to run it
3. **What It Does** - Detailed explanation
4. **Code Structure** - Key components
5. **Example Output** - What to expect
6. **Customization** - How to modify
7. **Troubleshooting** - Common issues
8. **API Reference** - Functions used
9. **Learning Points** - Key takeaways
10. **Related Examples** - What to explore next

### Code Comments

All examples include:
- Inline comments explaining key concepts
- Function documentation
- Usage examples
- Error handling patterns

## 🎨 Example Combinations

### Monitor While Processing

```bash
# Terminal 1: Monitor events
./bin/websocket

# Terminal 2: Submit workflows
./bin/basic
```

### Batch with Progress

Combine batch processing logic from `advanced` with progress tracking from `progress`:

```go
// See examples/advanced/main.go - Example 2
// See examples/progress/main.go - MonitorProgress function
```

### Upload and Track

Combine image upload from `advanced` with progress tracking:

```go
// Upload image (from advanced example)
uploaded, _ := client.UploadImage(ctx, "input.png", options)

// Build workflow
workflow := buildImg2ImgWorkflow(uploaded.Name)

// Submit and track progress (from progress example)
result, _ := client.QueuePrompt(ctx, workflow, nil)
MonitorProgress(ctx, client, result.PromptID)
```

## 🐛 Troubleshooting

### All Examples Fail to Connect

```
Failed to connect: connection refused
```

**Solution:** Ensure ComfyUI is running at `http://127.0.0.1:8188`

```bash
# Check if ComfyUI is running
curl http://127.0.0.1:8188/system_stats

# Start ComfyUI if needed
cd /path/to/ComfyUI
python main.py
```

### Build Errors

```
go: module not found
```

**Solution:** Initialize Go module

```bash
cd /data/comfyui-go-sdk
go mod tidy
make build
```

### Model Not Found

```
Failed to queue prompt: model not found
```

**Solution:** Update checkpoint name in examples to match your installed models

```bash
# Check available models
curl http://127.0.0.1:8188/object_info | jq '.CheckpointLoaderSimple'
```

### Permission Denied

```
permission denied: ./bin/basic
```

**Solution:** Make binaries executable

```bash
chmod +x bin/*
```

## 💡 Tips & Best Practices

### 1. Start Simple

Begin with the basic example and gradually explore advanced features.

### 2. Use WebSocket for Debugging

Run the WebSocket example in a separate terminal while developing.

### 3. Check Prerequisites

Ensure ComfyUI is running and models are installed before running examples.

### 4. Read the Code

Each example is well-commented. Reading the source code is highly recommended.

### 5. Experiment

Modify examples to fit your needs. They're designed to be starting points.

### 6. Combine Features

Mix and match code from different examples for your use case.

### 7. Handle Errors

All examples demonstrate proper error handling. Follow these patterns.

### 8. Use Contexts

Examples show proper context usage for timeouts and cancellation.

## 📚 Additional Resources

### SDK Documentation

- [Main README](../README.md) - Complete SDK documentation
- [Quick Start Guide](../QUICKSTART.md) - Getting started tutorial
- [Progress Tracking Guide](../PROGRESS_TRACKING.md) - Progress patterns
- [API Reference](../README.md#api-reference) - Complete API docs

### ComfyUI Resources

- [ComfyUI GitHub](https://github.com/comfyanonymous/ComfyUI)
- [ComfyUI API Documentation](https://github.com/comfyanonymous/ComfyUI/wiki/API)
- [ComfyUI Custom Nodes](https://github.com/ltdrdata/ComfyUI-Manager)

### Go Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go by Example](https://gobyexample.com/)

## 🎓 Next Steps

After exploring the examples:

1. **Build Your Own Application** - Use examples as templates
2. **Contribute** - Share your improvements or new examples
3. **Integrate** - Add ComfyUI to your existing projects
4. **Automate** - Create batch processing pipelines
5. **Monitor** - Build monitoring dashboards
6. **Extend** - Develop custom workflows and nodes

## 📊 Example Statistics

| Example | README Lines | Code Lines | Features | Complexity |
|---------|--------------|------------|----------|------------|
| Basic | 397 | ~180 | 6 | ⭐⭐ |
| WebSocket | 513 | ~130 | 5 | ⭐⭐ |
| Advanced | 654 | ~280 | 7 | ⭐⭐⭐⭐ |
| Progress | 178 | ~340 | 8 | ⭐⭐⭐ |

**Total:** 1,742 lines of documentation, ~930 lines of example code

## 🎉 Summary

The ComfyUI Go SDK examples provide:

- ✅ **4 comprehensive examples** covering all SDK features
- ✅ **1,742 lines of documentation** with detailed explanations
- ✅ **~930 lines of example code** ready to run
- ✅ **Complete API coverage** from basic to advanced
- ✅ **Production-ready patterns** for error handling and resource management
- ✅ **Visual progress tracking** for better UX
- ✅ **Real-time monitoring** for debugging and integration

**Start exploring now!** 🚀

---

For questions or issues, please refer to the main [README](../README.md) or open an issue on GitHub.

**Happy coding!** 💻

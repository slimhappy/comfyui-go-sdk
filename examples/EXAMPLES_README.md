# ComfyUI Go SDK Examples

This directory contains comprehensive examples demonstrating all features of the ComfyUI Go SDK.

## Quick Start

### Build All Examples

```bash
make build
```

### Run Specific Examples

```bash
# Queue management
make run-queue

# History operations
make run-history

# Model and node information
make run-model-info

# Image operations
make run-image-ops

# Error handling
make run-error-handling

# Integration test suite
make run-integration

# Run all examples
make run-all-examples
```

## Examples Overview

### 1. Basic (`examples/basic/`)
**Purpose**: Introduction to basic SDK usage

**Features**:
- System statistics retrieval
- Model discovery
- Simple workflow creation
- Workflow submission
- WebSocket monitoring
- Result retrieval

**Run**: `./bin/basic`

---

### 2. WebSocket (`examples/websocket/`)
**Purpose**: Real-time event monitoring

**Features**:
- WebSocket connection
- Event type handling
- Progress tracking
- Error detection

**Run**: `./bin/websocket`

---

### 3. Advanced (`examples/advanced/`)
**Purpose**: Advanced SDK features

**Features**:
- Image upload
- Batch processing
- Queue management
- History management
- Node information queries
- Workflow from file
- Concurrent execution

**Run**: `./bin/advanced`

---

### 4. Progress (`examples/progress/`)
**Purpose**: Visual progress tracking

**Features**:
- Real-time progress bar
- Step-by-step monitoring
- Execution time tracking

**Run**: `./bin/progress`

---

### 5. Execute from JSON (`examples/execute_from_json/`)
**Purpose**: Load and execute workflows from JSON files

**Features**:
- JSON workflow loading
- Workflow modification
- Parameter customization
- Execution and monitoring

**Run**: `./bin/execute_from_json workflow_example.json`

---

### 6. Queue Management (`examples/queue_management/`) ⭐ NEW
**Purpose**: Comprehensive queue operations

**Features**:
- Queue status monitoring
- Multiple workflow submission
- Queue item deletion
- Execution interruption
- Queue clearing
- Metadata handling
- Queue monitoring until completion

**Run**: `./bin/queue_management`

**Key Functions Demonstrated**:
- `GetQueue()` - Get current queue status
- `QueuePrompt()` - Submit workflows
- `DeleteFromQueue()` - Remove specific items
- `Interrupt()` - Stop current execution
- `ClearQueue()` - Clear all pending items

---

### 7. History Operations (`examples/history_operations/`) ⭐ NEW
**Purpose**: History management and analysis

**Features**:
- History retrieval (all and specific)
- Execution tracking
- History statistics analysis
- Image download from history
- History deletion
- History export

**Run**: `./bin/history_operations`

**Key Functions Demonstrated**:
- `GetHistory()` - Retrieve execution history
- `WaitForCompletion()` - Wait for workflow completion
- `DeleteHistory()` - Remove history items
- `ClearHistory()` - Clear all history
- History data analysis and statistics

---

### 8. Model Info (`examples/model_info/`) ⭐ NEW
**Purpose**: Model and node information queries

**Features**:
- System statistics
- Model folder discovery
- Model listing by type
- Embeddings retrieval
- Node class information
- Node input/output inspection
- Node search by category
- Node compatibility analysis
- Server features

**Run**: `./bin/model_info`

**Key Functions Demonstrated**:
- `GetSystemStats()` - System and device information
- `GetModels()` - List available models
- `GetEmbeddings()` - List embeddings
- `GetObjectInfo()` - Node class information
- `GetFeatures()` - Server capabilities

---

### 9. Image Operations (`examples/image_operations/`) ⭐ NEW
**Purpose**: Image upload and download operations

**Features**:
- Test image creation
- Image upload (file and bytes)
- Upload to subfolders
- Image-to-image workflows
- Image download
- Batch image processing
- Image format handling

**Run**: `./bin/image_operations`

**Key Functions Demonstrated**:
- `UploadImage()` - Upload image file
- `UploadImageBytes()` - Upload from memory
- `GetImage()` - Download image
- `SaveImage()` - Save image to file
- Image workflow integration

---

### 10. Error Handling (`examples/error_handling/`) ⭐ NEW
**Purpose**: Error handling and retry strategies

**Features**:
- Connection error handling
- Retry with exponential backoff
- Invalid workflow detection
- Missing model handling
- Context timeout handling
- Graceful degradation
- Workflow validation
- Safe execution with recovery
- Execution monitoring with error detection
- Batch execution with error handling

**Run**: `./bin/error_handling`

**Key Patterns Demonstrated**:
- Retry logic with backoff
- Context cancellation
- Panic recovery
- Error type checking
- Validation before execution
- Timeout handling

---

### 11. Integration Test (`examples/integration_test/`) ⭐ NEW
**Purpose**: Comprehensive SDK test suite

**Features**:
- Server connectivity test
- System information test
- Model discovery test
- Node information test
- Queue operations test
- Workflow execution test
- WebSocket monitoring test
- History management test
- Image operations test
- Error handling test

**Run**: `./bin/integration_test`

**Environment Variables**:
- `COMFYUI_URL` - ComfyUI server URL (default: http://127.0.0.1:8188)

**Output**: Test results with pass/fail status for each test

---

## Prerequisites

1. **ComfyUI Server**: Ensure ComfyUI is running on `http://127.0.0.1:8188` (or set `COMFYUI_URL`)
2. **Model**: At least one checkpoint model installed (e.g., `v1-5-pruned-emaonly.safetensors`)
3. **Go**: Go 1.21 or later

## Building Examples

### Build All
```bash
make build
```

### Build Specific Example
```bash
cd examples/queue_management
go build -o ../../bin/queue_management
```

## Running Examples

### Direct Execution
```bash
./bin/queue_management
./bin/history_operations
./bin/model_info
./bin/image_operations
./bin/error_handling
./bin/integration_test
```

### With Make
```bash
make run-queue
make run-history
make run-model-info
make run-image-ops
make run-error-handling
make run-integration
```

### Run All Examples
```bash
make run-all-examples
```

## Example Output

Each example provides detailed output showing:
- ✓ Successful operations
- ✗ Failed operations
- Progress information
- Results and statistics
- Error messages (if any)

## Customization

All examples can be customized by:
1. Modifying the workflow parameters
2. Changing the server URL
3. Adjusting timeouts and retry logic
4. Adding custom error handling

## Common Issues

### Connection Refused
- Ensure ComfyUI server is running
- Check the server URL and port
- Verify firewall settings

### Model Not Found
- Install required models in ComfyUI
- Check model names in the examples
- Update checkpoint names to match your installation

### Timeout Errors
- Increase timeout values in the code
- Check server performance
- Reduce workflow complexity (fewer steps)

## Learning Path

Recommended order for learning:

1. **basic** - Start here to understand fundamentals
2. **websocket** - Learn real-time monitoring
3. **queue_management** - Understand queue operations
4. **history_operations** - Learn result retrieval
5. **model_info** - Explore available resources
6. **image_operations** - Work with images
7. **error_handling** - Implement robust error handling
8. **advanced** - Combine multiple features
9. **integration_test** - See complete test suite

## Contributing

When adding new examples:
1. Create a new directory under `examples/`
2. Add `main.go` with clear comments
3. Update this README
4. Add build target to Makefile
5. Test thoroughly

## Support

For issues or questions:
- Check the main README.md
- Review the SDK documentation
- Open an issue on GitHub

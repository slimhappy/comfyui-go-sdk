# ComfyUI Go SDK

A comprehensive Go SDK for interacting with ComfyUI API.

## Features

- ✅ HTTP REST API client
- ✅ WebSocket client for real-time updates
- ✅ Workflow submission and management
- ✅ **Load and execute workflows from JSON files**
- ✅ Queue management
- ✅ History tracking
- ✅ Image upload and download
- ✅ System information queries
- ✅ Type-safe API


## Installation

```bash
go get github.com/slimhappy/comfyui-go-sdk
```

## Quick Start

### Method 1: Execute Workflow from JSON File (Easiest!)

```go
package main

import (
    "context"
    "log"
    
    comfyui "github.com/slimhappy/comfyui-go-sdk"
)

func main() {
    // Create client
    client := comfyui.NewClient("http://127.0.0.1:8188")
    ctx := context.Background()
    
    // Execute workflow directly from JSON file
    result, err := client.QueuePromptFromFile(ctx, "workflow.json", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Workflow queued: %s", result.PromptID)
}
```

### Method 2: Load, Modify, and Execute

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    comfyui "github.com/slimhappy/comfyui-go-sdk"
)

func main() {
    // Create client
    client := comfyui.NewClient("http://127.0.0.1:8188")
    
    // Load workflow from API format JSON
    workflow, err := comfyui.LoadWorkflowFromFile("workflow_api.json")
    if err != nil {
        log.Fatal(err)
    }
    
    // Modify workflow parameters
    workflow.SetNodeInput("3", "seed", 12345)
    workflow.SetNodeInput("3", "steps", 30)
    
    // Submit workflow
    ctx := context.Background()
    result, err := client.QueuePrompt(ctx, workflow, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Submitted workflow: %s\n", result.PromptID)
    
    // Wait for completion and get images
    images, err := client.WaitForCompletion(ctx, result.PromptID)
    if err != nil {
        log.Fatal(err)
    }
    
    // Save images
    for _, img := range images {
        if err := client.SaveImage(img, "output_"+img.Filename); err != nil {
            log.Printf("Failed to save image: %v", err)
        }
    }
}
```


## Usage Examples

### Execute Workflow from JSON File

```go
// Method 1: Direct execution
client := comfyui.NewClient("http://127.0.0.1:8188")
result, err := client.QueuePromptFromFile(context.Background(), "workflow.json", nil)

// Method 2: Load and modify before execution
workflow, err := comfyui.LoadWorkflowFromFile("workflow.json")
if err != nil {
    log.Fatal(err)
}

// Modify parameters
workflow.SetNodeInput("3", "seed", 42)
workflow.SetNodeInput("6", "text", "beautiful landscape")

// Execute
result, err := client.QueuePrompt(context.Background(), workflow, nil)
```

**How to get workflow JSON:**
1. Open ComfyUI web interface
2. Create your workflow
3. Click **File → Export (API Format)**
4. Save as `workflow.json`

See [examples/execute_from_json](examples/execute_from_json/README.md) for a complete CLI tool with progress tracking!

### Basic Workflow Submission


```go
client := comfyui.NewClient("http://127.0.0.1:8188")

workflow := comfyui.Workflow{
    "3": comfyui.Node{
        ClassType: "KSampler",
        Inputs: map[string]interface{}{
            "seed": 123456,
            "steps": 20,
            "cfg": 8.0,
            // ... more inputs
        },
    },
    // ... more nodes
}

result, err := client.QueuePrompt(context.Background(), workflow, nil)
```

### WebSocket Monitoring

```go
client := comfyui.NewClient("http://127.0.0.1:8188")

// Connect WebSocket
ws, err := client.ConnectWebSocket(context.Background())
if err != nil {
    log.Fatal(err)
}
defer ws.Close()

// Listen for messages
for msg := range ws.Messages() {
    switch msg.Type {
    case comfyui.MessageTypeExecuting:
        fmt.Printf("Executing node: %s\n", msg.Data.Node)
    case comfyui.MessageTypeProgress:
        fmt.Printf("Progress: %d/%d\n", msg.Data.Value, msg.Data.Max)
    case comfyui.MessageTypeExecuted:
        fmt.Println("Node completed!")
    }
}
```

### Progress Tracking with Visual Progress Bar

```go
// Submit workflow
result, err := client.QueuePrompt(ctx, workflow, nil)
if err != nil {
    log.Fatal(err)
}

// Connect WebSocket for progress monitoring
ws, err := client.ConnectWebSocket(ctx)
if err != nil {
    log.Fatal(err)
}
defer ws.Close()

// Track progress with visual feedback
for msg := range ws.Messages() {
    if msg.Type == string(comfyui.MessageTypeProgress) {
        data, _ := msg.GetProgressData()
        percentage := float64(data.Value) / float64(data.Max) * 100
        
        // Draw progress bar
        bar := strings.Repeat("█", int(percentage/2.5)) + 
               strings.Repeat("░", 40-int(percentage/2.5))
        fmt.Printf("\r⏳ [%s] %.1f%% | Step %d/%d", 
            bar, percentage, data.Value, data.Max)
    }
    
    if msg.Type == string(comfyui.MessageTypeExecuting) {
        data, _ := msg.GetExecutingData()
        if data.Node == nil {
            fmt.Println("\n✅ Completed!")
            break
        }
    }
}

// See examples/progress/ for a complete implementation
```

### Upload and Use Image


```go
// Upload image
imageRef, err := client.UploadImage(ctx, "input.png", comfyui.UploadOptions{
    Type: "input",
    Subfolder: "",
})
if err != nil {
    log.Fatal(err)
}

// Use in workflow
workflow["10"] = comfyui.Node{
    ClassType: "LoadImage",
    Inputs: map[string]interface{}{
        "image": imageRef.Name,
    },
}
```

### Query History

```go
// Get all history
history, err := client.GetHistory(ctx, "")
if err != nil {
    log.Fatal(err)
}

// Get specific prompt history
promptHistory, err := client.GetHistory(ctx, "prompt-id-here")
if err != nil {
    log.Fatal(err)
}
```

### Queue Management

```go
// Get queue status
queue, err := client.GetQueue(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Running: %d, Pending: %d\n", 
    len(queue.QueueRunning), len(queue.QueuePending))

// Clear queue
err = client.ClearQueue(ctx)

// Delete specific items
err = client.DeleteFromQueue(ctx, []string{"prompt-id-1", "prompt-id-2"})

// Interrupt execution
err = client.Interrupt(ctx, "")
```

### Get System Info

```go
// Get system stats
stats, err := client.GetSystemStats(ctx)
if err != nil {
    log.Fatal(err)
}

for _, device := range stats.Devices {
    fmt.Printf("Device: %s, VRAM: %d MB free\n", 
        device.Name, device.VRAMFree/1024/1024)
}

// Get object info (node definitions)
objectInfo, err := client.GetObjectInfo(ctx, "")
if err != nil {
    log.Fatal(err)
}
```

## API Documentation

### Client Methods

#### Workflow Operations
- `QueuePrompt(ctx, workflow, options)` - Submit workflow
- `QueuePromptFromFile(ctx, filepath, options)` - **Load and execute workflow from JSON file**
- `LoadWorkflowFromFile(filepath)` - Load workflow from JSON file
- `SaveWorkflowToFile(workflow, filepath)` - Save workflow to JSON file
- `WaitForCompletion(ctx, promptID)` - Wait for workflow completion
- `GetHistory(ctx, promptID)` - Get execution history
- `ClearHistory(ctx)` - Clear all history
- `DeleteHistory(ctx, promptIDs)` - Delete specific history


#### Queue Management
- `GetQueue(ctx)` - Get queue status
- `ClearQueue(ctx)` - Clear queue
- `DeleteFromQueue(ctx, promptIDs)` - Delete from queue
- `Interrupt(ctx, promptID)` - Interrupt execution

#### File Operations
- `UploadImage(ctx, filepath, options)` - Upload image
- `UploadImageBytes(ctx, data, filename, options)` - Upload image from bytes
- `GetImage(ctx, filename, subfolder, folderType)` - Download image
- `SaveImage(imageInfo, outputPath)` - Save image to file

#### System Information
- `GetSystemStats(ctx)` - Get system statistics
- `GetObjectInfo(ctx, nodeClass)` - Get node definitions
- `GetEmbeddings(ctx)` - Get embeddings list
- `GetModels(ctx, folder)` - Get models list

#### WebSocket
- `ConnectWebSocket(ctx)` - Connect WebSocket
- `ws.Messages()` - Get message channel
- `ws.Close()` - Close connection

## Project Structure

```
comfyui-go-sdk/
├── client.go          # Main client implementation
├── types.go           # Type definitions
├── websocket.go       # WebSocket client
├── workflow.go        # Workflow utilities
├── errors.go          # Error definitions
├── examples/          # Usage examples
│   ├── basic/              # Basic workflow submission
│   ├── websocket/          # WebSocket event monitoring
│   ├── advanced/           # Advanced features
│   ├── progress/           # Progress tracking with visual bar
│   └── execute_from_json/  # Execute workflow from JSON file ⭐
└── README.md
```

## Examples

All examples are located in the `examples/` directory:

| Example | Description | Documentation |
|---------|-------------|---------------|
| **basic** | Basic workflow submission and system info | [README](examples/basic/README.md) |
| **websocket** | Real-time event monitoring via WebSocket | [README](examples/websocket/README.md) |
| **advanced** | Advanced features (batch, upload, etc.) | [README](examples/advanced/README.md) |
| **progress** | Progress tracking with visual progress bar | [README](examples/progress/README.md) |
| **execute_from_json** | Execute workflows from JSON files ⭐ | [README](examples/execute_from_json/README.md) |

### Build and Run Examples

```bash
# Build all examples
make build

# Or build specific example
make build-execute-json

# Run examples
./bin/basic
./bin/websocket
./bin/advanced
./bin/progress
./bin/execute_from_json workflow.json
```


## Changelog

### 2025-10-22

#### Bug Fixes
- **Fixed Queue API JSON parsing error** ([QUEUE_API_BUGFIX.md](QUEUE_API_BUGFIX.md))
  - Fixed "cannot unmarshal array into Go struct field QueueStatus.queue_running" error
  - Added custom JSON unmarshaling for `QueueStatus` to handle array-based queue items
  - All queue monitoring features now work correctly
  
- **Fixed History API JSON parsing error** ([BUGFIX_SUMMARY.md](BUGFIX_SUMMARY.md))
  - Fixed "cannot unmarshal array into Go struct field HistoryItem.prompt" error
  - Added custom JSON unmarshaling for `PromptArray` to handle array-based prompt data
  - History retrieval now works correctly

#### New Features
- **Execute workflows from JSON files** ([EXECUTE_FROM_JSON_SUMMARY.md](EXECUTE_FROM_JSON_SUMMARY.md))
  - Added `QueuePromptFromFile()` method to SDK
  - Created complete CLI tool with progress monitoring
  - Support for dynamic parameter modification
  - Automatic result retrieval and image download
  - Comprehensive documentation and examples

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License


## Related Links

- [ComfyUI](https://github.com/comfyanonymous/ComfyUI)
- [ComfyUI Documentation](https://docs.comfy.org/)
- [API Examples](https://github.com/comfyanonymous/ComfyUI/tree/master/script_examples)

# Basic Example - ComfyUI Go SDK

This example demonstrates the fundamental features of the ComfyUI Go SDK, including system information retrieval, workflow building, submission, and result handling.

## 📋 Features

This example covers the following core functionalities:

1. **System Stats** - Get ComfyUI system information (OS, Python version, GPU info)
2. **Model Listing** - Retrieve available models (checkpoints, VAEs, etc.)
3. **Workflow Building** - Construct a simple text-to-image workflow programmatically
4. **Workflow Submission** - Queue a workflow for execution
5. **Real-time Monitoring** - Monitor execution progress via WebSocket
6. **Result Retrieval** - Get execution history and download generated images

## 🚀 Quick Start

### Prerequisites

- ComfyUI server running at `http://127.0.0.1:8188`
- At least one checkpoint model installed (e.g., `v1-5-pruned-emaonly.safetensors`)

### Build and Run

```bash
# From the SDK root directory
cd /data/comfyui-go-sdk
make build
./bin/basic
```

Or run directly:

```bash
cd examples/basic
go run main.go
```

## 📖 What This Example Does

### 1. Get System Stats

Retrieves and displays ComfyUI system information:

```go
stats, err := client.GetSystemStats(ctx)
// Displays: OS, Python version, GPU devices, VRAM usage
```

**Output:**
```
=== System Stats ===
OS: Linux
Python: 3.10.12
Device: NVIDIA GeForce RTX 3090 (cuda)
  VRAM: 22.5 GB / 24.0 GB
```

### 2. List Available Models

Fetches available checkpoint models:

```go
checkpoints, err := client.GetModels(ctx, "checkpoints")
// Shows first 5 checkpoints
```

**Output:**
```
=== Available Checkpoints ===
  - v1-5-pruned-emaonly.safetensors
  - sd_xl_base_1.0.safetensors
  - dreamshaper_8.safetensors
  ... and 12 more
```

### 3. Build a Workflow

Constructs a complete text-to-image workflow programmatically:

```go
workflow := buildSimpleWorkflow()
```

The workflow includes:
- **CheckpointLoaderSimple** - Load the model
- **EmptyLatentImage** - Create blank latent (512x512)
- **CLIPTextEncode** (x2) - Encode positive and negative prompts
- **KSampler** - Generate image (20 steps, euler sampler)
- **VAEDecode** - Decode latent to image
- **SaveImage** - Save the result

### 4. Submit Workflow

Queues the workflow for execution:

```go
result, err := client.QueuePrompt(ctx, workflow, nil)
fmt.Printf("Prompt ID: %s\n", result.PromptID)
```

**Output:**
```
=== Building Workflow ===
Submitting workflow...
Prompt ID: abc123-def456-ghi789
Queue Number: 1
```

### 5. Monitor Execution

Uses WebSocket to monitor real-time progress:

```go
ws, err := client.ConnectWebSocket(ctx)
for msg := range ws.Messages() {
    switch msg.Type {
    case comfyui.MessageTypeExecuting:
        // Handle node execution
    case comfyui.MessageTypeProgress:
        // Handle progress updates
    case comfyui.MessageTypeError:
        // Handle errors
    }
}
```

**Output:**
```
=== Monitoring Execution ===
Executing node: 4
Executing node: 5
Executing node: 6
Executing node: 7
Executing node: 3
Progress: 5/20
Progress: 10/20
Progress: 15/20
Progress: 20/20
Executing node: 8
Executing node: 9
✓ Execution completed!
```

### 6. Download Results

Retrieves execution history and saves generated images:

```go
history, err := client.GetHistory(ctx, result.PromptID)
for nodeID, output := range item.Outputs {
    for i, img := range output.Images {
        client.SaveImage(ctx, img, outputPath)
    }
}
```

**Output:**
```
=== Getting Results ===
Status: success
Completed: true

Node 9 produced 1 image(s):
  ✓ Saved: output_9_0.png

=== Done ===
```

## 🔧 Customization

### Change the Prompt

Modify the text in `buildSimpleWorkflow()`:

```go
positiveID := builder.AddNode("CLIPTextEncode", map[string]interface{}{
    "text": "your custom prompt here",
})
```

### Adjust Image Size

Change the latent dimensions:

```go
latentID := builder.AddNode("EmptyLatentImage", map[string]interface{}{
    "width":      768,  // Change from 512
    "height":     768,  // Change from 512
    "batch_size": 1,
})
```

### Modify Sampling Parameters

Adjust KSampler settings:

```go
samplerID := builder.AddNode("KSampler", map[string]interface{}{
    "seed":         54321,        // Different seed
    "steps":        30,            // More steps
    "cfg":          7.5,           // Lower CFG
    "sampler_name": "dpmpp_2m",    // Different sampler
    "scheduler":    "karras",      // Different scheduler
    "denoise":      1.0,
})
```

### Use a Different Model

Change the checkpoint:

```go
checkpointID := builder.AddNode("CheckpointLoaderSimple", map[string]interface{}{
    "ckpt_name": "dreamshaper_8.safetensors",
})
```

## 📊 Complete Workflow Structure

The example builds this workflow:

```
CheckpointLoaderSimple (node 1)
    ├─[0:MODEL]──────────────┐
    ├─[1:CLIP]───────────┐   │
    └─[2:VAE]────────┐   │   │
                     │   │   │
EmptyLatentImage     │   │   │
    └─[0:LATENT]─────┼───┼───┼──┐
                     │   │   │  │
CLIPTextEncode (+)   │   │   │  │
    ├─clip ◄─────────┘   │   │  │
    └─[0:CONDITIONING]───┼───┼──┼──┐
                         │   │  │  │
CLIPTextEncode (-)       │   │  │  │
    ├─clip ◄─────────────┘   │  │  │
    └─[0:CONDITIONING]───────┼──┼──┼──┐
                             │  │  │  │
KSampler                     │  │  │  │
    ├─model ◄────────────────┘  │  │  │
    ├─positive ◄────────────────┘  │  │
    ├─negative ◄───────────────────┘  │
    ├─latent_image ◄──────────────────┘
    └─[0:LATENT]────────────────────────┐
                                        │
VAEDecode                               │
    ├─samples ◄─────────────────────────┘
    ├─vae ◄─────────────┘
    └─[0:IMAGE]──────────────────────────┐
                                         │
SaveImage                                │
    └─images ◄───────────────────────────┘
```

## 🎯 Learning Points

This example teaches you:

1. **Client Creation** - How to initialize the ComfyUI client
2. **API Calls** - Making REST API calls to ComfyUI
3. **Workflow Builder** - Using the builder pattern to construct workflows
4. **Node Connections** - Connecting nodes with proper input/output slots
5. **WebSocket Usage** - Real-time event monitoring
6. **Error Handling** - Proper error handling patterns
7. **Image Saving** - Downloading and saving generated images

## 🔗 Related Examples

- **[websocket](../websocket/)** - Detailed WebSocket event monitoring
- **[advanced](../advanced/)** - Advanced features (batch processing, queue management)
- **[progress](../progress/)** - Visual progress tracking with progress bars

## 📚 API Reference

### Key Functions Used

- `NewClient(url string)` - Create a new client
- `GetSystemStats(ctx)` - Get system information
- `GetModels(ctx, type)` - List available models
- `QueuePrompt(ctx, workflow, options)` - Submit workflow
- `ConnectWebSocket(ctx)` - Connect to WebSocket
- `GetHistory(ctx, promptID)` - Get execution history
- `SaveImage(ctx, image, path)` - Save image to disk

### Workflow Builder Methods

- `NewWorkflowBuilder()` - Create a new builder
- `AddNode(classType, inputs)` - Add a node
- `ConnectNodes(fromID, fromSlot, toID, toInput)` - Connect nodes
- `Build()` - Build the final workflow

## 🐛 Troubleshooting

### Connection Refused

```
Failed to connect: connection refused
```

**Solution:** Make sure ComfyUI is running at `http://127.0.0.1:8188`

### Model Not Found

```
Failed to queue prompt: model not found
```

**Solution:** Update the checkpoint name in `buildSimpleWorkflow()` to match an installed model

### Timeout

```
Timeout waiting for completion
```

**Solution:** Increase the timeout duration or check if ComfyUI is processing the workflow

### WebSocket Closed

```
WebSocket closed
```

**Solution:** This is normal after execution completes. The example exits gracefully.

## 💡 Tips

1. **Check Available Models First** - Run the example to see which models are available
2. **Monitor System Resources** - Watch VRAM usage to avoid OOM errors
3. **Use Appropriate Timeouts** - Adjust the 5-minute timeout based on your workflow complexity
4. **Save Prompt IDs** - Store prompt IDs for later retrieval of results
5. **Handle Errors Gracefully** - Always check for errors in production code

## 📝 Example Output

Complete output from a successful run:

```
=== System Stats ===
OS: Linux
Python: 3.10.12
Device: NVIDIA GeForce RTX 3090 (cuda)
  VRAM: 22.5 GB / 24.0 GB

=== Available Checkpoints ===
  - v1-5-pruned-emaonly.safetensors
  - sd_xl_base_1.0.safetensors
  - dreamshaper_8.safetensors
  - realisticVisionV51_v51VAE.safetensors
  - deliberate_v2.safetensors
  ... and 8 more

=== Building Workflow ===
Submitting workflow...
Prompt ID: 12345678-90ab-cdef-1234-567890abcdef
Queue Number: 1

=== Monitoring Execution ===
Executing node: 4
Executing node: 5
Executing node: 6
Executing node: 7
Executing node: 3
Progress: 1/20
Progress: 5/20
Progress: 10/20
Progress: 15/20
Progress: 20/20
Executing node: 8
Executing node: 9
✓ Execution completed!

=== Getting Results ===
Status: success
Completed: true

Node 9 produced 1 image(s):
  ✓ Saved: output_9_0.png

=== Done ===
```

## 🎓 Next Steps

After mastering this basic example:

1. Explore **[websocket example](../websocket/)** for detailed event handling
2. Try **[advanced example](../advanced/)** for batch processing and queue management
3. Check **[progress example](../progress/)** for visual progress tracking
4. Read the main [SDK documentation](../../README.md) for complete API reference

---

**Happy coding!** 🚀

For questions or issues, please refer to the main [README](../../README.md) or check the [ComfyUI documentation](https://github.com/comfyanonymous/ComfyUI).

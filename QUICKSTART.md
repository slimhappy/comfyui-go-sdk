# Quick Start Guide

## Installation

```bash
go get github.com/yourusername/comfyui-go-sdk
```

## Prerequisites

1. ComfyUI server running at `http://127.0.0.1:8188`
2. At least one checkpoint model installed

## Your First Workflow

Create a file `main.go`:

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    comfyui "github.com/yourusername/comfyui-go-sdk"
)

func main() {
    // 1. Create client
    client := comfyui.NewClient("http://127.0.0.1:8188")
    ctx := context.Background()
    
    // 2. Build workflow
    builder := comfyui.NewWorkflowBuilder()
    
    // Load checkpoint
    ckptID := builder.AddNode("CheckpointLoaderSimple", map[string]interface{}{
        "ckpt_name": "v1-5-pruned-emaonly.safetensors",
    })
    
    // Create empty latent
    latentID := builder.AddNode("EmptyLatentImage", map[string]interface{}{
        "width": 512, "height": 512, "batch_size": 1,
    })
    
    // Positive prompt
    posID := builder.AddNode("CLIPTextEncode", map[string]interface{}{
        "text": "a beautiful sunset over mountains",
    })
    builder.ConnectNodes(ckptID, 1, posID, "clip")
    
    // Negative prompt
    negID := builder.AddNode("CLIPTextEncode", map[string]interface{}{
        "text": "bad quality, blurry",
    })
    builder.ConnectNodes(ckptID, 1, negID, "clip")
    
    // Sampler
    samplerID := builder.AddNode("KSampler", map[string]interface{}{
        "seed": 12345, "steps": 20, "cfg": 8.0,
        "sampler_name": "euler", "scheduler": "normal", "denoise": 1.0,
    })
    builder.ConnectNodes(ckptID, 0, samplerID, "model")
    builder.ConnectNodes(posID, 0, samplerID, "positive")
    builder.ConnectNodes(negID, 0, samplerID, "negative")
    builder.ConnectNodes(latentID, 0, samplerID, "latent_image")
    
    // VAE Decode
    decodeID := builder.AddNode("VAEDecode", map[string]interface{}{})
    builder.ConnectNodes(samplerID, 0, decodeID, "samples")
    builder.ConnectNodes(ckptID, 2, decodeID, "vae")
    
    // Save Image
    saveID := builder.AddNode("SaveImage", map[string]interface{}{
        "filename_prefix": "ComfyUI",
    })
    builder.ConnectNodes(decodeID, 0, saveID, "images")
    
    workflow := builder.Build()
    
    // 3. Submit workflow
    result, err := client.QueuePrompt(ctx, workflow, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Submitted! Prompt ID: %s\n", result.PromptID)
    
    // 4. Wait for completion
    execResult, err := client.WaitForCompletion(ctx, result.PromptID)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Completed in %v\n", execResult.Duration)
    
    // 5. Save images
    for i, img := range execResult.Images {
        outputPath := fmt.Sprintf("output_%d.png", i)
        if err := client.SaveImage(ctx, img, outputPath); err != nil {
            log.Printf("Failed to save: %v", err)
        } else {
            fmt.Printf("Saved: %s\n", outputPath)
        }
    }
}
```

Run it:

```bash
go run main.go
```

## Using Existing Workflows

If you have a workflow exported from ComfyUI (File → Export API Format):

```go
// Load workflow
workflow, err := comfyui.LoadWorkflowFromFile("workflow_api.json")
if err != nil {
    log.Fatal(err)
}

// Modify parameters
workflow.SetNodeInput("6", "text", "your custom prompt")
workflow.SetNodeInput("3", "seed", 99999)

// Submit
result, err := client.QueuePrompt(ctx, workflow, nil)
```

## WebSocket Monitoring

For real-time progress updates:

```go
ws, err := client.ConnectWebSocket(ctx)
if err != nil {
    log.Fatal(err)
}
defer ws.Close()

for msg := range ws.Messages() {
    switch msg.Type {
    case string(comfyui.MessageTypeProgress):
        data, _ := msg.GetProgressData()
        fmt.Printf("Progress: %d/%d\n", data.Value, data.Max)
        
    case string(comfyui.MessageTypeExecuting):
        data, _ := msg.GetExecutingData()
        if data.Node == nil {
            fmt.Println("Completed!")
            return
        }
    }
}
```

## Next Steps

- Check out the [examples](examples/) directory
- Read the [full documentation](README.md)
- Explore the [API reference](https://pkg.go.dev/github.com/yourusername/comfyui-go-sdk)

## Common Issues

### Connection Refused

Make sure ComfyUI is running:
```bash
cd /data/ComfyUI
python main.py
```

### Model Not Found

Check available models:
```go
models, _ := client.GetModels(ctx, "checkpoints")
fmt.Println(models)
```

### Workflow Errors

Validate your workflow:
```go
if err := workflow.Validate(); err != nil {
    log.Fatal(err)
}
```

Get detailed node information:
```go
info, _ := client.GetObjectInfo(ctx, "KSampler")
fmt.Printf("%+v\n", info)
```

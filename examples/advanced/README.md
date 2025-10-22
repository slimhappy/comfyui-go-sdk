# Advanced Example - ComfyUI Go SDK

This example demonstrates advanced features of the ComfyUI Go SDK, including image upload, batch processing, queue management, workflow manipulation, and concurrent execution.

## 📋 Features

This comprehensive example covers:

1. **Image Upload** - Upload images and use them in img2img workflows
2. **Batch Processing** - Process multiple workflows with different parameters
3. **Queue Management** - Monitor and manage the execution queue
4. **History Management** - Retrieve and analyze execution history
5. **Node Information** - Query node definitions and capabilities
6. **Workflow File Operations** - Load, modify, and save workflow files
7. **Concurrent Execution** - Run multiple workflows in parallel with timeouts

## 🚀 Quick Start

### Prerequisites

- ComfyUI server running at `http://127.0.0.1:8188`
- At least one checkpoint model installed
- (Optional) `input.png` file for image upload example
- (Optional) `workflow_api.json` file for workflow loading example

### Build and Run

```bash
# From the SDK root directory
cd /data/comfyui-go-sdk
make build
./bin/advanced
```

Or run directly:

```bash
cd examples/advanced
go run main.go
```

## 📖 What This Example Does

### Example 1: Image Upload & Img2Img

Upload an image and use it in an image-to-image workflow:

```go
uploadedImage, err := client.UploadImage(ctx, "input.png", comfyui.UploadOptions{
    Type:      "input",
    Subfolder: "",
    Overwrite: true,
})

workflow := buildImg2ImgWorkflow(uploadedImage.Name)
result, _ := client.QueuePrompt(ctx, workflow, nil)
```

**Output:**
```
=== Example 1: Image Upload ===
Uploaded: input_20240121_123456.png
Queued prompt: abc123-def456-ghi789
```

**Workflow Changes:**
- Replaces `EmptyLatentImage` with `LoadImage`
- Adds `VAEEncode` to encode the input image
- Sets denoise to 0.75 for img2img effect

### Example 2: Batch Processing

Process multiple workflows with different seeds:

```go
seeds := []int{12345, 67890, 11111, 22222, 33333}

for _, seed := range seeds {
    workflow, _ := baseWorkflow.Clone()
    workflow.SetNodeInput("6", "seed", seed)
    
    result, _ := client.QueuePrompt(ctx, workflow, nil)
    promptIDs = append(promptIDs, result.PromptID)
}

// Wait for all to complete
for i, promptID := range promptIDs {
    result, _ := client.WaitForCompletion(ctx, promptID)
    // Save images...
}
```

**Output:**
```
=== Example 2: Batch Processing ===
Queued seed 12345: abc123...
Queued seed 67890: def456...
Queued seed 11111: ghi789...
Queued seed 22222: jkl012...
Queued seed 33333: mno345...

Waiting for batch completion...
Waiting for 1/5...
  ✓ Completed in 12.3s, 1 images
Waiting for 2/5...
  ✓ Completed in 11.8s, 1 images
Waiting for 3/5...
  ✓ Completed in 12.1s, 1 images
Waiting for 4/5...
  ✓ Completed in 12.5s, 1 images
Waiting for 5/5...
  ✓ Completed in 12.0s, 1 images
```

**Generated Files:**
```
batch_0_seed_12345_0.png
batch_1_seed_67890_0.png
batch_2_seed_11111_0.png
batch_3_seed_22222_0.png
batch_4_seed_33333_0.png
```

### Example 3: Queue Management

Monitor and manage the execution queue:

```go
queue, err := client.GetQueue(ctx)
fmt.Printf("Running: %d, Pending: %d\n", len(queue.QueueRunning), len(queue.QueuePending))

for _, item := range queue.QueuePending {
    fmt.Printf("  - Prompt ID: %s, Number: %d\n", item.PromptID, item.Number)
}
```

**Output:**
```
=== Example 3: Queue Management ===
Running: 1, Pending: 3
Pending items:
  - Prompt ID: abc123..., Number: 2
  - Prompt ID: def456..., Number: 3
  - Prompt ID: ghi789..., Number: 4
```

**Use Cases:**
- Check queue status before submitting
- Estimate wait time
- Cancel pending items if needed

### Example 4: History Management

Retrieve and analyze execution history:

```go
history, err := client.GetHistory(ctx, "")
fmt.Printf("Total history items: %d\n", len(history))

for promptID, item := range history {
    fmt.Printf("  - %s: %s (completed: %v)\n", 
        promptID, item.Status.StatusStr, item.Status.Completed)
}
```

**Output:**
```
=== Example 4: History Management ===
Total history items: 47
  - abc123...: success (completed: true)
  - def456...: success (completed: true)
  - ghi789...: error (completed: false)
  - jkl012...: success (completed: true)
  - mno345...: success (completed: true)
```

**Use Cases:**
- Audit workflow executions
- Retrieve old results
- Analyze success/failure rates
- Debug failed workflows

### Example 5: Node Information

Query node definitions and capabilities:

```go
objectInfo, err := client.GetObjectInfo(ctx, "KSampler")

if info, ok := objectInfo["KSampler"]; ok {
    fmt.Printf("Node: %s\n", info.DisplayName)
    fmt.Printf("Category: %s\n", info.Category)
    fmt.Printf("Description: %s\n", info.Description)
    
    for name := range info.Input.Required {
        fmt.Printf("  - %s (required)\n", name)
    }
}
```

**Output:**
```
=== Example 5: Node Information ===
Node: KSampler
Category: sampling
Description: Sample latent images using various samplers
Inputs:
  - model (required)
  - seed (required)
  - steps (required)
  - cfg (required)
  - sampler_name (required)
  - scheduler (required)
  - positive (required)
  - negative (required)
  - latent_image (required)
  - denoise (required)
```

**Use Cases:**
- Discover available nodes
- Validate workflow inputs
- Build dynamic workflows
- Generate documentation

### Example 6: Workflow File Operations

Load, modify, and save workflow files:

```go
// Load from file
workflow, err := comfyui.LoadWorkflowFromFile("workflow_api.json")
fmt.Printf("Loaded workflow with %d nodes\n", len(workflow))

// Modify workflow
workflow.SetNodeInput("6", "text", "a beautiful sunset over mountains")

// Save modified workflow
err = comfyui.SaveWorkflowToFile(workflow, "workflow_modified.json")
```

**Output:**
```
=== Example 6: Load Workflow from File ===
Loaded workflow with 9 nodes
Saved modified workflow to workflow_modified.json
```

**Use Cases:**
- Template-based workflow generation
- Workflow version control
- Batch modifications
- Workflow sharing

### Example 7: Concurrent Execution

Run multiple workflows in parallel with timeout:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

results := make(chan *comfyui.ExecutionResult, 3)
errors := make(chan error, 3)

for i := 0; i < 3; i++ {
    go func(index int) {
        workflow := buildSimpleWorkflow()
        workflow.SetNodeInput("6", "seed", 10000+index)
        
        result, _ := client.QueuePrompt(ctx, workflow, nil)
        execResult, _ := client.WaitForCompletion(ctx, result.PromptID)
        
        results <- execResult
    }(i)
}

// Collect results
for completed := 0; completed < 3; {
    select {
    case result := <-results:
        completed++
        fmt.Printf("  ✓ Workflow %d completed in %v\n", completed, result.Duration)
    case <-ctx.Done():
        fmt.Println("  ✗ Timeout reached")
        return
    }
}
```

**Output:**
```
=== Example 7: Concurrent Execution ===
  ✓ Workflow 1 completed in 12.3s
  ✓ Workflow 2 completed in 11.9s
  ✓ Workflow 3 completed in 12.5s
```

**Benefits:**
- Faster batch processing
- Better resource utilization
- Timeout protection
- Error isolation

## 🔧 Advanced Techniques

### Custom Upload Options

```go
uploadedImage, err := client.UploadImage(ctx, "input.png", comfyui.UploadOptions{
    Type:      "input",        // or "temp"
    Subfolder: "my_images",    // organize in subfolders
    Overwrite: true,           // replace existing
})
```

### Workflow Cloning

```go
baseWorkflow := buildSimpleWorkflow()

// Create variations
for i := 0; i < 10; i++ {
    workflow, _ := baseWorkflow.Clone()
    workflow.SetNodeInput("6", "seed", 1000+i)
    workflow.SetNodeInput("6", "steps", 20+i*2)
    // Submit...
}
```

### Dynamic Node Modification

```go
workflow := buildSimpleWorkflow()

// Change checkpoint
workflow.SetNodeInput("4", "ckpt_name", "dreamshaper_8.safetensors")

// Adjust sampling
workflow.SetNodeInput("3", "steps", 30)
workflow.SetNodeInput("3", "cfg", 7.5)
workflow.SetNodeInput("3", "sampler_name", "dpmpp_2m")

// Update prompts
workflow.SetNodeInput("6", "text", "new positive prompt")
workflow.SetNodeInput("7", "text", "new negative prompt")
```

### Error Recovery

```go
for i := 0; i < maxRetries; i++ {
    result, err := client.QueuePrompt(ctx, workflow, nil)
    if err == nil {
        break
    }
    
    log.Printf("Attempt %d failed: %v", i+1, err)
    time.Sleep(time.Second * time.Duration(i+1))
}
```

### Progress Tracking

```go
result, _ := client.QueuePrompt(ctx, workflow, nil)

// Monitor with timeout
ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
defer cancel()

execResult, err := client.WaitForCompletion(ctx, result.PromptID)
if err != nil {
    log.Printf("Execution failed or timed out: %v", err)
}
```

## 📊 Workflow Structures

### Simple Text-to-Image

```go
workflow := comfyui.Workflow{
    "4": comfyui.Node{
        ClassType: "CheckpointLoaderSimple",
        Inputs: map[string]interface{}{
            "ckpt_name": "v1-5-pruned-emaonly.safetensors",
        },
    },
    "5": comfyui.Node{
        ClassType: "EmptyLatentImage",
        Inputs: map[string]interface{}{
            "width": 512, "height": 512, "batch_size": 1,
        },
    },
    // ... more nodes
}
```

### Image-to-Image

```go
workflow := buildSimpleWorkflow()

// Add image loading
workflow["10"] = comfyui.Node{
    ClassType: "LoadImage",
    Inputs: map[string]interface{}{
        "image": uploadedImageName,
    },
}

// Add VAE encoding
workflow["11"] = comfyui.Node{
    ClassType: "VAEEncode",
    Inputs: map[string]interface{}{
        "pixels": []interface{}{"10", 0},
        "vae":    []interface{}{"4", 2},
    },
}

// Update KSampler
workflow.SetNodeInput("3", "latent_image", []interface{}{"11", 0})
workflow.SetNodeInput("3", "denoise", 0.75)
```

## 🎯 Use Cases

### 1. Batch Image Generation

Generate multiple variations of an image:

```go
prompts := []string{
    "a cat in a garden",
    "a dog on a beach",
    "a bird in the sky",
}

for _, prompt := range prompts {
    workflow := buildSimpleWorkflow()
    workflow.SetNodeInput("6", "text", prompt)
    client.QueuePrompt(ctx, workflow, nil)
}
```

### 2. A/B Testing

Compare different models or settings:

```go
models := []string{
    "v1-5-pruned-emaonly.safetensors",
    "dreamshaper_8.safetensors",
    "realisticVisionV51_v51VAE.safetensors",
}

for _, model := range models {
    workflow := buildSimpleWorkflow()
    workflow.SetNodeInput("4", "ckpt_name", model)
    // Submit and compare results...
}
```

### 3. Automated Processing Pipeline

Process uploaded images automatically:

```go
imageFiles := []string{"img1.png", "img2.png", "img3.png"}

for _, imgFile := range imageFiles {
    // Upload
    uploaded, _ := client.UploadImage(ctx, imgFile, options)
    
    // Process
    workflow := buildImg2ImgWorkflow(uploaded.Name)
    result, _ := client.QueuePrompt(ctx, workflow, nil)
    
    // Wait and save
    execResult, _ := client.WaitForCompletion(ctx, result.PromptID)
    // Save results...
}
```

### 4. Queue Monitoring Service

Monitor queue and send alerts:

```go
ticker := time.NewTicker(30 * time.Second)
for range ticker.C {
    queue, _ := client.GetQueue(ctx)
    
    if len(queue.QueuePending) > 10 {
        sendAlert("Queue is getting long: %d items", len(queue.QueuePending))
    }
}
```

## 🐛 Troubleshooting

### Image Upload Failed

```
Failed to upload image: file not found
```

**Solution:** Ensure the image file exists and path is correct.

### Workflow Clone Error

```
Failed to clone workflow
```

**Solution:** Make sure the workflow is valid JSON-serializable.

### Timeout in Batch Processing

```
context deadline exceeded
```

**Solution:** Increase timeout or process in smaller batches.

### Node Not Found

```
Failed to get object info: node not found
```

**Solution:** Check node name spelling and ensure ComfyUI has the node installed.

### Concurrent Execution Issues

```
Too many open connections
```

**Solution:** Limit concurrent goroutines or use a worker pool.

## 💡 Best Practices

### 1. Use Context Timeouts

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
```

### 2. Handle Errors Gracefully

```go
result, err := client.QueuePrompt(ctx, workflow, nil)
if err != nil {
    log.Printf("Failed to queue: %v", err)
    // Retry or skip
    continue
}
```

### 3. Limit Concurrent Executions

```go
semaphore := make(chan struct{}, 3) // Max 3 concurrent

for _, workflow := range workflows {
    semaphore <- struct{}{} // Acquire
    go func(wf comfyui.Workflow) {
        defer func() { <-semaphore }() // Release
        // Process workflow...
    }(workflow)
}
```

### 4. Clean Up Resources

```go
defer ws.Close()
defer cancel()
defer file.Close()
```

### 5. Validate Before Submission

```go
if len(workflow) == 0 {
    return errors.New("empty workflow")
}

// Check required nodes exist
if _, ok := workflow["4"]; !ok {
    return errors.New("missing checkpoint loader")
}
```

## 📚 API Reference

### Key Functions Used

#### Image Operations
- `UploadImage(ctx, path, options)` - Upload image file
- `SaveImage(ctx, image, path)` - Save generated image

#### Workflow Operations
- `QueuePrompt(ctx, workflow, options)` - Submit workflow
- `WaitForCompletion(ctx, promptID)` - Wait for execution
- `LoadWorkflowFromFile(path)` - Load workflow from JSON
- `SaveWorkflowToFile(workflow, path)` - Save workflow to JSON
- `workflow.Clone()` - Clone workflow
- `workflow.SetNodeInput(nodeID, input, value)` - Modify node input

#### Queue & History
- `GetQueue(ctx)` - Get current queue
- `GetHistory(ctx, promptID)` - Get execution history

#### Node Information
- `GetObjectInfo(ctx, nodeClass)` - Get node definition

## 🎓 Learning Points

This example teaches:

1. **File Operations** - Upload and save images
2. **Workflow Manipulation** - Clone and modify workflows
3. **Batch Processing** - Process multiple items efficiently
4. **Concurrency** - Run workflows in parallel
5. **Context Management** - Use timeouts and cancellation
6. **Error Handling** - Robust error handling patterns
7. **Resource Management** - Proper cleanup and resource limits
8. **API Integration** - Complete API usage patterns

## 🔗 Related Examples

- **[basic](../basic/)** - Basic workflow submission
- **[websocket](../websocket/)** - Real-time event monitoring
- **[progress](../progress/)** - Visual progress tracking

## 📝 Next Steps

After mastering advanced features:

1. Build your own automation pipelines
2. Integrate with web services or APIs
3. Create batch processing tools
4. Develop monitoring dashboards
5. Explore custom node development

---

**Happy coding!** 🚀

For questions or issues, please refer to the main [README](../../README.md) or check the [ComfyUI documentation](https://github.com/comfyanonymous/ComfyUI).

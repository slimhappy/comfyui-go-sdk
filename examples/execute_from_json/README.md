# Execute Workflow from JSON File

This example demonstrates how to load a ComfyUI workflow from a JSON file and execute it using the Go SDK.

## Features

✅ **Load workflows from JSON files** - Support standard ComfyUI API format  
✅ **Dynamic parameter modification** - Change parameters at runtime  
✅ **Real-time progress monitoring** - Track execution status  
✅ **Automatic result retrieval** - Download and save generated images  
✅ **Error handling** - Comprehensive error reporting  
✅ **Command-line interface** - Easy to use CLI tool  

---

## Quick Start

### 1. Build the Example

```bash
cd /data/comfyui-go-sdk
make build-execute-json
```

Or manually:

```bash
go build -o bin/execute_from_json examples/execute_from_json/main.go
```

### 2. Prepare Your Workflow JSON

You can get a workflow JSON file in two ways:

#### Method 1: Export from ComfyUI Web Interface

1. Open ComfyUI web interface
2. Create or load your workflow
3. Click **File → Export (API Format)**
4. Save the JSON file

#### Method 2: Use the Example Workflow

```bash
# Use the provided example workflow
cp examples/execute_from_json/workflow_example.json my_workflow.json
```

### 3. Run the Example

```bash
# Basic execution
./bin/execute_from_json my_workflow.json

# With custom parameters
./bin/execute_from_json my_workflow.json seed=12345 steps=30 cfg=7.5

# With custom prompts
./bin/execute_from_json my_workflow.json prompt="beautiful landscape" negative="blurry"
```

---

## Usage

### Command Syntax

```bash
./execute_from_json <workflow.json> [parameters...]
```

### Parameters

You can override workflow parameters using `key=value` format:

| Parameter | Description | Example |
|-----------|-------------|---------|
| `seed=<number>` | Set random seed for reproducibility | `seed=12345` |
| `steps=<number>` | Set number of sampling steps | `steps=30` |
| `cfg=<number>` | Set CFG (Classifier Free Guidance) scale | `cfg=7.5` |
| `prompt=<text>` | Set positive prompt text | `prompt="beautiful sunset"` |
| `negative=<text>` | Set negative prompt text | `negative="blurry, low quality"` |

---

## Examples

### Example 1: Basic Execution

```bash
./bin/execute_from_json workflow.json
```

**Output:**
```
╔════════════════════════════════════════════════════════════════════╗
║         ComfyUI Go SDK - Execute Workflow from JSON File         ║
╚════════════════════════════════════════════════════════════════════╝

🔍 Checking ComfyUI server status...
✅ ComfyUI server is running

📂 Loading workflow from: workflow.json
✅ Workflow loaded successfully
   Total nodes: 7
   Node types:
     - KSampler: 1
     - CheckpointLoaderSimple: 1
     - EmptyLatentImage: 1
     - CLIPTextEncode: 2
     - VAEDecode: 1
     - SaveImage: 1

🚀 Submitting workflow to ComfyUI...
✅ Workflow queued successfully!
   Prompt ID: abc123-def456-ghi789
   Queue Position: 1

⏳ Monitoring execution progress...
✅ Completed in 15.3 seconds

📥 Retrieving execution results...
📊 Execution ID: abc123-def456-ghi789
   Status: success
   Outputs:
     Node 9:
       Images: 1
         [1] ComfyUI_00001.png (type: output, subfolder: )
         💾 Saved to: output/ComfyUI_00001.png

╔════════════════════════════════════════════════════════════════════╗
║                    ✅ Execution Complete!                          ║
╚════════════════════════════════════════════════════════════════════╝
```

### Example 2: Custom Parameters

```bash
./bin/execute_from_json workflow.json seed=42 steps=25 cfg=8.0
```

**Output:**
```
...
🔧 Applying custom parameters...
   ✓ Set seed=42 for node 3
   ✓ Set steps=25 for node 3
   ✓ Set cfg=8.0 for node 3
...
```

### Example 3: Custom Prompts

```bash
./bin/execute_from_json workflow.json \
  prompt="a beautiful mountain landscape at sunset, highly detailed" \
  negative="blurry, low quality, distorted"
```

### Example 4: Batch Processing

```bash
#!/bin/bash
# Generate multiple images with different seeds

for seed in 100 200 300 400 500; do
  echo "Generating image with seed $seed..."
  ./bin/execute_from_json workflow.json seed=$seed
  sleep 2
done
```

---

## Workflow JSON Format

The workflow JSON file follows ComfyUI's API format. Each node is identified by a unique ID (string) and contains:

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

### Example Workflow Structure

```json
{
  "3": {
    "class_type": "KSampler",
    "inputs": {
      "seed": 12345,
      "steps": 20,
      "cfg": 8.0,
      "sampler_name": "euler",
      "scheduler": "normal",
      "denoise": 1.0,
      "model": ["4", 0],
      "positive": ["6", 0],
      "negative": ["7", 0],
      "latent_image": ["5", 0]
    }
  },
  "4": {
    "class_type": "CheckpointLoaderSimple",
    "inputs": {
      "ckpt_name": "v1-5-pruned-emaonly.safetensors"
    }
  },
  "5": {
    "class_type": "EmptyLatentImage",
    "inputs": {
      "width": 512,
      "height": 512,
      "batch_size": 1
    }
  },
  "6": {
    "class_type": "CLIPTextEncode",
    "inputs": {
      "text": "beautiful landscape",
      "clip": ["4", 1]
    }
  },
  "7": {
    "class_type": "CLIPTextEncode",
    "inputs": {
      "text": "blurry, bad quality",
      "clip": ["4", 1]
    }
  },
  "8": {
    "class_type": "VAEDecode",
    "inputs": {
      "samples": ["3", 0],
      "vae": ["4", 2]
    }
  },
  "9": {
    "class_type": "SaveImage",
    "inputs": {
      "filename_prefix": "ComfyUI",
      "images": ["8", 0]
    }
  }
}
```

---

## Programmatic Usage

You can also use the SDK programmatically in your own Go code:

### Method 1: Using QueuePromptFromFile

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
    
    // Execute workflow directly from file
    resp, err := client.QueuePromptFromFile(ctx, "workflow.json", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Queued with ID: %s", resp.PromptID)
}
```

### Method 2: Load, Modify, and Execute

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
    
    // Load workflow from file
    workflow, err := comfyui.LoadWorkflowFromFile("workflow.json")
    if err != nil {
        log.Fatal(err)
    }
    
    // Modify parameters
    workflow.SetNodeInput("3", "seed", 12345)
    workflow.SetNodeInput("3", "steps", 30)
    workflow.SetNodeInput("6", "text", "beautiful sunset")
    
    // Execute modified workflow
    resp, err := client.QueuePrompt(ctx, workflow, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Queued with ID: %s", resp.PromptID)
}
```

### Method 3: Advanced - Clone and Batch Process

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    comfyui "github.com/yourusername/comfyui-go-sdk"
)

func main() {
    client := comfyui.NewClient("http://127.0.0.1:8188")
    ctx := context.Background()
    
    // Load base workflow
    baseWorkflow, err := comfyui.LoadWorkflowFromFile("workflow.json")
    if err != nil {
        log.Fatal(err)
    }
    
    // Generate multiple variations
    seeds := []int{100, 200, 300, 400, 500}
    
    for i, seed := range seeds {
        // Clone workflow for each variation
        workflow, err := baseWorkflow.Clone()
        if err != nil {
            log.Fatal(err)
        }
        
        // Modify seed
        workflow.SetNodeInput("3", "seed", seed)
        
        // Queue execution
        resp, err := client.QueuePrompt(ctx, workflow, nil)
        if err != nil {
            log.Printf("Failed to queue variation %d: %v", i, err)
            continue
        }
        
        fmt.Printf("Variation %d queued with ID: %s\n", i+1, resp.PromptID)
    }
}
```

---

## Output

Generated images are automatically saved to the `output/` directory in the current working directory.

```
output/
├── ComfyUI_00001.png
├── ComfyUI_00002.png
└── ComfyUI_00003.png
```

---

## Error Handling

The example includes comprehensive error handling:

### Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| `Workflow file not found` | File path is incorrect | Check file path and permissions |
| `Failed to unmarshal workflow` | Invalid JSON format | Validate JSON syntax |
| `Server check failed` | ComfyUI not running | Start ComfyUI server |
| `Node errors` | Invalid node configuration | Check node inputs and connections |
| `Execution timeout` | Workflow takes too long | Increase timeout or optimize workflow |

### Example Error Output

```
❌ Failed to load workflow: failed to unmarshal workflow: invalid character '}' looking for beginning of object key string
```

---

## Advanced Features

### 1. Workflow Validation

```go
workflow, _ := comfyui.LoadWorkflowFromFile("workflow.json")

if err := workflow.Validate(); err != nil {
    log.Fatalf("Invalid workflow: %v", err)
}
```

### 2. Node Inspection

```go
// Get all node IDs
nodeIDs := workflow.NodeIDs()
fmt.Printf("Nodes: %v\n", nodeIDs)

// Find nodes by type
samplers := workflow.NodesByClass("KSampler")
fmt.Printf("Found %d KSampler nodes\n", len(samplers))

// Get specific node
if node, ok := workflow.GetNode("3"); ok {
    fmt.Printf("Node 3 class: %s\n", node.ClassType)
}
```

### 3. Save Modified Workflow

```go
workflow, _ := comfyui.LoadWorkflowFromFile("workflow.json")

// Modify workflow
workflow.SetNodeInput("3", "seed", 99999)

// Save to new file
if err := comfyui.SaveWorkflowToFile(workflow, "modified_workflow.json"); err != nil {
    log.Fatal(err)
}
```

---

## Integration with Other Examples

This example can be combined with other SDK features:

### With WebSocket Monitoring

```go
// Start WebSocket connection for real-time updates
ws, _ := client.NewWebSocket(ctx)
defer ws.Close()

events := ws.Subscribe()
go func() {
    for event := range events {
        fmt.Printf("Event: %s\n", event.Type)
    }
}()

// Execute workflow
resp, _ := client.QueuePromptFromFile(ctx, "workflow.json", nil)
```

### With Progress Tracking

See the [progress example](../progress/README.md) for detailed progress tracking implementation.

---

## Troubleshooting

### Issue: "Connection refused"

**Solution:** Make sure ComfyUI is running:
```bash
# Check if ComfyUI is running
curl http://127.0.0.1:8188/system_stats

# Start ComfyUI if not running
cd /data/ComfyUI
python main.py
```

### Issue: "Node errors: missing checkpoint"

**Solution:** Ensure the checkpoint file exists in ComfyUI's models directory:
```bash
ls /data/ComfyUI/models/checkpoints/
```

Update the workflow JSON with the correct checkpoint name.

### Issue: Parameters not applied

**Solution:** Check node IDs in your workflow. The example assumes:
- Node "3" is KSampler
- Node "6" is positive prompt (CLIPTextEncode)
- Node "7" is negative prompt (CLIPTextEncode)

Adjust the parameter application logic if your workflow uses different node IDs.

---

## Performance Tips

1. **Reuse workflows**: Load once, clone for variations
2. **Batch processing**: Queue multiple workflows before waiting
3. **Optimize parameters**: Reduce steps/resolution for faster generation
4. **Monitor queue**: Check queue status to avoid overloading

---

## Related Documentation

- [Main README](../../README.md) - SDK overview
- [Basic Example](../basic/README.md) - Basic SDK usage
- [WebSocket Example](../websocket/README.md) - Real-time event monitoring
- [Progress Example](../progress/README.md) - Detailed progress tracking
- [Advanced Example](../advanced/README.md) - Advanced features

---

## API Reference

### Client Methods

```go
// Load and execute workflow from file
func (c *Client) QueuePromptFromFile(ctx context.Context, filepath string, extraData map[string]interface{}) (*QueuePromptResponse, error)

// Execute workflow
func (c *Client) QueuePrompt(ctx context.Context, workflow Workflow, extraData map[string]interface{}) (*QueuePromptResponse, error)

// Get execution history
func (c *Client) GetHistory(ctx context.Context, promptID string) (map[string]HistoryItem, error)

// Download generated image
func (c *Client) DownloadImage(ctx context.Context, image ImageOutput, outputPath string) error
```

### Workflow Methods

```go
// Load workflow from JSON file
func LoadWorkflowFromFile(filepath string) (Workflow, error)

// Save workflow to JSON file
func SaveWorkflowToFile(workflow Workflow, filepath string) error

// Modify node input
func (w Workflow) SetNodeInput(nodeID string, inputName string, value interface{}) error

// Get node input
func (w Workflow) GetNodeInput(nodeID string, inputName string) (interface{}, error)

// Clone workflow
func (w Workflow) Clone() (Workflow, error)

// Validate workflow
func (w Workflow) Validate() error

// Get all node IDs
func (w Workflow) NodeIDs() []string

// Get nodes by class type
func (w Workflow) NodesByClass(classType string) map[string]Node
```

---

## Contributing

Found a bug or have a feature request? Please open an issue on GitHub!

---

## License

This example is part of the ComfyUI Go SDK project.

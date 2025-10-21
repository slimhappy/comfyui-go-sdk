# Progress Tracking Example

This example demonstrates how to track the real-time progress of ComfyUI workflow execution with a visual progress bar.

## Features

- **Visual Progress Bar**: ASCII-based progress bar showing completion percentage
- **Real-time Updates**: Updates progress in the same terminal line
- **Detailed Statistics**: Shows elapsed time, ETA, current node, and step count
- **Complete Workflow**: Submits a sample workflow and monitors its execution
- **Error Handling**: Displays errors with detailed information
- **Result Retrieval**: Automatically saves generated images after completion

## Usage

### Build and Run

```bash
# Build the example
cd /data/comfyui-go-sdk
make build

# Run the progress tracking demo
./bin/progress
```

Or run directly:

```bash
cd examples/progress
go run main.go
```

## What It Does

1. **Creates a Sample Workflow**: Generates a text-to-image workflow with:
   - Checkpoint loader
   - Positive and negative prompts
   - KSampler (20 steps)
   - VAE decoder
   - Image saver

2. **Submits to ComfyUI**: Queues the workflow for execution

3. **Monitors Progress**: Connects via WebSocket and displays:
   ```
   ⏳ [████████████████████░░░░░░░░░░░░░░░░░░░░] 50.0% | Step 10/20 | Node: KSampler | Time: 5s ETA: 5s
   ```

4. **Shows Completion**: Displays final statistics and saves images:
   ```
   ✅ Completed in 10.5s (Processed 7 nodes)
   ```

## Progress Display Format

The progress bar shows:
- **Visual Bar**: `[████████░░░░░░░░]` - 40 characters wide
- **Percentage**: Current completion percentage
- **Steps**: Current step / Total steps
- **Node**: Currently executing node name
- **Time**: Elapsed time since start
- **ETA**: Estimated time to completion (calculated dynamically)

## Example Output

```
╔════════════════════════════════════════════════════════════╗
║        ComfyUI Go SDK - Progress Tracking Demo           ║
╚════════════════════════════════════════════════════════════╝

📝 Creating sample workflow...
📤 Submitting workflow to ComfyUI...
✓ Workflow queued successfully (ID: abc123...)

🚀 Monitoring workflow execution: abc123...

⏳ [████████████████████████████████████████] 100.0% | Step 20/20 | Node: SaveImage | Time: 12s
✅ Completed in 12.3s (Processed 7 nodes)

📥 Retrieving execution results...

📊 Execution Summary:
  • Status: success
  • Generated Images: 1

💾 Saving generated images...
  ✓ Saved: output_9_0.png

✨ Demo completed!
```

## Code Structure

### ProgressTracker

The `ProgressTracker` struct maintains the state of workflow execution:

```go
type ProgressTracker struct {
    PromptID       string        // Workflow ID
    StartTime      time.Time     // Execution start time
    CurrentNode    string        // Currently executing node
    TotalNodes     int           // Total number of nodes
    CompletedNodes int           // Number of completed nodes
    CurrentStep    int           // Current step in sampling
    TotalSteps     int           // Total steps in sampling
    LastUpdate     time.Time     // Last update timestamp
    IsCompleted    bool          // Completion flag
    HasError       bool          // Error flag
    ErrorMessage   string        // Error details
}
```

### Key Functions

- **`MonitorProgress()`**: Main monitoring loop that connects to WebSocket and tracks progress
- **`handleProgressMessage()`**: Processes WebSocket messages and updates tracker
- **`DrawProgressBar()`**: Renders the visual progress bar
- **`PrintProgress()`**: Displays the current progress state
- **`CreateSampleWorkflow()`**: Builds a sample workflow for demonstration

## Integration in Your Code

To use progress tracking in your own application:

```go
import (
    "context"
    comfyui "github.com/yourusername/comfyui-go-sdk"
)

func main() {
    client := comfyui.NewClient("http://127.0.0.1:8188")
    
    // Submit your workflow
    result, err := client.QueuePrompt(context.Background(), workflow, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Monitor progress (copy the MonitorProgress function from this example)
    if err := MonitorProgress(context.Background(), client, result.PromptID); err != nil {
        log.Printf("Error: %v", err)
    }
}
```

## Requirements

- ComfyUI server running on `http://127.0.0.1:8188`
- A checkpoint model (e.g., `v1-5-pruned-emaonly.safetensors`)
- Terminal with ANSI escape code support (for progress bar updates)

## Customization

You can customize the progress display by modifying:

- **Progress bar width**: Change the `width` parameter in `DrawProgressBar()`
- **Update frequency**: Adjust the WebSocket message handling
- **Display format**: Modify the `PrintProgress()` function
- **Workflow**: Replace `CreateSampleWorkflow()` with your own workflow

## Tips

1. **Terminal Support**: The progress bar uses ANSI escape codes. If your terminal doesn't support them, you'll see raw escape sequences.

2. **Long-running Workflows**: For workflows with many steps, the ETA calculation becomes more accurate over time.

3. **Multiple Nodes**: The tracker counts completed nodes to give you an overall sense of progress beyond just the current sampling step.

4. **Error Recovery**: If an error occurs, the tracker will display it and stop monitoring.

## See Also

- [Basic Example](../basic/) - Simple workflow submission
- [WebSocket Example](../websocket/) - Raw WebSocket message monitoring
- [Advanced Example](../advanced/) - Advanced SDK features

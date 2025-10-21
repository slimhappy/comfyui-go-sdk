# Progress Tracking Quick Reference

## Quick Start

```bash
# Build and run
cd /data/comfyui-go-sdk
make build
./bin/progress

# Or use the demo script
./run_progress_demo.sh
```

## Visual Output

The progress tracker displays real-time information in a single line:

```
⏳ [████████████████████░░░░░░░░░░░░░░░░░░░░] 50.0% | Step 10/20 | Node: KSampler | Time: 5s ETA: 5s
```

### Progress Bar Elements

| Element | Description | Example |
|---------|-------------|---------|
| 🔄 Icon | Status indicator | ⏳ (running), ✅ (done), ❌ (error) |
| Bar | Visual progress (40 chars) | `[████████░░░░░░░░]` |
| Percentage | Completion percentage | `50.0%` |
| Steps | Current/Total steps | `Step 10/20` |
| Node | Current node name | `Node: KSampler` |
| Time | Elapsed time | `Time: 5s` |
| ETA | Estimated time remaining | `ETA: 5s` |

## Integration Example

### Minimal Integration

```go
package main

import (
    "context"
    "fmt"
    comfyui "github.com/yourusername/comfyui-go-sdk"
)

func main() {
    client := comfyui.NewClient("http://127.0.0.1:8188")
    
    // Submit workflow
    result, _ := client.QueuePrompt(context.Background(), workflow, nil)
    
    // Monitor with progress bar
    MonitorProgress(context.Background(), client, result.PromptID)
}

// Copy MonitorProgress, ProgressTracker, and helper functions 
// from examples/progress/main.go
```

### Custom Progress Handler

```go
type CustomProgressHandler struct {
    onProgress func(current, total int, percentage float64)
    onComplete func(duration time.Duration)
    onError    func(err error)
}

func (h *CustomProgressHandler) HandleMessage(msg comfyui.WebSocketMessage) {
    switch msg.Type {
    case string(comfyui.MessageTypeProgress):
        data, _ := msg.GetProgressData()
        percentage := float64(data.Value) / float64(data.Max) * 100
        h.onProgress(data.Value, data.Max, percentage)
        
    case string(comfyui.MessageTypeExecuting):
        data, _ := msg.GetExecutingData()
        if data.Node == nil {
            h.onComplete(time.Since(startTime))
        }
        
    case string(comfyui.MessageTypeError):
        data, _ := msg.GetErrorData()
        h.onError(fmt.Errorf("%s: %s", data.ExceptionType, data.ExceptionMessage))
    }
}
```

## Progress Tracking Patterns

### Pattern 1: Simple Progress Bar

```go
for msg := range ws.Messages() {
    if msg.Type == string(comfyui.MessageTypeProgress) {
        data, _ := msg.GetProgressData()
        fmt.Printf("\rProgress: %d/%d", data.Value, data.Max)
    }
}
```

### Pattern 2: Percentage Display

```go
for msg := range ws.Messages() {
    if msg.Type == string(comfyui.MessageTypeProgress) {
        data, _ := msg.GetProgressData()
        pct := float64(data.Value) / float64(data.Max) * 100
        fmt.Printf("\r%.1f%% complete", pct)
    }
}
```

### Pattern 3: Visual Bar

```go
func drawBar(current, total, width int) string {
    filled := int(float64(current) / float64(total) * float64(width))
    return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

for msg := range ws.Messages() {
    if msg.Type == string(comfyui.MessageTypeProgress) {
        data, _ := msg.GetProgressData()
        bar := drawBar(data.Value, data.Max, 40)
        fmt.Printf("\r[%s]", bar)
    }
}
```

### Pattern 4: Multi-line Status

```go
type Status struct {
    CurrentNode string
    Progress    string
    Elapsed     time.Duration
}

func printStatus(s Status) {
    fmt.Print("\033[2J\033[H") // Clear screen
    fmt.Printf("Node: %s\n", s.CurrentNode)
    fmt.Printf("Progress: %s\n", s.Progress)
    fmt.Printf("Time: %s\n", s.Elapsed)
}
```

## WebSocket Message Types

### Progress Messages

```go
type ProgressData struct {
    Value int `json:"value"` // Current step
    Max   int `json:"max"`   // Total steps
}
```

### Executing Messages

```go
type ExecutingData struct {
    Node     *string `json:"node"`      // Current node (nil when done)
    PromptID string  `json:"prompt_id"` // Workflow ID
}
```

### Executed Messages

```go
type ExecutedData struct {
    Node     string                 `json:"node"`
    PromptID string                 `json:"prompt_id"`
    Output   map[string]interface{} `json:"output"`
}
```

### Error Messages

```go
type ErrorData struct {
    PromptID         string   `json:"prompt_id"`
    NodeID           string   `json:"node_id"`
    NodeType         string   `json:"node_type"`
    ExceptionType    string   `json:"exception_type"`
    ExceptionMessage string   `json:"exception_message"`
    Traceback        []string `json:"traceback"`
}
```

## Tips & Tricks

### 1. Clear Line Before Update

```go
func clearLine() {
    fmt.Print("\r\033[K")
}

clearLine()
fmt.Printf("New progress: %d%%", percentage)
```

### 2. Calculate ETA

```go
elapsed := time.Since(startTime)
if percentage > 0 {
    totalEstimated := float64(elapsed) / (percentage / 100)
    eta := time.Duration(totalEstimated) - elapsed
    fmt.Printf("ETA: %s", eta.Round(time.Second))
}
```

### 3. Track Multiple Metrics

```go
type Metrics struct {
    NodesCompleted int
    StepsCompleted int
    ImagesGenerated int
    StartTime      time.Time
}

func (m *Metrics) Speed() float64 {
    elapsed := time.Since(m.StartTime).Seconds()
    return float64(m.StepsCompleted) / elapsed
}
```

### 4. Graceful Shutdown

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

go func() {
    <-sigChan
    fmt.Println("\nShutting down...")
    cancel()
}()
```

### 5. Progress Persistence

```go
type ProgressLog struct {
    PromptID  string
    Timestamp time.Time
    Progress  int
    Total     int
}

func saveProgress(log ProgressLog) {
    data, _ := json.Marshal(log)
    ioutil.WriteFile("progress.json", data, 0644)
}
```

## Troubleshooting

### Progress Not Updating

**Problem**: Progress bar doesn't update
**Solution**: Ensure you're using `\r` (carriage return) not `\n` (newline)

```go
// Wrong
fmt.Printf("Progress: %d%%\n", pct)

// Correct
fmt.Printf("\rProgress: %d%%", pct)
```

### Terminal Doesn't Support ANSI

**Problem**: See escape codes like `^[[K` in output
**Solution**: Use simple text output without escape codes

```go
// Instead of clearing line
fmt.Printf("\rProgress: %d/%d", current, total)
```

### WebSocket Disconnects

**Problem**: Connection drops during long workflows
**Solution**: Implement reconnection logic

```go
func connectWithRetry(client *comfyui.Client, ctx context.Context) (*comfyui.WebSocketClient, error) {
    for i := 0; i < 3; i++ {
        ws, err := client.ConnectWebSocket(ctx)
        if err == nil {
            return ws, nil
        }
        time.Sleep(time.Second * time.Duration(i+1))
    }
    return nil, fmt.Errorf("failed to connect after retries")
}
```

### Missing Progress Events

**Problem**: Some progress updates are skipped
**Solution**: This is normal - ComfyUI sends updates at intervals

```go
// Don't rely on every single step
// Instead, track overall progress
```

## Performance Considerations

### 1. Update Frequency

Limit update frequency to avoid terminal flicker:

```go
lastUpdate := time.Now()
updateInterval := 100 * time.Millisecond

if time.Since(lastUpdate) > updateInterval {
    printProgress()
    lastUpdate = time.Now()
}
```

### 2. Buffer Output

Use buffered output for better performance:

```go
buf := bufio.NewWriter(os.Stdout)
defer buf.Flush()

fmt.Fprintf(buf, "\rProgress: %d%%", pct)
buf.Flush()
```

### 3. Minimize Allocations

Reuse strings and buffers:

```go
var progressBuf strings.Builder
progressBuf.Grow(100)

progressBuf.Reset()
progressBuf.WriteString("\r[")
progressBuf.WriteString(bar)
progressBuf.WriteString("]")
fmt.Print(progressBuf.String())
```

## See Also

- [examples/progress/main.go](examples/progress/main.go) - Complete implementation
- [examples/progress/README.md](examples/progress/README.md) - Detailed documentation
- [examples/websocket/main.go](examples/websocket/main.go) - WebSocket basics
- [README.md](README.md) - Full SDK documentation

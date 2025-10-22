# WebSocket Example - ComfyUI Go SDK

This example demonstrates real-time event monitoring using WebSocket connections to ComfyUI. It listens to all events and displays them in a human-readable format.

## 📋 Features

This example showcases:

1. **WebSocket Connection** - Establish persistent connection to ComfyUI
2. **Real-time Event Monitoring** - Listen to all ComfyUI events
3. **Event Type Handling** - Parse and display different message types
4. **Graceful Shutdown** - Handle Ctrl+C interrupts properly
5. **Error Handling** - Robust error handling for WebSocket operations

## 🚀 Quick Start

### Prerequisites

- ComfyUI server running at `http://127.0.0.1:8188`
- No workflow submission required - just monitors events

### Build and Run

```bash
# From the SDK root directory
cd /data/comfyui-go-sdk
make build
./bin/websocket
```

Or run directly:

```bash
cd examples/websocket
go run main.go
```

### Stop the Monitor

Press `Ctrl+C` to gracefully shutdown the WebSocket connection.

## 📖 What This Example Does

### Connection

Establishes a WebSocket connection and listens for all events:

```go
ws, err := client.ConnectWebSocket(ctx)
defer ws.Close()

for msg := range ws.Messages() {
    handleMessage(msg)
}
```

### Event Types Monitored

The example handles these ComfyUI event types:

#### 1. **STATUS** - Queue Status Updates

Displays the number of items remaining in the queue.

```
[STATUS] Queue remaining: 3
```

#### 2. **EXECUTING** - Node Execution Events

Shows which node is currently being executed.

```
[EXECUTING] Prompt abc123, Node 4
[EXECUTING] Prompt abc123 completed
```

#### 3. **PROGRESS** - Sampling Progress

Displays progress during sampling operations (e.g., KSampler).

```
[PROGRESS] 5/20 (25.0%)
[PROGRESS] 10/20 (50.0%)
[PROGRESS] 20/20 (100.0%)
```

#### 4. **EXECUTED** - Node Completion

Indicates when a node finishes execution and shows output info.

```
[EXECUTED] Node 9 in prompt abc123
  → Produced 1 image(s)
```

#### 5. **CACHED** - Cached Execution

Shows when a node's execution is cached (skipped).

```
[CACHED] Node execution cached
```

#### 6. **ERROR** - Execution Errors

Displays detailed error information including traceback.

```
[ERROR] Prompt abc123, Node 3 (KSampler)
  Type: RuntimeError
  Message: CUDA out of memory
  Traceback:
    File "/path/to/file.py", line 123, in function
    ...
```

## 🎯 Use Cases

### 1. Development & Debugging

Monitor ComfyUI events while developing workflows:

```bash
# Terminal 1: Run the monitor
./bin/websocket

# Terminal 2: Submit workflows
./bin/basic
```

### 2. System Monitoring

Keep track of ComfyUI activity in production:

```bash
# Run in background
nohup ./bin/websocket > comfyui-events.log 2>&1 &
```

### 3. Integration Testing

Verify that workflows execute correctly:

```go
// Start monitoring in a goroutine
go func() {
    ws, _ := client.ConnectWebSocket(ctx)
    for msg := range ws.Messages() {
        // Verify expected events
    }
}()

// Submit test workflow
client.QueuePrompt(ctx, testWorkflow, nil)
```

### 4. Performance Analysis

Track execution times and identify bottlenecks:

```
[EXECUTING] Prompt abc123, Node 4  // Checkpoint loading
[EXECUTING] Prompt abc123, Node 3  // KSampler (slow)
[PROGRESS] 1/20 (5.0%)
[PROGRESS] 2/20 (10.0%)
...
```

## 🔧 Code Structure

### Main Function

```go
func main() {
    client := comfyui.NewClient("http://127.0.0.1:8188")
    
    // Setup context with cancellation
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Handle Ctrl+C
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-sigChan
        cancel()
    }()
    
    // Connect and listen
    ws, _ := client.ConnectWebSocket(ctx)
    defer ws.Close()
    
    for {
        select {
        case msg := <-ws.Messages():
            handleMessage(msg)
        case err := <-ws.Errors():
            log.Printf("Error: %v", err)
        case <-ctx.Done():
            return
        }
    }
}
```

### Message Handler

```go
func handleMessage(msg comfyui.WebSocketMessage) {
    switch msg.Type {
    case string(comfyui.MessageTypeStatus):
        // Handle status updates
    case string(comfyui.MessageTypeExecuting):
        // Handle execution events
    case string(comfyui.MessageTypeProgress):
        // Handle progress updates
    case string(comfyui.MessageTypeExecuted):
        // Handle completion events
    case string(comfyui.MessageTypeCached):
        // Handle cached executions
    case string(comfyui.MessageTypeError):
        // Handle errors
    default:
        // Handle unknown events
    }
}
```

## 📊 Example Output

### Typical Workflow Execution

```
Connecting to ComfyUI WebSocket...
Connected! Listening for events...
Press Ctrl+C to exit

[STATUS] Queue remaining: 1
[EXECUTING] Prompt 12345678-90ab-cdef-1234-567890abcdef, Node 4
[EXECUTING] Prompt 12345678-90ab-cdef-1234-567890abcdef, Node 5
[EXECUTING] Prompt 12345678-90ab-cdef-1234-567890abcdef, Node 6
[EXECUTING] Prompt 12345678-90ab-cdef-1234-567890abcdef, Node 7
[EXECUTING] Prompt 12345678-90ab-cdef-1234-567890abcdef, Node 3
[PROGRESS] 1/20 (5.0%)
[PROGRESS] 2/20 (10.0%)
[PROGRESS] 5/20 (25.0%)
[PROGRESS] 10/20 (50.0%)
[PROGRESS] 15/20 (75.0%)
[PROGRESS] 20/20 (100.0%)
[EXECUTED] Node 3 in prompt 12345678-90ab-cdef-1234-567890abcdef
[EXECUTING] Prompt 12345678-90ab-cdef-1234-567890abcdef, Node 8
[EXECUTED] Node 8 in prompt 12345678-90ab-cdef-1234-567890abcdef
[EXECUTING] Prompt 12345678-90ab-cdef-1234-567890abcdef, Node 9
[EXECUTED] Node 9 in prompt 12345678-90ab-cdef-1234-567890abcdef
  → Produced 1 image(s)
[EXECUTING] Prompt 12345678-90ab-cdef-1234-567890abcdef completed
[STATUS] Queue remaining: 0
```

### Error Scenario

```
[STATUS] Queue remaining: 1
[EXECUTING] Prompt abc123, Node 4
[EXECUTING] Prompt abc123, Node 3
[ERROR] Prompt abc123, Node 3 (KSampler)
  Type: RuntimeError
  Message: CUDA out of memory. Tried to allocate 2.00 GiB
  Traceback:
    File "/ComfyUI/execution.py", line 151, in recursive_execute
    File "/ComfyUI/nodes.py", line 1456, in sample
    RuntimeError: CUDA out of memory
[STATUS] Queue remaining: 0
```

### Graceful Shutdown

```
^C
Shutting down...
WebSocket closed
```

## 🎨 Customization

### Filter Specific Events

Only show progress updates:

```go
func handleMessage(msg comfyui.WebSocketMessage) {
    if msg.Type == string(comfyui.MessageTypeProgress) {
        data, _ := msg.GetProgressData()
        percentage := float64(data.Value) / float64(data.Max) * 100
        fmt.Printf("[PROGRESS] %d/%d (%.1f%%)\n", data.Value, data.Max, percentage)
    }
}
```

### Track Specific Prompt

Monitor only a specific workflow:

```go
targetPromptID := "abc123"

func handleMessage(msg comfyui.WebSocketMessage) {
    switch msg.Type {
    case string(comfyui.MessageTypeExecuting):
        data, _ := msg.GetExecutingData()
        if data.PromptID == targetPromptID {
            fmt.Printf("Target workflow executing: %s\n", *data.Node)
        }
    }
}
```

### Add Timestamps

Include timestamps in output:

```go
func handleMessage(msg comfyui.WebSocketMessage) {
    timestamp := time.Now().Format("15:04:05")
    fmt.Printf("[%s] ", timestamp)
    
    // ... rest of handling
}
```

### Log to File

Write events to a log file:

```go
logFile, _ := os.OpenFile("comfyui-events.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
defer logFile.Close()
logger := log.New(logFile, "", log.LstdFlags)

func handleMessage(msg comfyui.WebSocketMessage) {
    logger.Printf("[%s] %v\n", msg.Type, msg.Data)
}
```

## 🔍 Message Data Structures

### ExecutingData

```go
type ExecutingData struct {
    PromptID string
    Node     *string  // nil when execution completes
}
```

### ProgressData

```go
type ProgressData struct {
    Value int  // Current step
    Max   int  // Total steps
}
```

### ExecutedData

```go
type ExecutedData struct {
    PromptID string
    Node     string
    Output   map[string]interface{}
}
```

### ErrorData

```go
type ErrorData struct {
    PromptID         string
    NodeID           string
    NodeType         string
    ExceptionType    string
    ExceptionMessage string
    Traceback        []string
}
```

## 🐛 Troubleshooting

### Connection Failed

```
Failed to connect: dial tcp 127.0.0.1:8188: connection refused
```

**Solution:** Ensure ComfyUI is running and WebSocket endpoint is accessible.

### WebSocket Closed Unexpectedly

```
WebSocket closed
```

**Solution:** This can happen if ComfyUI restarts. The example will exit gracefully. Restart the monitor.

### No Events Received

**Solution:** Submit a workflow using the basic example or ComfyUI web interface to generate events.

### High CPU Usage

**Solution:** The example is very lightweight. If you see high CPU, check ComfyUI server itself.

## 💡 Tips

1. **Run in Separate Terminal** - Keep the monitor running while working with ComfyUI
2. **Pipe to grep** - Filter specific events: `./bin/websocket | grep PROGRESS`
3. **Use with tmux/screen** - Run in a persistent session
4. **Combine with Basic Example** - Monitor events while submitting workflows
5. **Check Error Details** - Error messages include full traceback for debugging

## 🔗 Integration Examples

### With Basic Example

```bash
# Terminal 1: Monitor events
./bin/websocket

# Terminal 2: Submit workflow
./bin/basic
```

### With Progress Example

```bash
# Terminal 1: Monitor raw events
./bin/websocket

# Terminal 2: Visual progress tracking
./bin/progress
```

### With Advanced Example

```bash
# Terminal 1: Monitor all events
./bin/websocket

# Terminal 2: Batch processing
./bin/advanced
```

## 📚 API Reference

### Key Functions Used

- `NewClient(url string)` - Create client
- `ConnectWebSocket(ctx)` - Establish WebSocket connection
- `ws.Messages()` - Channel for incoming messages
- `ws.Errors()` - Channel for WebSocket errors
- `ws.Close()` - Close the connection
- `msg.GetExecutingData()` - Parse executing event
- `msg.GetProgressData()` - Parse progress event
- `msg.GetExecutedData()` - Parse executed event
- `msg.GetErrorData()` - Parse error event

### Message Types

```go
const (
    MessageTypeStatus    = "status"
    MessageTypeExecuting = "executing"
    MessageTypeProgress  = "progress"
    MessageTypeExecuted  = "executed"
    MessageTypeCached    = "execution_cached"
    MessageTypeError     = "execution_error"
)
```

## 🎓 Learning Points

This example teaches:

1. **WebSocket Basics** - How to establish and maintain WebSocket connections
2. **Event-Driven Programming** - Handling asynchronous events
3. **Context Management** - Using context for cancellation
4. **Signal Handling** - Graceful shutdown on Ctrl+C
5. **Channel Operations** - Working with Go channels
6. **Type Assertions** - Parsing dynamic JSON data
7. **Error Handling** - Robust error handling patterns

## 🔗 Related Examples

- **[basic](../basic/)** - Submit workflows to generate events
- **[progress](../progress/)** - Visual progress tracking
- **[advanced](../advanced/)** - Advanced workflow management

## 📝 Next Steps

After understanding WebSocket monitoring:

1. Try **[progress example](../progress/)** for visual progress bars
2. Explore **[advanced example](../advanced/)** for complex workflows
3. Build your own event handlers for custom use cases
4. Integrate WebSocket monitoring into your applications

---

**Happy monitoring!** 🎧

For questions or issues, please refer to the main [README](../../README.md) or check the [ComfyUI documentation](https://github.com/comfyanonymous/ComfyUI).

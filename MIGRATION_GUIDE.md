# Migration Guide - PromptArray Update

## Overview

The `HistoryItem.Prompt` field type has been changed from `Workflow` to `PromptArray` to correctly parse ComfyUI's history API response format.

---

## What Changed

### Before (Old Code)
```go
type HistoryItem struct {
    Prompt  Workflow              `json:"prompt"`  // ❌ Incorrect - causes unmarshal error
    Outputs map[string]NodeOutput `json:"outputs"`
    Status  HistoryStatus         `json:"status"`
}
```

### After (New Code)
```go
type HistoryItem struct {
    Prompt  PromptArray           `json:"prompt"`  // ✅ Correct - parses array properly
    Outputs map[string]NodeOutput `json:"outputs"`
    Status  HistoryStatus         `json:"status"`
}
```

---

## PromptArray Structure

```go
type PromptArray struct {
    Number              float64                // Priority/order number
    PromptID            string                 // Unique prompt identifier
    Workflow            Workflow               // The actual workflow (what you probably want)
    ExtraData           map[string]interface{} // Additional metadata
    OutputsToExecute    []string               // List of node IDs to execute
}
```

---

## Migration Examples

### Example 1: Accessing the Workflow

**Before**:
```go
history, err := client.GetHistory(ctx, promptID)
if err != nil {
    log.Fatal(err)
}

for id, item := range history {
    workflow := item.Prompt  // ❌ This was directly a Workflow
    
    // Access nodes
    for nodeID, node := range workflow {
        fmt.Printf("Node %s: %s\n", nodeID, node.ClassType)
    }
}
```

**After**:
```go
history, err := client.GetHistory(ctx, promptID)
if err != nil {
    log.Fatal(err)
}

for id, item := range history {
    workflow := item.Prompt.Workflow  // ✅ Now access via .Workflow
    
    // Access nodes (same as before)
    for nodeID, node := range workflow {
        fmt.Printf("Node %s: %s\n", nodeID, node.ClassType)
    }
}
```

### Example 2: Accessing Additional Metadata

**New capability** - You can now access additional information:

```go
history, err := client.GetHistory(ctx, promptID)
if err != nil {
    log.Fatal(err)
}

for id, item := range history {
    // Access the workflow
    workflow := item.Prompt.Workflow
    
    // Access additional metadata (NEW!)
    fmt.Printf("Prompt ID: %s\n", item.Prompt.PromptID)
    fmt.Printf("Priority: %.0f\n", item.Prompt.Number)
    fmt.Printf("Client ID: %v\n", item.Prompt.ExtraData["client_id"])
    fmt.Printf("Nodes to execute: %v\n", item.Prompt.OutputsToExecute)
    
    // Access outputs and status (unchanged)
    fmt.Printf("Status: %s\n", item.Status.StatusStr)
    fmt.Printf("Outputs: %d\n", len(item.Outputs))
}
```

### Example 3: No Changes Needed

If you only use `Outputs` and `Status`, **no changes are required**:

```go
// This code works exactly the same as before ✅
history, err := client.GetHistory(ctx, promptID)
if err != nil {
    log.Fatal(err)
}

for id, item := range history {
    // These fields are unchanged
    fmt.Printf("Status: %s\n", item.Status.StatusStr)
    
    for nodeID, output := range item.Outputs {
        for _, img := range output.Images {
            client.SaveImage(ctx, img, fmt.Sprintf("%s.png", nodeID))
        }
    }
}
```

---

## Quick Reference

### Accessing Common Fields

| What you want | Old way | New way |
|--------------|---------|---------|
| The workflow | `item.Prompt` | `item.Prompt.Workflow` |
| Prompt ID | N/A | `item.Prompt.PromptID` |
| Priority number | N/A | `item.Prompt.Number` |
| Extra metadata | N/A | `item.Prompt.ExtraData` |
| Outputs | `item.Outputs` | `item.Outputs` (unchanged) |
| Status | `item.Status` | `item.Status` (unchanged) |

---

## Complete Example

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
    
    // Queue a workflow
    workflow := comfyui.NewWorkflowBuilder().
        AddNode("1", "KSampler", map[string]interface{}{
            "seed": 12345,
        }).
        Build()
    
    result, err := client.QueuePrompt(ctx, workflow, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Wait for completion...
    // (WebSocket monitoring code here)
    
    // Get history
    history, err := client.GetHistory(ctx, result.PromptID)
    if err != nil {
        log.Fatal(err)
    }
    
    // Access the result
    if item, ok := history[result.PromptID]; ok {
        // Access workflow (UPDATED)
        workflow := item.Prompt.Workflow
        fmt.Printf("Workflow has %d nodes\n", len(workflow))
        
        // Access metadata (NEW)
        fmt.Printf("Prompt ID: %s\n", item.Prompt.PromptID)
        fmt.Printf("Priority: %.0f\n", item.Prompt.Number)
        
        // Access outputs (UNCHANGED)
        fmt.Printf("Status: %s\n", item.Status.StatusStr)
        for nodeID, output := range item.Outputs {
            fmt.Printf("Node %s generated %d images\n", nodeID, len(output.Images))
        }
    }
}
```

---

## Testing Your Migration

After updating your code, verify it works:

```bash
# Build your project
go build

# Run tests
go test ./...

# Test with actual ComfyUI server
go run your_main.go
```

---

## Need Help?

- See [BUGFIX_SUMMARY.md](BUGFIX_SUMMARY.md) for technical details
- Check [examples/progress/main.go](examples/progress/main.go) for a working example
- Review [types.go](types.go) for the complete PromptArray definition

---

## Summary

✅ **Simple change**: Add `.Workflow` when accessing the workflow from history  
✅ **Bonus**: Access to additional metadata (prompt ID, priority, extra data)  
✅ **No impact**: If you only use `Outputs` and `Status`, no changes needed  
✅ **Better accuracy**: Correctly represents ComfyUI's internal structure

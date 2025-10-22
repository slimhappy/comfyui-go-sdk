# Queue API JSON Parsing Bugfix

## 问题描述

在运行 `execute_from_json` 示例时，出现以下错误：

```
2025/10/22 11:32:20 ❌ Execution monitoring failed: failed to get queue: 
failed to get queue: failed to decode response: json: cannot unmarshal array 
into Go struct field QueueStatus.queue_running of type comfyui.QueueItem
```

## 根本原因

ComfyUI 的队列 API (`/queue`) 返回的数据结构与 Go SDK 中定义的类型不匹配。

### ComfyUI API 实际返回格式

```json
{
  "queue_running": [
    [number, prompt_id, workflow, extra_data, outputs_to_execute],
    [number, prompt_id, workflow, extra_data, outputs_to_execute]
  ],
  "queue_pending": [
    [number, prompt_id, workflow, extra_data, outputs_to_execute],
    [number, prompt_id, workflow, extra_data, outputs_to_execute]
  ]
}
```

每个队列项是一个**数组**（元组），而不是对象。

### 原始 Go SDK 定义

```go
type QueueStatus struct {
    QueueRunning []QueueItem `json:"queue_running"`
    QueuePending []QueueItem `json:"queue_pending"`
}

type QueueItem struct {
    Number    int                    `json:"0"`
    PromptID  string                 `json:"1"`
    Prompt    Workflow               `json:"2"`
    ExtraData map[string]interface{} `json:"3"`
    Outputs   []string               `json:"4"`
}
```

**问题**：`QueueItem` 使用数字字符串作为 JSON 标签（`"0"`, `"1"`, `"2"` 等），这只适用于直接解析单个数组，但当数组嵌套在对象的字段中时，Go 的 JSON 解析器无法正确处理。

## 修复方案

### 1. 为 `QueueStatus` 添加自定义 JSON 解析

```go
// UnmarshalJSON implements custom JSON unmarshaling for QueueStatus
func (q *QueueStatus) UnmarshalJSON(data []byte) error {
    var raw struct {
        QueueRunning [][]interface{} `json:"queue_running"`
        QueuePending [][]interface{} `json:"queue_pending"`
    }
    
    if err := json.Unmarshal(data, &raw); err != nil {
        return err
    }
    
    // Parse running queue
    q.QueueRunning = make([]QueueItem, 0, len(raw.QueueRunning))
    for _, arr := range raw.QueueRunning {
        item, err := parseQueueItemArray(arr)
        if err != nil {
            return fmt.Errorf("failed to parse running queue item: %w", err)
        }
        q.QueueRunning = append(q.QueueRunning, item)
    }
    
    // Parse pending queue
    q.QueuePending = make([]QueueItem, 0, len(raw.QueuePending))
    for _, arr := range raw.QueuePending {
        item, err := parseQueueItemArray(arr)
        if err != nil {
            return fmt.Errorf("failed to parse pending queue item: %w", err)
        }
        q.QueuePending = append(q.QueuePending, item)
    }
    
    return nil
}
```

### 2. 更新 `QueueItem` 结构体

移除 JSON 标签，因为我们使用自定义解析：

```go
type QueueItem struct {
    Number    int                    `json:"-"`
    PromptID  string                 `json:"-"`
    Prompt    Workflow               `json:"-"`
    ExtraData map[string]interface{} `json:"-"`
    Outputs   []string               `json:"-"`
}
```

### 3. 创建辅助解析函数

```go
// parseQueueItemArray parses a queue item from an array
func parseQueueItemArray(arr []interface{}) (QueueItem, error) {
    if len(arr) < 3 {
        return QueueItem{}, fmt.Errorf("queue item array too short: expected at least 3 elements, got %d", len(arr))
    }
    
    item := QueueItem{}
    
    // Parse number (index 0)
    if num, ok := arr[0].(float64); ok {
        item.Number = int(num)
    }
    
    // Parse prompt_id (index 1)
    if id, ok := arr[1].(string); ok {
        item.PromptID = id
    }
    
    // Parse workflow (index 2)
    if promptMap, ok := arr[2].(map[string]interface{}); ok {
        workflow := make(Workflow)
        for k, v := range promptMap {
            if nodeMap, ok := v.(map[string]interface{}); ok {
                node := Node{}
                if classType, ok := nodeMap["class_type"].(string); ok {
                    node.ClassType = classType
                }
                if inputs, ok := nodeMap["inputs"].(map[string]interface{}); ok {
                    node.Inputs = inputs
                }
                workflow[k] = node
            }
        }
        item.Prompt = workflow
    }
    
    // Parse extra_data (index 3)
    if len(arr) > 3 {
        if extraData, ok := arr[3].(map[string]interface{}); ok {
            item.ExtraData = extraData
        }
    }
    
    // Parse outputs_to_execute (index 4)
    if len(arr) > 4 {
        if outputs, ok := arr[4].([]interface{}); ok {
            item.Outputs = make([]string, 0, len(outputs))
            for _, output := range outputs {
                if str, ok := output.(string); ok {
                    item.Outputs = append(item.Outputs, str)
                }
            }
        }
    }
    
    return item, nil
}
```

## 修复效果

### 修复前

```
❌ Execution monitoring failed: failed to get queue: failed to get queue: 
failed to decode response: json: cannot unmarshal array into Go struct 
field QueueStatus.queue_running of type comfyui.QueueItem
```

### 修复后

```
✅ Completed in 0.2 seconds
📥 Retrieving execution results...
📊 Execution ID: 2d2b79be-aa32-44e7-97d4-04cc6c601e83
   Status: success
```

## 测试结果

### 单元测试

```bash
$ go test -v .
=== RUN   TestNewClient
--- PASS: TestNewClient (0.00s)
=== RUN   TestClientIDManagement
--- PASS: TestClientIDManagement (0.00s)
=== RUN   TestWorkflowOperations
--- PASS: TestWorkflowOperations (0.00s)
=== RUN   TestWorkflowBuilder
--- PASS: TestWorkflowBuilder (0.00s)
=== RUN   TestWorkflowValidation
--- PASS: TestWorkflowValidation (0.00s)
=== RUN   TestWebSocketMessageParsing
--- PASS: TestWebSocketMessageParsing (0.00s)
=== RUN   TestErrorTypes
--- PASS: TestErrorTypes (0.00s)
=== RUN   TestAPICallsIntegration
--- PASS: TestAPICallsIntegration (0.08s)
PASS
ok      github.com/yourusername/comfyui-go-sdk  0.081s
```

### 功能测试

```bash
$ ./bin/execute_from_json examples/execute_from_json/workflow_example.json

╔════════════════════════════════════════════════════════════════════╗
║         ComfyUI Go SDK - Execute Workflow from JSON File         ║
╚════════════════════════════════════════════════════════════════════╝

🔍 Checking ComfyUI server status...
✅ ComfyUI server is running

📂 Loading workflow from: examples/execute_from_json/workflow_example.json
✅ Workflow loaded successfully

🚀 Submitting workflow to ComfyUI...
✅ Workflow queued successfully!

⏳ Monitoring execution progress...
✅ Completed in 0.2 seconds

📥 Retrieving execution results...
📊 Execution ID: 2d2b79be-aa32-44e7-97d4-04cc6c601e83
   Status: success
```

## 相关问题

这是继 [BUGFIX_SUMMARY.md](BUGFIX_SUMMARY.md) 中修复的历史 API 问题之后的第二个类似问题。两个问题都源于 ComfyUI API 返回数组格式的数据，而 Go 的标准 JSON 解析器无法直接处理嵌套在对象字段中的数组到结构体的映射。

### 已修复的 API

1. ✅ **History API** (`/history/{prompt_id}`) - 修复于 BUGFIX_SUMMARY.md
   - 问题：`prompt` 字段返回数组而非对象
   - 解决：为 `PromptArray` 添加自定义 `UnmarshalJSON`

2. ✅ **Queue API** (`/queue`) - 本次修复
   - 问题：`queue_running` 和 `queue_pending` 返回数组的数组
   - 解决：为 `QueueStatus` 添加自定义 `UnmarshalJSON`

## 向后兼容性

✅ **完全向后兼容**

- 所有现有的 API 方法签名保持不变
- `QueueItem` 结构体的字段保持不变
- 只是改变了内部的 JSON 解析实现
- 所有现有代码无需修改即可继续工作

## 文件变更

- **修改文件**：`types.go`
  - 添加 `QueueStatus.UnmarshalJSON()` 方法
  - 更新 `QueueItem` 结构体的 JSON 标签
  - 添加 `parseQueueItemArray()` 辅助函数

## 提交信息

```
Fix queue API JSON parsing error

- Add custom UnmarshalJSON for QueueStatus to handle array-based queue items
- Update QueueItem struct to remove JSON tags (use custom parsing)
- Add parseQueueItemArray helper function
- Fix "cannot unmarshal array into Go struct field" error
- All tests passing
- Backward compatible
```

## 总结

通过为 `QueueStatus` 添加自定义 JSON 解析方法，成功修复了队列 API 的数据解析问题。现在 SDK 可以正确处理 ComfyUI 返回的数组格式的队列数据，`execute_from_json` 示例可以正常运行并监控执行进度。

---

**修复日期**：2025-10-22  
**影响范围**：Queue API (`/queue`)  
**向后兼容**：✅ 是  
**测试状态**：✅ 全部通过

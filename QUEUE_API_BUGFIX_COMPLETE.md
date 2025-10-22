# Queue API Bugfix - Complete Summary

## 问题报告

用户在运行 `execute_from_json` 示例时遇到以下错误：

```bash
./execute_from_json ../examples/execute_from_json/workflow_example.json 
报错：
⏳ Monitoring execution progress...
2025/10/22 11:32:20 ❌ Execution monitoring failed: failed to get queue: 
failed to get queue: failed to decode response: json: cannot unmarshal array 
into Go struct field QueueStatus.queue_running of type comfyui.QueueItem
```

## 问题分析

### 根本原因

ComfyUI 的队列 API (`GET /queue`) 返回的数据格式与 Go SDK 中定义的类型不匹配。

**ComfyUI API 实际返回格式**：
```json
{
  "queue_running": [
    [number, prompt_id, workflow, extra_data, outputs_to_execute],
    [number, prompt_id, workflow, extra_data, outputs_to_execute]
  ],
  "queue_pending": [
    [number, prompt_id, workflow, extra_data, outputs_to_execute]
  ]
}
```

每个队列项是一个**数组**（元组），索引含义：
- `[0]` = number (int) - 队列编号
- `[1]` = prompt_id (string) - 提示词 ID
- `[2]` = workflow (object) - 工作流定义
- `[3]` = extra_data (object) - 额外数据
- `[4]` = outputs_to_execute (array) - 要执行的输出

### 原始代码问题

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

**问题**：Go 的标准 JSON 解析器无法将嵌套在对象字段中的数组直接映射到使用数字字符串作为 JSON 标签的结构体。

## 解决方案

### 1. 为 QueueStatus 添加自定义 JSON 解析

```go
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

### 2. 更新 QueueItem 结构体

```go
type QueueItem struct {
    Number    int                    `json:"-"`
    PromptID  string                 `json:"-"`
    Prompt    Workflow               `json:"-"`
    ExtraData map[string]interface{} `json:"-"`
    Outputs   []string               `json:"-"`
}
```

移除 JSON 标签，改用自定义解析。

### 3. 创建辅助解析函数

```go
func parseQueueItemArray(arr []interface{}) (QueueItem, error) {
    if len(arr) < 3 {
        return QueueItem{}, fmt.Errorf("queue item array too short: expected at least 3 elements, got %d", len(arr))
    }
    
    item := QueueItem{}
    
    // Parse each field from array indices
    if num, ok := arr[0].(float64); ok {
        item.Number = int(num)
    }
    
    if id, ok := arr[1].(string); ok {
        item.PromptID = id
    }
    
    // Parse workflow (complex nested structure)
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
    
    // Parse optional fields
    if len(arr) > 3 {
        if extraData, ok := arr[3].(map[string]interface{}); ok {
            item.ExtraData = extraData
        }
    }
    
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

### 修复前 ❌

```
⏳ Monitoring execution progress...
2025/10/22 11:32:20 ❌ Execution monitoring failed: failed to get queue: 
failed to get queue: failed to decode response: json: cannot unmarshal array 
into Go struct field QueueStatus.queue_running of type comfyui.QueueItem
```

### 修复后 ✅

```
╔════════════════════════════════════════════════════════════════════╗
║         ComfyUI Go SDK - Execute Workflow from JSON File         ║
╚════════════════════════════════════════════════════════════════════╝

🔍 Checking ComfyUI server status...
✅ ComfyUI server is running

📂 Loading workflow from: examples/execute_from_json/workflow_example.json
✅ Workflow loaded successfully
   Total nodes: 7

🚀 Submitting workflow to ComfyUI...
✅ Workflow queued successfully!
   Prompt ID: 2d2b79be-aa32-44e7-97d4-04cc6c601e83
   Queue Position: 18

⏳ Monitoring execution progress...
✅ Completed in 0.2 seconds

📥 Retrieving execution results...
📊 Execution ID: 2d2b79be-aa32-44e7-97d4-04cc6c601e83
   Status: success
```

## 测试验证

### 单元测试 ✅

```bash
$ cd /data/comfyui-go-sdk && go test -v .

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
=== RUN   TestAPICallsIntegration/GetQueue
--- PASS: TestAPICallsIntegration (0.08s)
    --- PASS: TestAPICallsIntegration/GetQueue (0.00s)
PASS
ok      github.com/yourusername/comfyui-go-sdk  0.081s
```

**结果**：8/8 测试全部通过 ✅

### 构建测试 ✅

```bash
$ cd /data/comfyui-go-sdk && make build

Building examples...
Done! Binaries in ./bin/
```

**结果**：所有示例编译成功 ✅

### 功能测试 ✅

```bash
$ ./bin/execute_from_json examples/execute_from_json/workflow_example.json

✅ ComfyUI server is running
✅ Workflow loaded successfully
✅ Workflow queued successfully!
⏳ Monitoring execution progress...
✅ Completed in 0.2 seconds
📊 Status: success
```

**结果**：功能完全正常 ✅

## 影响范围

### 受影响的组件

1. **Queue API** (`GET /queue`)
   - `GetQueue()` 方法
   - 队列状态查询

2. **Progress Monitoring**
   - 实时进度跟踪
   - 队列位置监控

3. **Examples**
   - `execute_from_json` 示例
   - `progress` 示例

### 向后兼容性

✅ **完全向后兼容**

- 所有公共 API 方法签名保持不变
- `QueueItem` 结构体字段保持不变
- 只修改了内部 JSON 解析实现
- 现有代码无需任何修改

## 相关问题

这是第二个类似的 JSON 解析问题：

1. **History API** - 已修复 ([BUGFIX_SUMMARY.md](BUGFIX_SUMMARY.md))
   - 问题：`prompt` 字段返回数组而非对象
   - 解决：为 `PromptArray` 添加自定义 `UnmarshalJSON`

2. **Queue API** - 本次修复 ([QUEUE_API_BUGFIX.md](QUEUE_API_BUGFIX.md))
   - 问题：`queue_running` 和 `queue_pending` 返回数组的数组
   - 解决：为 `QueueStatus` 添加自定义 `UnmarshalJSON`

两个问题的共同点：
- ComfyUI API 使用数组格式存储结构化数据
- Go 标准 JSON 解析器无法直接处理
- 需要自定义 `UnmarshalJSON` 方法

## 文件变更

### 修改的文件

1. **types.go**
   - 添加 `QueueStatus.UnmarshalJSON()` 方法
   - 更新 `QueueItem` 结构体（移除 JSON 标签）
   - 添加 `parseQueueItemArray()` 辅助函数

### 新增的文档

1. **QUEUE_API_BUGFIX.md**
   - 详细的问题分析
   - 完整的解决方案
   - 测试结果和验证

2. **README.md** (更新)
   - 添加 Changelog 部分
   - 记录所有 bug 修复

3. **QUEUE_API_BUGFIX_COMPLETE.md** (本文档)
   - 完整的修复总结
   - 从问题到解决的全过程

## Git 提交信息

```
Fix queue API JSON parsing error

Problem:
- execute_from_json example failed with JSON unmarshal error
- Queue API returns array-based queue items
- Go JSON parser cannot handle nested arrays to structs

Solution:
- Add custom UnmarshalJSON for QueueStatus
- Update QueueItem to use custom parsing
- Add parseQueueItemArray helper function

Testing:
- All unit tests passing (8/8)
- All examples building successfully
- execute_from_json working correctly
- Queue monitoring functional

Impact:
- Fixes queue monitoring
- Fixes progress tracking
- Backward compatible
- No breaking changes

Related: BUGFIX_SUMMARY.md (History API fix)
```

## 使用示例

### 基本使用

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
    
    // Get queue status (now works correctly!)
    queue, err := client.GetQueue(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Running: %d, Pending: %d", 
        len(queue.QueueRunning), 
        len(queue.QueuePending))
    
    // Access queue items
    for _, item := range queue.QueueRunning {
        log.Printf("Prompt ID: %s, Number: %d", 
            item.PromptID, item.Number)
    }
}
```

### 使用 CLI 工具

```bash
# 执行工作流并监控进度
./bin/execute_from_json workflow.json

# 带参数执行
./bin/execute_from_json workflow.json seed=12345 steps=30

# 自定义提示词
./bin/execute_from_json workflow.json \
  prompt="beautiful landscape" \
  negative="blurry"
```

## 总结

✅ **问题已完全解决**

通过为 `QueueStatus` 添加自定义 JSON 解析方法，成功修复了队列 API 的数据解析问题。现在：

- ✅ 队列 API 正常工作
- ✅ 进度监控功能正常
- ✅ `execute_from_json` 示例完全可用
- ✅ 所有测试通过
- ✅ 向后兼容
- ✅ 无破坏性变更

用户现在可以正常使用所有队列相关功能，包括实时进度监控和工作流执行跟踪。

---

**修复日期**：2025-10-22  
**修复人员**：AI Assistant  
**测试状态**：✅ 全部通过  
**生产状态**：✅ 可用  
**向后兼容**：✅ 是

package comfyui

import (
	"encoding/json"
	"fmt"
	"time"
)


// Workflow represents a ComfyUI workflow (prompt)
type Workflow map[string]Node

// Node represents a single node in the workflow
type Node struct {
	ClassType string                 `json:"class_type"`
	Inputs    map[string]interface{} `json:"inputs"`
}

// QueuePromptRequest represents the request to queue a prompt
type QueuePromptRequest struct {
	Prompt    Workflow               `json:"prompt"`
	ClientID  string                 `json:"client_id,omitempty"`
	PromptID  string                 `json:"prompt_id,omitempty"`
	ExtraData map[string]interface{} `json:"extra_data,omitempty"`
	FrontEnd  bool                   `json:"front,omitempty"`
	Number    int                    `json:"number,omitempty"`
}

// QueuePromptResponse represents the response from queueing a prompt
type QueuePromptResponse struct {
	PromptID   string                 `json:"prompt_id"`
	Number     int                    `json:"number"`
	NodeErrors map[string]interface{} `json:"node_errors"`
}

// QueueStatus represents the current queue status
type QueueStatus struct {
	QueueRunning []QueueItem `json:"queue_running"`
	QueuePending []QueueItem `json:"queue_pending"`
}

// QueueItem represents an item in the queue
type QueueItem struct {
	Number    int                    `json:"0"`
	PromptID  string                 `json:"1"`
	Prompt    Workflow               `json:"2"`
	ExtraData map[string]interface{} `json:"3"`
	Outputs   []string               `json:"4"`
}

// History represents execution history
type History map[string]HistoryItem

// PromptArray represents the prompt array structure returned by ComfyUI history API
// Structure: [number, prompt_id, workflow, extra_data, outputs_to_execute]
type PromptArray struct {
	Number              float64                `json:"-"`
	PromptID            string                 `json:"-"`
	Workflow            Workflow               `json:"-"`
	ExtraData           map[string]interface{} `json:"-"`
	OutputsToExecute    []string               `json:"-"`
	rawArray            []interface{}
}

// UnmarshalJSON implements custom JSON unmarshaling for PromptArray
func (p *PromptArray) UnmarshalJSON(data []byte) error {
	var arr []interface{}
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}

	p.rawArray = arr
	if len(arr) < 3 {
		return fmt.Errorf("prompt array too short: expected at least 3 elements, got %d", len(arr))
	}

	// Parse number (index 0)
	if num, ok := arr[0].(float64); ok {
		p.Number = num
	}

	// Parse prompt_id (index 1)
	if id, ok := arr[1].(string); ok {
		p.PromptID = id
	}

	// Parse workflow (index 2)
	if workflowData, ok := arr[2].(map[string]interface{}); ok {
		workflow := make(Workflow)
		for k, v := range workflowData {
			if nodeData, ok := v.(map[string]interface{}); ok {
				node := Node{}
				if classType, ok := nodeData["class_type"].(string); ok {
					node.ClassType = classType
				}
				if inputs, ok := nodeData["inputs"].(map[string]interface{}); ok {
					node.Inputs = inputs
				}
				workflow[k] = node
			}
		}
		p.Workflow = workflow
	}

	// Parse extra_data (index 3)
	if len(arr) > 3 {
		if extraData, ok := arr[3].(map[string]interface{}); ok {
			p.ExtraData = extraData
		}
	}

	// Parse outputs_to_execute (index 4)
	if len(arr) > 4 {
		if outputs, ok := arr[4].([]interface{}); ok {
			p.OutputsToExecute = make([]string, 0, len(outputs))
			for _, output := range outputs {
				if str, ok := output.(string); ok {
					p.OutputsToExecute = append(p.OutputsToExecute, str)
				}
			}
		}
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for PromptArray
func (p *PromptArray) MarshalJSON() ([]byte, error) {
	if p.rawArray != nil {
		return json.Marshal(p.rawArray)
	}
	arr := []interface{}{p.Number, p.PromptID, p.Workflow, p.ExtraData, p.OutputsToExecute}
	return json.Marshal(arr)
}

// HistoryItem represents a single history entry
type HistoryItem struct {
	Prompt  PromptArray           `json:"prompt"`
	Outputs map[string]NodeOutput `json:"outputs"`
	Status  HistoryStatus         `json:"status"`
}


// HistoryStatus represents the status of a history item
type HistoryStatus struct {
	StatusStr string        `json:"status_str"`
	Completed bool          `json:"completed"`
	Messages  []interface{} `json:"messages"`
}

// NodeOutput represents the output of a node
type NodeOutput struct {
	Images []ImageInfo            `json:"images,omitempty"`
	Text   []string               `json:"text,omitempty"`
	Data   map[string]interface{} `json:"-"`
}

// ImageInfo represents information about an output image
type ImageInfo struct {
	Filename  string `json:"filename"`
	Subfolder string `json:"subfolder"`
	Type      string `json:"type"`
}

// SystemStats represents system statistics
type SystemStats struct {
	System  SystemInfo   `json:"system"`
	Devices []DeviceInfo `json:"devices"`
}

// SystemInfo represents system information
type SystemInfo struct {
	OS             string `json:"os"`
	PythonVersion  string `json:"python_version"`
	EmbeddedPython bool   `json:"embedded_python"`
}

// DeviceInfo represents device information
type DeviceInfo struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Index          int    `json:"index"`
	VRAMTotal      int64  `json:"vram_total"`
	VRAMFree       int64  `json:"vram_free"`
	TorchVRAMTotal int64  `json:"torch_vram_total"`
	TorchVRAMFree  int64  `json:"torch_vram_free"`
}

// ObjectInfo represents node class information
type ObjectInfo map[string]NodeClassInfo

// NodeClassInfo represents information about a node class
type NodeClassInfo struct {
	Input       NodeInputInfo `json:"input"`
	Output      []string      `json:"output"`
	OutputName  []string      `json:"output_name"`
	Name        string        `json:"name"`
	DisplayName string        `json:"display_name"`
	Description string        `json:"description"`
	Category    string        `json:"category"`
	OutputNode  bool          `json:"output_node"`
}

// NodeInputInfo represents input information for a node
type NodeInputInfo struct {
	Required map[string]interface{} `json:"required"`
	Optional map[string]interface{} `json:"optional,omitempty"`
	Hidden   map[string]interface{} `json:"hidden,omitempty"`
}

// UploadImageResponse represents the response from uploading an image
type UploadImageResponse struct {
	Name      string `json:"name"`
	Subfolder string `json:"subfolder"`
	Type      string `json:"type"`
}

// UploadOptions represents options for uploading files
type UploadOptions struct {
	Subfolder string
	Type      string // "input", "temp", "output"
	Overwrite bool
}

// WebSocketMessage represents a message received via WebSocket
type WebSocketMessage struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// MessageType represents the type of WebSocket message
type MessageType string

const (
	MessageTypeStatus    MessageType = "status"
	MessageTypeExecuting MessageType = "executing"
	MessageTypeProgress  MessageType = "progress"
	MessageTypeExecuted  MessageType = "executed"
	MessageTypeError     MessageType = "execution_error"
	MessageTypeCached    MessageType = "execution_cached"
)

// ExecutingData represents data for executing message
type ExecutingData struct {
	Node     *string `json:"node"`
	PromptID string  `json:"prompt_id"`
}

// ProgressData represents data for progress message
type ProgressData struct {
	Value int `json:"value"`
	Max   int `json:"max"`
}

// ExecutedData represents data for executed message
type ExecutedData struct {
	Node     string                 `json:"node"`
	Output   map[string]interface{} `json:"output"`
	PromptID string                 `json:"prompt_id"`
}

// StatusData represents data for status message
type StatusData struct {
	Status struct {
		ExecInfo struct {
			QueueRemaining int `json:"queue_remaining"`
		} `json:"exec_info"`
	} `json:"status"`
	SID string `json:"sid"`
}

// ErrorData represents data for error message
type ErrorData struct {
	PromptID         string   `json:"prompt_id"`
	NodeID           string   `json:"node_id"`
	NodeType         string   `json:"node_type"`
	ExceptionType    string   `json:"exception_type"`
	ExceptionMessage string   `json:"exception_message"`
	Traceback        []string `json:"traceback"`
}

// QueueManagementRequest represents a request to manage the queue
type QueueManagementRequest struct {
	Clear  bool     `json:"clear,omitempty"`
	Delete []string `json:"delete,omitempty"`
}

// HistoryManagementRequest represents a request to manage history
type HistoryManagementRequest struct {
	Clear  bool     `json:"clear,omitempty"`
	Delete []string `json:"delete,omitempty"`
}

// InterruptRequest represents a request to interrupt execution
type InterruptRequest struct {
	PromptID string `json:"prompt_id,omitempty"`
}

// FreeMemoryRequest represents a request to free memory
type FreeMemoryRequest struct {
	UnloadModels bool `json:"unload_models,omitempty"`
	FreeMemory   bool `json:"free_memory,omitempty"`
}

// Features represents server features
type Features struct {
	Features []string `json:"features"`
}

// ExecutionResult represents the result of workflow execution
type ExecutionResult struct {
	PromptID  string
	Images    []ImageInfo
	Outputs   map[string]NodeOutput
	Status    HistoryStatus
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

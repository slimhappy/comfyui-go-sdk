package comfyui

import "fmt"

// Error types
var (
	ErrNodeNotFound     = fmt.Errorf("node not found")
	ErrInvalidWorkflow  = fmt.Errorf("invalid workflow")
	ErrExecutionFailed  = fmt.Errorf("execution failed")
	ErrConnectionFailed = fmt.Errorf("connection failed")
	ErrTimeout          = fmt.Errorf("timeout")
	ErrInvalidResponse  = fmt.Errorf("invalid response")
)

// APIError represents an API error
type APIError struct {
	StatusCode int
	Message    string
	Details    interface{}
}

func (e *APIError) Error() string {
	if e.Details != nil {
		return fmt.Sprintf("API error (status %d): %s - %v", e.StatusCode, e.Message, e.Details)
	}
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// NodeError represents a node execution error
type NodeError struct {
	NodeID    string
	NodeType  string
	Message   string
	Traceback []string
}

func (e *NodeError) Error() string {
	return fmt.Sprintf("node error in %s (%s): %s", e.NodeID, e.NodeType, e.Message)
}

// ValidationError represents a workflow validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in %s: %s", e.Field, e.Message)
}

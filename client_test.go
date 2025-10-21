package comfyui

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("http://127.0.0.1:8188")
	if client == nil {
		t.Fatal("Expected client to be created")
	}

	if client.baseURL != "http://127.0.0.1:8188" {
		t.Errorf("Expected baseURL to be http://127.0.0.1:8188, got %s", client.baseURL)
	}

	if client.clientID == "" {
		t.Error("Expected clientID to be set")
	}
}

func TestClientIDManagement(t *testing.T) {
	client := NewClient("http://127.0.0.1:8188")

	customID := "test-client-id"
	client.SetClientID(customID)

	if client.GetClientID() != customID {
		t.Errorf("Expected clientID to be %s, got %s", customID, client.GetClientID())
	}
}

func TestWorkflowOperations(t *testing.T) {
	workflow := Workflow{
		"1": Node{
			ClassType: "TestNode",
			Inputs: map[string]interface{}{
				"param1": "value1",
			},
		},
	}

	// Test SetNodeInput
	err := workflow.SetNodeInput("1", "param2", "value2")
	if err != nil {
		t.Errorf("Failed to set node input: %v", err)
	}

	// Test GetNodeInput
	value, err := workflow.GetNodeInput("1", "param2")
	if err != nil {
		t.Errorf("Failed to get node input: %v", err)
	}
	if value != "value2" {
		t.Errorf("Expected value2, got %v", value)
	}

	// Test AddNode
	workflow.AddNode("2", "AnotherNode", map[string]interface{}{
		"test": true,
	})

	if len(workflow) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(workflow))
	}

	// Test GetNode
	node, ok := workflow.GetNode("2")
	if !ok {
		t.Error("Expected to find node 2")
	}
	if node.ClassType != "AnotherNode" {
		t.Errorf("Expected AnotherNode, got %s", node.ClassType)
	}

	// Test RemoveNode
	workflow.RemoveNode("2")
	if len(workflow) != 1 {
		t.Errorf("Expected 1 node after removal, got %d", len(workflow))
	}

	// Test NodeIDs
	ids := workflow.NodeIDs()
	if len(ids) != 1 {
		t.Errorf("Expected 1 node ID, got %d", len(ids))
	}

	// Test NodesByClass
	nodes := workflow.NodesByClass("TestNode")
	if len(nodes) != 1 {
		t.Errorf("Expected 1 TestNode, got %d", len(nodes))
	}

	// Test Validate
	if err := workflow.Validate(); err != nil {
		t.Errorf("Validation failed: %v", err)
	}

	// Test Clone
	clone, err := workflow.Clone()
	if err != nil {
		t.Errorf("Failed to clone workflow: %v", err)
	}
	if len(clone) != len(workflow) {
		t.Error("Clone has different number of nodes")
	}
}

func TestWorkflowBuilder(t *testing.T) {
	builder := NewWorkflowBuilder()

	// Add nodes
	id1 := builder.AddNode("CheckpointLoader", map[string]interface{}{
		"ckpt_name": "test.safetensors",
	})

	id2 := builder.AddNode("KSampler", map[string]interface{}{
		"steps": 20,
	})

	// Connect nodes
	err := builder.ConnectNodes(id1, 0, id2, "model")
	if err != nil {
		t.Errorf("Failed to connect nodes: %v", err)
	}

	// Build workflow
	workflow := builder.Build()

	if len(workflow) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(workflow))
	}

	// Check connection
	node2 := workflow[id2]
	if model, ok := node2.Inputs["model"].([]interface{}); ok {
		if model[0] != id1 {
			t.Error("Connection not set correctly")
		}
	} else {
		t.Error("Model input not set")
	}
}

func TestWorkflowValidation(t *testing.T) {
	// Empty workflow
	emptyWorkflow := Workflow{}
	if err := emptyWorkflow.Validate(); err == nil {
		t.Error("Expected validation to fail for empty workflow")
	}

	// Workflow with node missing class_type
	invalidWorkflow := Workflow{
		"1": Node{
			Inputs: map[string]interface{}{},
		},
	}
	if err := invalidWorkflow.Validate(); err == nil {
		t.Error("Expected validation to fail for node without class_type")
	}

	// Valid workflow
	validWorkflow := Workflow{
		"1": Node{
			ClassType: "TestNode",
			Inputs:    map[string]interface{}{},
		},
	}
	if err := validWorkflow.Validate(); err != nil {
		t.Errorf("Expected validation to pass: %v", err)
	}
}

func TestWebSocketMessageParsing(t *testing.T) {
	// Test ExecutingData
	execMsg := WebSocketMessage{
		Type: string(MessageTypeExecuting),
		Data: map[string]interface{}{
			"prompt_id": "test-prompt",
			"node":      "3",
		},
	}

	execData, err := execMsg.GetExecutingData()
	if err != nil {
		t.Errorf("Failed to parse executing data: %v", err)
	}
	if execData.PromptID != "test-prompt" {
		t.Errorf("Expected prompt_id test-prompt, got %s", execData.PromptID)
	}
	if execData.Node == nil || *execData.Node != "3" {
		t.Error("Node not parsed correctly")
	}

	// Test ProgressData
	progMsg := WebSocketMessage{
		Type: string(MessageTypeProgress),
		Data: map[string]interface{}{
			"value": float64(5),
			"max":   float64(20),
		},
	}

	progData, err := progMsg.GetProgressData()
	if err != nil {
		t.Errorf("Failed to parse progress data: %v", err)
	}
	if progData.Value != 5 || progData.Max != 20 {
		t.Errorf("Expected 5/20, got %d/%d", progData.Value, progData.Max)
	}

	// Test ExecutedData
	execedMsg := WebSocketMessage{
		Type: string(MessageTypeExecuted),
		Data: map[string]interface{}{
			"prompt_id": "test-prompt",
			"node":      "9",
			"output": map[string]interface{}{
				"images": []interface{}{},
			},
		},
	}

	execedData, err := execedMsg.GetExecutedData()
	if err != nil {
		t.Errorf("Failed to parse executed data: %v", err)
	}
	if execedData.Node != "9" {
		t.Errorf("Expected node 9, got %s", execedData.Node)
	}
}

func TestErrorTypes(t *testing.T) {
	// Test APIError
	apiErr := &APIError{
		StatusCode: 500,
		Message:    "Internal Server Error",
		Details:    map[string]string{"error": "test"},
	}
	if apiErr.Error() == "" {
		t.Error("APIError should have error message")
	}

	// Test NodeError
	nodeErr := &NodeError{
		NodeID:   "3",
		NodeType: "KSampler",
		Message:  "Test error",
	}
	if nodeErr.Error() == "" {
		t.Error("NodeError should have error message")
	}

	// Test ValidationError
	valErr := &ValidationError{
		Field:   "test_field",
		Message: "Invalid value",
	}
	if valErr.Error() == "" {
		t.Error("ValidationError should have error message")
	}
}

// Mock test for API calls (requires running ComfyUI server)
func TestAPICallsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewClient("http://127.0.0.1:8188")
	ctx := context.Background()

	// Test GetSystemStats
	t.Run("GetSystemStats", func(t *testing.T) {
		stats, err := client.GetSystemStats(ctx)
		if err != nil {
			t.Skipf("ComfyUI not available: %v", err)
		}
		if stats == nil {
			t.Error("Expected stats to be returned")
		}
	})

	// Test GetObjectInfo
	t.Run("GetObjectInfo", func(t *testing.T) {
		info, err := client.GetObjectInfo(ctx, "")
		if err != nil {
			t.Skipf("ComfyUI not available: %v", err)
		}
		if len(info) == 0 {
			t.Error("Expected object info to be returned")
		}
	})

	// Test GetQueue
	t.Run("GetQueue", func(t *testing.T) {
		queue, err := client.GetQueue(ctx)
		if err != nil {
			t.Skipf("ComfyUI not available: %v", err)
		}
		if queue == nil {
			t.Error("Expected queue to be returned")
		}
	})
}

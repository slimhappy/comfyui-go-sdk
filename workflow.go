package comfyui

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadWorkflowFromFile loads a workflow from a JSON file
func LoadWorkflowFromFile(filepath string) (Workflow, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var workflow Workflow
	if err := json.Unmarshal(data, &workflow); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow: %w", err)
	}

	return workflow, nil
}

// SaveWorkflowToFile saves a workflow to a JSON file
func SaveWorkflowToFile(workflow Workflow, filepath string) error {
	data, err := json.MarshalIndent(workflow, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal workflow: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// SetNodeInput sets an input value for a node
func (w Workflow) SetNodeInput(nodeID string, inputName string, value interface{}) error {
	node, ok := w[nodeID]
	if !ok {
		return fmt.Errorf("node %s not found", nodeID)
	}

	if node.Inputs == nil {
		node.Inputs = make(map[string]interface{})
	}

	node.Inputs[inputName] = value
	w[nodeID] = node

	return nil
}

// GetNodeInput gets an input value from a node
func (w Workflow) GetNodeInput(nodeID string, inputName string) (interface{}, error) {
	node, ok := w[nodeID]
	if !ok {
		return nil, fmt.Errorf("node %s not found", nodeID)
	}

	value, ok := node.Inputs[inputName]
	if !ok {
		return nil, fmt.Errorf("input %s not found in node %s", inputName, nodeID)
	}

	return value, nil
}

// AddNode adds a new node to the workflow
func (w Workflow) AddNode(nodeID string, classType string, inputs map[string]interface{}) {
	w[nodeID] = Node{
		ClassType: classType,
		Inputs:    inputs,
	}
}

// RemoveNode removes a node from the workflow
func (w Workflow) RemoveNode(nodeID string) {
	delete(w, nodeID)
}

// GetNode gets a node from the workflow
func (w Workflow) GetNode(nodeID string) (Node, bool) {
	node, ok := w[nodeID]
	return node, ok
}

// Clone creates a deep copy of the workflow
func (w Workflow) Clone() (Workflow, error) {
	data, err := json.Marshal(w)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow: %w", err)
	}

	var clone Workflow
	if err := json.Unmarshal(data, &clone); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow: %w", err)
	}

	return clone, nil
}

// NodeIDs returns all node IDs in the workflow
func (w Workflow) NodeIDs() []string {
	ids := make([]string, 0, len(w))
	for id := range w {
		ids = append(ids, id)
	}
	return ids
}

// NodesByClass returns all nodes of a specific class type
func (w Workflow) NodesByClass(classType string) map[string]Node {
	nodes := make(map[string]Node)
	for id, node := range w {
		if node.ClassType == classType {
			nodes[id] = node
		}
	}
	return nodes
}

// Validate performs basic validation on the workflow
func (w Workflow) Validate() error {
	if len(w) == 0 {
		return fmt.Errorf("workflow is empty")
	}

	for id, node := range w {
		if node.ClassType == "" {
			return fmt.Errorf("node %s has no class_type", id)
		}
	}

	return nil
}

// WorkflowBuilder helps build workflows programmatically
type WorkflowBuilder struct {
	workflow Workflow
	nextID   int
}

// NewWorkflowBuilder creates a new workflow builder
func NewWorkflowBuilder() *WorkflowBuilder {
	return &WorkflowBuilder{
		workflow: make(Workflow),
		nextID:   1,
	}
}

// AddNode adds a node and returns its ID
func (wb *WorkflowBuilder) AddNode(classType string, inputs map[string]interface{}) string {
	id := fmt.Sprintf("%d", wb.nextID)
	wb.nextID++

	wb.workflow[id] = Node{
		ClassType: classType,
		Inputs:    inputs,
	}

	return id
}

// AddNodeWithID adds a node with a specific ID
func (wb *WorkflowBuilder) AddNodeWithID(id string, classType string, inputs map[string]interface{}) {
	wb.workflow[id] = Node{
		ClassType: classType,
		Inputs:    inputs,
	}
}

// Build returns the built workflow
func (wb *WorkflowBuilder) Build() Workflow {
	return wb.workflow
}

// ConnectNodes creates a connection between two nodes
// The output of sourceNode will be used as input for targetNode
func (wb *WorkflowBuilder) ConnectNodes(sourceNodeID string, sourceOutput int, targetNodeID string, targetInput string) error {
	targetNode, ok := wb.workflow[targetNodeID]
	if !ok {
		return fmt.Errorf("target node %s not found", targetNodeID)
	}

	if targetNode.Inputs == nil {
		targetNode.Inputs = make(map[string]interface{})
	}

	targetNode.Inputs[targetInput] = []interface{}{sourceNodeID, sourceOutput}
	wb.workflow[targetNodeID] = targetNode

	return nil
}

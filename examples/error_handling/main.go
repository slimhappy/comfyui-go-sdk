package main

import (
	"context"
	"fmt"
	"log"
	"time"

	comfyui "github.com/yourusername/comfyui-go-sdk"
)

func main() {
	// Create client
	client := comfyui.NewClient("http://127.0.0.1:8188")
	ctx := context.Background()

	fmt.Println("=== ComfyUI Error Handling & Retry Example ===\n")

	// Example 1: Handle connection errors
	fmt.Println("1. Testing connection error handling...")
	badClient := comfyui.NewClient("http://127.0.0.1:9999") // Wrong port
	_, err := badClient.GetSystemStats(ctx)
	if err != nil {
		fmt.Printf("   ✓ Correctly caught connection error: %v\n", err)
	}

	// Example 2: Retry with exponential backoff
	fmt.Println("\n2. Retry with exponential backoff...")
	stats, err := retryWithBackoff(ctx, func() (interface{}, error) {
		return client.GetSystemStats(ctx)
	}, 3, time.Second)
	
	if err != nil {
		log.Printf("   Failed after retries: %v", err)
	} else {
		fmt.Printf("   ✓ Successfully retrieved stats after retry\n")
		if s, ok := stats.(*comfyui.SystemStats); ok {
			fmt.Printf("   OS: %s\n", s.System.OS)
		}
	}

	// Example 3: Handle invalid workflow errors
	fmt.Println("\n3. Testing invalid workflow handling...")
	invalidWorkflow := comfyui.Workflow{
		"1": comfyui.Node{
			ClassType: "NonExistentNode",
			Inputs:    map[string]interface{}{},
		},
	}
	
	_, err = client.QueuePrompt(ctx, invalidWorkflow, nil)
	if err != nil {
		fmt.Printf("   ✓ Correctly caught invalid workflow error: %v\n", err)
	}

	// Example 4: Handle missing model errors
	fmt.Println("\n4. Testing missing model handling...")
	workflowWithMissingModel := comfyui.Workflow{
		"1": comfyui.Node{
			ClassType: "CheckpointLoaderSimple",
			Inputs: map[string]interface{}{
				"ckpt_name": "non_existent_model.safetensors",
			},
		},
	}
	
	result, err := client.QueuePrompt(ctx, workflowWithMissingModel, nil)
	if err != nil {
		fmt.Printf("   ✓ Correctly caught missing model error: %v\n", err)
	} else if len(result.NodeErrors) > 0 {
		fmt.Printf("   ✓ Node errors detected: %v\n", result.NodeErrors)
	}

	// Example 5: Context timeout handling
	fmt.Println("\n5. Testing context timeout...")
	timeoutCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	
	workflow := buildLongRunningWorkflow()
	_, err = client.QueuePrompt(timeoutCtx, workflow, nil)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Printf("   ✓ Correctly handled timeout: %v\n", err)
		} else {
			fmt.Printf("   Request completed or failed: %v\n", err)
		}
	}

	// Example 6: Graceful degradation
	fmt.Println("\n6. Testing graceful degradation...")
	operations := []struct {
		name string
		fn   func() error
	}{
		{"Get System Stats", func() error {
			_, err := client.GetSystemStats(ctx)
			return err
		}},
		{"Get Queue", func() error {
			_, err := client.GetQueue(ctx)
			return err
		}},
		{"Get History", func() error {
			_, err := client.GetHistory(ctx, "")
			return err
		}},
		{"Get Models", func() error {
			_, err := client.GetModels(ctx, "checkpoints")
			return err
		}},
	}
	
	successCount := 0
	for _, op := range operations {
		err := op.fn()
		if err != nil {
			fmt.Printf("   ✗ %s failed: %v\n", op.name, err)
		} else {
			fmt.Printf("   ✓ %s succeeded\n", op.name)
			successCount++
		}
	}
	fmt.Printf("   Success rate: %d/%d\n", successCount, len(operations))

	// Example 7: Workflow validation before submission
	fmt.Println("\n7. Testing workflow validation...")
	testWorkflows := []struct {
		name     string
		workflow comfyui.Workflow
	}{
		{"Empty workflow", comfyui.Workflow{}},
		{"Valid workflow", buildSimpleWorkflow()},
		{"Workflow with missing inputs", comfyui.Workflow{
			"1": comfyui.Node{
				ClassType: "KSampler",
				Inputs:    map[string]interface{}{}, // Missing required inputs
			},
		}},
	}
	
	for _, test := range testWorkflows {
		err := test.workflow.Validate()
		if err != nil {
			fmt.Printf("   ✗ %s: %v\n", test.name, err)
		} else {
			fmt.Printf("   ✓ %s: valid\n", test.name)
		}
	}

	// Example 8: Safe workflow execution with recovery
	fmt.Println("\n8. Testing safe workflow execution...")
	err = safeExecuteWorkflow(ctx, client, buildSimpleWorkflow())
	if err != nil {
		fmt.Printf("   Workflow execution failed: %v\n", err)
	} else {
		fmt.Printf("   ✓ Workflow executed safely\n")
	}

	// Example 9: Monitor execution with error detection
	fmt.Println("\n9. Testing execution monitoring with error detection...")
	workflow = buildSimpleWorkflow()
	result, err = client.QueuePrompt(ctx, workflow, nil)
	if err != nil {
		log.Printf("   Failed to queue: %v", err)
	} else {
		fmt.Printf("   Queued: %s\n", result.PromptID)
		
		// Monitor with error detection
		err = monitorExecution(ctx, client, result.PromptID, 2*time.Minute)
		if err != nil {
			fmt.Printf("   ✗ Execution failed: %v\n", err)
		} else {
			fmt.Printf("   ✓ Execution completed successfully\n")
		}
	}

	// Example 10: Batch execution with error handling
	fmt.Println("\n10. Testing batch execution with error handling...")
	results := executeBatch(ctx, client, 3)
	fmt.Printf("   Batch results: %d succeeded, %d failed\n", 
		results.Succeeded, results.Failed)

	fmt.Println("\n=== Error Handling Example Complete ===")
}

// retryWithBackoff retries a function with exponential backoff
func retryWithBackoff(ctx context.Context, fn func() (interface{}, error), maxRetries int, initialDelay time.Duration) (interface{}, error) {
	var lastErr error
	delay := initialDelay
	
	for i := 0; i < maxRetries; i++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}
		
		lastErr = err
		fmt.Printf("   Attempt %d failed: %v\n", i+1, err)
		
		if i < maxRetries-1 {
			fmt.Printf("   Retrying in %v...\n", delay)
			select {
			case <-time.After(delay):
				delay *= 2 // Exponential backoff
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
	
	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// safeExecuteWorkflow executes a workflow with panic recovery
func safeExecuteWorkflow(ctx context.Context, client *comfyui.Client, workflow comfyui.Workflow) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during execution: %v", r)
		}
	}()
	
	// Validate workflow first
	if err := workflow.Validate(); err != nil {
		return fmt.Errorf("workflow validation failed: %w", err)
	}
	
	// Queue workflow
	result, err := client.QueuePrompt(ctx, workflow, nil)
	if err != nil {
		return fmt.Errorf("failed to queue: %w", err)
	}
	
	// Wait for completion with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	
	_, err = client.WaitForCompletion(timeoutCtx, result.PromptID)
	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}
	
	return nil
}

// monitorExecution monitors workflow execution and detects errors
func monitorExecution(ctx context.Context, client *comfyui.Client, promptID string, timeout time.Duration) error {
	ws, err := client.ConnectWebSocket(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect websocket: %w", err)
	}
	defer ws.Close()
	
	timeoutTimer := time.After(timeout)
	
	for {
		select {
		case <-timeoutTimer:
			return fmt.Errorf("execution timeout after %v", timeout)
			
		case <-ctx.Done():
			return ctx.Err()
			
		case err := <-ws.Errors():
			return fmt.Errorf("websocket error: %w", err)
			
		case msg, ok := <-ws.Messages():
			if !ok {
				return fmt.Errorf("websocket closed unexpectedly")
			}
			
			switch msg.Type {
			case string(comfyui.MessageTypeError):
				data, _ := msg.GetErrorData()
				if data.PromptID == promptID {
					return fmt.Errorf("execution error in node %s: %s", 
						data.NodeID, data.ExceptionMessage)
				}
				
			case string(comfyui.MessageTypeExecuting):
				data, _ := msg.GetExecutingData()
				if data.PromptID == promptID && data.Node == nil {
					// Execution completed
					return nil
				}
			}
		}
	}
}

// BatchResults contains results of batch execution
type BatchResults struct {
	Succeeded int
	Failed    int
	Errors    []error
}

// executeBatch executes multiple workflows with error handling
func executeBatch(ctx context.Context, client *comfyui.Client, count int) BatchResults {
	results := BatchResults{
		Errors: make([]error, 0),
	}
	
	for i := 0; i < count; i++ {
		workflow := buildSimpleWorkflow()
		workflow.SetNodeInput("3", "seed", 10000+i)
		
		result, err := client.QueuePrompt(ctx, workflow, nil)
		if err != nil {
			results.Failed++
			results.Errors = append(results.Errors, err)
			fmt.Printf("   Batch %d failed to queue: %v\n", i+1, err)
			continue
		}
		
		// Try to wait for completion with timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		_, err = client.WaitForCompletion(timeoutCtx, result.PromptID)
		cancel()
		
		if err != nil {
			results.Failed++
			results.Errors = append(results.Errors, err)
			fmt.Printf("   Batch %d failed to complete: %v\n", i+1, err)
		} else {
			results.Succeeded++
			fmt.Printf("   Batch %d completed successfully\n", i+1)
		}
	}
	
	return results
}

// buildSimpleWorkflow creates a simple workflow for testing
func buildSimpleWorkflow() comfyui.Workflow {
	return comfyui.Workflow{
		"3": comfyui.Node{
			ClassType: "KSampler",
			Inputs: map[string]interface{}{
				"seed":         12345,
				"steps":        10,
				"cfg":          7.0,
				"sampler_name": "euler",
				"scheduler":    "normal",
				"denoise":      1.0,
				"model":        []interface{}{"4", 0},
				"positive":     []interface{}{"6", 0},
				"negative":     []interface{}{"7", 0},
				"latent_image": []interface{}{"5", 0},
			},
		},
		"4": comfyui.Node{
			ClassType: "CheckpointLoaderSimple",
			Inputs: map[string]interface{}{
				"ckpt_name": "v1-5-pruned-emaonly.safetensors",
			},
		},
		"5": comfyui.Node{
			ClassType: "EmptyLatentImage",
			Inputs: map[string]interface{}{
				"width":      512,
				"height":     512,
				"batch_size": 1,
			},
		},
		"6": comfyui.Node{
			ClassType: "CLIPTextEncode",
			Inputs: map[string]interface{}{
				"text": "test image",
				"clip": []interface{}{"4", 1},
			},
		},
		"7": comfyui.Node{
			ClassType: "CLIPTextEncode",
			Inputs: map[string]interface{}{
				"text": "bad quality",
				"clip": []interface{}{"4", 1},
			},
		},
		"8": comfyui.Node{
			ClassType: "VAEDecode",
			Inputs: map[string]interface{}{
				"samples": []interface{}{"3", 0},
				"vae":     []interface{}{"4", 2},
			},
		},
		"9": comfyui.Node{
			ClassType: "SaveImage",
			Inputs: map[string]interface{}{
				"filename_prefix": "error_test",
				"images":          []interface{}{"8", 0},
			},
		},
	}
}

// buildLongRunningWorkflow creates a workflow that takes longer to execute
func buildLongRunningWorkflow() comfyui.Workflow {
	workflow := buildSimpleWorkflow()
	workflow.SetNodeInput("3", "steps", 50) // More steps = longer execution
	return workflow
}

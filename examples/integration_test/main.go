package main

import (
	"context"
	"fmt"
	"os"
	"time"

	comfyui "github.com/yourusername/comfyui-go-sdk"
)


func main() {
	// Create client
	serverURL := os.Getenv("COMFYUI_URL")
	if serverURL == "" {
		serverURL = "http://127.0.0.1:8188"
	}
	
	client := comfyui.NewClient(serverURL)
	ctx := context.Background()

	fmt.Println("=== ComfyUI Integration Test ===")
	fmt.Printf("Server: %s\n", serverURL)
	fmt.Printf("Client ID: %s\n\n", client.GetClientID())

	// Test suite
	tests := []struct {
		name string
		fn   func(context.Context, *comfyui.Client) error
	}{
		{"Server Connection", testServerConnection},
		{"System Information", testSystemInfo},
		{"Model Discovery", testModelDiscovery},
		{"Node Information", testNodeInfo},
		{"Queue Operations", testQueueOperations},
		{"Workflow Execution", testWorkflowExecution},
		{"WebSocket Monitoring", testWebSocketMonitoring},
		{"History Management", testHistoryManagement},
		{"Image Operations", testImageOperations},
		{"Error Handling", testErrorHandling},
	}

	results := make(map[string]bool)
	startTime := time.Now()

	for i, test := range tests {
		fmt.Printf("[%d/%d] Running: %s...\n", i+1, len(tests), test.name)
		
		testCtx, cancel := context.WithTimeout(ctx, 3*time.Minute)
		err := test.fn(testCtx, client)
		cancel()
		
		if err != nil {
			fmt.Printf("  ✗ FAILED: %v\n\n", err)
			results[test.name] = false
		} else {
			fmt.Printf("  ✓ PASSED\n\n")
			results[test.name] = true
		}
	}

	// Print summary
	duration := time.Since(startTime)
	fmt.Println("=== Test Summary ===")
	fmt.Printf("Duration: %v\n\n", duration)

	passed := 0
	failed := 0
	for name, result := range results {
		status := "✗ FAILED"
		if result {
			status = "✓ PASSED"
			passed++
		} else {
			failed++
		}
		fmt.Printf("%s: %s\n", status, name)
	}

	fmt.Printf("\nTotal: %d, Passed: %d, Failed: %d\n", len(tests), passed, failed)
	
	if failed > 0 {
		os.Exit(1)
	}
}

// testServerConnection tests basic server connectivity
func testServerConnection(ctx context.Context, client *comfyui.Client) error {
	_, err := client.GetSystemStats(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	fmt.Println("  - Server is reachable")
	return nil
}

// testSystemInfo tests system information retrieval
func testSystemInfo(ctx context.Context, client *comfyui.Client) error {
	stats, err := client.GetSystemStats(ctx)
	if err != nil {
		return fmt.Errorf("failed to get system stats: %w", err)
	}
	
	fmt.Printf("  - OS: %s\n", stats.System.OS)
	fmt.Printf("  - Python: %s\n", stats.System.PythonVersion)
	fmt.Printf("  - Devices: %d\n", len(stats.Devices))
	
	for _, device := range stats.Devices {
		fmt.Printf("    • %s (%s): %.2f GB VRAM\n", 
			device.Name, device.Type, float64(device.VRAMTotal)/1024/1024/1024)
	}
	
	return nil
}

// testModelDiscovery tests model discovery
func testModelDiscovery(ctx context.Context, client *comfyui.Client) error {
	checkpoints, err := client.GetModels(ctx, "checkpoints")
	if err != nil {
		return fmt.Errorf("failed to get checkpoints: %w", err)
	}
	
	fmt.Printf("  - Found %d checkpoints\n", len(checkpoints))
	
	if len(checkpoints) == 0 {
		return fmt.Errorf("no checkpoints found - please install at least one model")
	}
	
	// Test other model types
	modelTypes := []string{"vae", "loras"}
	for _, modelType := range modelTypes {
		models, err := client.GetModels(ctx, modelType)
		if err != nil {
			fmt.Printf("  - Warning: failed to get %s: %v\n", modelType, err)
		} else {
			fmt.Printf("  - Found %d %s\n", len(models), modelType)
		}
	}
	
	return nil
}

// testNodeInfo tests node information retrieval
func testNodeInfo(ctx context.Context, client *comfyui.Client) error {
	objectInfo, err := client.GetObjectInfo(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get object info: %w", err)
	}
	
	fmt.Printf("  - Total node classes: %d\n", len(objectInfo))
	
	// Check for essential nodes
	essentialNodes := []string{"KSampler", "CheckpointLoaderSimple", "CLIPTextEncode", "SaveImage"}
	for _, nodeName := range essentialNodes {
		if _, ok := objectInfo[nodeName]; !ok {
			return fmt.Errorf("essential node not found: %s", nodeName)
		}
	}
	fmt.Printf("  - All essential nodes present\n")
	
	return nil
}

// testQueueOperations tests queue management
func testQueueOperations(ctx context.Context, client *comfyui.Client) error {
	// Get initial queue status
	queue, err := client.GetQueue(ctx)
	if err != nil {
		return fmt.Errorf("failed to get queue: %w", err)
	}
	
	initialPending := len(queue.QueuePending)
	fmt.Printf("  - Initial queue: %d running, %d pending\n", 
		len(queue.QueueRunning), initialPending)
	
	// Queue a test workflow
	workflow := buildTestWorkflow(99999)
	result, err := client.QueuePrompt(ctx, workflow, map[string]interface{}{
		"test": "queue_operations",
	})
	if err != nil {
		return fmt.Errorf("failed to queue prompt: %w", err)
	}
	
	fmt.Printf("  - Queued test workflow: %s\n", result.PromptID)
	
	// Verify queue increased
	time.Sleep(500 * time.Millisecond)
	queue, err = client.GetQueue(ctx)
	if err != nil {
		return fmt.Errorf("failed to get queue after queueing: %w", err)
	}
	
	currentTotal := len(queue.QueueRunning) + len(queue.QueuePending)
	fmt.Printf("  - Queue after queueing: %d total items\n", currentTotal)
	
	// Delete the queued item
	err = client.DeleteFromQueue(ctx, []string{result.PromptID})
	if err != nil {
		fmt.Printf("  - Warning: failed to delete from queue: %v\n", err)
	} else {
		fmt.Printf("  - Successfully deleted test item from queue\n")
	}
	
	return nil
}

// testWorkflowExecution tests workflow execution
func testWorkflowExecution(ctx context.Context, client *comfyui.Client) error {
	workflow := buildTestWorkflow(12345)
	
	result, err := client.QueuePrompt(ctx, workflow, nil)
	if err != nil {
		return fmt.Errorf("failed to queue workflow: %w", err)
	}
	
	fmt.Printf("  - Queued workflow: %s\n", result.PromptID)
	
	// Wait for completion
	execResult, err := client.WaitForCompletion(ctx, result.PromptID)
	if err != nil {
		return fmt.Errorf("workflow execution failed: %w", err)
	}
	
	fmt.Printf("  - Execution completed in %v\n", execResult.Duration)
	fmt.Printf("  - Generated %d images\n", len(execResult.Images))
	
	if len(execResult.Images) == 0 {
		return fmt.Errorf("no images generated")
	}
	
	return nil
}

// testWebSocketMonitoring tests WebSocket functionality
func testWebSocketMonitoring(ctx context.Context, client *comfyui.Client) error {
	ws, err := client.ConnectWebSocket(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect websocket: %w", err)
	}
	defer ws.Close()
	
	fmt.Println("  - WebSocket connected")
	
	// Queue a workflow
	workflow := buildTestWorkflow(54321)
	result, err := client.QueuePrompt(ctx, workflow, nil)
	if err != nil {
		return fmt.Errorf("failed to queue workflow: %w", err)
	}
	
	fmt.Printf("  - Monitoring execution: %s\n", result.PromptID)
	
	// Monitor execution
	timeout := time.After(2 * time.Minute)
	progressCount := 0
	
	for {
		select {
		case <-timeout:
			return fmt.Errorf("monitoring timeout")
			
		case err := <-ws.Errors():
			return fmt.Errorf("websocket error: %w", err)
			
		case msg, ok := <-ws.Messages():
			if !ok {
				return fmt.Errorf("websocket closed")
			}
			
			switch msg.Type {
			case string(comfyui.MessageTypeProgress):
				progressCount++
				
			case string(comfyui.MessageTypeExecuting):
				data, _ := msg.GetExecutingData()
				if data.PromptID == result.PromptID && data.Node == nil {
					fmt.Printf("  - Received %d progress updates\n", progressCount)
					return nil
				}
				
			case string(comfyui.MessageTypeError):
				data, _ := msg.GetErrorData()
				if data.PromptID == result.PromptID {
					return fmt.Errorf("execution error: %s", data.ExceptionMessage)
				}
			}
		}
	}
}

// testHistoryManagement tests history operations
func testHistoryManagement(ctx context.Context, client *comfyui.Client) error {
	// Get all history
	history, err := client.GetHistory(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get history: %w", err)
	}
	
	fmt.Printf("  - Total history items: %d\n", len(history))
	
	// Queue a test workflow
	workflow := buildTestWorkflow(77777)
	result, err := client.QueuePrompt(ctx, workflow, nil)
	if err != nil {
		return fmt.Errorf("failed to queue workflow: %w", err)
	}
	
	// Wait for completion
	_, err = client.WaitForCompletion(ctx, result.PromptID)
	if err != nil {
		return fmt.Errorf("workflow execution failed: %w", err)
	}
	
	// Get specific history
	specificHistory, err := client.GetHistory(ctx, result.PromptID)
	if err != nil {
		return fmt.Errorf("failed to get specific history: %w", err)
	}
	
	if _, ok := specificHistory[result.PromptID]; !ok {
		return fmt.Errorf("history item not found")
	}
	
	fmt.Printf("  - Successfully retrieved history for: %s\n", result.PromptID)
	
	return nil
}

// testImageOperations tests image upload and download
func testImageOperations(ctx context.Context, client *comfyui.Client) error {
	// Create a test image
	testImagePath := "integration_test_image.png"
	err := createSimpleTestImage(testImagePath)
	if err != nil {
		return fmt.Errorf("failed to create test image: %w", err)
	}
	defer os.Remove(testImagePath)
	
	fmt.Println("  - Created test image")
	
	// Upload image
	uploadResult, err := client.UploadImage(ctx, testImagePath, comfyui.UploadOptions{
		Type:      "input",
		Overwrite: true,
	})
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}
	
	fmt.Printf("  - Uploaded image: %s\n", uploadResult.Name)
	
	// Download image
	imageData, err := client.GetImage(ctx, uploadResult.Name, uploadResult.Subfolder, uploadResult.Type)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	
	fmt.Printf("  - Downloaded image: %.2f KB\n", float64(len(imageData))/1024)
	
	if len(imageData) == 0 {
		return fmt.Errorf("downloaded image is empty")
	}
	
	return nil
}

// testErrorHandling tests error handling
func testErrorHandling(ctx context.Context, client *comfyui.Client) error {
	// Test invalid workflow
	invalidWorkflow := comfyui.Workflow{
		"1": comfyui.Node{
			ClassType: "NonExistentNode",
			Inputs:    map[string]interface{}{},
		},
	}
	
	_, err := client.QueuePrompt(ctx, invalidWorkflow, nil)
	if err == nil {
		return fmt.Errorf("expected error for invalid workflow, got nil")
	}
	
	fmt.Println("  - Correctly handled invalid workflow")
	
	// Test context cancellation
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel() // Cancel immediately
	
	_, err = client.GetSystemStats(cancelCtx)
	if err == nil {
		return fmt.Errorf("expected error for cancelled context, got nil")
	}
	
	fmt.Println("  - Correctly handled context cancellation")
	
	return nil
}

// buildTestWorkflow creates a simple test workflow
func buildTestWorkflow(seed int) comfyui.Workflow {
	return comfyui.Workflow{
		"3": comfyui.Node{
			ClassType: "KSampler",
			Inputs: map[string]interface{}{
				"seed":         seed,
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
				"text": "integration test image",
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
				"filename_prefix": "integration_test",
				"images":          []interface{}{"8", 0},
			},
		},
	}
}

// createSimpleTestImage creates a minimal test image
func createSimpleTestImage(path string) error {
	// Create a simple 1x1 PNG
	data := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0xCF, 0xC0, 0x00,
		0x00, 0x03, 0x01, 0x01, 0x00, 0x18, 0xDD, 0x8D,
		0xB4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E,
		0x44, 0xAE, 0x42, 0x60, 0x82,
	}
	
	return os.WriteFile(path, data, 0644)
}

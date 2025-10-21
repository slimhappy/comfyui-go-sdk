package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	comfyui "github.com/yourusername/comfyui-go-sdk"
)


// ProgressTracker tracks the progress of a workflow execution
type ProgressTracker struct {
	PromptID      string
	StartTime     time.Time
	CurrentNode   string
	TotalNodes    int
	CompletedNodes int
	CurrentStep   int
	TotalSteps    int
	LastUpdate    time.Time
	IsCompleted   bool
	HasError      bool
	ErrorMessage  string
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(promptID string) *ProgressTracker {
	return &ProgressTracker{
		PromptID:   promptID,
		StartTime:  time.Now(),
		LastUpdate: time.Now(),
	}
}

// Update updates the progress information
func (pt *ProgressTracker) Update(currentStep, totalSteps int, node string) {
	pt.CurrentStep = currentStep
	pt.TotalSteps = totalSteps
	pt.CurrentNode = node
	pt.LastUpdate = time.Now()
}

// CompleteNode marks a node as completed
func (pt *ProgressTracker) CompleteNode() {
	pt.CompletedNodes++
}

// SetError sets an error state
func (pt *ProgressTracker) SetError(msg string) {
	pt.HasError = true
	pt.ErrorMessage = msg
}

// Complete marks the workflow as completed
func (pt *ProgressTracker) Complete() {
	pt.IsCompleted = true
}

// GetElapsedTime returns the elapsed time since start
func (pt *ProgressTracker) GetElapsedTime() time.Duration {
	return time.Since(pt.StartTime)
}

// GetProgressPercentage returns the overall progress percentage
func (pt *ProgressTracker) GetProgressPercentage() float64 {
	if pt.TotalSteps == 0 {
		return 0
	}
	return float64(pt.CurrentStep) / float64(pt.TotalSteps) * 100
}

// DrawProgressBar draws a visual progress bar
func DrawProgressBar(current, total int, width int) string {
	if total == 0 {
		return strings.Repeat("░", width)
	}
	
	percentage := float64(current) / float64(total)
	filled := int(percentage * float64(width))
	
	if filled > width {
		filled = width
	}
	
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return bar
}

// ClearLine clears the current terminal line
func ClearLine() {
	fmt.Print("\r\033[K")
}

// PrintProgress prints the current progress with a visual bar
func PrintProgress(tracker *ProgressTracker) {
	ClearLine()
	
	if tracker.HasError {
		fmt.Printf("❌ Error: %s\n", tracker.ErrorMessage)
		return
	}
	
	if tracker.IsCompleted {
		elapsed := tracker.GetElapsedTime()
		fmt.Printf("✅ Completed in %s (Processed %d nodes)\n", 
			elapsed.Round(time.Millisecond), tracker.CompletedNodes)
		return
	}
	
	// Draw progress bar
	bar := DrawProgressBar(tracker.CurrentStep, tracker.TotalSteps, 40)
	percentage := tracker.GetProgressPercentage()
	elapsed := tracker.GetElapsedTime()
	
	// Calculate ETA
	var eta string
	if tracker.CurrentStep > 0 && percentage > 0 {
		totalEstimated := float64(elapsed) / (percentage / 100)
		remaining := time.Duration(totalEstimated) - elapsed
		eta = fmt.Sprintf(" ETA: %s", remaining.Round(time.Second))
	}
	
	fmt.Printf("⏳ [%s] %.1f%% | Step %d/%d | Node: %s | Time: %s%s",
		bar,
		percentage,
		tracker.CurrentStep,
		tracker.TotalSteps,
		tracker.CurrentNode,
		elapsed.Round(time.Second),
		eta,
	)
}

// MonitorProgress monitors the progress of a workflow execution
func MonitorProgress(ctx context.Context, client *comfyui.Client, promptID string) error {
	tracker := NewProgressTracker(promptID)
	
	// Connect to WebSocket
	ws, err := client.ConnectWebSocket(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	defer ws.Close()
	
	fmt.Printf("🚀 Monitoring workflow execution: %s\n\n", promptID)
	
	// Monitor messages
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
			
		case msg, ok := <-ws.Messages():
			if !ok {
				return fmt.Errorf("WebSocket closed")
			}
			
			handleProgressMessage(msg, tracker, promptID)
			PrintProgress(tracker)
			
			if tracker.IsCompleted || tracker.HasError {
				fmt.Println() // New line after completion
				return nil
			}
			
		case err := <-ws.Errors():
			tracker.SetError(err.Error())
			PrintProgress(tracker)
			fmt.Println()
			return err
		}
	}
}

// handleProgressMessage handles WebSocket messages for progress tracking
func handleProgressMessage(msg comfyui.WebSocketMessage, tracker *ProgressTracker, promptID string) {
	switch msg.Type {
	case string(comfyui.MessageTypeProgress):
		data, err := msg.GetProgressData()
		if err == nil {
			tracker.Update(data.Value, data.Max, tracker.CurrentNode)
		}
		
	case string(comfyui.MessageTypeExecuting):
		data, err := msg.GetExecutingData()
		if err == nil {
			if data.PromptID == promptID {
				if data.Node == nil {
					// Execution completed
					tracker.Complete()
				} else {
					// New node started
					tracker.CurrentNode = *data.Node
				}
			}
		}
		
	case string(comfyui.MessageTypeExecuted):
		data, err := msg.GetExecutedData()
		if err == nil && data.PromptID == promptID {
			tracker.CompleteNode()
		}
		
	case string(comfyui.MessageTypeError):
		data, err := msg.GetErrorData()
		if err == nil && data.PromptID == promptID {
			tracker.SetError(fmt.Sprintf("%s: %s", data.ExceptionType, data.ExceptionMessage))
		}
	}
}

// CreateSampleWorkflow creates a sample workflow for testing
func CreateSampleWorkflow() comfyui.Workflow {
	builder := comfyui.NewWorkflowBuilder()
	
	// Load checkpoint
	ckptID := builder.AddNode("CheckpointLoaderSimple", map[string]interface{}{
		"ckpt_name": "v1-5-pruned-emaonly.safetensors",
	})
	
	// Positive prompt
	posPromptID := builder.AddNode("CLIPTextEncode", map[string]interface{}{
		"text": "beautiful landscape, mountains, sunset, highly detailed, 8k",
	})
	builder.ConnectNodes(ckptID, 1, posPromptID, "clip")
	
	// Negative prompt
	negPromptID := builder.AddNode("CLIPTextEncode", map[string]interface{}{
		"text": "ugly, blurry, low quality",
	})
	builder.ConnectNodes(ckptID, 1, negPromptID, "clip")
	
	// Empty latent
	latentID := builder.AddNode("EmptyLatentImage", map[string]interface{}{
		"width":       512,
		"height":      512,
		"batch_size":  1,
	})
	
	// KSampler
	samplerID := builder.AddNode("KSampler", map[string]interface{}{
		"seed":          int(time.Now().Unix()),
		"steps":         20,
		"cfg":           7.0,
		"sampler_name":  "euler",
		"scheduler":     "normal",
		"denoise":       1.0,
	})
	builder.ConnectNodes(ckptID, 0, samplerID, "model")
	builder.ConnectNodes(posPromptID, 0, samplerID, "positive")
	builder.ConnectNodes(negPromptID, 0, samplerID, "negative")
	builder.ConnectNodes(latentID, 0, samplerID, "latent_image")
	
	// VAE Decode
	decodeID := builder.AddNode("VAEDecode", map[string]interface{}{})
	builder.ConnectNodes(samplerID, 0, decodeID, "samples")
	builder.ConnectNodes(ckptID, 2, decodeID, "vae")
	
	// Save image
	saveID := builder.AddNode("SaveImage", map[string]interface{}{
		"filename_prefix": "ComfyUI_Progress_Demo",
	})
	builder.ConnectNodes(decodeID, 0, saveID, "images")
	
	return builder.Build()
}

func main() {
	// Create client
	client := comfyui.NewClient("http://127.0.0.1:8188")
	
	ctx := context.Background()
	
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║        ComfyUI Go SDK - Progress Tracking Demo           ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
	
	// Create a sample workflow
	fmt.Println("📝 Creating sample workflow...")
	workflow := CreateSampleWorkflow()
	
	// Queue the workflow
	fmt.Println("📤 Submitting workflow to ComfyUI...")
	result, err := client.QueuePrompt(ctx, workflow, nil)
	if err != nil {
		log.Fatalf("Failed to queue prompt: %v", err)
	}
	
	fmt.Printf("✓ Workflow queued successfully (ID: %s)\n", result.PromptID)
	fmt.Println()
	
	// Monitor progress
	if err := MonitorProgress(ctx, client, result.PromptID); err != nil {
		log.Printf("Error monitoring progress: %v", err)
	}
	
	// Get the final result
	fmt.Println("\n📥 Retrieving execution results...")
	execResult, err := client.GetHistory(ctx, result.PromptID)
	if err != nil {
		log.Printf("Failed to get history: %v", err)
		return
	}
	
	// Display results
	if history, ok := execResult[result.PromptID]; ok {
		fmt.Println("\n📊 Execution Summary:")
		fmt.Printf("  • Status: %s\n", history.Status.StatusStr)
		
		if len(history.Outputs) > 0 {
			totalImages := 0
			for _, output := range history.Outputs {
				totalImages += len(output.Images)
			}
			fmt.Printf("  • Generated Images: %d\n", totalImages)

			
			// Save images
			if totalImages > 0 {
				fmt.Println("\n💾 Saving generated images...")
				for nodeID, output := range history.Outputs {
					for i, imgInfo := range output.Images {
						filename := fmt.Sprintf("output_%s_%d.png", nodeID, i)
						
						if err := client.SaveImage(ctx, imgInfo, filename); err != nil {
							log.Printf("Failed to save image: %v", err)
						} else {
							fmt.Printf("  ✓ Saved: %s\n", filename)
						}
					}
				}
			}

		}
	}
	
	fmt.Println("\n✨ Demo completed!")
}

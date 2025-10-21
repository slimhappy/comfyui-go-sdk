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

	// Example 1: Get system stats
	fmt.Println("=== System Stats ===")
	stats, err := client.GetSystemStats(ctx)
	if err != nil {
		log.Printf("Failed to get system stats: %v", err)
	} else {
		fmt.Printf("OS: %s\n", stats.System.OS)
		fmt.Printf("Python: %s\n", stats.System.PythonVersion)
		for _, device := range stats.Devices {
			fmt.Printf("Device: %s (%s)\n", device.Name, device.Type)
			fmt.Printf("  VRAM: %.2f GB / %.2f GB\n",
				float64(device.VRAMFree)/1024/1024/1024,
				float64(device.VRAMTotal)/1024/1024/1024)
		}
	}

	// Example 2: Get available models
	fmt.Println("\n=== Available Checkpoints ===")
	checkpoints, err := client.GetModels(ctx, "checkpoints")
	if err != nil {
		log.Printf("Failed to get checkpoints: %v", err)
	} else {
		for i, ckpt := range checkpoints {
			if i < 5 { // Show first 5
				fmt.Printf("  - %s\n", ckpt)
			}
		}
		if len(checkpoints) > 5 {
			fmt.Printf("  ... and %d more\n", len(checkpoints)-5)
		}
	}

	// Example 3: Build a simple workflow
	fmt.Println("\n=== Building Workflow ===")
	workflow := buildSimpleWorkflow()

	// Example 4: Submit workflow
	fmt.Println("Submitting workflow...")
	result, err := client.QueuePrompt(ctx, workflow, nil)
	if err != nil {
		log.Fatalf("Failed to queue prompt: %v", err)
	}
	fmt.Printf("Prompt ID: %s\n", result.PromptID)
	fmt.Printf("Queue Number: %d\n", result.Number)

	// Example 5: Monitor execution with WebSocket
	fmt.Println("\n=== Monitoring Execution ===")
	ws, err := client.ConnectWebSocket(ctx)
	if err != nil {
		log.Fatalf("Failed to connect websocket: %v", err)
	}
	defer ws.Close()

	// Monitor progress
	for {
		select {
		case msg, ok := <-ws.Messages():
			if !ok {
				fmt.Println("WebSocket closed")
				return
			}

			switch msg.Type {
			case string(comfyui.MessageTypeExecuting):
				data, _ := msg.GetExecutingData()
				if data.PromptID == result.PromptID {
					if data.Node == nil {
						fmt.Println("✓ Execution completed!")
						goto done
					} else {
						fmt.Printf("Executing node: %s\n", *data.Node)
					}
				}

			case string(comfyui.MessageTypeProgress):
				data, _ := msg.GetProgressData()
				fmt.Printf("Progress: %d/%d\n", data.Value, data.Max)

			case string(comfyui.MessageTypeError):
				data, _ := msg.GetErrorData()
				if data.PromptID == result.PromptID {
					fmt.Printf("✗ Error: %s\n", data.ExceptionMessage)
					return
				}
			}

		case err := <-ws.Errors():
			log.Printf("WebSocket error: %v", err)
			return

		case <-time.After(5 * time.Minute):
			log.Println("Timeout waiting for completion")
			return
		}
	}

done:
	// Example 6: Get results
	fmt.Println("\n=== Getting Results ===")
	history, err := client.GetHistory(ctx, result.PromptID)
	if err != nil {
		log.Fatalf("Failed to get history: %v", err)
	}

	if item, ok := history[result.PromptID]; ok {
		fmt.Printf("Status: %s\n", item.Status.StatusStr)
		fmt.Printf("Completed: %v\n", item.Status.Completed)

		// Download images
		for nodeID, output := range item.Outputs {
			if len(output.Images) > 0 {
				fmt.Printf("\nNode %s produced %d image(s):\n", nodeID, len(output.Images))
				for i, img := range output.Images {
					outputPath := fmt.Sprintf("output_%s_%d.png", nodeID, i)
					if err := client.SaveImage(ctx, img, outputPath); err != nil {
						log.Printf("Failed to save image: %v", err)
					} else {
						fmt.Printf("  ✓ Saved: %s\n", outputPath)
					}
				}
			}
		}
	}

	fmt.Println("\n=== Done ===")
}

// buildSimpleWorkflow creates a simple text-to-image workflow
func buildSimpleWorkflow() comfyui.Workflow {
	builder := comfyui.NewWorkflowBuilder()

	// Load checkpoint
	checkpointID := builder.AddNode("CheckpointLoaderSimple", map[string]interface{}{
		"ckpt_name": "v1-5-pruned-emaonly.safetensors",
	})

	// Create empty latent
	latentID := builder.AddNode("EmptyLatentImage", map[string]interface{}{
		"width":      512,
		"height":     512,
		"batch_size": 1,
	})

	// Positive prompt
	positiveID := builder.AddNode("CLIPTextEncode", map[string]interface{}{
		"text": "beautiful landscape, mountains, sunset, high quality",
	})
	builder.ConnectNodes(checkpointID, 1, positiveID, "clip")

	// Negative prompt
	negativeID := builder.AddNode("CLIPTextEncode", map[string]interface{}{
		"text": "bad quality, blurry, low resolution",
	})
	builder.ConnectNodes(checkpointID, 1, negativeID, "clip")

	// KSampler
	samplerID := builder.AddNode("KSampler", map[string]interface{}{
		"seed":         12345,
		"steps":        20,
		"cfg":          8.0,
		"sampler_name": "euler",
		"scheduler":    "normal",
		"denoise":      1.0,
	})
	builder.ConnectNodes(checkpointID, 0, samplerID, "model")
	builder.ConnectNodes(positiveID, 0, samplerID, "positive")
	builder.ConnectNodes(negativeID, 0, samplerID, "negative")
	builder.ConnectNodes(latentID, 0, samplerID, "latent_image")

	// VAE Decode
	decodeID := builder.AddNode("VAEDecode", map[string]interface{}{})
	builder.ConnectNodes(samplerID, 0, decodeID, "samples")
	builder.ConnectNodes(checkpointID, 2, decodeID, "vae")

	// Save Image
	builder.AddNode("SaveImage", map[string]interface{}{
		"filename_prefix": "ComfyUI",
	})
	builder.ConnectNodes(decodeID, 0, "9", "images")

	return builder.Build()
}

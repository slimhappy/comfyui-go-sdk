package main

import (
	"context"
	"fmt"
	"log"

	comfyui "github.com/yourusername/comfyui-go-sdk"
)


func main() {
	// Create client
	client := comfyui.NewClient("http://127.0.0.1:8188")
	ctx := context.Background()

	fmt.Println("=== ComfyUI History Operations Example ===\n")

	// Example 1: Get all history
	fmt.Println("1. Retrieving all history...")
	history, err := client.GetHistory(ctx, "")
	if err != nil {
		log.Fatalf("Failed to get history: %v", err)
	}
	fmt.Printf("   Total history items: %d\n", len(history))

	// Display recent history items
	if len(history) > 0 {
		fmt.Println("\n   Recent history items:")
		count := 0
		for promptID, item := range history {
			if count >= 5 {
				break
			}
			fmt.Printf("     - Prompt ID: %s\n", promptID)
			fmt.Printf("       Status: %s (Completed: %v)\n", item.Status.StatusStr, item.Status.Completed)
			fmt.Printf("       Outputs: %d nodes\n", len(item.Outputs))
			
			// Count images
			imageCount := 0
			for _, output := range item.Outputs {
				imageCount += len(output.Images)
			}
			if imageCount > 0 {
				fmt.Printf("       Images: %d\n", imageCount)
			}
			count++
		}
	}

	// Example 2: Queue a new workflow and track it
	fmt.Println("\n2. Queueing a new workflow to track...")
	workflow := buildSimpleWorkflow(42424)
	result, err := client.QueuePrompt(ctx, workflow, map[string]interface{}{
		"test_name": "history_tracking_test",
	})
	if err != nil {
		log.Fatalf("Failed to queue prompt: %v", err)
	}
	fmt.Printf("   ✓ Queued workflow: %s\n", result.PromptID)
	trackedPromptID := result.PromptID

	// Wait for completion
	fmt.Println("\n3. Waiting for workflow completion...")
	execResult, err := client.WaitForCompletion(ctx, trackedPromptID)
	if err != nil {
		log.Printf("   Failed to wait for completion: %v", err)
	} else {
		fmt.Printf("   ✓ Workflow completed in %v\n", execResult.Duration)
		fmt.Printf("   Generated %d images\n", len(execResult.Images))
	}

	// Example 4: Get specific history item
	fmt.Println("\n4. Retrieving specific history item...")
	specificHistory, err := client.GetHistory(ctx, trackedPromptID)
	if err != nil {
		log.Printf("   Failed to get specific history: %v", err)
	} else {
		if item, ok := specificHistory[trackedPromptID]; ok {
			fmt.Printf("   ✓ Found history for prompt: %s\n", trackedPromptID)
			fmt.Printf("   Status: %s\n", item.Status.StatusStr)
			fmt.Printf("   Completed: %v\n", item.Status.Completed)
			
			// Display output details
			fmt.Println("\n   Output details:")
			for nodeID, output := range item.Outputs {
				fmt.Printf("     Node %s:\n", nodeID)
				if len(output.Images) > 0 {
					fmt.Printf("       Images: %d\n", len(output.Images))
					for i, img := range output.Images {
						fmt.Printf("         [%d] %s (subfolder: %s, type: %s)\n", 
							i+1, img.Filename, img.Subfolder, img.Type)
					}
				}
				if len(output.Text) > 0 {
					fmt.Printf("       Text outputs: %d\n", len(output.Text))
				}
			}

			// Display workflow details
			fmt.Println("\n   Workflow details:")
			fmt.Printf("     Prompt ID: %s\n", item.Prompt.PromptID)
			fmt.Printf("     Number: %.0f\n", item.Prompt.Number)
			fmt.Printf("     Nodes in workflow: %d\n", len(item.Prompt.Workflow))
		}
	}

	// Example 5: Analyze history statistics
	fmt.Println("\n5. Analyzing history statistics...")
	history, err = client.GetHistory(ctx, "")
	if err != nil {
		log.Printf("   Failed to get history: %v", err)
	} else {
		stats := analyzeHistory(history)
		fmt.Printf("   Total executions: %d\n", stats.TotalExecutions)
		fmt.Printf("   Completed: %d\n", stats.CompletedExecutions)
		fmt.Printf("   Failed: %d\n", stats.FailedExecutions)
		fmt.Printf("   Total images generated: %d\n", stats.TotalImages)
		fmt.Printf("   Average images per execution: %.2f\n", stats.AvgImagesPerExecution)
		
		if len(stats.NodeClassUsage) > 0 {
			fmt.Println("\n   Most used node classes:")
			count := 0
			for class, usage := range stats.NodeClassUsage {
				if count >= 5 {
					break
				}
				fmt.Printf("     - %s: %d times\n", class, usage)
				count++
			}
		}
	}

	// Example 6: Download images from history
	fmt.Println("\n6. Downloading images from history...")
	if execResult != nil && len(execResult.Images) > 0 {
		for i, img := range execResult.Images {
			outputPath := fmt.Sprintf("history_image_%d.png", i+1)
			err := client.SaveImage(ctx, img, outputPath)
			if err != nil {
				log.Printf("   Failed to save image %d: %v", i+1, err)
			} else {
				fmt.Printf("   ✓ Saved image %d: %s\n", i+1, outputPath)
			}
		}
	} else {
		fmt.Println("   No images to download")
	}

	// Example 7: Delete specific history items
	fmt.Println("\n7. Deleting specific history items...")
	fmt.Println("   Note: Skipping actual deletion to preserve history")
	fmt.Println("   To delete history, uncomment the following code:")
	fmt.Printf("   // err = client.DeleteHistory(ctx, []string{\"%s\"})\n", trackedPromptID)
	
	// Uncomment to actually delete:
	// if trackedPromptID != "" {
	// 	err = client.DeleteHistory(ctx, []string{trackedPromptID})
	// 	if err != nil {
	// 		log.Printf("   Failed to delete history: %v", err)
	// 	} else {
	// 		fmt.Println("   ✓ History item deleted")
	// 	}
	// }

	// Example 8: Clear all history
	fmt.Println("\n8. Clear history demonstration...")
	fmt.Println("   Note: Skipping actual clear to preserve history")
	fmt.Println("   To clear all history, uncomment the following code:")
	fmt.Println("   // err = client.ClearHistory(ctx)")
	
	// Uncomment to actually clear:
	// err = client.ClearHistory(ctx)
	// if err != nil {
	// 	log.Printf("   Failed to clear history: %v", err)
	// } else {
	// 	fmt.Println("   ✓ All history cleared")
	// }

	// Example 9: Export history to JSON
	fmt.Println("\n9. Exporting history data...")
	if len(history) > 0 {
		// Get a sample history item
		var samplePromptID string
		for id := range history {
			samplePromptID = id
			break
		}
		
		_, err := client.GetHistory(ctx, samplePromptID)
		if err != nil {
			log.Printf("   Failed to get sample history: %v", err)
		} else {
			fmt.Printf("   ✓ Retrieved history for export: %s\n", samplePromptID)
			fmt.Println("   History data can be saved to JSON for analysis")
			// In a real application, you would marshal this to JSON and save to file
		}
	}


	fmt.Println("\n=== History Operations Example Complete ===")
}

// HistoryStats contains statistics about execution history
type HistoryStats struct {
	TotalExecutions       int
	CompletedExecutions   int
	FailedExecutions      int
	TotalImages           int
	AvgImagesPerExecution float64
	NodeClassUsage        map[string]int
}

// analyzeHistory analyzes history and returns statistics
func analyzeHistory(history comfyui.History) HistoryStats {
	stats := HistoryStats{
		NodeClassUsage: make(map[string]int),
	}

	stats.TotalExecutions = len(history)

	for _, item := range history {
		// Count completion status
		if item.Status.Completed {
			stats.CompletedExecutions++
		} else {
			stats.FailedExecutions++
		}

		// Count images
		for _, output := range item.Outputs {
			stats.TotalImages += len(output.Images)
		}

		// Count node class usage
		for _, node := range item.Prompt.Workflow {
			stats.NodeClassUsage[node.ClassType]++
		}
	}

	// Calculate average
	if stats.TotalExecutions > 0 {
		stats.AvgImagesPerExecution = float64(stats.TotalImages) / float64(stats.TotalExecutions)
	}

	return stats
}

// buildSimpleWorkflow creates a simple text-to-image workflow
func buildSimpleWorkflow(seed int) comfyui.Workflow {
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
				"text": "a beautiful test image",
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
				"filename_prefix": "history_test",
				"images":          []interface{}{"8", 0},
			},
		},
	}
}

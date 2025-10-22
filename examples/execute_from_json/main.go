package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	comfyui "github.com/yourusername/comfyui-go-sdk"
)

const (
	comfyUIURL = "http://127.0.0.1:8188"
)

func main() {
	// Parse command line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	workflowFile := os.Args[1]

	// Check if file exists
	if _, err := os.Stat(workflowFile); os.IsNotExist(err) {
		log.Fatalf("❌ Workflow file not found: %s", workflowFile)
	}

	fmt.Println("╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║         ComfyUI Go SDK - Execute Workflow from JSON File         ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Create client
	client := comfyui.NewClient(comfyUIURL)
	ctx := context.Background()

	// Check server status
	fmt.Println("🔍 Checking ComfyUI server status...")
	if err := checkServerStatus(ctx, client); err != nil {
		log.Fatalf("❌ Server check failed: %v", err)
	}
	fmt.Println("✅ ComfyUI server is running")
	fmt.Println()

	// Load and display workflow info
	fmt.Printf("📂 Loading workflow from: %s\n", workflowFile)
	workflow, err := comfyui.LoadWorkflowFromFile(workflowFile)
	if err != nil {
		log.Fatalf("❌ Failed to load workflow: %v", err)
	}

	displayWorkflowInfo(workflow)
	fmt.Println()

	// Optional: Modify workflow parameters
	if len(os.Args) > 2 {
		fmt.Println("🔧 Applying custom parameters...")
		if err := applyCustomParameters(workflow, os.Args[2:]); err != nil {
			log.Printf("⚠️  Warning: %v", err)
		}
		fmt.Println()
	}

	// Queue the workflow
	fmt.Println("🚀 Submitting workflow to ComfyUI...")
	resp, err := client.QueuePromptFromFile(ctx, workflowFile, nil)
	if err != nil {
		log.Fatalf("❌ Failed to queue workflow: %v", err)
	}

	fmt.Printf("✅ Workflow queued successfully!\n")
	fmt.Printf("   Prompt ID: %s\n", resp.PromptID)
	fmt.Printf("   Queue Position: %d\n", resp.Number)
	fmt.Println()

	// Monitor execution
	fmt.Println("⏳ Monitoring execution progress...")
	if err := monitorExecution(ctx, client, resp.PromptID); err != nil {
		log.Fatalf("❌ Execution monitoring failed: %v", err)
	}

	// Retrieve results
	fmt.Println()
	fmt.Println("📥 Retrieving execution results...")
	if err := retrieveResults(ctx, client, resp.PromptID); err != nil {
		log.Fatalf("❌ Failed to retrieve results: %v", err)
	}

	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    ✅ Execution Complete!                          ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════╝")
}

func printUsage() {
	fmt.Println("Usage: execute_from_json <workflow.json> [parameters...]")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Execute workflow from JSON file")
	fmt.Println("  ./execute_from_json workflow.json")
	fmt.Println()
	fmt.Println("  # Execute with custom parameters")
	fmt.Println("  ./execute_from_json workflow.json seed=12345 steps=30")
	fmt.Println()
	fmt.Println("Parameters format: key=value")
	fmt.Println("  seed=<number>      - Set random seed")
	fmt.Println("  steps=<number>     - Set sampling steps")
	fmt.Println("  cfg=<number>       - Set CFG scale")
	fmt.Println("  prompt=<text>      - Set positive prompt")
	fmt.Println("  negative=<text>    - Set negative prompt")
}

func checkServerStatus(ctx context.Context, client *comfyui.Client) error {
	queue, err := client.GetQueue(ctx)
	if err != nil {
		return err
	}
	_ = queue
	return nil
}

func displayWorkflowInfo(workflow comfyui.Workflow) {
	fmt.Printf("✅ Workflow loaded successfully\n")
	fmt.Printf("   Total nodes: %d\n", len(workflow))

	// Count nodes by type
	nodeTypes := make(map[string]int)
	for _, node := range workflow {
		nodeTypes[node.ClassType]++
	}

	fmt.Println("   Node types:")
	for classType, count := range nodeTypes {
		fmt.Printf("     - %s: %d\n", classType, count)
	}
}

func applyCustomParameters(workflow comfyui.Workflow, params []string) error {
	for _, param := range params {
		// Parse key=value format
		var key, value string
		if _, err := fmt.Sscanf(param, "%s=%s", &key, &value); err != nil {
			return fmt.Errorf("invalid parameter format: %s", param)
		}

		// Apply parameter based on key
		switch key {
		case "seed":
			var seed int
			if _, err := fmt.Sscanf(value, "%d", &seed); err != nil {
				return fmt.Errorf("invalid seed value: %s", value)
			}
			// Find KSampler nodes and update seed
			for id, node := range workflow {
				if node.ClassType == "KSampler" {
					workflow.SetNodeInput(id, "seed", seed)
					fmt.Printf("   ✓ Set seed=%d for node %s\n", seed, id)
				}
			}

		case "steps":
			var steps int
			if _, err := fmt.Sscanf(value, "%d", &steps); err != nil {
				return fmt.Errorf("invalid steps value: %s", value)
			}
			for id, node := range workflow {
				if node.ClassType == "KSampler" {
					workflow.SetNodeInput(id, "steps", steps)
					fmt.Printf("   ✓ Set steps=%d for node %s\n", steps, id)
				}
			}

		case "cfg":
			var cfg float64
			if _, err := fmt.Sscanf(value, "%f", &cfg); err != nil {
				return fmt.Errorf("invalid cfg value: %s", value)
			}
			for id, node := range workflow {
				if node.ClassType == "KSampler" {
					workflow.SetNodeInput(id, "cfg", cfg)
					fmt.Printf("   ✓ Set cfg=%.1f for node %s\n", cfg, id)
				}
			}

		case "prompt":
			for id, node := range workflow {
				if node.ClassType == "CLIPTextEncode" {
					// Assume first CLIPTextEncode is positive prompt
					workflow.SetNodeInput(id, "text", value)
					fmt.Printf("   ✓ Set prompt='%s' for node %s\n", value, id)
					break
				}
			}

		case "negative":
			count := 0
			for id, node := range workflow {
				if node.ClassType == "CLIPTextEncode" {
					count++
					if count == 2 {
						// Assume second CLIPTextEncode is negative prompt
						workflow.SetNodeInput(id, "text", value)
						fmt.Printf("   ✓ Set negative='%s' for node %s\n", value, id)
						break
					}
				}
			}

		default:
			fmt.Printf("   ⚠️  Unknown parameter: %s\n", key)
		}
	}

	return nil
}

func monitorExecution(ctx context.Context, client *comfyui.Client, promptID string) error {
	startTime := time.Now()
	lastStatus := ""

	for {
		// Check queue status
		queue, err := client.GetQueue(ctx)
		if err != nil {
			return fmt.Errorf("failed to get queue: %w", err)
		}

		// Check if still in queue
		inQueue := false
		for _, item := range queue.QueuePending {
			if item.PromptID == promptID {
				inQueue = true
				break
			}
		}

		// Check if currently running
		running := false
		for _, item := range queue.QueueRunning {
			if item.PromptID == promptID {
				running = true
				break
			}
		}


		// Update status
		var status string
		if inQueue {
			status = "⏳ In queue..."
		} else if running {
			status = "🔄 Executing..."
		} else {
			// Check if completed
			history, err := client.GetHistory(ctx, promptID)
			if err == nil && len(history) > 0 {
				elapsed := time.Since(startTime)
				fmt.Printf("\r✅ Completed in %.1f seconds\n", elapsed.Seconds())
				return nil
			}
			status = "⏳ Waiting..."
		}

		// Print status if changed
		if status != lastStatus {
			fmt.Printf("\r%s", status)
			lastStatus = status
		}

		// Wait before next check
		time.Sleep(500 * time.Millisecond)

		// Timeout after 5 minutes
		if time.Since(startTime) > 5*time.Minute {
			return fmt.Errorf("execution timeout")
		}
	}
}

func retrieveResults(ctx context.Context, client *comfyui.Client, promptID string) error {
	history, err := client.GetHistory(ctx, promptID)
	if err != nil {
		return fmt.Errorf("failed to get history: %w", err)
	}

	if len(history) == 0 {
		return fmt.Errorf("no history found for prompt ID: %s", promptID)
	}

	for id, item := range history {
		fmt.Printf("\n📊 Execution ID: %s\n", id)
		fmt.Printf("   Status: %s\n", item.Status.StatusStr)

		if len(item.Status.Messages) > 0 {
			fmt.Println("   Messages:")
			for _, msg := range item.Status.Messages {
				msgJSON, _ := json.MarshalIndent(msg, "     ", "  ")
				fmt.Printf("     %s\n", string(msgJSON))
			}
		}

		if len(item.Outputs) > 0 {
			fmt.Println("   Outputs:")
			for nodeID, output := range item.Outputs {
				fmt.Printf("     Node %s:\n", nodeID)

				if len(output.Images) > 0 {
					fmt.Printf("       Images: %d\n", len(output.Images))
					for i, img := range output.Images {
						fmt.Printf("         [%d] %s (type: %s, subfolder: %s)\n",
							i+1, img.Filename, img.Type, img.Subfolder)

						// Save image
						outputDir := "output"
						os.MkdirAll(outputDir, 0755)
						outputPath := filepath.Join(outputDir, img.Filename)

					if err := client.SaveImage(ctx, img, outputPath); err != nil {
						log.Printf("⚠️  Failed to save image: %v", err)
					} else {
						fmt.Printf("         💾 Saved to: %s\n", outputPath)
					}

					}
				}



			}
		}
	}

	return nil
}

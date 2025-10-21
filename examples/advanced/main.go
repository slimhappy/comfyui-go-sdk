package main

import (
	"context"
	"fmt"
	"log"
	"time"

	comfyui "github.com/yourusername/comfyui-go-sdk"
)

func main() {
	client := comfyui.NewClient("http://127.0.0.1:8188")
	ctx := context.Background()

	// Example 1: Upload image and use in workflow
	fmt.Println("=== Example 1: Image Upload ===")
	uploadedImage, err := client.UploadImage(ctx, "input.png", comfyui.UploadOptions{
		Type:      "input",
		Subfolder: "",
		Overwrite: true,
	})
	if err != nil {
		log.Printf("Failed to upload image: %v", err)
	} else {
		fmt.Printf("Uploaded: %s\n", uploadedImage.Name)

		// Use uploaded image in workflow
		workflow := buildImg2ImgWorkflow(uploadedImage.Name)
		result, err := client.QueuePrompt(ctx, workflow, nil)
		if err != nil {
			log.Printf("Failed to queue prompt: %v", err)
		} else {
			fmt.Printf("Queued prompt: %s\n", result.PromptID)
		}
	}

	// Example 2: Batch processing with different seeds
	fmt.Println("\n=== Example 2: Batch Processing ===")
	baseWorkflow := buildSimpleWorkflow()

	seeds := []int{12345, 67890, 11111, 22222, 33333}
	promptIDs := make([]string, 0, len(seeds))

	for _, seed := range seeds {
		workflow, _ := baseWorkflow.Clone()
		workflow.SetNodeInput("6", "seed", seed) // Assuming node 6 is KSampler

		result, err := client.QueuePrompt(ctx, workflow, nil)
		if err != nil {
			log.Printf("Failed to queue seed %d: %v", seed, err)
			continue
		}

		fmt.Printf("Queued seed %d: %s\n", seed, result.PromptID)
		promptIDs = append(promptIDs, result.PromptID)
	}

	// Wait for all to complete
	fmt.Println("\nWaiting for batch completion...")
	for i, promptID := range promptIDs {
		fmt.Printf("Waiting for %d/%d...\n", i+1, len(promptIDs))
		result, err := client.WaitForCompletion(ctx, promptID)
		if err != nil {
			log.Printf("Failed to wait for %s: %v", promptID, err)
			continue
		}

		fmt.Printf("  ✓ Completed in %v, %d images\n", result.Duration, len(result.Images))

		// Save images
		for j, img := range result.Images {
			outputPath := fmt.Sprintf("batch_%d_seed_%d_%d.png", i, seeds[i], j)
			if err := client.SaveImage(ctx, img, outputPath); err != nil {
				log.Printf("Failed to save: %v", err)
			}
		}
	}

	// Example 3: Queue management
	fmt.Println("\n=== Example 3: Queue Management ===")
	queue, err := client.GetQueue(ctx)
	if err != nil {
		log.Printf("Failed to get queue: %v", err)
	} else {
		fmt.Printf("Running: %d, Pending: %d\n", len(queue.QueueRunning), len(queue.QueuePending))

		if len(queue.QueuePending) > 0 {
			fmt.Println("Pending items:")
			for _, item := range queue.QueuePending {
				fmt.Printf("  - Prompt ID: %s, Number: %d\n", item.PromptID, item.Number)
			}
		}
	}

	// Example 4: History management
	fmt.Println("\n=== Example 4: History Management ===")
	history, err := client.GetHistory(ctx, "")
	if err != nil {
		log.Printf("Failed to get history: %v", err)
	} else {
		fmt.Printf("Total history items: %d\n", len(history))

		// Show recent items
		count := 0
		for promptID, item := range history {
			if count >= 5 {
				break
			}
			fmt.Printf("  - %s: %s (completed: %v)\n", promptID, item.Status.StatusStr, item.Status.Completed)
			count++
		}
	}

	// Example 5: Node information
	fmt.Println("\n=== Example 5: Node Information ===")
	objectInfo, err := client.GetObjectInfo(ctx, "KSampler")
	if err != nil {
		log.Printf("Failed to get object info: %v", err)
	} else {
		if info, ok := objectInfo["KSampler"]; ok {
			fmt.Printf("Node: %s\n", info.DisplayName)
			fmt.Printf("Category: %s\n", info.Category)
			fmt.Printf("Description: %s\n", info.Description)
			fmt.Printf("Inputs:\n")
			for name := range info.Input.Required {
				fmt.Printf("  - %s (required)\n", name)
			}
			for name := range info.Input.Optional {
				fmt.Printf("  - %s (optional)\n", name)
			}
		}
	}

	// Example 6: Workflow from file
	fmt.Println("\n=== Example 6: Load Workflow from File ===")
	workflow, err := comfyui.LoadWorkflowFromFile("workflow_api.json")
	if err != nil {
		log.Printf("Failed to load workflow: %v (make sure workflow_api.json exists)", err)
	} else {
		fmt.Printf("Loaded workflow with %d nodes\n", len(workflow))

		// Modify workflow
		workflow.SetNodeInput("6", "text", "a beautiful sunset over mountains")

		// Save modified workflow
		if err := comfyui.SaveWorkflowToFile(workflow, "workflow_modified.json"); err != nil {
			log.Printf("Failed to save workflow: %v", err)
		} else {
			fmt.Println("Saved modified workflow to workflow_modified.json")
		}
	}

	// Example 7: Concurrent execution with context timeout
	fmt.Println("\n=== Example 7: Concurrent Execution ===")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	results := make(chan *comfyui.ExecutionResult, 3)
	errors := make(chan error, 3)

	for i := 0; i < 3; i++ {
		go func(index int) {
			workflow := buildSimpleWorkflow()
			workflow.SetNodeInput("6", "seed", 10000+index)

			result, err := client.QueuePrompt(ctx, workflow, nil)
			if err != nil {
				errors <- err
				return
			}

			execResult, err := client.WaitForCompletion(ctx, result.PromptID)
			if err != nil {
				errors <- err
				return
			}

			results <- execResult
		}(i)
	}

	// Collect results
	completed := 0
	for completed < 3 {
		select {
		case result := <-results:
			completed++
			fmt.Printf("  ✓ Workflow %d completed in %v\n", completed, result.Duration)

		case err := <-errors:
			completed++
			fmt.Printf("  ✗ Workflow failed: %v\n", err)

		case <-ctx.Done():
			fmt.Println("  ✗ Timeout reached")
			return
		}
	}

	fmt.Println("\n=== All Examples Completed ===")
}

func buildSimpleWorkflow() comfyui.Workflow {
	return comfyui.Workflow{
		"3": comfyui.Node{
			ClassType: "KSampler",
			Inputs: map[string]interface{}{
				"seed":         12345,
				"steps":        20,
				"cfg":          8.0,
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
				"text": "beautiful landscape",
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
				"filename_prefix": "ComfyUI",
				"images":          []interface{}{"8", 0},
			},
		},
	}
}

func buildImg2ImgWorkflow(imageName string) comfyui.Workflow {
	workflow := buildSimpleWorkflow()

	// Replace EmptyLatentImage with LoadImage
	workflow["10"] = comfyui.Node{
		ClassType: "LoadImage",
		Inputs: map[string]interface{}{
			"image": imageName,
		},
	}

	// Add VAEEncode
	workflow["11"] = comfyui.Node{
		ClassType: "VAEEncode",
		Inputs: map[string]interface{}{
			"pixels": []interface{}{"10", 0},
			"vae":    []interface{}{"4", 2},
		},
	}

	// Update KSampler to use encoded image
	workflow.SetNodeInput("3", "latent_image", []interface{}{"11", 0})
	workflow.SetNodeInput("3", "denoise", 0.75)

	return workflow
}

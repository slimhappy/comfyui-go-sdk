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

	fmt.Println("=== ComfyUI Queue Management Example ===\n")

	// Example 1: Check current queue status
	fmt.Println("1. Checking current queue status...")
	queue, err := client.GetQueue(ctx)
	if err != nil {
		log.Fatalf("Failed to get queue: %v", err)
	}
	fmt.Printf("   Running: %d items\n", len(queue.QueueRunning))
	fmt.Printf("   Pending: %d items\n", len(queue.QueuePending))

	// Display queue details
	if len(queue.QueueRunning) > 0 {
		fmt.Println("\n   Running items:")
		for _, item := range queue.QueueRunning {
			fmt.Printf("     - Prompt ID: %s, Number: %d\n", item.PromptID, item.Number)
		}
	}

	if len(queue.QueuePending) > 0 {
		fmt.Println("\n   Pending items:")
		for i, item := range queue.QueuePending {
			if i < 5 { // Show first 5
				fmt.Printf("     - Prompt ID: %s, Number: %d\n", item.PromptID, item.Number)
			}
		}
		if len(queue.QueuePending) > 5 {
			fmt.Printf("     ... and %d more\n", len(queue.QueuePending)-5)
		}
	}

	// Example 2: Queue multiple workflows
	fmt.Println("\n2. Queueing multiple workflows...")
	promptIDs := []string{}
	
	for i := 0; i < 3; i++ {
		workflow := buildSimpleWorkflow(12345 + i)
		result, err := client.QueuePrompt(ctx, workflow, map[string]interface{}{
			"batch_name": fmt.Sprintf("test_batch_%d", i),
		})
		if err != nil {
			log.Printf("   Failed to queue workflow %d: %v", i, err)
			continue
		}
		promptIDs = append(promptIDs, result.PromptID)
		fmt.Printf("   ✓ Queued workflow %d: %s (number: %d)\n", i+1, result.PromptID, result.Number)
		time.Sleep(100 * time.Millisecond) // Small delay between submissions
	}

	// Example 3: Monitor queue changes
	fmt.Println("\n3. Monitoring queue status...")
	time.Sleep(1 * time.Second)
	queue, err = client.GetQueue(ctx)
	if err != nil {
		log.Printf("   Failed to get queue: %v", err)
	} else {
		fmt.Printf("   Running: %d items\n", len(queue.QueueRunning))
		fmt.Printf("   Pending: %d items\n", len(queue.QueuePending))
	}

	// Example 4: Delete specific items from queue
	if len(promptIDs) > 0 && len(queue.QueuePending) > 0 {
		fmt.Println("\n4. Deleting specific items from queue...")
		// Delete the last queued item
		toDelete := promptIDs[len(promptIDs)-1:]
		fmt.Printf("   Deleting prompt ID: %s\n", toDelete[0])
		
		err = client.DeleteFromQueue(ctx, toDelete)
		if err != nil {
			log.Printf("   Failed to delete from queue: %v", err)
		} else {
			fmt.Println("   ✓ Successfully deleted item from queue")
			
			// Verify deletion
			time.Sleep(500 * time.Millisecond)
			queue, err = client.GetQueue(ctx)
			if err == nil {
				fmt.Printf("   Queue status after deletion - Pending: %d items\n", len(queue.QueuePending))
			}
		}
	}

	// Example 5: Interrupt current execution
	fmt.Println("\n5. Testing interrupt functionality...")
	if len(queue.QueueRunning) > 0 {
		runningItem := queue.QueueRunning[0]
		fmt.Printf("   Interrupting prompt: %s\n", runningItem.PromptID)
		
		err = client.Interrupt(ctx, runningItem.PromptID)
		if err != nil {
			log.Printf("   Failed to interrupt: %v", err)
		} else {
			fmt.Println("   ✓ Interrupt signal sent")
		}
	} else {
		fmt.Println("   No running items to interrupt")
	}

	// Example 6: Clear entire queue
	fmt.Println("\n6. Queue clearing demonstration...")
	fmt.Println("   Note: Skipping actual clear to preserve other work")
	fmt.Println("   To clear queue, uncomment the following code:")
	fmt.Println("   // err = client.ClearQueue(ctx)")
	
	// Uncomment to actually clear the queue:
	// err = client.ClearQueue(ctx)
	// if err != nil {
	// 	log.Printf("   Failed to clear queue: %v", err)
	// } else {
	// 	fmt.Println("   ✓ Queue cleared successfully")
	// }

	// Example 7: Queue with priority (using extra_data)
	fmt.Println("\n7. Queueing with metadata...")
	workflow := buildSimpleWorkflow(99999)
	result, err := client.QueuePrompt(ctx, workflow, map[string]interface{}{
		"priority":    "high",
		"user":        "test_user",
		"description": "Test workflow with metadata",
		"timestamp":   time.Now().Unix(),
	})
	if err != nil {
		log.Printf("   Failed to queue with metadata: %v", err)
	} else {
		fmt.Printf("   ✓ Queued with metadata: %s\n", result.PromptID)
	}

	// Example 8: Monitor queue until empty (with timeout)
	fmt.Println("\n8. Monitoring queue until completion...")
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("   ⏱ Timeout reached")
			goto done
		case <-ticker.C:
			queue, err := client.GetQueue(ctx)
			if err != nil {
				log.Printf("   Failed to get queue: %v", err)
				continue
			}
			
			total := len(queue.QueueRunning) + len(queue.QueuePending)
			fmt.Printf("   Queue status: %d total items (Running: %d, Pending: %d)\n", 
				total, len(queue.QueueRunning), len(queue.QueuePending))
			
			if total == 0 {
				fmt.Println("   ✓ Queue is empty")
				goto done
			}
		}
	}

done:
	// Final status
	fmt.Println("\n9. Final queue status...")
	queue, err = client.GetQueue(ctx)
	if err != nil {
		log.Printf("   Failed to get final queue status: %v", err)
	} else {
		fmt.Printf("   Running: %d items\n", len(queue.QueueRunning))
		fmt.Printf("   Pending: %d items\n", len(queue.QueuePending))
	}

	fmt.Println("\n=== Queue Management Example Complete ===")
}

// buildSimpleWorkflow creates a simple text-to-image workflow
func buildSimpleWorkflow(seed int) comfyui.Workflow {
	return comfyui.Workflow{
		"3": comfyui.Node{
			ClassType: "KSampler",
			Inputs: map[string]interface{}{
				"seed":         seed,
				"steps":        10, // Reduced steps for faster testing
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
				"text": "a simple test image",
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
				"filename_prefix": "queue_test",
				"images":          []interface{}{"8", 0},
			},
		},
	}
}

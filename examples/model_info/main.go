package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	comfyui "github.com/yourusername/comfyui-go-sdk"
)

func main() {
	// Create client
	client := comfyui.NewClient("http://127.0.0.1:8188")
	ctx := context.Background()

	fmt.Println("=== ComfyUI Model & Node Information Example ===\n")

	// Example 1: Get system statistics
	fmt.Println("1. System Statistics...")
	stats, err := client.GetSystemStats(ctx)
	if err != nil {
		log.Printf("Failed to get system stats: %v", err)
	} else {
		fmt.Printf("   OS: %s\n", stats.System.OS)
		fmt.Printf("   Python Version: %s\n", stats.System.PythonVersion)
		fmt.Printf("   Embedded Python: %v\n", stats.System.EmbeddedPython)
		
		fmt.Println("\n   Devices:")
		for i, device := range stats.Devices {
			fmt.Printf("     [%d] %s (%s)\n", i, device.Name, device.Type)
			fmt.Printf("         Index: %d\n", device.Index)
			fmt.Printf("         VRAM Total: %.2f GB\n", float64(device.VRAMTotal)/1024/1024/1024)
			fmt.Printf("         VRAM Free: %.2f GB\n", float64(device.VRAMFree)/1024/1024/1024)
			fmt.Printf("         VRAM Usage: %.1f%%\n", 
				float64(device.VRAMTotal-device.VRAMFree)/float64(device.VRAMTotal)*100)
			
			if device.TorchVRAMTotal > 0 {
				fmt.Printf("         Torch VRAM Total: %.2f GB\n", 
					float64(device.TorchVRAMTotal)/1024/1024/1024)
				fmt.Printf("         Torch VRAM Free: %.2f GB\n", 
					float64(device.TorchVRAMFree)/1024/1024/1024)
			}
		}
	}

	// Example 2: Get available model folders
	fmt.Println("\n2. Available Model Folders...")
	folders, err := client.GetModels(ctx, "")
	if err != nil {
		log.Printf("Failed to get model folders: %v", err)
	} else {
		fmt.Printf("   Found %d model folders:\n", len(folders))
		for _, folder := range folders {
			fmt.Printf("     - %s\n", folder)
		}
	}

	// Example 3: Get models in specific folders
	fmt.Println("\n3. Models in Specific Folders...")
	modelFolders := []string{"checkpoints", "vae", "loras", "embeddings", "upscale_models"}
	
	for _, folder := range modelFolders {
		models, err := client.GetModels(ctx, folder)
		if err != nil {
			log.Printf("   Failed to get models for %s: %v", folder, err)
			continue
		}
		
		fmt.Printf("\n   %s (%d models):\n", folder, len(models))
		if len(models) > 0 {
			// Show first 5 models
			for i, model := range models {
				if i < 5 {
					fmt.Printf("     - %s\n", model)
				}
			}
			if len(models) > 5 {
				fmt.Printf("     ... and %d more\n", len(models)-5)
			}
		} else {
			fmt.Println("     (no models found)")
		}
	}

	// Example 4: Get embeddings
	fmt.Println("\n4. Available Embeddings...")
	embeddings, err := client.GetEmbeddings(ctx)
	if err != nil {
		log.Printf("Failed to get embeddings: %v", err)
	} else {
		fmt.Printf("   Found %d embeddings:\n", len(embeddings))
		for i, emb := range embeddings {
			if i < 10 {
				fmt.Printf("     - %s\n", emb)
			}
		}
		if len(embeddings) > 10 {
			fmt.Printf("     ... and %d more\n", len(embeddings)-10)
		}
	}

	// Example 5: Get all node classes
	fmt.Println("\n5. Available Node Classes...")
	objectInfo, err := client.GetObjectInfo(ctx, "")
	if err != nil {
		log.Printf("Failed to get object info: %v", err)
	} else {
		fmt.Printf("   Total node classes: %d\n", len(objectInfo))
		
		// Categorize nodes
		categories := make(map[string][]string)
		for className, info := range objectInfo {
			categories[info.Category] = append(categories[info.Category], className)
		}
		
		fmt.Printf("\n   Node categories (%d):\n", len(categories))
		for category, nodes := range categories {
			fmt.Printf("     - %s: %d nodes\n", category, len(nodes))
		}
	}

	// Example 6: Get specific node information
	fmt.Println("\n6. Detailed Node Information...")
	nodesToInspect := []string{"KSampler", "CheckpointLoaderSimple", "CLIPTextEncode", "SaveImage", "LoadImage"}
	
	for _, nodeName := range nodesToInspect {
		nodeInfo, err := client.GetObjectInfo(ctx, nodeName)
		if err != nil {
			log.Printf("   Failed to get info for %s: %v", nodeName, err)
			continue
		}
		
		if info, ok := nodeInfo[nodeName]; ok {
			fmt.Printf("\n   Node: %s\n", nodeName)
			fmt.Printf("     Display Name: %s\n", info.DisplayName)
			fmt.Printf("     Category: %s\n", info.Category)
			if info.Description != "" {
				fmt.Printf("     Description: %s\n", info.Description)
			}
			fmt.Printf("     Output Node: %v\n", info.OutputNode)
			
			// Input information
			fmt.Println("     Required Inputs:")
			for inputName, inputSpec := range info.Input.Required {
				fmt.Printf("       - %s: %v\n", inputName, inputSpec)
			}
			
			if len(info.Input.Optional) > 0 {
				fmt.Println("     Optional Inputs:")
				for inputName, inputSpec := range info.Input.Optional {
					fmt.Printf("       - %s: %v\n", inputName, inputSpec)
				}
			}
			
			// Output information
			if len(info.Output) > 0 {
				fmt.Println("     Outputs:")
				for i, output := range info.Output {
					outputName := output
					if i < len(info.OutputName) {
						outputName = info.OutputName[i]
					}
					fmt.Printf("       [%d] %s (%s)\n", i, outputName, output)
				}
			}
		}
	}

	// Example 7: Search for nodes by category
	fmt.Println("\n7. Searching Nodes by Category...")
	searchCategories := []string{"sampling", "loaders", "conditioning"}
	
	for _, searchCat := range searchCategories {
		fmt.Printf("\n   Nodes in '%s' category:\n", searchCat)
		count := 0
		for className, info := range objectInfo {
			if strings.Contains(strings.ToLower(info.Category), strings.ToLower(searchCat)) {
				fmt.Printf("     - %s (%s)\n", className, info.DisplayName)
				count++
				if count >= 5 {
					break
				}
			}
		}
		if count == 0 {
			fmt.Println("     (no nodes found)")
		}
	}

	// Example 8: Find nodes with specific output types
	fmt.Println("\n8. Finding Nodes by Output Type...")
	outputTypes := []string{"MODEL", "CLIP", "VAE", "LATENT", "IMAGE"}
	
	for _, outputType := range outputTypes {
		fmt.Printf("\n   Nodes that output '%s':\n", outputType)
		count := 0
		for className, info := range objectInfo {
			for _, output := range info.Output {
				if output == outputType {
					fmt.Printf("     - %s\n", className)
					count++
					if count >= 5 {
						break
					}
				}
			}
			if count >= 5 {
				break
			}
		}
		if count == 0 {
			fmt.Println("     (no nodes found)")
		}
	}

	// Example 9: Get server features
	fmt.Println("\n9. Server Features...")
	features, err := client.GetFeatures(ctx)
	if err != nil {
		log.Printf("Failed to get features: %v", err)
	} else {
		if len(features.Features) > 0 {
			fmt.Println("   Supported features:")
			for _, feature := range features.Features {
				fmt.Printf("     - %s\n", feature)
			}
		} else {
			fmt.Println("   No special features reported")
		}
	}

	// Example 10: Memory management information
	fmt.Println("\n10. Memory Management...")
	fmt.Println("   Available memory management operations:")
	fmt.Println("     - Free memory: client.FreeMemory(ctx, false, true)")
	fmt.Println("     - Unload models: client.FreeMemory(ctx, true, false)")
	fmt.Println("     - Both: client.FreeMemory(ctx, true, true)")
	fmt.Println("\n   Note: Skipping actual memory operations")
	
	// Uncomment to actually free memory:
	// err = client.FreeMemory(ctx, true, true)
	// if err != nil {
	// 	log.Printf("   Failed to free memory: %v", err)
	// } else {
	// 	fmt.Println("   ✓ Memory freed successfully")
	// }

	// Example 11: Generate node compatibility report
	fmt.Println("\n11. Node Compatibility Analysis...")
	analyzeNodeCompatibility(objectInfo)

	fmt.Println("\n=== Model & Node Information Example Complete ===")
}

// analyzeNodeCompatibility analyzes which nodes can connect to each other
func analyzeNodeCompatibility(objectInfo comfyui.ObjectInfo) {
	// Find common connection patterns
	connectionPatterns := make(map[string]map[string]int) // output_type -> input_type -> count
	
	for _, info := range objectInfo {
		for _, outputType := range info.Output {
			if connectionPatterns[outputType] == nil {
				connectionPatterns[outputType] = make(map[string]int)
			}
			
			// Check what this output can connect to
			for inputName := range info.Input.Required {
				connectionPatterns[outputType][inputName]++
			}
		}
	}
	
	fmt.Println("   Common data flow patterns:")
	count := 0
	for outputType, inputs := range connectionPatterns {
		if count >= 5 {
			break
		}
		fmt.Printf("     %s can connect to:\n", outputType)
		inputCount := 0
		for inputName := range inputs {
			if inputCount >= 3 {
				break
			}
			fmt.Printf("       - %s\n", inputName)
			inputCount++
		}
		count++
	}
}

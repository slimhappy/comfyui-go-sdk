package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	comfyui "github.com/yourusername/comfyui-go-sdk"
)

func main() {
	// Create client
	client := comfyui.NewClient("http://127.0.0.1:8188")
	ctx := context.Background()

	fmt.Println("=== ComfyUI Image Operations Example ===\n")

	// Example 1: Create a test image
	fmt.Println("1. Creating test image...")
	testImagePath := "test_input.png"
	err := createTestImage(testImagePath, 512, 512)
	if err != nil {
		log.Fatalf("Failed to create test image: %v", err)
	}
	fmt.Printf("   ✓ Created test image: %s\n", testImagePath)

	// Example 2: Upload image with default options
	fmt.Println("\n2. Uploading image with default options...")
	uploadResult, err := client.UploadImage(ctx, testImagePath, comfyui.UploadOptions{})
	if err != nil {
		log.Fatalf("Failed to upload image: %v", err)
	}
	fmt.Printf("   ✓ Uploaded successfully\n")
	fmt.Printf("   Name: %s\n", uploadResult.Name)
	fmt.Printf("   Subfolder: %s\n", uploadResult.Subfolder)
	fmt.Printf("   Type: %s\n", uploadResult.Type)

	// Example 3: Upload image to specific subfolder
	fmt.Println("\n3. Uploading image to subfolder...")
	uploadResult2, err := client.UploadImage(ctx, testImagePath, comfyui.UploadOptions{
		Subfolder: "test_folder",
		Type:      "input",
		Overwrite: true,
	})
	if err != nil {
		log.Printf("Failed to upload to subfolder: %v", err)
	} else {
		fmt.Printf("   ✓ Uploaded to subfolder\n")
		fmt.Printf("   Name: %s\n", uploadResult2.Name)
		fmt.Printf("   Subfolder: %s\n", uploadResult2.Subfolder)
	}

	// Example 4: Upload image bytes directly
	fmt.Println("\n4. Uploading image from bytes...")
	imageData, err := os.ReadFile(testImagePath)
	if err != nil {
		log.Printf("Failed to read image: %v", err)
	} else {
		uploadResult3, err := client.UploadImageBytes(ctx, imageData, "from_bytes.png", comfyui.UploadOptions{
			Type: "input",
		})
		if err != nil {
			log.Printf("Failed to upload bytes: %v", err)
		} else {
			fmt.Printf("   ✓ Uploaded from bytes\n")
			fmt.Printf("   Name: %s\n", uploadResult3.Name)
		}
	}

	// Example 5: Use uploaded image in workflow
	fmt.Println("\n5. Using uploaded image in img2img workflow...")
	workflow := buildImg2ImgWorkflow(uploadResult.Name)
	queueResult, err := client.QueuePrompt(ctx, workflow, nil)
	if err != nil {
		log.Fatalf("Failed to queue workflow: %v", err)
	}
	fmt.Printf("   ✓ Queued workflow: %s\n", queueResult.PromptID)

	// Wait for completion
	fmt.Println("\n6. Waiting for workflow completion...")
	execResult, err := client.WaitForCompletion(ctx, queueResult.PromptID)
	if err != nil {
		log.Fatalf("Failed to wait for completion: %v", err)
	}
	fmt.Printf("   ✓ Workflow completed in %v\n", execResult.Duration)
	fmt.Printf("   Generated %d images\n", len(execResult.Images))

	// Example 7: Download generated images
	fmt.Println("\n7. Downloading generated images...")
	for i, img := range execResult.Images {
		outputPath := fmt.Sprintf("output_image_%d.png", i+1)
		err := client.SaveImage(ctx, img, outputPath)
		if err != nil {
			log.Printf("   Failed to save image %d: %v", i+1, err)
		} else {
			fmt.Printf("   ✓ Saved image %d: %s\n", i+1, outputPath)
			
			// Get file size
			if fileInfo, err := os.Stat(outputPath); err == nil {
				fmt.Printf("     Size: %.2f KB\n", float64(fileInfo.Size())/1024)
			}
		}
	}

	// Example 8: Download image directly (without SaveImage helper)
	fmt.Println("\n8. Downloading image using GetImage...")
	if len(execResult.Images) > 0 {
		img := execResult.Images[0]
		imageData, err := client.GetImage(ctx, img.Filename, img.Subfolder, img.Type)
		if err != nil {
			log.Printf("   Failed to get image: %v", err)
		} else {
			fmt.Printf("   ✓ Downloaded image data\n")
			fmt.Printf("   Size: %.2f KB\n", float64(len(imageData))/1024)
			
			// Save manually
			manualPath := "manual_download.png"
			err = os.WriteFile(manualPath, imageData, 0644)
			if err != nil {
				log.Printf("   Failed to write file: %v", err)
			} else {
				fmt.Printf("   ✓ Saved manually: %s\n", manualPath)
			}
		}
	}

	// Example 9: Batch image processing
	fmt.Println("\n9. Batch image processing...")
	batchImages := []string{}
	
	// Create multiple test images
	for i := 0; i < 3; i++ {
		imagePath := fmt.Sprintf("batch_input_%d.png", i)
		err := createTestImage(imagePath, 256, 256)
		if err != nil {
			log.Printf("   Failed to create batch image %d: %v", i, err)
			continue
		}
		batchImages = append(batchImages, imagePath)
	}
	
	fmt.Printf("   Created %d test images\n", len(batchImages))
	
	// Upload all images
	uploadedNames := []string{}
	for i, imagePath := range batchImages {
		result, err := client.UploadImage(ctx, imagePath, comfyui.UploadOptions{
			Type:      "input",
			Overwrite: true,
		})
		if err != nil {
			log.Printf("   Failed to upload batch image %d: %v", i, err)
			continue
		}
		uploadedNames = append(uploadedNames, result.Name)
		fmt.Printf("   ✓ Uploaded batch image %d: %s\n", i+1, result.Name)
	}

	// Example 10: Process each uploaded image
	fmt.Println("\n10. Processing batch images...")
	for i, imageName := range uploadedNames {
		workflow := buildImg2ImgWorkflow(imageName)
		result, err := client.QueuePrompt(ctx, workflow, map[string]interface{}{
			"batch_index": i,
		})
		if err != nil {
			log.Printf("   Failed to queue batch %d: %v", i, err)
			continue
		}
		fmt.Printf("   ✓ Queued batch %d: %s\n", i+1, result.PromptID)
	}

	// Example 11: Image format handling
	fmt.Println("\n11. Image format information...")
	fmt.Println("   Supported operations:")
	fmt.Println("     - Upload: PNG, JPG, JPEG, WEBP")
	fmt.Println("     - Download: Based on ComfyUI output format")
	fmt.Println("     - Default output: PNG")

	// Example 12: Cleanup
	fmt.Println("\n12. Cleanup...")
	filesToClean := []string{testImagePath}
	filesToClean = append(filesToClean, batchImages...)
	
	for _, file := range filesToClean {
		if err := os.Remove(file); err == nil {
			fmt.Printf("   ✓ Removed: %s\n", file)
		}
	}

	fmt.Println("\n=== Image Operations Example Complete ===")
}

// createTestImage creates a simple test image with gradient
func createTestImage(path string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Create a gradient pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create a colorful gradient
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(((x + y) * 255) / (width + height))
			
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	
	// Add some patterns
	// Horizontal lines
	for y := 0; y < height; y += 50 {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}
	
	// Vertical lines
	for x := 0; x < width; x += 50 {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}
	
	// Save to file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	return png.Encode(file, img)
}

// buildImg2ImgWorkflow creates an image-to-image workflow
func buildImg2ImgWorkflow(imageName string) comfyui.Workflow {
	return comfyui.Workflow{
		"1": comfyui.Node{
			ClassType: "LoadImage",
			Inputs: map[string]interface{}{
				"image": imageName,
			},
		},
		"2": comfyui.Node{
			ClassType: "CheckpointLoaderSimple",
			Inputs: map[string]interface{}{
				"ckpt_name": "v1-5-pruned-emaonly.safetensors",
			},
		},
		"3": comfyui.Node{
			ClassType: "VAEEncode",
			Inputs: map[string]interface{}{
				"pixels": []interface{}{"1", 0},
				"vae":    []interface{}{"2", 2},
			},
		},
		"4": comfyui.Node{
			ClassType: "CLIPTextEncode",
			Inputs: map[string]interface{}{
				"text": "high quality, detailed, enhanced",
				"clip": []interface{}{"2", 1},
			},
		},
		"5": comfyui.Node{
			ClassType: "CLIPTextEncode",
			Inputs: map[string]interface{}{
				"text": "low quality, blurry",
				"clip": []interface{}{"2", 1},
			},
		},
		"6": comfyui.Node{
			ClassType: "KSampler",
			Inputs: map[string]interface{}{
				"seed":         54321,
				"steps":        15,
				"cfg":          7.0,
				"sampler_name": "euler",
				"scheduler":    "normal",
				"denoise":      0.6,
				"model":        []interface{}{"2", 0},
				"positive":     []interface{}{"4", 0},
				"negative":     []interface{}{"5", 0},
				"latent_image": []interface{}{"3", 0},
			},
		},
		"7": comfyui.Node{
			ClassType: "VAEDecode",
			Inputs: map[string]interface{}{
				"samples": []interface{}{"6", 0},
				"vae":     []interface{}{"2", 2},
			},
		},
		"8": comfyui.Node{
			ClassType: "SaveImage",
			Inputs: map[string]interface{}{
				"filename_prefix": "img2img_test",
				"images":          []interface{}{"7", 0},
			},
		},
	}
}

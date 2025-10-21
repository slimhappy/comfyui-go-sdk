package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	comfyui "github.com/yourusername/comfyui-go-sdk"
)

func main() {
	// Create client
	client := comfyui.NewClient("http://127.0.0.1:8188")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		cancel()
	}()

	// Connect WebSocket
	fmt.Println("Connecting to ComfyUI WebSocket...")
	ws, err := client.ConnectWebSocket(ctx)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer ws.Close()

	fmt.Println("Connected! Listening for events...")
	fmt.Println("Press Ctrl+C to exit")

	// Listen for all messages
	for {
		select {
		case <-ctx.Done():
			return

		case msg, ok := <-ws.Messages():
			if !ok {
				fmt.Println("WebSocket closed")
				return
			}

			handleMessage(msg)

		case err := <-ws.Errors():
			log.Printf("WebSocket error: %v", err)
			return
		}
	}
}

func handleMessage(msg comfyui.WebSocketMessage) {
	switch msg.Type {
	case string(comfyui.MessageTypeStatus):
		data, err := extractStatusData(msg.Data)
		if err == nil {
			fmt.Printf("[STATUS] Queue remaining: %d\n", data.QueueRemaining)
		}

	case string(comfyui.MessageTypeExecuting):
		data, err := msg.GetExecutingData()
		if err == nil {
			if data.Node == nil {
				fmt.Printf("[EXECUTING] Prompt %s completed\n", data.PromptID)
			} else {
				fmt.Printf("[EXECUTING] Prompt %s, Node %s\n", data.PromptID, *data.Node)
			}
		}

	case string(comfyui.MessageTypeProgress):
		data, err := msg.GetProgressData()
		if err == nil {
			percentage := float64(data.Value) / float64(data.Max) * 100
			fmt.Printf("[PROGRESS] %d/%d (%.1f%%)\n", data.Value, data.Max, percentage)
		}

	case string(comfyui.MessageTypeExecuted):
		data, err := msg.GetExecutedData()
		if err == nil {
			fmt.Printf("[EXECUTED] Node %s in prompt %s\n", data.Node, data.PromptID)
			if images, ok := data.Output["images"].([]interface{}); ok {
				fmt.Printf("  → Produced %d image(s)\n", len(images))
			}
		}

	case string(comfyui.MessageTypeCached):
		fmt.Printf("[CACHED] Node execution cached\n")

	case string(comfyui.MessageTypeError):
		data, err := msg.GetErrorData()
		if err == nil {
			fmt.Printf("[ERROR] Prompt %s, Node %s (%s)\n", data.PromptID, data.NodeID, data.NodeType)
			fmt.Printf("  Type: %s\n", data.ExceptionType)
			fmt.Printf("  Message: %s\n", data.ExceptionMessage)
			if len(data.Traceback) > 0 {
				fmt.Println("  Traceback:")
				for _, line := range data.Traceback {
					fmt.Printf("    %s\n", line)
				}
			}
		}

	default:
		fmt.Printf("[%s] %v\n", msg.Type, msg.Data)
	}
}

func extractStatusData(data map[string]interface{}) (*struct{ QueueRemaining int }, error) {
	result := &struct{ QueueRemaining int }{}

	if status, ok := data["status"].(map[string]interface{}); ok {
		if execInfo, ok := status["exec_info"].(map[string]interface{}); ok {
			if qr, ok := execInfo["queue_remaining"].(float64); ok {
				result.QueueRemaining = int(qr)
				return result, nil
			}
		}
	}

	return nil, fmt.Errorf("invalid status data")
}

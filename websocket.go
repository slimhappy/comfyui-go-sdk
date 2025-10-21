package comfyui

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketClient represents a WebSocket connection to ComfyUI
type WebSocketClient struct {
	conn     *websocket.Conn
	messages chan WebSocketMessage
	errors   chan error
	done     chan struct{}
	once     sync.Once
	clientID string
}

// ConnectWebSocket establishes a WebSocket connection
func (c *Client) ConnectWebSocket(ctx context.Context) (*WebSocketClient, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	// Convert http/https to ws/wss
	scheme := "ws"
	if u.Scheme == "https" {
		scheme = "wss"
	}

	wsURL := fmt.Sprintf("%s://%s/ws?clientId=%s", scheme, u.Host, c.clientID)

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect websocket: %w", err)
	}

	ws := &WebSocketClient{
		conn:     conn,
		messages: make(chan WebSocketMessage, 100),
		errors:   make(chan error, 10),
		done:     make(chan struct{}),
		clientID: c.clientID,
	}

	go ws.readLoop()

	return ws, nil
}

// readLoop reads messages from the WebSocket
func (ws *WebSocketClient) readLoop() {
	defer close(ws.messages)
	defer close(ws.errors)

	for {
		select {
		case <-ws.done:
			return
		default:
			_, message, err := ws.conn.ReadMessage()
			if err != nil {
				select {
				case ws.errors <- fmt.Errorf("read error: %w", err):
				case <-ws.done:
				}
				return
			}

			var msg WebSocketMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				select {
				case ws.errors <- fmt.Errorf("unmarshal error: %w", err):
				case <-ws.done:
				}
				continue
			}

			select {
			case ws.messages <- msg:
			case <-ws.done:
				return
			}
		}
	}
}

// Messages returns the channel for receiving messages
func (ws *WebSocketClient) Messages() <-chan WebSocketMessage {
	return ws.messages
}

// Errors returns the channel for receiving errors
func (ws *WebSocketClient) Errors() <-chan error {
	return ws.errors
}

// Close closes the WebSocket connection
func (ws *WebSocketClient) Close() error {
	var err error
	ws.once.Do(func() {
		close(ws.done)
		err = ws.conn.Close()
	})
	return err
}

// SendMessage sends a message through the WebSocket
func (ws *WebSocketClient) SendMessage(msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := ws.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// WaitForPromptCompletion waits for a specific prompt to complete
func (ws *WebSocketClient) WaitForPromptCompletion(ctx context.Context, promptID string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-ws.errors:
			return err
		case msg, ok := <-ws.messages:
			if !ok {
				return fmt.Errorf("websocket closed")
			}

			if msg.Type == string(MessageTypeExecuting) {
				if pid, ok := msg.Data["prompt_id"].(string); ok && pid == promptID {
					if node, ok := msg.Data["node"].(string); !ok || node == "" {
						return nil // Execution completed
					}
				}
			}

			if msg.Type == string(MessageTypeError) {
				if pid, ok := msg.Data["prompt_id"].(string); ok && pid == promptID {
					return fmt.Errorf("execution error: %v", msg.Data)
				}
			}
		}
	}
}

// GetExecutingData extracts executing data from a message
func (msg *WebSocketMessage) GetExecutingData() (*ExecutingData, error) {
	if msg.Type != string(MessageTypeExecuting) {
		return nil, fmt.Errorf("not an executing message")
	}

	data := &ExecutingData{}
	if pid, ok := msg.Data["prompt_id"].(string); ok {
		data.PromptID = pid
	}
	if node, ok := msg.Data["node"].(string); ok {
		data.Node = &node
	}

	return data, nil
}

// GetProgressData extracts progress data from a message
func (msg *WebSocketMessage) GetProgressData() (*ProgressData, error) {
	if msg.Type != string(MessageTypeProgress) {
		return nil, fmt.Errorf("not a progress message")
	}

	data := &ProgressData{}
	if val, ok := msg.Data["value"].(float64); ok {
		data.Value = int(val)
	}
	if max, ok := msg.Data["max"].(float64); ok {
		data.Max = int(max)
	}

	return data, nil
}

// GetExecutedData extracts executed data from a message
func (msg *WebSocketMessage) GetExecutedData() (*ExecutedData, error) {
	if msg.Type != string(MessageTypeExecuted) {
		return nil, fmt.Errorf("not an executed message")
	}

	data := &ExecutedData{}
	if node, ok := msg.Data["node"].(string); ok {
		data.Node = node
	}
	if pid, ok := msg.Data["prompt_id"].(string); ok {
		data.PromptID = pid
	}
	if output, ok := msg.Data["output"].(map[string]interface{}); ok {
		data.Output = output
	}

	return data, nil
}

// GetErrorData extracts error data from a message
func (msg *WebSocketMessage) GetErrorData() (*ErrorData, error) {
	if msg.Type != string(MessageTypeError) {
		return nil, fmt.Errorf("not an error message")
	}

	data := &ErrorData{}
	if pid, ok := msg.Data["prompt_id"].(string); ok {
		data.PromptID = pid
	}
	if nodeID, ok := msg.Data["node_id"].(string); ok {
		data.NodeID = nodeID
	}
	if nodeType, ok := msg.Data["node_type"].(string); ok {
		data.NodeType = nodeType
	}
	if excType, ok := msg.Data["exception_type"].(string); ok {
		data.ExceptionType = excType
	}
	if excMsg, ok := msg.Data["exception_message"].(string); ok {
		data.ExceptionMessage = excMsg
	}
	if tb, ok := msg.Data["traceback"].([]interface{}); ok {
		for _, line := range tb {
			if str, ok := line.(string); ok {
				data.Traceback = append(data.Traceback, str)
			}
		}
	}

	return data, nil
}

package comfyui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Client represents a ComfyUI API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	clientID   string
}

// NewClient creates a new ComfyUI client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		clientID: uuid.New().String(),
	}
}

// NewClientWithHTTPClient creates a new client with custom HTTP client
func NewClientWithHTTPClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: httpClient,
		clientID:   uuid.New().String(),
	}
}

// SetClientID sets the client ID
func (c *Client) SetClientID(clientID string) {
	c.clientID = clientID
}

// GetClientID returns the client ID
func (c *Client) GetClientID() string {
	return c.clientID
}

// QueuePrompt queues a workflow for execution
func (c *Client) QueuePrompt(ctx context.Context, workflow Workflow, extraData map[string]interface{}) (*QueuePromptResponse, error) {
	req := QueuePromptRequest{
		Prompt:    workflow,
		ClientID:  c.clientID,
		ExtraData: extraData,
	}

	var resp QueuePromptResponse
	if err := c.doRequest(ctx, "POST", "/prompt", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to queue prompt: %w", err)
	}

	if len(resp.NodeErrors) > 0 {
		return &resp, fmt.Errorf("node errors: %v", resp.NodeErrors)
	}

	return &resp, nil
}

// QueuePromptFromFile loads a workflow from a JSON file and queues it for execution
func (c *Client) QueuePromptFromFile(ctx context.Context, filepath string, extraData map[string]interface{}) (*QueuePromptResponse, error) {
	workflow, err := LoadWorkflowFromFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to load workflow from file: %w", err)
	}

	return c.QueuePrompt(ctx, workflow, extraData)
}

// GetQueue retrieves the current queue status
func (c *Client) GetQueue(ctx context.Context) (*QueueStatus, error) {
	var queue QueueStatus
	if err := c.doRequest(ctx, "GET", "/queue", nil, &queue); err != nil {
		return nil, fmt.Errorf("failed to get queue: %w", err)
	}
	return &queue, nil
}

// ClearQueue clears all items from the queue
func (c *Client) ClearQueue(ctx context.Context) error {
	req := QueueManagementRequest{Clear: true}
	return c.doRequest(ctx, "POST", "/queue", req, nil)
}

// DeleteFromQueue deletes specific items from the queue
func (c *Client) DeleteFromQueue(ctx context.Context, promptIDs []string) error {
	req := QueueManagementRequest{Delete: promptIDs}
	return c.doRequest(ctx, "POST", "/queue", req, nil)
}

// Interrupt interrupts the current execution
func (c *Client) Interrupt(ctx context.Context, promptID string) error {
	req := InterruptRequest{PromptID: promptID}
	return c.doRequest(ctx, "POST", "/interrupt", req, nil)
}

// GetHistory retrieves execution history
// If promptID is empty, returns all history
func (c *Client) GetHistory(ctx context.Context, promptID string) (History, error) {
	path := "/history"
	if promptID != "" {
		path = fmt.Sprintf("/history/%s", promptID)
	}

	var history History
	if err := c.doRequest(ctx, "GET", path, nil, &history); err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}
	return history, nil
}

// ClearHistory clears all history
func (c *Client) ClearHistory(ctx context.Context) error {
	req := HistoryManagementRequest{Clear: true}
	return c.doRequest(ctx, "POST", "/history", req, nil)
}

// DeleteHistory deletes specific history items
func (c *Client) DeleteHistory(ctx context.Context, promptIDs []string) error {
	req := HistoryManagementRequest{Delete: promptIDs}
	return c.doRequest(ctx, "POST", "/history", req, nil)
}

// GetSystemStats retrieves system statistics
func (c *Client) GetSystemStats(ctx context.Context) (*SystemStats, error) {
	var stats SystemStats
	if err := c.doRequest(ctx, "GET", "/system_stats", nil, &stats); err != nil {
		return nil, fmt.Errorf("failed to get system stats: %w", err)
	}
	return &stats, nil
}

// GetObjectInfo retrieves node class information
// If nodeClass is empty, returns all node classes
func (c *Client) GetObjectInfo(ctx context.Context, nodeClass string) (ObjectInfo, error) {
	path := "/object_info"
	if nodeClass != "" {
		path = fmt.Sprintf("/object_info/%s", nodeClass)
	}

	var info ObjectInfo
	if err := c.doRequest(ctx, "GET", path, nil, &info); err != nil {
		return nil, fmt.Errorf("failed to get object info: %w", err)
	}
	return info, nil
}

// GetEmbeddings retrieves the list of available embeddings
func (c *Client) GetEmbeddings(ctx context.Context) ([]string, error) {
	var embeddings []string
	if err := c.doRequest(ctx, "GET", "/embeddings", nil, &embeddings); err != nil {
		return nil, fmt.Errorf("failed to get embeddings: %w", err)
	}
	return embeddings, nil
}

// GetModels retrieves the list of models
// If folder is empty, returns list of model folders
func (c *Client) GetModels(ctx context.Context, folder string) ([]string, error) {
	path := "/models"
	if folder != "" {
		path = fmt.Sprintf("/models/%s", folder)
	}

	var models []string
	if err := c.doRequest(ctx, "GET", path, nil, &models); err != nil {
		return nil, fmt.Errorf("failed to get models: %w", err)
	}
	return models, nil
}

// FreeMemory requests the server to free memory
func (c *Client) FreeMemory(ctx context.Context, unloadModels, freeMemory bool) error {
	req := FreeMemoryRequest{
		UnloadModels: unloadModels,
		FreeMemory:   freeMemory,
	}
	return c.doRequest(ctx, "POST", "/free", req, nil)
}

// GetFeatures retrieves server features
func (c *Client) GetFeatures(ctx context.Context) (*Features, error) {
	var features Features
	if err := c.doRequest(ctx, "GET", "/features", nil, &features); err != nil {
		return nil, fmt.Errorf("failed to get features: %w", err)
	}
	return &features, nil
}

// UploadImage uploads an image file
func (c *Client) UploadImage(ctx context.Context, filepath string, opts UploadOptions) (*UploadImageResponse, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	filename := filepath[strings.LastIndex(filepath, "/")+1:]
	return c.UploadImageBytes(ctx, data, filename, opts)
}

// UploadImageBytes uploads an image from bytes
func (c *Client) UploadImageBytes(ctx context.Context, data []byte, filename string, opts UploadOptions) (*UploadImageResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add image file
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write file data: %w", err)
	}

	// Add options
	if opts.Subfolder != "" {
		writer.WriteField("subfolder", opts.Subfolder)
	}
	if opts.Type != "" {
		writer.WriteField("type", opts.Type)
	} else {
		writer.WriteField("type", "input")
	}
	if opts.Overwrite {
		writer.WriteField("overwrite", "true")
	}

	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/upload/image", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var uploadResp UploadImageResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &uploadResp, nil
}

// GetImage downloads an image
func (c *Client) GetImage(ctx context.Context, filename, subfolder, folderType string) ([]byte, error) {
	params := url.Values{}
	params.Add("filename", filename)
	params.Add("subfolder", subfolder)
	params.Add("type", folderType)

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/view?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get image: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return data, nil
}

// SaveImage saves an image to a file
func (c *Client) SaveImage(ctx context.Context, img ImageInfo, outputPath string) error {
	data, err := c.GetImage(ctx, img.Filename, img.Subfolder, img.Type)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// WaitForCompletion waits for a workflow to complete and returns the results
func (c *Client) WaitForCompletion(ctx context.Context, promptID string) (*ExecutionResult, error) {
	ws, err := c.ConnectWebSocket(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect websocket: %w", err)
	}
	defer ws.Close()

	result := &ExecutionResult{
		PromptID:  promptID,
		StartTime: time.Now(),
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case msg, ok := <-ws.Messages():
			if !ok {
				return nil, fmt.Errorf("websocket closed unexpectedly")
			}

			if msg.Type == string(MessageTypeExecuting) {
				data := msg.Data
				if pid, ok := data["prompt_id"].(string); ok && pid == promptID {
					if node, ok := data["node"].(string); !ok || node == "" {

						// Execution completed
						result.EndTime = time.Now()
						result.Duration = result.EndTime.Sub(result.StartTime)

						// Get history to retrieve outputs
						history, err := c.GetHistory(ctx, promptID)
						if err != nil {
							return nil, fmt.Errorf("failed to get history: %w", err)
						}

						if item, ok := history[promptID]; ok {
							result.Outputs = item.Outputs
							result.Status = item.Status

							// Collect all images
							for _, output := range item.Outputs {
								result.Images = append(result.Images, output.Images...)
							}
						}

						return result, nil
					}
				}
			}
		}
	}
}

// doRequest performs an HTTP request
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

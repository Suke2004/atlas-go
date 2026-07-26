package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OllamaProvider implements Provider using the Ollama local HTTP API.
// Compatible with Ollama v0.1+ /api/chat endpoint.
type OllamaProvider struct {
	baseURL string
	model   string
	client  *http.Client
}

// NewOllamaProvider creates an Ollama provider.
// baseURL is typically "http://localhost:11434".
func NewOllamaProvider(baseURL, model string) *OllamaProvider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3.2"
	}
	return &OllamaProvider{
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		client:  &http.Client{Timeout: 120 * time.Second},
	}
}

// Name returns the provider identifier.
func (o *OllamaProvider) Name() string { return "ollama:" + o.model }

// IsAvailable pings Ollama to check if it is running.
func (o *OllamaProvider) IsAvailable(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, o.baseURL+"/api/tags", nil)
	if err != nil {
		return false
	}
	resp, err := o.client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// ollama REST types -------------------------------------------------------

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type ollamaChatResponse struct {
	Message ollamaMessage `json:"message"`
	Error   string        `json:"error"`
	Done    bool          `json:"done"`
}

// Complete sends messages to Ollama and returns the full assistant reply.
func (o *OllamaProvider) Complete(ctx context.Context, messages []Message) (string, error) {
	ollamaMsgs := make([]ollamaMessage, 0, len(messages))
	for _, m := range messages {
		role := string(m.Role)
		// Ollama uses "assistant" not "assistant"
		if m.Role == RoleAssistant {
			role = "assistant"
		}
		ollamaMsgs = append(ollamaMsgs, ollamaMessage{
			Role:    role,
			Content: m.Content,
		})
	}

	reqBody := ollamaChatRequest{
		Model:    o.model,
		Messages: ollamaMsgs,
		Stream:   false,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("ollama: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		o.baseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("ollama: build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("ollama: HTTP request failed — is Ollama running? Error: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ollama: read response: %w", err)
	}

	var chatResp ollamaChatResponse
	if err := json.Unmarshal(raw, &chatResp); err != nil {
		return "", fmt.Errorf("ollama: parse response: %w", err)
	}

	if chatResp.Error != "" {
		return "", fmt.Errorf("ollama error: %s", chatResp.Error)
	}

	return chatResp.Message.Content, nil
}

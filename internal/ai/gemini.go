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

// GeminiProvider implements Provider using the Google Gemini REST API.
// It uses the generateContent endpoint directly over HTTPS — no SDK dependency.
type GeminiProvider struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// NewGeminiProvider creates a Gemini provider. model should be e.g. "gemini-2.0-flash".
func NewGeminiProvider(apiKey, model string) *GeminiProvider {
	if model == "" {
		model = "gemini-2.0-flash"
	}
	return &GeminiProvider{
		apiKey:  apiKey,
		model:   model,
		baseURL: "https://generativelanguage.googleapis.com/v1beta",
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

// Name returns the provider identifier.
func (g *GeminiProvider) Name() string { return "gemini:" + g.model }

// IsAvailable checks whether the Gemini API key is set and the endpoint reachable.
func (g *GeminiProvider) IsAvailable(ctx context.Context) bool {
	return g.apiKey != ""
}

// gemini REST request / response types -----------------------------------

type geminiRequest struct {
	SystemInstruction *geminiContent  `json:"systemInstruction,omitempty"`
	Contents          []geminiContent `json:"contents"`
	GenerationConfig  geminiGenConfig `json:"generationConfig"`
}

type geminiContent struct {
	Role  string        `json:"role,omitempty"`
	Parts []geminiPart  `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenConfig struct {
	Temperature     float64 `json:"temperature"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
}

type geminiResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}

// Complete sends messages to Gemini and returns the assistant response text.
func (g *GeminiProvider) Complete(ctx context.Context, messages []Message) (string, error) {
	if g.apiKey == "" {
		return "", fmt.Errorf("gemini: API key not configured — set it in Settings → AI Provider")
	}

	req := geminiRequest{
		GenerationConfig: geminiGenConfig{
			Temperature:     0.7,
			MaxOutputTokens: 2048,
		},
	}

	// Separate system message from the conversation
	for _, m := range messages {
		switch m.Role {
		case RoleSystem:
			req.SystemInstruction = &geminiContent{
				Parts: []geminiPart{{Text: m.Content}},
			}
		case RoleUser:
			req.Contents = append(req.Contents, geminiContent{
				Role:  "user",
				Parts: []geminiPart{{Text: m.Content}},
			})
		case RoleAssistant:
			req.Contents = append(req.Contents, geminiContent{
				Role:  "model",
				Parts: []geminiPart{{Text: m.Content}},
			})
		}
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("gemini: marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", g.baseURL, g.model, g.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("gemini: build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("gemini: HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gemini: read response body: %w", err)
	}

	var gemResp geminiResponse
	if err := json.Unmarshal(raw, &gemResp); err != nil {
		return "", fmt.Errorf("gemini: parse response: %w", err)
	}

	if gemResp.Error != nil {
		return "", fmt.Errorf("gemini API error %d: %s", gemResp.Error.Code, gemResp.Error.Message)
	}

	if len(gemResp.Candidates) == 0 || len(gemResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini: empty response from model")
	}

	var sb strings.Builder
	for _, part := range gemResp.Candidates[0].Content.Parts {
		sb.WriteString(part.Text)
	}
	return sb.String(), nil
}

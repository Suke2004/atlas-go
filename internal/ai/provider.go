package ai

import (
	"context"
	"fmt"
	"strings"
)

// Role represents the speaker in a conversation.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message is one turn in a conversation.
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

// Provider is the interface every AI backend must implement.
// All implementations must be safe for concurrent use.
type Provider interface {
	// Complete returns the full assistant response for a conversation.
	Complete(ctx context.Context, messages []Message) (string, error)
	// IsAvailable checks whether the provider is reachable right now.
	IsAvailable(ctx context.Context) bool
	// Name returns a human-readable provider identifier.
	Name() string
}

// ProviderType identifies which backend to use.
type ProviderType string

const (
	ProviderOllama ProviderType = "ollama"
	ProviderGemini ProviderType = "gemini"
	ProviderOpenAI ProviderType = "openai"
)

// NewFromSettings constructs the appropriate Provider from the settings map
// returned by settings.Service.GetAll(). Returns nil if the provider type is
// unsupported — callers must handle a nil provider gracefully.
func NewFromSettings(s map[string]string) Provider {
	switch ProviderType(s["ai_provider"]) {
	case ProviderGemini:
		return NewGeminiProvider(s["ai_gemini_key"], s["ai_gemini_model"])
	case ProviderOllama:
		return NewOllamaProvider(s["ai_ollama_url"], s["ai_ollama_model"])
	default:
		// Default to Ollama (local-first)
		return NewOllamaProvider(s["ai_ollama_url"], s["ai_ollama_model"])
	}
}

// ChatRequest is the canonical input to the AI workspace.
type ChatRequest struct {
	Mode    string `json:"mode"`    // "chat" | "ask_atlas" | "explain" | "summarise"
	Message string `json:"message"` // user message
	Context string `json:"context"` // optional injected workspace context
}

// AtlasContext holds live workspace data injected into Ask Atlas prompts.
type AtlasContext struct {
	ActiveProjects []string
	TodayTasks     []string
	RecentNotes    []string
	LastJournal    string
}

// BuildAtlasSystemPrompt constructs the system message for Ask Atlas mode.
func BuildAtlasSystemPrompt(ac AtlasContext) string {
	var sb strings.Builder
	sb.WriteString("You are Atlas AI, an intelligent assistant embedded inside the user's personal workspace.\n")
	sb.WriteString("You have access to the following live context from the user's workspace:\n\n")

	if len(ac.ActiveProjects) > 0 {
		sb.WriteString(fmt.Sprintf("ACTIVE PROJECTS (%d):\n", len(ac.ActiveProjects)))
		for _, p := range ac.ActiveProjects {
			sb.WriteString("  • " + p + "\n")
		}
		sb.WriteString("\n")
	}

	if len(ac.TodayTasks) > 0 {
		sb.WriteString(fmt.Sprintf("TODAY'S TASKS (%d):\n", len(ac.TodayTasks)))
		for _, t := range ac.TodayTasks {
			sb.WriteString("  • " + t + "\n")
		}
		sb.WriteString("\n")
	}

	if len(ac.RecentNotes) > 0 {
		sb.WriteString(fmt.Sprintf("RECENT NOTES (%d):\n", len(ac.RecentNotes)))
		for _, n := range ac.RecentNotes {
			sb.WriteString("  • " + n + "\n")
		}
		sb.WriteString("\n")
	}

	if ac.LastJournal != "" {
		sb.WriteString("LATEST JOURNAL SNIPPET:\n")
		sb.WriteString(ac.LastJournal + "\n\n")
	}

	sb.WriteString("Answer the user's question using the above context wherever relevant. ")
	sb.WriteString("Be concise, actionable, and grounded in their actual data. ")
	sb.WriteString("If asked about something not in the context, say so honestly.")
	return sb.String()
}

// SystemPromptFor returns the system prompt for a given mode.
func SystemPromptFor(mode string) string {
	switch mode {
	case "explain":
		return "You are an expert code reviewer and teacher. When given code, provide a clear, structured explanation: what it does, how it works, potential issues, and suggestions for improvement. Use markdown formatting with headers and code blocks."
	case "summarise":
		return "You are a precise summarisation assistant. Condense the provided text into a clear, bullet-pointed summary preserving all key information. Use markdown formatting."
	default:
		return "You are Atlas AI, a helpful, concise assistant for a developer's personal workspace. Respond in markdown. Be direct and practical."
	}
}

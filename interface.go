package yacloud_gpt

import "context"

type YandexGpt interface {
	Completion(req CompletionRequest) (CompletionResponse, error)
	CompletionWithContext(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
}

type YandexGptTokenizer interface {
	Tokenize()
	TokenizeCompletion()
}

type MessageRole string

const (
	ModeSystem    MessageRole = "system"
	ModeAssistant MessageRole = "assistant"
	ModeUser      MessageRole = "user"
)

type CompletionRequest struct {
	ModelUri          string              `json:"modelUri,omitempty"`
	CompletionOptions *CompletionOptions  `json:"completionOptions,omitempty"`
	Messages          []CompletionMessage `json:"messages"`
}

type CompletionOptions struct {
	Stream      bool    `json:"stream,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"maxTokens,omitempty"`
}

type CompletionMessage struct {
	Role MessageRole `json:"role"`
	Text string      `json:"text"`
}

type CompletionResponse struct {
	Alternatives []CompetionAlternative `json:"alternatives"`
}

type CompetionAlternative struct {
	Message CompletionMessage `json:"message"`
	Status  string            `json:"status"`
}

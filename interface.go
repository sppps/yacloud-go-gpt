package yacloud_gpt

type YandexGpt interface {
	Completion(req CompletionRequest) (CompletionResponse, error)
}

type YandexGptTokenizer interface {
	Tokenize()
	TokenizeCompletion()
}

type ModelUri string

const (
	YandexGptPro           ModelUri = "yandexgpt"
	YandexGptLite          ModelUri = "yandexgpt-lite"
	YandexGptSummarization ModelUri = "summarization"
)

type MessageRole string

const (
	ModeSystem    MessageRole = "system"
	ModeAssistant MessageRole = "assistant"
	ModeUser      MessageRole = "user"
)

type CompletionRequest struct {
	ModelUri          ModelUri            `json:"modelUri,omitempty"`
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

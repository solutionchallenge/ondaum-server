package llm

type Statistics struct {
	TotalTokens      int64
	PromptTokens     int64
	CompletionTokens int64
	CachedTokens     int64
}

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	ConversationID string
	ID             string
	Role           Role
	Content        string
	Metadata       map[string]any
}

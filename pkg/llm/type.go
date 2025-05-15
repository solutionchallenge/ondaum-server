package llm

type Statistics struct {
	TotalTokens      int64
	PromptTokens     int64
	CompletionTokens int64
	ThoughtsTokens   int64
	CachedTokens     int64
}

type Role string

const (
	RoleUser  Role = "user"
	RoleModel Role = "model"
)

type Message struct {
	ConversationID string
	ID             string
	Role           Role
	Content        string
	Metadata       map[string]any
}

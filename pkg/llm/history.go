package llm

type HistoryRole string

const (
	HistoryRoleUser      HistoryRole = "user"
	HistoryRoleSystem    HistoryRole = "system"
	HistoryRoleAssistant HistoryRole = "assistant"
)

type HistoryManager interface {
	Add(messages ...Message)
	Get(conversationID string) []Message
}

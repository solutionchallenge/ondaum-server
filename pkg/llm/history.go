package llm

import "context"

type HistoryRole string

const (
	HistoryRoleUser      HistoryRole = "user"
	HistoryRoleSystem    HistoryRole = "system"
	HistoryRoleAssistant HistoryRole = "assistant"
)

type HistoryManager interface {
	Add(ctx context.Context, messages ...Message)
	Get(ctx context.Context, conversationID string) []Message
}

package llm

import "context"

type HistoryManager interface {
	Add(ctx context.Context, messages ...Message)
	Get(ctx context.Context, conversationID string) []Message
}

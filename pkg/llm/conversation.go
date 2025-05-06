package llm

import "context"

type Conversation interface {
	Request(ctx context.Context, request Message) (Message, error)
	GetHistory(ctx context.Context) []Message
	GetStatistics() Statistics
	End()
}

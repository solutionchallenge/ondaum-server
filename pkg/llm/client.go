package llm

import "context"

type Client interface {
	StartConversation(ctx context.Context, historyManager HistoryManager, instructionIdentifier string, id ...string) (Conversation, error)
	RunActionPrompt(ctx context.Context, instructionIdentifier string, promptIdentifier string, histories ...Message) (Message, error)
	GetStatistics() Statistics
	Close() error
}

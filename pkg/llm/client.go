package llm

type Client interface {
	StartConversation(manager HistoryManager, id ...string) Conversation
	GetStatistics() Statistics
	Close() error
}

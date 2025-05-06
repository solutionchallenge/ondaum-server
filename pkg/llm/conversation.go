package llm

type Conversation interface {
	Request(request Message) (Message, error)
	GetStatistics() Statistics
	End()
}

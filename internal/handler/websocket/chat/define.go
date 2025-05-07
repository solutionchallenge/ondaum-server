package chat

import "time"

const (
	ChatAutoFinishAfter = 30 * time.Minute
)

const (
	ChatPayloadNotifyConversationFinished = "conversation_finished"
	ChatPayloadNotifyConversationArchived = "conversation_archived"
	ChatPayloadNotifyNewConversation      = "new_conversation"
	ChatPayloadNotifyExistingConversation = "existing_conversation"
)

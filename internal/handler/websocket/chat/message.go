package chat

import (
	"context"
	"encoding/json"

	"github.com/benbjohnson/clock"
	domain "github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	ftpkg "github.com/solutionchallenge/ondaum-server/pkg/future"
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
	llmpkg "github.com/solutionchallenge/ondaum-server/pkg/llm"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	wspkg "github.com/solutionchallenge/ondaum-server/pkg/websocket"
	"github.com/uptrace/bun"
)

func HandleMessage(
	db *bun.DB, clk clock.Clock, llm llm.Client, future *ftpkg.Scheduler, request wspkg.MessageWrapper,
) (wspkg.ResponseWrapper, bool, error) {
	if !request.Authorized || !checkAuthorization(db, request.UserID) {
		utils.Log(utils.InfoLevel).CID(request.SessionID).RID(request.MessageID).BT().Send("Unauthorized")
		return wspkg.BuildRejectResponse(request), false, nil
	}

	var response wspkg.ResponseWrapper
	var shouldClose bool
	err := db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		chat := &domain.Chat{}
		err := tx.NewSelect().
			Model(chat).
			Where("session_id = ?", request.SessionID).
			Where("user_id = ?", request.UserID).
			Scan(ctx)
		if err != nil {
			utils.Log(utils.ErrorLevel).CID(request.SessionID).RID(request.MessageID).Err(err).BT().Send("Failed to query chat")
			return utils.WrapError(err, "failed to query chat")
		}

		if !chat.ArchivedAt.IsZero() {
			utils.Log(utils.InfoLevel).CID(request.SessionID).RID(request.MessageID).BT().Send("Cannot reactivate archived chat")
			response = wspkg.BuildRejectResponse(request)
			return nil
		}

		if !chat.FinishedAt.IsZero() {
			_, err = tx.NewUpdate().
				Model(chat).
				Set("finished_at = NULL").
				Where("id = ?", chat.ID).
				Exec(ctx)
			if err != nil {
				utils.Log(utils.ErrorLevel).CID(request.SessionID).RID(request.MessageID).Err(err).BT().Send("Failed to reactivate chat")
				return utils.WrapError(err, "failed to reactivate chat")
			}
		}
		return nil
	})
	if err != nil {
		utils.Log(utils.ErrorLevel).CID(request.SessionID).RID(request.MessageID).Err(err).BT().Send("Failed to run transaction")
		return wspkg.ResponseWrapper{}, false, utils.WrapError(err, "failed to run transaction")
	}

	_, err = UpsertChattingEndFutureJob(context.Background(), future, request.SessionID, request.UserID, ChatAutoFinishAfter)
	if err != nil {
		utils.Log(utils.ErrorLevel).CID(request.SessionID).RID(request.MessageID).Err(err).BT().Send("Failed to upsert future job")
		return wspkg.ResponseWrapper{}, false, utils.WrapError(err, "failed to upsert future job")
	}

	manager := NewChatHistoryManager(db, request.SessionID)
	conversation, err := llm.StartConversation(context.Background(), manager, "interactive_chat", request.SessionID)
	if err != nil {
		utils.Log(utils.ErrorLevel).CID(request.SessionID).RID(request.MessageID).Err(err).BT().Send("Failed to start conversation")
		return wspkg.ResponseWrapper{}, false, utils.WrapError(err, "failed to start conversation")
	}
	// TODO: apply user addition to the conversation
	llmResponse, err := conversation.Request(context.Background(), llmpkg.Message{
		ConversationID: request.SessionID,
		ID:             request.MessageID,
		Role:           llmpkg.RoleUser,
		Content:        request.Payload.(string),
	})
	if err != nil {
		utils.Log(utils.ErrorLevel).CID(request.SessionID).RID(request.MessageID).Err(err).BT().Send("Failed to request to conversation")
		return wspkg.ResponseWrapper{}, false, utils.WrapError(err, "failed to request to conversation")
	}
	marshaled, _ := json.Marshal(llmResponse.Metadata)
	utils.Log(utils.DebugLevel).CID(request.SessionID).RID(request.MessageID).BT().Send("Message ID: %v", llmResponse.ID)
	utils.Log(utils.DebugLevel).CID(request.SessionID).RID(request.MessageID).BT().Send("Message Metadata: %v", string(marshaled))

	response = wspkg.BuildResponseFrom(
		request, llmResponse.ID,
		wspkg.PredefinedActionData, llmResponse.Content,
	)

	return response, shouldClose, nil
}

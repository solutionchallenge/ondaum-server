package chat

import (
	"context"

	"github.com/google/uuid"
	domain "github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	wspkg "github.com/solutionchallenge/ondaum-server/pkg/websocket"
	"github.com/uptrace/bun"
)

func HandlePing(db *bun.DB, request wspkg.PingWrapper) (wspkg.ResponseWrapper, bool, error) {
	if !request.Authorized || !checkAuthorization(db, request.UserID) {
		utils.Log(utils.InfoLevel).CID(request.SessionID).BT().Send("Unauthorized")
		return wspkg.BuildRejectResponse(request), false, nil
	}

	chat := &domain.Chat{}
	err := db.NewSelect().
		Model(chat).
		Where("user_id = ? AND session_id = ?", request.UserID, request.SessionID).
		Scan(context.Background())
	if err != nil {
		utils.Log(utils.ErrorLevel).CID(request.SessionID).Err(err).BT().Send("Failed to query chat")
		return wspkg.ResponseWrapper{}, false, err
	}

	if !chat.FinishedAt.IsZero() && chat.ArchivedAt.IsZero() {
		return wspkg.BuildResponseFrom(
			request, uuid.New().String(),
			wspkg.PredefinedActionNotify, ChatPayloadNotifyConversationFinished,
		), false, nil
	}

	if !chat.ArchivedAt.IsZero() {
		return wspkg.BuildCloseResponse(
			request, wspkg.PredefinedActionNotify, ChatPayloadNotifyConversationArchived,
		), false, nil
	}

	return wspkg.BuildNoopResponse(request), false, nil
}

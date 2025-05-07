package chat

import (
	"context"
	"database/sql"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/google/uuid"
	domain "github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	wspkg "github.com/solutionchallenge/ondaum-server/pkg/websocket"
	"github.com/uptrace/bun"
)

func HandleConnect(db *bun.DB, clk clock.Clock, request wspkg.ConnectWrapper) (wspkg.ResponseWrapper, string, error) {
	if !request.Authorized || !checkAuthorization(db, request.UserID) {
		utils.Log(utils.InfoLevel).CID(request.ConnectID).BT().Send("Unauthorized")
		return wspkg.BuildRejectResponse(request), "", nil
	}

	chat := &domain.Chat{}
	err := db.NewSelect().
		Model(chat).
		Where("user_id = ?", request.UserID).
		Where("archived_at IS NULL").
		Order("created_at DESC").
		Limit(1).
		Scan(context.Background())

	if err == sql.ErrNoRows {
		chat = &domain.Chat{
			UserID:       request.UserID,
			SessionID:    request.ConnectID,
			StartedDate:  clk.Now().Truncate(24 * time.Hour),
			UserTimezone: clk.Now().Location().String(),
		}
		_, err = db.NewInsert().Model(chat).Exec(context.Background())
		if err != nil {
			utils.Log(utils.ErrorLevel).CID(request.ConnectID).Err(err).BT().Send("Failed to create chat")
			return wspkg.ResponseWrapper{}, "", err
		}
		return wspkg.BuildResponseFrom(
			request, uuid.New().String(),
			wspkg.PredefinedActionNotify, ChatPayloadNotifyNewConversation,
		), request.ConnectID, nil
	}

	request.ConnectID = chat.SessionID
	if !chat.FinishedAt.IsZero() {
		return wspkg.BuildResponseFrom(
			request, uuid.New().String(),
			wspkg.PredefinedActionNotify, ChatPayloadNotifyConversationFinished,
		), request.ConnectID, nil
	}
	return wspkg.BuildResponseFrom(
		request, uuid.New().String(),
		wspkg.PredefinedActionNotify, ChatPayloadNotifyExistingConversation,
	), request.ConnectID, nil
}

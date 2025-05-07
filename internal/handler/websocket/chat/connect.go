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

func HandleConnect(db *bun.DB, clk clock.Clock, request wspkg.ConnectWrapper) (wspkg.ResponseWrapper, error) {
	if !request.Authorized || !checkAuthorization(db, request.UserID) {
		utils.Log(utils.InfoLevel).CID(request.SessionID).BT().Send("Unauthorized")
		return wspkg.BuildRejectResponse(request), nil
	}

	chat := &domain.Chat{}
	err := db.NewSelect().
		Model(chat).
		Where("user_id = ? AND session_id = ?", request.UserID, request.SessionID).
		Scan(context.Background())
	if err != nil && err != sql.ErrNoRows {
		utils.Log(utils.ErrorLevel).CID(request.SessionID).Err(err).BT().Send("Failed to query chat")
		return wspkg.ResponseWrapper{}, err
	}

	if err == sql.ErrNoRows {
		chat = &domain.Chat{
			UserID:       request.UserID,
			SessionID:    request.SessionID,
			StartedDate:  clk.Now().Truncate(24 * time.Hour),
			UserTimezone: clk.Now().Location().String(),
		}
		_, err = db.NewInsert().Model(chat).Exec(context.Background())
		if err != nil {
			utils.Log(utils.ErrorLevel).CID(request.SessionID).Err(err).BT().Send("Failed to create chat")
			return wspkg.ResponseWrapper{}, err
		}
		return wspkg.BuildResponseFrom(
			request, uuid.New().String(),
			wspkg.PredefinedActionNotify, ChatPayloadNotifyNewConversation,
		), nil
	}

	if !chat.FinishedAt.IsZero() && chat.ArchivedAt.IsZero() {
		return wspkg.BuildResponseFrom(
			request, uuid.New().String(),
			wspkg.PredefinedActionNotify, ChatPayloadNotifyConversationFinished,
		), nil
	}

	if !chat.ArchivedAt.IsZero() {
		return wspkg.BuildCloseResponse(
			request, wspkg.PredefinedActionNotify, ChatPayloadNotifyConversationArchived,
		), nil
	}

	return wspkg.BuildResponseFrom(
		request, uuid.New().String(),
		wspkg.PredefinedActionNotify, ChatPayloadNotifyExistingConversation,
	), nil
}

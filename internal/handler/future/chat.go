package future

import (
	"context"
	"encoding/json"

	"github.com/benbjohnson/clock"
	domain "github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	"github.com/solutionchallenge/ondaum-server/pkg/future"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

const (
	ChatJobType = future.JobType("chat")
)

type ChatFutureHandlerDependencies struct {
	fx.In
	DB    *bun.DB
	Clock clock.Clock
}

type ChatFutureHandlerParams struct {
	UserID         int64
	ConversationID string
}

type ChatFutureHandler struct {
	deps ChatFutureHandlerDependencies
}

func NewChatFutureHandler(deps ChatFutureHandlerDependencies) (*ChatFutureHandler, error) {
	return &ChatFutureHandler{
		deps: deps,
	}, nil
}

func (h *ChatFutureHandler) Handle(ctx context.Context, job *future.Job) error {
	var input ChatFutureHandlerParams
	err := json.Unmarshal([]byte(job.ActionParams), &input)
	if err != nil {
		return utils.WrapError(err, "failed to unmarshal future job params")
	}

	tx, err := h.deps.DB.BeginTx(ctx, nil)
	if err != nil {
		return utils.WrapError(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	chat := &domain.Chat{}
	err = tx.NewSelect().
		Model(chat).
		Where("session_id = ?", input.ConversationID).
		Where("user_id = ?", input.UserID).
		Scan(ctx)
	if err != nil {
		return utils.WrapError(err, "failed to select chat (%v:%v)", input.UserID, input.ConversationID)
	}

	if chat.FinishedAt.IsZero() {
		_, err = tx.NewUpdate().
			Model(chat).
			Set("finished_at = CURRENT_TIMESTAMP").
			WherePK().
			Exec(ctx)
		if err != nil {
			return utils.WrapError(err, "failed to update chat (%v:%v)", input.UserID, input.ConversationID)
		}
	}

	if err := tx.Commit(); err != nil {
		return utils.WrapError(err, "failed to commit transaction (%v:%v)", input.UserID, input.ConversationID)
	}

	return nil
}

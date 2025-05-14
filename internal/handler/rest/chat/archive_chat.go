package chat

import (
	"database/sql"
	"errors"

	"github.com/benbjohnson/clock"
	"github.com/gofiber/fiber/v2"
	domain "github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type ArchiveChatHandlerDependencies struct {
	fx.In
	DB    *bun.DB
	Clock clock.Clock
}

type ArchiveChatHandler struct {
	deps ArchiveChatHandlerDependencies
}

type ArchiveChatHandlerResponse struct {
	Success  bool `json:"success"`
	Finished bool `json:"finished"`
}

func NewArchiveChatHandler(deps ArchiveChatHandlerDependencies) (*ArchiveChatHandler, error) {
	return &ArchiveChatHandler{deps: deps}, nil
}

// @ID ArchiveChat
// @Summary Archive a chat
// @Description Archive a chat to prevent it from being accessed again and allow to summarize it.
// @Tags chat
// @Accept json
// @Produce json
// @Param session_id path string true "Session ID"
// @Success 200 {object} ArchiveChatHandlerResponse
// @Failure 401 {object} http.Error
// @Failure 404 {object} http.Error
// @Failure 500 {object} http.Error
// @Router /chats/{session_id}/archive [post]
// @Security BearerAuth
func (h *ArchiveChatHandler) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userID, err := http.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			http.NewError(ctx, err, "Unauthorized"),
		)
	}

	user := &user.User{ID: userID}
	if err := h.deps.DB.NewSelect().Model(user).Where("id = ?", userID).Scan(ctx); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(ctx, err, "User not found"),
		)
	}

	sessionID := c.Params("session_id")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(ctx, errors.New("session_id is required"), "Bad Request"),
		)
	}

	tx, err := h.deps.DB.BeginTx(ctx, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to begin transaction"),
		)
	}
	defer tx.Rollback()

	chat := &domain.Chat{}
	if err := tx.NewSelect().
		Model(chat).
		Where("session_id = ?", sessionID).
		Where("user_id = ?", userID).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(
				http.NewError(ctx, err, "Chat not found for session_id: %v", sessionID),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to get chat for session_id: %v", sessionID),
		)
	}

	now := h.deps.Clock.Now().UTC()
	endedAt := now
	if !chat.FinishedAt.IsZero() && !chat.FinishedAt.Time.Before(chat.CreatedAt) {
		endedAt = chat.FinishedAt.Time
	}
	chatDuration := common.NewNullableDuration(endedAt.Sub(chat.CreatedAt))

	wasFinished := chat.FinishedAt.IsZero()
	mutation := tx.NewUpdate().
		Model(chat).
		Set("archived_at = CURRENT_TIMESTAMP").
		Set("chat_duration = ?", chatDuration).
		Where("session_id = ?", sessionID).
		Where("user_id = ?", userID)

	if wasFinished {
		mutation = mutation.Set("finished_at = CURRENT_TIMESTAMP")
	}

	if _, err := mutation.Exec(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to archive chat"),
		)
	}

	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to commit transaction"),
		)
	}

	return c.JSON(ArchiveChatHandlerResponse{
		Success:  true,
		Finished: wasFinished,
	})
}

func (h *ArchiveChatHandler) Identify() string {
	return "archive-chat"
}

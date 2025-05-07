package chat

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	domain "github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type ArchiveChatHandlerDependencies struct {
	fx.In
	DB *bun.DB
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
// @Router /chats/:session_id/archive [post]
func (h *ArchiveChatHandler) Handle(c *fiber.Ctx) error {
	userID, err := http.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			http.NewError(c.UserContext(), err, "Unauthorized"),
		)
	}
	user := &user.User{ID: userID}
	if err := h.deps.DB.NewSelect().Model(user).Where("id = ?", userID).Scan(c.UserContext()); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(c.UserContext(), err, "User not found"),
		)
	}

	sessionID := c.Params("session_id")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), errors.New("session_id is required"), "Bad Request"),
		)
	}

	tx, err := h.deps.DB.BeginTx(c.UserContext(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to begin transaction"),
		)
	}
	defer tx.Rollback()

	chat := &domain.Chat{}
	if err := tx.NewSelect().
		Model(chat).
		Where("session_id = ?", sessionID).
		Where("user_id = ?", userID).
		Scan(c.UserContext()); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(c.UserContext(), err, "Chat not found"),
		)
	}

	wasFinished := chat.FinishedAt.IsZero()
	mutation := tx.NewUpdate().
		Model(chat).
		Set("archived_at = CURRENT_TIMESTAMP").
		Where("session_id = ?", sessionID).
		Where("user_id = ?", userID)
	if !wasFinished {
		mutation = mutation.Set("finished_at = CURRENT_TIMESTAMP")
	}
	if _, err := mutation.Exec(c.UserContext()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to archive chat"),
		)
	}

	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to commit transaction"),
		)
	}

	response := ArchiveChatHandlerResponse{
		Success:  true,
		Finished: wasFinished,
	}
	return c.JSON(response)
}

func (h *ArchiveChatHandler) Identify() string {
	return "archive-chat"
}

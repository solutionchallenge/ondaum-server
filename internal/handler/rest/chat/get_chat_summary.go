package chat

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type GetChatSummaryHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type GetChatSummaryHandler struct {
	deps GetChatSummaryHandlerDependencies
}

func NewGetChatSummaryHandler(deps GetChatSummaryHandlerDependencies) (*GetChatSummaryHandler, error) {
	return &GetChatSummaryHandler{deps: deps}, nil
}

// @ID GetSummary
// @Summary Get summary
// @Description Get summary of the chat
// @Tags chat
// @Accept json
// @Produce json
// @Param session_id path string true "Session ID"
// @Success 200 {object} chat.SummaryDTO
// @Failure 401 {object} http.Error
// @Failure 404 {object} http.Error
// @Failure 500 {object} http.Error
// @Router /chats/{session_id}/summary [get]
// @Security BearerAuth
func (h *GetChatSummaryHandler) Handle(c *fiber.Ctx) error {
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
	chat := &chat.Chat{}
	if err := h.deps.DB.NewSelect().
		Model(chat).
		Relation("Summary").
		Relation("Histories", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("inserted_at ASC")
		}).
		Where("session_id = ?", sessionID).
		Where("user_id = ?", userID).
		Scan(ctx); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(ctx, err, "Failed to get summary"),
		)
	}
	if chat.Summary == nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(ctx, errors.New("summary not found"), "Failed to get summary"),
		)
	}
	response := chat.Summary.ToSummaryDTO(chat.Histories)
	return c.JSON(response)
}

func (h *GetChatSummaryHandler) Identify() string {
	return "get-chat-summary"
}

package chat

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	domain "github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type GetChatHandlerDependencies struct {
	fx.In
	DB *bun.DB
}
type GetChatHandler struct {
	deps GetChatHandlerDependencies
}

func NewGetChatHandler(deps GetChatHandlerDependencies) (*GetChatHandler, error) {
	return &GetChatHandler{deps: deps}, nil
}

// @ID GetChat
// @Summary Get chat
// @Description Get chat with histories and summary
// @Tags chat
// @Accept json
// @Produce json
// @Param session_id path string true "Session ID"
// @Success 200 {object} domain.ChatWithSimplifiedSummaryAndHistoriesDTO
// @Failure 401 {object} http.Error
// @Failure 404 {object} http.Error
// @Failure 500 {object} http.Error
// @Router /chat/{session_id} [get]
// @Security BearerAuth
func (h *GetChatHandler) Handle(c *fiber.Ctx) error {
	userID, err := http.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			http.NewError(c.UserContext(), err, "Unauthorized"),
		)
	}
	user := &user.User{ID: userID}
	if err := h.deps.DB.NewSelect().Model(user).Where("id = ?", userID).Scan(context.Background()); err != nil {
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

	chat := &domain.Chat{}
	err = h.deps.DB.NewSelect().
		Model(chat).
		Relation("Histories").
		Relation("Summary").
		Where("session_id = ?", sessionID).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Scan(context.Background())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(c.UserContext(), err, "Chat not found"),
		)
	}

	response := chat.ToChatWithSimplifiedSummaryAndHistoriesDTO()
	return c.JSON(response)
}

func (h *GetChatHandler) Identify() string {
	return "get-chat"
}

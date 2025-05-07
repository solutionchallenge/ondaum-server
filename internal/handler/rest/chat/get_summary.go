package chat

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type GetSummaryHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type GetSummaryHandler struct {
	deps GetSummaryHandlerDependencies
}

func NewGetSummaryHandler(deps GetSummaryHandlerDependencies) (*GetSummaryHandler, error) {
	return &GetSummaryHandler{deps: deps}, nil
}

// @ID GetSummary
// @Summary Get summary
// @Description Get summary of the chat
// @Tags chat
// @Accept json
// @Produce json
// @Param session_id path string true "Session ID"
// @Success 200 {object} chat.SimplifiedSummaryDTO
// @Failure 401 {object} http.Error
// @Failure 404 {object} http.Error
// @Failure 500 {object} http.Error
// @Router /chat/{session_id}/summary [get]
// @Security BearerAuth
func (h *GetSummaryHandler) Handle(c *fiber.Ctx) error {
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
	chat := &chat.Chat{}
	if err := h.deps.DB.NewSelect().
		Model(chat).
		Relation("Summary").
		Where("session_id = ?", sessionID).
		Where("user_id = ?", userID).
		Scan(context.Background()); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(c.UserContext(), err, "Failed to get summary"),
		)
	}
	if chat.Summary == nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(c.UserContext(), errors.New("summary not found"), "Failed to get summary"),
		)
	}
	response := chat.Summary.ToSimplifiedSummaryDTO()
	return c.JSON(response)
}

func (h *GetSummaryHandler) Identify() string {
	return "get-summary"
}

package chat

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type ListChatHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type ListChatHandlerResponse struct {
	Chats []chat.ChatWithSimplifiedSummaryDTO `json:"chats"`
}

type ListChatHandler struct {
	deps ListChatHandlerDependencies
}

func NewListChatHandler(deps ListChatHandlerDependencies) (*ListChatHandler, error) {
	return &ListChatHandler{deps: deps}, nil
}

func (h *ListChatHandler) Handle(c *fiber.Ctx) error {
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
	chats := []*chat.Chat{}
	if err := h.deps.DB.NewSelect().
		Model(&chats).
		Relation("Summary").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Scan(context.Background()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to list chats"),
		)
	}

	dtos := utils.Map(chats, func(c *chat.Chat) chat.ChatWithSimplifiedSummaryDTO {
		return c.ToChatWithSimplifiedSummaryDTO()
	})
	response := ListChatHandlerResponse{
		Chats: dtos,
	}
	return c.JSON(response)
}

func (h *ListChatHandler) Identify() string {
	return "list-chat"
}

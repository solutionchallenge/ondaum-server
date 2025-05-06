package debug

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type ListChatHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type ListHistoryHandlerResponse struct {
	Chats []chat.Chat `json:"chats"`
}

type ListChatHandler struct {
	deps ListChatHandlerDependencies
}

func NewListChatHandler(deps ListChatHandlerDependencies) (*ListChatHandler, error) {
	return &ListChatHandler{deps: deps}, nil
}

// Must not be documented. (Debugging purpose only!)
func (h *ListChatHandler) Handle(c *fiber.Ctx) error {
	if os.Getenv("FLAG_DEBUGGING_FEATURES_ENABLED") != "true" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	cid := c.Params("chat_id")
	sid := c.Params("session_id")

	chats := []chat.Chat{}
	query := h.deps.DB.NewSelect().
		Model(&chats).
		Relation("User").
		Relation("Histories")
	if cid != "" {
		query = query.Where("id = ?", cid)
	}
	if sid != "" {
		query = query.Where("session_id = ?", sid)
	}
	err := query.Scan(c.UserContext(), &chats)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(c.UserContext(), err, "Chat not found"),
		)
	}
	return c.JSON(ListHistoryHandlerResponse{Chats: chats})
}

func (h *ListChatHandler) Identify() string {
	return "list_chat"
}

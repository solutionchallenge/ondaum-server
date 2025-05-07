package chat

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	domain "github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type UpsertSummaryHandlerDependencies struct {
	fx.In
	DB  *bun.DB
	LLM llm.Client
}

type UpsertSummaryHandlerResponse struct {
	Success   bool                        `json:"success"`
	Created   bool                        `json:"created"`
	Returning domain.SimplifiedSummaryDTO `json:"returning"`
}

type UpsertSummaryHandler struct {
	deps UpsertSummaryHandlerDependencies
}

func NewUpsertSummaryHandler(deps UpsertSummaryHandlerDependencies) (*UpsertSummaryHandler, error) {
	return &UpsertSummaryHandler{deps: deps}, nil
}

func (h *UpsertSummaryHandler) Handle(c *fiber.Ctx) error {
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
		Where("session_id = ?", sessionID).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Scan(context.Background())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(c.UserContext(), err, "Chat not found"),
		)
	}
	if chat.ArchivedAt.IsZero() {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(c.UserContext(), errors.New("chat is not archived"), "Chat is not archived"),
		)
	}
	if len(chat.Histories) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(c.UserContext(), errors.New("chat history not found"), "Chat history not found"),
		)
	}

	histories := utils.Map(chat.Histories, func(h *domain.History) llm.Message {
		return llm.Message{
			Role:    llm.Role(h.Role),
			Content: h.Content,
		}
	})

	resolved, err := h.deps.LLM.ResolvePrompt(c.UserContext(), "interactive_chat", "summary_chat", histories...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to resolve prompt"),
		)
	}

	summary := struct {
		Title    string   `json:"title"`
		Text     string   `json:"text"`
		Keywords []string `json:"keywords"`
		Emotions []struct {
			Emotion common.Emotion `json:"emotion"`
			Rate    float64        `json:"rate"`
		} `json:"emotions"`
	}{}
	if err := json.Unmarshal([]byte(resolved.Content), &summary); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to unmarshal response"),
		)
	}
	var emotions common.EmotionRateList = utils.Map(summary.Emotions, func(e struct {
		Emotion common.Emotion `json:"emotion"`
		Rate    float64        `json:"rate"`
	}) common.EmotionRate {
		return common.EmotionRate{
			Emotion: e.Emotion,
			Rate:    e.Rate,
		}
	})
	model := &domain.Summary{
		ChatID:   chat.ID,
		Title:    summary.Title,
		Text:     summary.Text,
		Keywords: summary.Keywords,
		Emotions: emotions,
	}
	result, err := h.deps.DB.NewInsert().
		Model(model).
		On("DUPLICATE KEY UPDATE").
		Set("title = ?", summary.Title).
		Set("text = ?", summary.Text).
		Set("keywords = ?", summary.Keywords).
		Set("emotions = ?", emotions.ToString()).
		Set("updated_at = CURRENT_TIMESTAMP").
		Exec(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to insert summary"),
		)
	}
	rowsAffected, _ := result.RowsAffected()
	response := &UpsertSummaryHandlerResponse{
		Success:   true,
		Created:   rowsAffected == 1,
		Returning: model.ToSimplifiedSummaryDTO(),
	}
	return c.JSON(response)
}

func (h *UpsertSummaryHandler) Identify() string {
	return "summary-chat"
}

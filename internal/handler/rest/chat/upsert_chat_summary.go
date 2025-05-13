package chat

import (
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

type UpsertChatSummaryHandlerDependencies struct {
	fx.In
	DB  *bun.DB
	LLM llm.Client
}

type UpsertChatSummaryHandlerResponse struct {
	Success   bool              `json:"success"`
	Created   bool              `json:"created"`
	Returning domain.SummaryDTO `json:"returning"`
}

type UpsertChatSummaryHandler struct {
	deps UpsertChatSummaryHandlerDependencies
}

func NewUpsertChatSummaryHandler(deps UpsertChatSummaryHandlerDependencies) (*UpsertChatSummaryHandler, error) {
	return &UpsertChatSummaryHandler{deps: deps}, nil
}

// @ID UpsertChatSummary
// @Summary Create or update chat summary
// @Description Create or update chat summary and return the created/updated chat summary
// @Tags chat
// @Accept json
// @Produce json
// @Param session_id path string true "Session ID"
// @Success 200 {object} UpsertChatSummaryHandlerResponse
// @Failure 401 {object} http.Error
// @Failure 404 {object} http.Error
// @Failure 500 {object} http.Error
// @Router /chats/{session_id}/summary [post]
// @Security BearerAuth
func (h *UpsertChatSummaryHandler) Handle(c *fiber.Ctx) error {
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

	chat := &domain.Chat{}
	err = h.deps.DB.NewSelect().
		Model(chat).
		Relation("Histories", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("inserted_at ASC")
		}).
		Where("session_id = ?", sessionID).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Scan(ctx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(ctx, err, "Chat not found"),
		)
	}
	if chat.ArchivedAt.IsZero() {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(ctx, errors.New("chat is not archived"), "Chat is not archived"),
		)
	}
	if len(chat.Histories) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(ctx, errors.New("chat history not found"), "Chat history not found"),
		)
	}

	histories := utils.Map(chat.Histories, func(h *domain.History) llm.Message {
		return llm.Message{
			Role:    llm.Role(h.Role),
			Content: h.Content,
		}
	})

	resolved, err := h.deps.LLM.RunActionPrompt(ctx, "interactive_chat", "summary_chat", histories...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to resolve prompt"),
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
		Recommendations []string `json:"recommendations"`
		PositiveScore   float64  `json:"positive_score"`
		NegativeScore   float64  `json:"negative_score"`
		NeutralScore    float64  `json:"neutral_score"`
		// We must fetch indicies and convert them to history IDs later because the LLM can't get corresponding history IDs.
		MainTopic struct {
			BeginHistoryIndex int `json:"begin_history_index"`
			EndHistoryIndex   int `json:"end_history_index"`
		} `json:"main_topic"`
	}{}
	if err := json.Unmarshal([]byte(resolved.Content), &summary); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to unmarshal response"),
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
	convertedMainTopic := func() domain.MainTopic {
		if summary.MainTopic.BeginHistoryIndex < 0 || summary.MainTopic.EndHistoryIndex < 0 {
			return domain.MainTopic{}
		}
		if summary.MainTopic.BeginHistoryIndex > summary.MainTopic.EndHistoryIndex {
			return domain.MainTopic{}
		}
		if summary.MainTopic.BeginHistoryIndex >= len(chat.Histories) || summary.MainTopic.EndHistoryIndex >= len(chat.Histories) {
			return domain.MainTopic{}
		}
		beginHistoryID := chat.Histories[summary.MainTopic.BeginHistoryIndex].MessageID
		endHistoryID := chat.Histories[summary.MainTopic.EndHistoryIndex].MessageID
		return domain.MainTopic{
			BeginMessageID: beginHistoryID,
			EndMessageID:   endHistoryID,
		}
	}()
	model := &domain.Summary{
		ChatID:          chat.ID,
		Title:           summary.Title,
		Text:            summary.Text,
		Keywords:        summary.Keywords,
		Emotions:        emotions,
		Recommendations: summary.Recommendations,
		PositiveScore:   summary.PositiveScore,
		NegativeScore:   summary.NegativeScore,
		NeutralScore:    summary.NeutralScore,
		MainTopic:       convertedMainTopic,
	}
	result, err := h.deps.DB.NewInsert().
		Model(model).
		On("DUPLICATE KEY UPDATE").
		Set("title = ?", summary.Title).
		Set("text = ?", summary.Text).
		Set("keywords = ?", summary.Keywords).
		Set("emotions = ?", emotions.ToString()).
		Set("recommendations = ?", summary.Recommendations).
		Set("positive_score = ?", summary.PositiveScore).
		Set("negative_score = ?", summary.NegativeScore).
		Set("neutral_score = ?", summary.NeutralScore).
		Set("main_topic = ?", convertedMainTopic.ToString()).
		Set("updated_at = CURRENT_TIMESTAMP").
		Exec(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to insert summary"),
		)
	}
	rowsAffected, _ := result.RowsAffected()
	response := &UpsertChatSummaryHandlerResponse{
		Success:   true,
		Created:   rowsAffected == 1,
		Returning: model.ToSummaryDTO(chat.Histories),
	}
	return c.JSON(response)
}

func (h *UpsertChatSummaryHandler) Identify() string {
	return "upsert-chat-summary"
}

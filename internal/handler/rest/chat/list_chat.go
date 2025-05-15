package chat

import (
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	domain "github.com/solutionchallenge/ondaum-server/internal/domain/chat"
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

type ListChatHandler struct {
	deps ListChatHandlerDependencies
}

type ListChatResponse struct {
	Chats []domain.ChatDTO `json:"chats"`
}

func NewListChatHandler(deps ListChatHandlerDependencies) (*ListChatHandler, error) {
	return &ListChatHandler{deps: deps}, nil
}

// @ID ListChat
// @Summary List chats
// @Description List chats with optional filters for datetime range and emotion
// @Tags chat
// @Accept json
// @Produce json
// @Param datetime_gte query string false "Filter by chat started datetime in ISO 8601 format (YYYY-MM-DDTHH:mm:ssZ)"
// @Param datetime_lte query string false "Filter by chat ended datetime in ISO 8601 format (YYYY-MM-DDTHH:mm:ssZ)"
// @Param dominant_emotions query string false "Filter by dominant emotions (comma separated, e.g. 'joy,sadness')"
// @Param matching_keyword query string false "Filter by sub-string matching keyword"
// @Param matching_content query string false "Filter by sub-string matching content (search raw-text from all contents, it could be slow for large data)"
// @Param message_id query string false "Filter by message ID"
// @Param only_archived query bool false "Filter only archived chats"
// @Success 200 {object} ListChatResponse
// @Failure 401 {object} http.Error
// @Failure 400 {object} http.Error
// @Failure 500 {object} http.Error
// @Router /chats [get]
// @Security BearerAuth
func (h *ListChatHandler) Handle(c *fiber.Ctx) error {
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

	datetimeGte := c.Query("datetime_gte")
	datetimeLte := c.Query("datetime_lte")
	dominantEmotions := c.Query("dominant_emotions")
	matchingKeyword := c.Query("matching_keyword")
	matchingContent := c.Query("matching_content")
	messageID := c.Query("message_id")
	onlyArchivedStr := c.Query("only_archived")

	query := h.deps.DB.NewSelect().
		Model((*domain.Chat)(nil)).
		Relation("Summary").
		Where("c.user_id = ?", userID)

	if datetimeGte != "" {
		localStartTime, err := time.Parse(time.RFC3339, datetimeGte)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				http.NewError(ctx, err, "Invalid datetime_gte format. Use YYYY-MM-DDTHH:mm:ssZ"),
			)
		}
		query = query.Where("c.created_at >= ?", localStartTime.UTC())
	}
	if datetimeLte != "" {
		localEndTime, err := time.Parse(time.RFC3339, datetimeLte)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				http.NewError(ctx, err, "Invalid datetime_lte format. Use YYYY-MM-DDTHH:mm:ssZ"),
			)
		}
		query = query.Where("(c.archived_at IS NULL OR c.archived_at <= ?)", localEndTime.UTC())
	}

	if onlyArchivedStr != "" {
		onlyArchived, err := strconv.ParseBool(onlyArchivedStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				http.NewError(ctx, err, "Invalid only_archived value. Use true or false"),
			)
		}
		if onlyArchived {
			query = query.Where("c.archived_at IS NOT NULL")
		}
	}

	if messageID != "" {
		query = query.
			Join("JOIN chat_histories ch").
			JoinOn("ch.chat_id = c.id").
			Where("ch.message_id = ?", messageID)
	}

	var chats []domain.Chat
	if err := query.Scan(ctx, &chats); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to list chats"),
		)
	}

	if matchingKeyword != "" {
		chats = utils.Filter(chats, func(chat domain.Chat) bool {
			if chat.Summary == nil {
				return false
			}
			return utils.OneOf(chat.Summary.Keywords, func(keyword string) bool {
				return strings.Contains(strings.ToLower(keyword), strings.ToLower(matchingKeyword))
			})
		})
	}

	if matchingContent != "" {
		chats = utils.Filter(chats, func(chat domain.Chat) bool {
			allContents := []string{}
			summaryTitle := chat.Summary.Title
			summaryText := chat.Summary.Text
			allContents = append(allContents, summaryTitle, summaryText)

			summaryDTO := chat.Summary.ToSummaryWithTopicMessages(chat.Histories)
			summaryContents := utils.Map(*summaryDTO.TopicMessages, func(topicMessage domain.HistoryDTO) string {
				return topicMessage.Content
			})
			allContents = append(allContents, summaryContents...)

			historyContents := utils.Map(chat.Histories, func(history *domain.History) string {
				return history.Content
			})
			allContents = append(allContents, historyContents...)

			return utils.OneOf(allContents, func(content string) bool {
				return strings.Contains(strings.ToLower(content), strings.ToLower(matchingContent))
			})
		})
	}

	if dominantEmotions != "" {
		emotions := strings.Split(dominantEmotions, ",")
		trimedEmotions := utils.Map(emotions, func(emotion string) string {
			return strings.TrimSpace(emotion)
		})
		chats = utils.Filter(chats, func(chat domain.Chat) bool {
			if chat.Summary == nil {
				return false
			}
			dominant := chat.Summary.Emotions.GetDominant()
			return slices.Contains(trimedEmotions, string(dominant))
		})
	}

	chatDTOs := utils.Map(chats, func(chat domain.Chat) domain.ChatDTO {
		return chat.ToChatDTO()
	})

	return c.JSON(ListChatResponse{
		Chats: chatDTOs,
	})
}

func (h *ListChatHandler) Identify() string {
	return "list-chat"
}

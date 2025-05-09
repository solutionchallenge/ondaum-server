package chat

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	domain "github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
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
	Chats []domain.ChatWithSummaryDTO `json:"chats"`
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
// @Param dominant_emotion query string false "Filter by dominant emotion"
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
	dominantEmotion := c.Query("dominant_emotion")
	onlyArchivedStr := c.Query("only_archived")

	var startTime, endTime time.Time
	if datetimeGte != "" {
		startTime, err = time.Parse(time.RFC3339, datetimeGte)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				http.NewError(ctx, err, "Invalid datetime_gte format. Use YYYY-MM-DDTHH:mm:ssZ"),
			)
		}
	}
	if datetimeLte != "" {
		endTime, err = time.Parse(time.RFC3339, datetimeLte)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				http.NewError(ctx, err, "Invalid datetime_lte format. Use YYYY-MM-DDTHH:mm:ssZ"),
			)
		}
	}

	onlyArchived := false
	if onlyArchivedStr != "" {
		onlyArchived, err = strconv.ParseBool(onlyArchivedStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				http.NewError(ctx, err, "Invalid only_archived value. Use true or false"),
			)
		}
	}

	query := h.deps.DB.NewSelect().
		Model((*domain.Chat)(nil)).
		Relation("Summary").
		Where("user_id = ?", userID)

	if onlyArchived {
		query = query.Where("archived_at IS NOT NULL")
	}

	if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("(finished_at IS NULL OR finished_at <= ?)", endTime)
	}

	var chats []domain.Chat
	if err := query.Scan(ctx, &chats); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to list chats"),
		)
	}

	if dominantEmotion != "" {
		filteredChats := make([]domain.Chat, 0, len(chats))
		for _, chat := range chats {
			if chat.Summary == nil {
				continue
			}
			var maxRate float64
			var dominant common.Emotion
			for _, emotion := range chat.Summary.Emotions {
				if emotion.Rate > maxRate {
					maxRate = emotion.Rate
					dominant = emotion.Emotion
				}
			}
			if dominant == common.Emotion(dominantEmotion) {
				filteredChats = append(filteredChats, chat)
			}
		}
		chats = filteredChats
	}

	chatDTOs := make([]domain.ChatWithSummaryDTO, len(chats))
	for i, chat := range chats {
		chatDTOs[i] = chat.ToChatWithSummaryDTO()
	}

	return c.JSON(ListChatResponse{
		Chats: chatDTOs,
	})
}

func (h *ListChatHandler) Identify() string {
	return "list-chat"
}

package chat

import (
	"time"

	"github.com/gofiber/fiber/v2"
	domain "github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type PredefinedStressLevel struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Threshold   float64 `json:"threshold"`
}

var PredefinedStressLevels = []PredefinedStressLevel{
	{
		Title:       "Your Stress Level seems High",
		Description: "Talking to someone you trust can lift the weight off your heart.",
		Threshold:   0.7,
	},
	{
		Title:       "Your stress level seems Moderate",
		Description: "It's important to take care of yourself and find ways to relax.",
		Threshold:   0.5,
	},
	{
		Title:       "Your stress level seems Low",
		Description: "You seem to be under a low amount of stress, keep it up!",
		Threshold:   0.3,
	},
}

type GetChatReportHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type GetChatReportHandlerResponse struct {
	DatetimeGte           time.Time               `json:"datetime_gte"`
	DatetimeLte           time.Time               `json:"datetime_lte"`
	EmotionCounts         common.EmotionCountList `json:"emotion_counts"`
	AveragePositiveScore  float64                 `json:"average_positive_score"`
	AverageNegativeScore  float64                 `json:"average_negative_score"`
	AverageNeutralScore   float64                 `json:"average_neutral_score"`
	TotalChatCount        int                     `json:"total_chat_count"`
	AverageChatDuration   string                  `json:"average_chat_duration"`
	StressLevelDescriptor PredefinedStressLevel   `json:"stress_level_descriptor"`
}

type GetChatReportHandler struct {
	deps GetChatReportHandlerDependencies
}

func NewGetChatReportHandler(deps GetChatReportHandlerDependencies) (*GetChatReportHandler, error) {
	return &GetChatReportHandler{deps: deps}, nil
}

// @Summary Get chat report
// @Description Get chat report
// @Tags chat
// @Accept json
// @Produce json
// @Param datetime_gte query string false "Filter by chat started datetime in ISO 8601 format (YYYY-MM-DDTHH:mm:ssZ)"
// @Param datetime_lte query string false "Filter by chat ended datetime in ISO 8601 format (YYYY-MM-DDTHH:mm:ssZ)"
// @Success 200 {object} GetChatReportHandlerResponse
// @Failure 401 {object} http.Error
// @Failure 404 {object} http.Error
// @Failure 500 {object} http.Error
// @Router /chat-reports [get]
func (h *GetChatReportHandler) Handle(c *fiber.Ctx) error {
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

	localStartTime := time.Time{}
	localEndTime := time.Time{}

	query := h.deps.DB.NewSelect().
		Model((*domain.Chat)(nil)).
		Relation("Summary").
		Where("c.user_id = ?", userID)

	if datetimeGte != "" {
		parsedTime, err := time.Parse(time.RFC3339, datetimeGte)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				http.NewError(ctx, err, "Invalid datetime_gte format. Use YYYY-MM-DDTHH:mm:ssZ"),
			)
		}
		localStartTime = parsedTime
		query = query.Where("c.created_at >= ?", parsedTime.UTC())
	}
	if datetimeLte != "" {
		parsedTime, err := time.Parse(time.RFC3339, datetimeLte)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				http.NewError(ctx, err, "Invalid datetime_lte format. Use YYYY-MM-DDTHH:mm:ssZ"),
			)
		}
		localEndTime = parsedTime
		query = query.Where("(c.archived_at IS NULL OR c.archived_at <= ?)", parsedTime.UTC())
	}

	var chats []domain.Chat
	if err := query.Scan(ctx, &chats); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to get chats"),
		)
	}

	chatCount := len(chats)
	if chatCount == 0 {
		return c.Status(fiber.StatusOK).JSON(GetChatReportHandlerResponse{
			DatetimeGte: localStartTime,
			DatetimeLte: localEndTime,
		})
	}

	validChats := utils.Filter(chats, func(chat domain.Chat) bool {
		return chat.Summary != nil && chat.ChatDuration.Valid
	})

	validChatCount := len(validChats)
	if validChatCount == 0 {
		return c.Status(fiber.StatusOK).JSON(GetChatReportHandlerResponse{
			DatetimeGte: localStartTime,
			DatetimeLte: localEndTime,
		})
	}

	type ChatStats struct {
		PositiveScore float64
		NegativeScore float64
		NeutralScore  float64
		Duration      time.Duration
		EmotionGroups map[common.Emotion][]domain.Chat
	}

	stats := utils.Reduce(validChats, func(acc ChatStats, chat domain.Chat) ChatStats {
		dominantEmotion := chat.Summary.Emotions.GetDominant()
		acc.EmotionGroups[dominantEmotion] = append(acc.EmotionGroups[dominantEmotion], chat)

		return ChatStats{
			PositiveScore: acc.PositiveScore + chat.Summary.PositiveScore,
			NegativeScore: acc.NegativeScore + chat.Summary.NegativeScore,
			NeutralScore:  acc.NeutralScore + chat.Summary.NeutralScore,
			Duration:      acc.Duration + time.Duration(chat.ChatDuration.Duration),
			EmotionGroups: acc.EmotionGroups,
		}
	}, ChatStats{
		EmotionGroups: make(map[common.Emotion][]domain.Chat),
	})

	emotionCountList := make(common.EmotionCountList, 0, len(stats.EmotionGroups))
	for emotion, chats := range stats.EmotionGroups {
		emotionCountList = append(emotionCountList, common.EmotionCount{
			Emotion: emotion,
			Count:   int64(len(chats)),
		})
	}

	avgPositiveScore := utils.RoundTo(stats.PositiveScore/float64(validChatCount), 3)
	avgNegativeScore := utils.RoundTo(stats.NegativeScore/float64(validChatCount), 3)
	avgNeutralScore := utils.RoundTo(stats.NeutralScore/float64(validChatCount), 3)

	total := avgPositiveScore + avgNegativeScore + avgNeutralScore
	if total != 1.00 {
		diff := 1.00 - total
		if diff > 0 {
			avgPositiveScore += diff
		} else {
			if avgNeutralScore > 0 {
				avgNeutralScore += diff
			} else if avgNegativeScore > 0 {
				avgNegativeScore += diff
			} else {
				avgPositiveScore += diff
			}
		}
	}

	avgConvertedDuration := common.NewDuration(stats.Duration)
	avgRedableDuration := avgConvertedDuration.ToString(common.DurationFormatMinutes)

	var stressLevel PredefinedStressLevel
	for _, level := range PredefinedStressLevels {
		if avgNegativeScore >= level.Threshold {
			stressLevel = level
			break
		}
	}

	return c.Status(fiber.StatusOK).JSON(GetChatReportHandlerResponse{
		DatetimeGte:           localStartTime,
		DatetimeLte:           localEndTime,
		EmotionCounts:         emotionCountList,
		AveragePositiveScore:  avgPositiveScore,
		AverageNegativeScore:  avgNegativeScore,
		AverageNeutralScore:   avgNeutralScore,
		TotalChatCount:        validChatCount,
		AverageChatDuration:   avgRedableDuration,
		StressLevelDescriptor: stressLevel,
	})
}

func (h *GetChatReportHandler) Identify() string {
	return "get-chat-report"
}

package chat

import (
	"math"
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
	AverageChatDuration   time.Duration           `json:"average_chat_duration"`
	StressLevelDescriptor PredefinedStressLevel   `json:"stress_level_descriptor"`
}

type GetChatReportHandler struct {
	deps GetChatReportHandlerDependencies
}

func NewGetChatReportHandler(deps GetChatReportHandlerDependencies) *GetChatReportHandler {
	return &GetChatReportHandler{deps: deps}
}

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

	var localStartTime, localEndTime time.Time
	if datetimeGte != "" {
		localStartTime, err = time.Parse(time.RFC3339, datetimeGte)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				http.NewError(ctx, err, "Invalid datetime_gte format. Use YYYY-MM-DDTHH:mm:ssZ"),
			)
		}
	}
	if datetimeLte != "" {
		localEndTime, err = time.Parse(time.RFC3339, datetimeLte)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(
				http.NewError(ctx, err, "Invalid datetime_lte format. Use YYYY-MM-DDTHH:mm:ssZ"),
			)
		}
	}

	chats := []domain.Chat{}
	if err := h.deps.DB.NewSelect().
		Model(&chats).
		Relation("Summary").
		Where("user_id = ?", userID).
		Where("created_at >= ?", localStartTime).
		Where("created_at <= ?", localEndTime).
		Scan(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to get chats"),
		)
	}

	emotionCounts := make(map[common.Emotion]int64)
	var totalPositiveScore, totalNegativeScore, totalNeutralScore float64
	var totalChatDuration time.Duration

	for _, chat := range chats {
		if chat.Summary == nil {
			continue
		}

		dominantEmotion := chat.Summary.Emotions.GetDominant()
		emotionCounts[dominantEmotion]++

		totalPositiveScore += chat.Summary.PositiveScore
		totalNegativeScore += chat.Summary.NegativeScore
		totalNeutralScore += chat.Summary.NeutralScore

		if chat.ChatDuration.Valid {
			totalChatDuration += time.Duration(chat.ChatDuration.Duration)
		}
	}

	chatCount := len(chats)
	if chatCount == 0 {
		return c.Status(fiber.StatusOK).JSON(GetChatReportHandlerResponse{
			DatetimeGte: localStartTime,
			DatetimeLte: localEndTime,
		})
	}

	avgPositiveScore := utils.RoundTo(totalPositiveScore/float64(chatCount), 3)
	avgNegativeScore := utils.RoundTo(totalNegativeScore/float64(chatCount), 3)
	avgNeutralScore := utils.RoundTo(totalNeutralScore/float64(chatCount), 3)

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

	avgDurationSeconds := int(math.Ceil(totalChatDuration.Seconds()))
	avgDurationMinutes := time.Duration(avgDurationSeconds/60+1) * time.Minute

	var stressLevel PredefinedStressLevel
	for _, level := range PredefinedStressLevels {
		if avgNegativeScore >= level.Threshold {
			stressLevel = level
			break
		}
	}

	emotionCountList := make(common.EmotionCountList, 0, len(emotionCounts))
	for emotion, count := range emotionCounts {
		emotionCountList = append(emotionCountList, common.EmotionCount{
			Emotion: emotion,
			Count:   count,
		})
	}

	return c.Status(fiber.StatusOK).JSON(GetChatReportHandlerResponse{
		DatetimeGte:           localStartTime,
		DatetimeLte:           localEndTime,
		EmotionCounts:         emotionCountList,
		AveragePositiveScore:  avgPositiveScore,
		AverageNegativeScore:  avgNegativeScore,
		AverageNeutralScore:   avgNeutralScore,
		TotalChatCount:        chatCount,
		AverageChatDuration:   avgDurationMinutes,
		StressLevelDescriptor: stressLevel,
	})
}

func (h *GetChatReportHandler) Identify() string {
	return "get-chat-report"
}

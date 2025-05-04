package schema

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type GetSupportedEmotionHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type GetSupportedEmotionResponse struct {
	SupportedEmotions common.EmotionList `json:"supported_emotions"`
}

type GetSupportedEmotionHandler struct {
	deps GetSupportedEmotionHandlerDependencies
}

func NewGetSupportedEmotionHandler(deps GetSupportedEmotionHandlerDependencies) (*GetSupportedEmotionHandler, error) {
	return &GetSupportedEmotionHandler{deps: deps}, nil
}

// @ID GetSupportedEmotions
// @Summary Get supported emotions
// @Description Get supported emotions
// @Tags schema
// @Accept json
// @Produce json
// @Success 200 {object} GetSupportedEmotionResponse
// @Router /_schema/supported-emotions [get]
func (h *GetSupportedEmotionHandler) Handle(c *fiber.Ctx) error {
	return c.JSON(GetSupportedEmotionResponse{
		SupportedEmotions: common.SupportedEmotions,
	})
}

func (h *GetSupportedEmotionHandler) Identify() string {
	return "get-supported-emotions"
}

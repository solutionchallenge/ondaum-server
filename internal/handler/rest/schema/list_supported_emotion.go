package schema

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type ListSupportedEmotionHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type ListSupportedEmotionHandler struct {
	deps ListSupportedEmotionHandlerDependencies
}

func NewListSupportedEmotionHandler(deps ListSupportedEmotionHandlerDependencies) (*ListSupportedEmotionHandler, error) {
	return &ListSupportedEmotionHandler{deps: deps}, nil
}

// @ID ListSupportedEmotion
// @Summary List supported emotions
// @Description List supported emotions
// @Tags schema
// @Accept json
// @Produce json
// @Success 200 {object} common.EmotionList
// @Router /_schema/supported-emotions [get]
func (h *ListSupportedEmotionHandler) Handle(c *fiber.Ctx) error {
	return c.JSON(common.SupportedEmotions)
}

func (h *ListSupportedEmotionHandler) Identify() string {
	return "list-supported-emotions"
}

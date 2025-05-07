package schema

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"go.uber.org/fx"
)

type ListSupportedFeatureHandlerDependencies struct {
	fx.In
}

type ListSupportedFeatureHandler struct {
	deps ListSupportedFeatureHandlerDependencies
}

func NewListSupportedFeatureHandler(deps ListSupportedFeatureHandlerDependencies) (*ListSupportedFeatureHandler, error) {
	return &ListSupportedFeatureHandler{deps: deps}, nil
}

// @ID ListSupportedFeature
// @Summary List supported features
// @Description List supported features
// @Tags schema
// @Accept json
// @Produce json
// @Success 200 {object} common.FeatureList
// @Router /_schema/supported-features [get]
func (h *ListSupportedFeatureHandler) Handle(c *fiber.Ctx) error {
	return c.JSON(common.SupportedFeatures)
}

func (h *ListSupportedFeatureHandler) Identify() string {
	return "list-supported-feature"
}

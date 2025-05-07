package schema

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"go.uber.org/fx"
)

type ListSupportedInspectionHandlerDependencies struct {
	fx.In
}

type ListSupportedInspectionHandler struct {
	deps ListSupportedInspectionHandlerDependencies
}

func NewListSupportedInspectionHandler(deps ListSupportedInspectionHandlerDependencies) (*ListSupportedInspectionHandler, error) {
	return &ListSupportedInspectionHandler{deps: deps}, nil
}

// @ID ListSupportedTest
// @Summary List supported tests
// @Description List supported tests
// @Tags schema
// @Accept json
// @Produce json
// @Success 200 {object} common.InspectionList
// @Router /_schema/supported-inspections [get]
func (h *ListSupportedInspectionHandler) Handle(c *fiber.Ctx) error {
	return c.JSON(common.SupportedInspections)
}

func (h *ListSupportedInspectionHandler) Identify() string {
	return "list-supported-inspection"
}

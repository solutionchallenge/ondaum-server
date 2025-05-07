package schema

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"go.uber.org/fx"
)

type ListSupportedDiagnosisHandlerDependencies struct {
	fx.In
}

type ListSupportedDiagnosisHandler struct {
	deps ListSupportedDiagnosisHandlerDependencies
}

func NewListSupportedDiagnosisHandler(deps ListSupportedDiagnosisHandlerDependencies) (*ListSupportedDiagnosisHandler, error) {
	return &ListSupportedDiagnosisHandler{deps: deps}, nil
}

// @ID ListSupportedDiagnosis
// @Summary List supported diagnoses
// @Description List supported diagnoses
// @Tags schema
// @Accept json
// @Produce json
// @Success 200 {object} common.DiagnosisList
// @Router /_schema/supported-diagnoses [get]
func (h *ListSupportedDiagnosisHandler) Handle(c *fiber.Ctx) error {
	return c.JSON(common.SupportedDiagnoses)
}

func (h *ListSupportedDiagnosisHandler) Identify() string {
	return "list-supported-diagnoses"
}

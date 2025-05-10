package diagnosis

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"go.uber.org/fx"
)

type ListDiagnosisPaperHandlerDependencies struct {
	fx.In
}

type ListDiagnosisPaperHandler struct {
	deps ListDiagnosisPaperHandlerDependencies
}

type ListDiagnosisPaperHandlerResponse struct {
	ID          common.Diagnosis `json:"id"`
	Description string           `json:"description"`
}

func NewListDiagnosisPaperHandler(deps ListDiagnosisPaperHandlerDependencies) (*ListDiagnosisPaperHandler, error) {
	return &ListDiagnosisPaperHandler{deps: deps}, nil
}

// @ID ListDiagnosisPaper
// @Summary List diagnosis papers
// @Description List diagnosis papers
// @Tags diagnosis
// @Accept json
// @Produce json
// @Success 200 {array} ListDiagnosisPaperHandlerResponse
// @Router /diagnoses/papers [get]
// @Security BearerAuth
func (h *ListDiagnosisPaperHandler) Handle(c *fiber.Ctx) error {
	diagnoses := utils.Map(common.SupportedDiagnoses, func(diagnosis common.Diagnosis) ListDiagnosisPaperHandlerResponse {
		return ListDiagnosisPaperHandlerResponse{
			ID:          diagnosis,
			Description: common.DiagnosisDescriptions[diagnosis],
		}
	})
	return c.JSON(diagnoses)
}

func (h *ListDiagnosisPaperHandler) Identify() string {
	return "list-diagnosis-paper"
}

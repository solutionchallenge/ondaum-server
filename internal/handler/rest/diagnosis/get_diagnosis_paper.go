package diagnosis

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"go.uber.org/fx"
)

type GetDiagnosisPaperHandlerDependencies struct {
	fx.In
}

type GetDiagnosisPaperHandler struct {
	deps GetDiagnosisPaperHandlerDependencies
}

func NewGetDiagnosisPaperHandler(deps GetDiagnosisPaperHandlerDependencies) (*GetDiagnosisPaperHandler, error) {
	return &GetDiagnosisPaperHandler{deps: deps}, nil
}

// @ID GetDiagnosisPaper
// @Summary Get diagnosis paper
// @Description Get diagnosis paper as JSON format
// @Accept json
// @Produce json
// @Param diagnosis_id path string true "Diagnosis ID"
// @Success 200 {object} common.DiagnosisPaper
// @Failure 404 {object} http.Error
// @Failure 500 {object} http.Error
// @Router /diagnosis/:diagnosis_id [get]
func (h *GetDiagnosisPaperHandler) Handle(c *fiber.Ctx) error {
	identifier := c.Params("diagnosis_id")
	switch common.Diagnosis(identifier) {
	case common.DiagnosisPHQ9:
		diagnosisPaper, err := common.ReadDiagnosisPaperFrom("resource/diagnosis/phq-9-en.json")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				http.NewError(c.UserContext(), err, "Failed to read diagnosis paper"),
			)
		}
		return c.JSON(diagnosisPaper)
	case common.DiagnosisGAD7:
		diagnosisPaper, err := common.ReadDiagnosisPaperFrom("resource/diagnosis/gad-7-en.json")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				http.NewError(c.UserContext(), err, "Failed to read diagnosis paper"),
			)
		}
		return c.JSON(diagnosisPaper)
	case common.DiagnosisPSS:
		diagnosisPaper, err := common.ReadDiagnosisPaperFrom("resource/diagnosis/pss-en.json")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				http.NewError(c.UserContext(), err, "Failed to read diagnosis paper"),
			)
		}
		return c.JSON(diagnosisPaper)
	default:
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(c.UserContext(), nil, "Diagnosis not found"),
		)
	}
}

func (h *GetDiagnosisPaperHandler) Identify() string {
	return "get-diagnosis-paper"
}

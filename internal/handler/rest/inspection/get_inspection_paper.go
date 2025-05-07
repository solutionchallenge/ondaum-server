package inspection

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"go.uber.org/fx"
)

type GetInspectionPaperHandlerDependencies struct {
	fx.In
}

type GetInspectionPaperHandler struct {
	deps GetInspectionPaperHandlerDependencies
}

func NewGetInspectionPaperHandler(deps GetInspectionPaperHandlerDependencies) (*GetInspectionPaperHandler, error) {
	return &GetInspectionPaperHandler{deps: deps}, nil
}

// @ID GetInspectionPaper
// @Summary Get inspection paper
// @Description Get inspection paper as JSON format
// @Accept json
// @Produce json
// @Param inspection_id path string true "Inspection ID"
// @Success 200 {object} common.InspectionPaper
// @Failure 404 {object} http.Error
// @Failure 500 {object} http.Error
// @Router /inspection/:inspection_id [get]
func (h *GetInspectionPaperHandler) Handle(c *fiber.Ctx) error {
	identifier := c.Params("inspection_id")
	switch common.Inspection(identifier) {
	case common.InspectionPHQ9:
		testPaper, err := common.ReadInspectionPaperFrom("resource/inspection/phq-9-en.json")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				http.NewError(c.UserContext(), err, "Failed to read test paper"),
			)
		}
		return c.JSON(testPaper)
	case common.InspectionGAD7:
		testPaper, err := common.ReadInspectionPaperFrom("resource/inspection/gad-7-en.json")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				http.NewError(c.UserContext(), err, "Failed to read test paper"),
			)
		}
		return c.JSON(testPaper)
	case common.InspectionPSS:
		testPaper, err := common.ReadInspectionPaperFrom("resource/inspection/pss-en.json")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				http.NewError(c.UserContext(), err, "Failed to read test paper"),
			)
		}
		return c.JSON(testPaper)
	default:
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(c.UserContext(), nil, "Test not found"),
		)
	}
}

func (h *GetInspectionPaperHandler) Identify() string {
	return "get-inspection-paper"
}

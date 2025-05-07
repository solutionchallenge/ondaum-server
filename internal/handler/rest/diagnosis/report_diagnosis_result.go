package diagnosis

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/diagnosis"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type ReportDiagnosisResultHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type ReportDiagnosisResultHandlerResponse struct {
	Success bool  `json:"success"`
	ID      int64 `json:"id"`
}

type ReportDiagnosisResultHandler struct {
	deps ReportDiagnosisResultHandlerDependencies
}

func NewReportDiagnosisResultHandler(deps ReportDiagnosisResultHandlerDependencies) (*ReportDiagnosisResultHandler, error) {
	return &ReportDiagnosisResultHandler{deps: deps}, nil
}

// @ID ReportDiagnosisResult
// @Summary Report diagnosis result
// @Description Report diagnosis result
// @Accept json
// @Produce json
// @Param request body diagnosis.DiagnosisDTO true "Diagnosis result"
// @Success 200 {object} ReportDiagnosisResultHandlerResponse
// @Failure 401 {object} http.Error
// @Failure 404 {object} http.Error
// @Failure 500 {object} http.Error
// @Router /diagnoses [post]
// @Security BearerAuth
func (h *ReportDiagnosisResultHandler) Handle(c *fiber.Ctx) error {
	userID, err := http.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			http.NewError(c.UserContext(), err, "Unauthorized"),
		)
	}
	user := &user.User{ID: userID}
	if err := h.deps.DB.NewSelect().Model(user).Where("id = ?", userID).Scan(c.UserContext()); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(c.UserContext(), err, "User not found"),
		)
	}

	var request diagnosis.DiagnosisDTO
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), err, "Invalid request"),
		)
	}

	diagnosis := &diagnosis.Diagnosis{
		UserID:            userID,
		Diagnosis:         request.Diagnosis,
		TotalScore:        request.TotalScore,
		ResultScore:       request.ResultScore,
		ResultName:        request.ResultName,
		ResultDescription: request.ResultDescription,
		ResultCritical:    request.ResultCritical,
	}
	result, err := h.deps.DB.NewInsert().Model(diagnosis).Exec(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to report diagnosis result"),
		)
	}

	inserted, err := result.LastInsertId()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to report diagnosis result"),
		)
	}
	response := ReportDiagnosisResultHandlerResponse{
		Success: true,
		ID:      inserted,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ReportDiagnosisResultHandler) Identify() string {
	return "report-diagnosis-result"
}

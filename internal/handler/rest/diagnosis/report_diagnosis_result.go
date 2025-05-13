package diagnosis

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	Success bool   `json:"success"`
	ID      string `json:"id"`
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
// @Tags diagnosis
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

	var request diagnosis.DiagnosisDTO
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(ctx, err, "Invalid request"),
		)
	}

	subID := uuid.New().String()
	diagnosis := &diagnosis.Diagnosis{
		UserID:            userID,
		SubID:             subID,
		Diagnosis:         request.Diagnosis,
		TotalScore:        request.TotalScore,
		ResultScore:       request.ResultScore,
		ResultName:        request.ResultName,
		ResultDescription: request.ResultDescription,
		ResultCritical:    request.ResultCritical,
	}
	_, err = h.deps.DB.NewInsert().Model(diagnosis).Exec(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to report diagnosis result"),
		)
	}

	response := ReportDiagnosisResultHandlerResponse{
		Success: true,
		ID:      subID,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ReportDiagnosisResultHandler) Identify() string {
	return "report-diagnosis-result"
}

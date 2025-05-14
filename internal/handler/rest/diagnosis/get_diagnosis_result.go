package diagnosis

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/diagnosis"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type GetDiagnosisResultHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type GetDiagnosisResultHandler struct {
	deps GetDiagnosisResultHandlerDependencies
}

func NewGetDiagnosisResultHandler(deps GetDiagnosisResultHandlerDependencies) (*GetDiagnosisResultHandler, error) {
	return &GetDiagnosisResultHandler{deps: deps}, nil
}

// @ID GetDiagnosisResult
// @Summary Get diagnosis result
// @Description Get diagnosis result
// @Tags diagnosis
// @Accept json
// @Produce json
// @Param diagnosis_id path string true "Diagnosis ID"
// @Success 200 {object} diagnosis.DiagnosisDTO
// @Failure 401 {object} http.Error
// @Failure 404 {object} http.Error
// @Failure 500 {object} http.Error
// @Router /diagnoses/{diagnosis_id} [get]
// @Security BearerAuth
func (h *GetDiagnosisResultHandler) Handle(c *fiber.Ctx) error {
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

	diagnosisID := c.Params("diagnosis_id")
	if diagnosisID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(ctx, errors.New("diagnosis_id is required"), "Bad Request"),
		)
	}

	diagnosis := &diagnosis.Diagnosis{}
	if err := h.deps.DB.NewSelect().
		Model(diagnosis).
		Where("user_id = ?", userID).
		Where("sub_id = ?", diagnosisID).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(
				http.NewError(ctx, err, "Diagnosis not found for id: %v", diagnosisID),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to get diagnosis for id: %v", diagnosisID),
		)
	}

	return c.Status(fiber.StatusOK).JSON(diagnosis.ToDiagnosisDTOWithSubID())
}

func (h *GetDiagnosisResultHandler) Identify() string {
	return "get-diagnosis-result"
}

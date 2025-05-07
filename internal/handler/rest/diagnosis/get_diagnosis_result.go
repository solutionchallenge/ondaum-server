package diagnosis

import (
	"errors"
	"strconv"

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
// @Accept json
// @Produce json
// @Param diagnosis_id path string true "Diagnosis ID"
// @Success 200 {object} diagnosis.DiagnosisDTO
// @Failure 401 {object} http.Error
// @Failure 404 {object} http.Error
// @Failure 500 {object} http.Error
// @Router /diagnosis-papers/{diagnosis_id} [get]
// @Security BearerAuth
func (h *GetDiagnosisResultHandler) Handle(c *fiber.Ctx) error {
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

	diagnosisID := c.Params("diagnosis_id")
	if diagnosisID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), errors.New("diagnosis_id is required"), "Bad Request"),
		)
	}

	convertedID, err := strconv.ParseInt(diagnosisID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), err, "Invalid diagnosis_id"),
		)
	}

	diagnosis := &diagnosis.Diagnosis{}
	if err := h.deps.DB.NewSelect().
		Model(diagnosis).
		Where("id = ?", convertedID).
		Where("user_id = ?", userID).
		Scan(c.UserContext()); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(c.UserContext(), err, "Diagnosis not found"),
		)
	}

	return c.Status(fiber.StatusOK).JSON(diagnosis.ToDiagnosisDTO())
}

func (h *GetDiagnosisResultHandler) Identify() string {
	return "get-diagnosis-result"
}

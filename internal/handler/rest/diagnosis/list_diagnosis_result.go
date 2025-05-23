package diagnosis

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/diagnosis"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type ListDiagnosisResultHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type ListDiagnosisResultHandler struct {
	deps ListDiagnosisResultHandlerDependencies
}

func NewListDiagnosisResultHandler(deps ListDiagnosisResultHandlerDependencies) (*ListDiagnosisResultHandler, error) {
	return &ListDiagnosisResultHandler{deps: deps}, nil
}

// @ID ListDiagnosisResult
// @Summary List diagnosis result
// @Description List diagnosis result
// @Tags diagnosis
// @Accept json
// @Produce json
// @Success 200 {object} diagnosis.DiagnosisDTOWithSubID
// @Failure 401 {object} http.Error
// @Failure 404 {object} http.Error
// @Failure 500 {object} http.Error
// @Router /diagnoses [get]
// @Security BearerAuth
func (h *ListDiagnosisResultHandler) Handle(c *fiber.Ctx) error {
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

	diagnoses := []*diagnosis.Diagnosis{}
	if err := h.deps.DB.NewSelect().
		Model(&diagnoses).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Scan(ctx); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(ctx, err, "Diagnosis not found"),
		)
	}

	dtos := utils.Map(diagnoses, func(d *diagnosis.Diagnosis) diagnosis.DiagnosisDTOWithSubID {
		return d.ToDiagnosisDTOWithSubID()
	})
	return c.Status(fiber.StatusOK).JSON(dtos)
}

func (h *ListDiagnosisResultHandler) Identify() string {
	return "list-diagnosis-result"
}

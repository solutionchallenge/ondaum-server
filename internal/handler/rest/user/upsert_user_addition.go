package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type UpsertUserAdditionHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type UpsertUserAdditionHandlerRequest struct {
	Concerns []string           `json:"concerns"`
	Emotions common.EmotionList `json:"emotions"`
}

type UpsertUserAdditionHandlerResponse struct {
	Success bool `json:"success"`
	Created bool `json:"created"`
}

type UpsertUserAdditionHandler struct {
	deps UpsertUserAdditionHandlerDependencies
}

func NewUpsertUserAdditionHandler(deps UpsertUserAdditionHandlerDependencies) (*UpsertUserAdditionHandler, error) {
	return &UpsertUserAdditionHandler{deps: deps}, nil
}

// @ID UpsertUserAddition
// @Summary      Update or Create User Additional Information
// @Description  Updates or creates the user's additional information including concerns and emotions.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param request body UpsertUserAdditionHandlerRequest true "User additional information"
// @Success      200 {object} UpsertUserAdditionHandlerResponse
// @Failure      400 {object} http.Error
// @Failure      401 {object} http.Error
// @Failure      500 {object} http.Error
// @Router       /user/addition [put]
// @Security     BearerAuth
func (h *UpsertUserAdditionHandler) Handle(c *fiber.Ctx) error {
	userID, err := http.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			http.NewError(c.UserContext(), err, "Unauthorized"),
		)
	}

	request := &UpsertUserAdditionHandlerRequest{}
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), err, "Invalid request"),
		)
	}

	if !request.Emotions.Validate() {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), nil, "Invalid emotions contained in request"),
		)
	}

	addition := &user.Addition{
		UserID:   userID,
		Concerns: request.Concerns,
		Emotions: request.Emotions,
	}

	result, err := h.deps.DB.NewInsert().
		Model(addition).
		On("DUPLICATE KEY UPDATE").
		Set("concerns = ?", request.Concerns).
		Set("emotions = ?", request.Emotions.ToString()).
		Set("updated_at = CURRENT_TIMESTAMP").
		Exec(c.UserContext())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to upsert user addition"),
		)
	}

	rowsAffected, _ := result.RowsAffected()
	return c.JSON(UpsertUserAdditionHandlerResponse{
		Success: true,
		Created: rowsAffected == 1,
	})
}

func (h *UpsertUserAdditionHandler) Identify() string {
	return "upsert-user-addition"
}

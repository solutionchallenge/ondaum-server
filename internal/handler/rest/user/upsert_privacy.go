package user

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type UpsertPrivacyHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type UpsertPrivacyHandlerRequest struct {
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
}

type UpsertPrivacyHandlerResponse struct {
	Success bool `json:"success"`
	Created bool `json:"created"`
}

type UpsertPrivacyHandler struct {
	deps UpsertPrivacyHandlerDependencies
}

func NewUpsertPrivacyHandler(deps UpsertPrivacyHandlerDependencies) (*UpsertPrivacyHandler, error) {
	return &UpsertPrivacyHandler{deps: deps}, nil
}

// @ID UpsertUserPrivacy
// @Summary      Update or Create User Privacy Information
// @Description  Updates or creates the user's privacy information including gender and birthday.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param request body UpsertPrivacyHandlerRequest true "User privacy information"
// @Success      200 {object} UpsertPrivacyHandlerResponse
// @Failure      400 {object} http.Error
// @Failure      401 {object} http.Error
// @Failure      500 {object} http.Error
// @Router       /user/privacy [put]
// @Security     BearerAuth
func (h *UpsertPrivacyHandler) Handle(c *fiber.Ctx) error {
	userID, err := http.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			http.NewError(c.UserContext(), err, "Unauthorized"),
		)
	}

	request := &UpsertPrivacyHandlerRequest{}
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), err, "Invalid request"),
		)
	}

	// Parse birthday string to time.Time
	birthday, err := time.Parse(utils.TIME_FORMAT_DATE, request.Birthday)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), err, "Invalid birthday format. Use YYYY-MM-DD"),
		)
	}

	// Validate gender
	gender := user.UserGender(request.Gender)
	if gender != user.UserGenderMale && gender != user.UserGenderFemale && gender != user.UserGenderOther {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), nil, "Invalid gender. Must be one of: male, female, other"),
		)
	}

	privacy := &user.Privacy{
		UserID:   userID,
		Gender:   gender,
		Birthday: birthday,
	}

	result, err := h.deps.DB.NewInsert().
		Model(privacy).
		On("DUPLICATE KEY UPDATE").
		Set("gender = ?", gender).
		Set("birthday = ?", birthday).
		Set("updated_at = CURRENT_TIMESTAMP").
		Exec(c.UserContext())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to upsert user privacy"),
		)
	}

	rowsAffected, _ := result.RowsAffected()
	return c.JSON(UpsertPrivacyHandlerResponse{
		Success: true,
		Created: rowsAffected == 1,
	})
}

func (h *UpsertPrivacyHandler) Identify() string {
	return "upsert-privacy"
}

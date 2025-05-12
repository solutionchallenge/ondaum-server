package user

import (
	"time"

	"github.com/gofiber/fiber/v2"
	domain "github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type UpsertUserPrivacyHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type UpsertUserPrivacyHandlerRequest struct {
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
}

type UpsertUserPrivacyHandlerResponse struct {
	Success bool `json:"success"`
	Created bool `json:"created"`
}

type UpsertUserPrivacyHandler struct {
	deps UpsertUserPrivacyHandlerDependencies
}

func NewUpsertUserPrivacyHandler(deps UpsertUserPrivacyHandlerDependencies) (*UpsertUserPrivacyHandler, error) {
	return &UpsertUserPrivacyHandler{deps: deps}, nil
}

// @ID UpsertUserPrivacy
// @Summary      Update or Create User Privacy Information
// @Description  Updates or creates the user's privacy information including gender and birthday.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param request body UpsertUserPrivacyHandlerRequest true "User privacy information"
// @Success      200 {object} UpsertUserPrivacyHandlerResponse
// @Failure      400 {object} http.Error
// @Failure      401 {object} http.Error
// @Failure      500 {object} http.Error
// @Router       /user/privacy [put]
// @Security     BearerAuth
func (h *UpsertUserPrivacyHandler) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()
	userID, err := http.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			http.NewError(ctx, err, "Unauthorized"),
		)
	}
	user := &domain.User{ID: userID}
	if err := h.deps.DB.NewSelect().Model(user).Where("id = ?", userID).Scan(ctx); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(ctx, err, "User not found"),
		)
	}

	request := &UpsertUserPrivacyHandlerRequest{}
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(ctx, err, "Invalid request"),
		)
	}

	birthday, err := time.Parse(utils.TIME_FORMAT_DATE, request.Birthday)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(ctx, err, "Invalid birthday format. Use YYYY-MM-DD"),
		)
	}

	// Validate gender
	gender := domain.UserGender(request.Gender)
	if gender != domain.UserGenderMale && gender != domain.UserGenderFemale && gender != domain.UserGenderOther {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(ctx, nil, "Invalid gender. Must be one of: male, female, other"),
		)
	}

	privacy := &domain.Privacy{
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
		Exec(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to upsert user privacy"),
		)
	}

	rowsAffected, _ := result.RowsAffected()
	return c.JSON(UpsertUserPrivacyHandlerResponse{
		Success: true,
		Created: rowsAffected == 1,
	})
}

func (h *UpsertUserPrivacyHandler) Identify() string {
	return "upsert-user-privacy"
}

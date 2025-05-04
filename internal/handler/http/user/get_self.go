package user

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type GetSelfHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type GetSelfHandlerPrivacyPartial struct {
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
}

type GetSelfHandlerAdditionPartial struct {
	Concerns []string           `json:"concerns"`
	Emotions common.EmotionList `json:"emotions"`
}

type GetSelfHandlerResponse struct {
	ID       int64                          `json:"id"`
	Email    string                         `json:"email"`
	Username string                         `json:"username"`
	Privacy  *GetSelfHandlerPrivacyPartial  `json:"privacy,omitempty"`
	Addition *GetSelfHandlerAdditionPartial `json:"addition,omitempty"`
}

type GetSelfHandler struct {
	deps GetSelfHandlerDependencies
}

var _ http.Handler = &GetSelfHandler{}

func NewGetSelfHandler(deps GetSelfHandlerDependencies) (*GetSelfHandler, error) {
	return &GetSelfHandler{
		deps: deps,
	}, nil
}

// @ID GetSelfUser
// @Summary      Get Current Authenticated User Profile
// @Description  Returns the authenticated user's profile, including basic information and optional onboarding data.
//   - The `privacy` and `addition` fields are optional and will be null if the user has not completed onboarding.
//   - This API requires a valid Bearer token and returns the profile of the authenticated user only.
//
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200 {object} GetSelfHandlerResponse
// @Failure      401 {object} http.Error
// @Failure      404 {object} http.Error
// @Failure      500 {object} http.Error
// @Router       /user/self [get]
// @Security     BearerAuth
func (h *GetSelfHandler) Handle(c *fiber.Ctx) error {
	id, err := http.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			http.NewError(c.UserContext(), err, "Unauthorized"),
		)
	}
	user := user.User{}
	err = h.deps.DB.NewSelect().
		Model(&user).
		Relation("Addition").
		Relation("Privacy").
		Where("id = ?", id).
		Scan(c.UserContext())
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(
				http.NewError(c.UserContext(), err, "User not found for id: %v", id),
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to get user for id: %v", id),
		)
	}

	response := GetSelfHandlerResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
	if user.Privacy != nil {
		response.Privacy = &GetSelfHandlerPrivacyPartial{
			Gender:   string(user.Privacy.Gender),
			Birthday: user.Privacy.Birthday.Format("2006-01-02"),
		}
	}
	if user.Addition != nil {
		response.Addition = &GetSelfHandlerAdditionPartial{
			Concerns: user.Addition.Concerns,
			Emotions: user.Addition.Emotions,
		}
	}
	return c.JSON(response)
}

func (h *GetSelfHandler) Identify() string {
	return "get-self"
}

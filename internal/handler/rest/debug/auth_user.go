package debug

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/solutionchallenge/ondaum-server/pkg/jwt"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type AuthUserHandlerDependencies struct {
	fx.In
	DB  *bun.DB
	JWT *jwt.Generator
}

type AuthUserHandlerRequest struct {
	ID int64 `json:"id"`
}

type AuthUserHandlerResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthUserHandler struct {
	deps AuthUserHandlerDependencies
}

func NewAuthUserHandler(deps AuthUserHandlerDependencies) (*AuthUserHandler, error) {
	return &AuthUserHandler{deps: deps}, nil
}

// Must not be documented. (Debugging purpose only!)
func (h *AuthUserHandler) Handle(c *fiber.Ctx) error {
	if os.Getenv("FLAG_DEBUGGING_FEATURES_ENABLED") != "true" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	request := &AuthUserHandlerRequest{}
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), err, "Invalid request"),
		)
	}

	user := &user.User{}
	err := h.deps.DB.NewSelect().
		Model(user).
		Where("id = ?", request.ID).
		Scan(c.UserContext())

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			http.NewError(c.UserContext(), err, "User not found"),
		)
	}

	tokenPair, err := h.deps.JWT.GenerateTokenPair(strconv.FormatInt(user.ID, 10))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to generate token pair"),
		)
	}

	response := AuthUserHandlerResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}
	return c.JSON(response)
}

func (h *AuthUserHandler) Identify() string {
	return "auth-user"
}

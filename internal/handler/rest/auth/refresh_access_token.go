package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/solutionchallenge/ondaum-server/pkg/jwt"
	"go.uber.org/fx"
)

type RefreshAccessTokenHandlerDependencies struct {
	fx.In
	JWT *jwt.Generator
}

type RefreshAccessTokenHandlerRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshAccessTokenHandlerResponse struct {
	AccessToken string `json:"access_token"`
}

type RefreshAccessTokenHandler struct {
	deps RefreshAccessTokenHandlerDependencies
}

func NewRefreshAccessTokenHandler(deps RefreshAccessTokenHandlerDependencies) (*RefreshAccessTokenHandler, error) {
	return &RefreshAccessTokenHandler{deps: deps}, nil
}

// @ID RefreshAccessToken
// @Summary Refresh access token
// @Description Refresh access token
// @Accept json
// @Produce json
// @Param request body RefreshAccessTokenHandlerRequest true "Refresh token"
// @Success 200 {object} RefreshAccessTokenHandlerResponse
// @Failure 401 {object} http.Error
// @Failure 404 {object} http.Error
// @Failure 500 {object} http.Error
// @Router /auth/refresh [post]
func (h *RefreshAccessTokenHandler) Handle(c *fiber.Ctx) error {
	var request RefreshAccessTokenHandlerRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	tokenType, err := h.deps.JWT.GetTokenType(request.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			http.NewError(c.UserContext(), err, "Invalid refresh token"),
		)
	}

	if tokenType != jwt.RefreshTokenType {
		return c.Status(fiber.StatusUnauthorized).JSON(
			http.NewError(c.UserContext(), errors.New("invalid token type"), "Invalid token type"),
		)
	}

	tokenPair, err := h.deps.JWT.RefreshTokenPair(request.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			http.NewError(c.UserContext(), err, "Failed to refresh token pair"),
		)
	}

	response := RefreshAccessTokenHandlerResponse{
		AccessToken: tokenPair.AccessToken,
	}

	return c.JSON(response)
}

func (h *RefreshAccessTokenHandler) Identify() string {
	return "refresh-access-token"
}

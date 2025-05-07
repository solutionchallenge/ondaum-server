package oauth

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/solutionchallenge/ondaum-server/pkg/jwt"

	"github.com/solutionchallenge/ondaum-server/pkg/oauth"
	"github.com/solutionchallenge/ondaum-server/pkg/oauth/google"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type AuthGoogleHandlerDependencies struct {
	fx.In
	DB    *bun.DB
	OAuth *oauth.Container
	JWT   *jwt.Generator
}

type AuthGoogleHandlerRequest struct {
	Code string `query:"code"`
}

type AuthGoogleHandlerResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthGoogleHandler struct {
	deps AuthGoogleHandlerDependencies
}

var _ http.Handler = &AuthGoogleHandler{}

func NewAuthGoogleHandler(deps AuthGoogleHandlerDependencies) (*AuthGoogleHandler, error) {
	return &AuthGoogleHandler{
		deps: deps,
	}, nil
}

// @ID ExchangeGoogleOAuthCode
// @Summary      Exchange Google OAuth Code for Tokens
// @Description  Receives the authorization code (obtained from Google OAuth) and exchanges it for access and refresh tokens.
// @Tags         oauth
// @Accept       json
// @Produce      json
// @Param request body AuthGoogleHandlerRequest true "Payload containing the authorization code received from Google OAuth"
// @Param redirect query string false "Redirect URI (optional, the client's callback URL where Google was redirect with the code)"
// @Success      200 {object} AuthGoogleHandlerResponse
// @Failure      400 {object} http.Error
// @Failure      500 {object} http.Error
// @Router       /oauth/google/auth [post]
func (h *AuthGoogleHandler) Handle(c *fiber.Ctx) error {
	redirectURI := c.Query("redirect")

	request := &AuthGoogleHandlerRequest{}
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), err, "Invalid request"),
		)
	}
	if request.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), nil, "Code is required"),
		)
	}

	var userInfo oauth.UserInfoOutput
	var err error
	if redirectURI != "" {
		userInfo, err = h.deps.OAuth.Use(google.Provider).GetUserInfo(request.Code, redirectURI)
	} else {
		userInfo, err = h.deps.OAuth.Use(google.Provider).GetUserInfo(request.Code)
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to get user info"),
		)
	}

	dbTx, err := h.deps.DB.BeginTx(c.Context(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to begin transaction"),
		)
	}
	defer dbTx.Rollback()

	userData := &user.User{
		Email:    userInfo.Email,
		Username: userInfo.Name,
	}
	_, err = dbTx.NewInsert().
		Model(userData).
		On("DUPLICATE KEY UPDATE").
		Set("username = VALUES(username)").
		Set("updated_at = CURRENT_TIMESTAMP").
		Exec(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to upsert user"),
		)
	}

	err = dbTx.NewSelect().
		Model(userData).
		Where("email = ?", userInfo.Email).
		Scan(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to get user"),
		)
	}

	userOAuth := &user.OAuth{
		UserID:       userData.ID,
		Provider:     google.Provider,
		ProviderCode: request.Code,
	}
	_, err = dbTx.NewInsert().
		Model(userOAuth).
		On("DUPLICATE KEY UPDATE").
		Set("provider_code = ?", request.Code).
		Set("updated_at = CURRENT_TIMESTAMP").
		Exec(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to create user oauth"),
		)
	}

	if err := dbTx.Commit(); err != nil {
		return c.Status(fiber.StatusConflict).JSON(
			http.NewError(c.UserContext(), err, "Failed to commit transaction"),
		)
	}

	tokenPair, err := h.deps.JWT.GenerateTokenPair(strconv.FormatInt(userData.ID, 10))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Failed to generate token pair"),
		)
	}

	response := AuthGoogleHandlerResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}
	return c.JSON(response)
}

func (h *AuthGoogleHandler) Identify() string {
	return "auth-google"
}

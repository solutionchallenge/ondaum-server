package debug

import (
	"database/sql"
	"os"
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

type OAuthCallbackHandlerDependencies struct {
	fx.In
	DB    *bun.DB
	OAuth *oauth.Container
	JWT   *jwt.Generator
}

type OAuthCallbackHandlerResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type OAuthCallbackHandler struct {
	deps OAuthCallbackHandlerDependencies
}

var _ http.Handler = &OAuthCallbackHandler{}

func NewOAuthCallbackHandler(deps OAuthCallbackHandlerDependencies) (*OAuthCallbackHandler, error) {
	return &OAuthCallbackHandler{
		deps: deps,
	}, nil
}

func (h *OAuthCallbackHandler) Handle(c *fiber.Ctx) error {
	if os.Getenv("FLAG_DEBUGGING_FEATURES_ENABLED") != "true" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Code is required",
		})
	}

	userInfo, err := h.deps.OAuth.Use(google.Provider).GetUserInfo(code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	dbTx, err := h.deps.DB.BeginTx(c.Context(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to begin transaction",
		})
	}
	defer dbTx.Rollback()

	foundUser := &user.User{
		Email: userInfo.Email,
	}
	err = dbTx.NewSelect().
		Model(foundUser).
		Where("email = ?", userInfo.Email).
		Scan(c.Context())
	if err != nil {
		if err == sql.ErrNoRows {
			newUser := &user.User{
				Email:    userInfo.Email,
				Username: userInfo.Name,
			}
			_, err = dbTx.NewInsert().
				Model(newUser).
				Exec(c.Context())
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to create user",
				})
			}
			foundUser = newUser
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get user",
			})
		}
	}

	userOAuth := &user.UserOAuth{
		UserID:       foundUser.ID,
		Provider:     google.Provider,
		ProviderCode: userInfo.ID,
	}
	_, err = dbTx.NewInsert().
		Model(userOAuth).
		On("DUPLICATE KEY UPDATE provider_code = VALUES(provider_code)").
		Exec(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user oauth",
		})
	}

	if err := dbTx.Commit(); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Failed to commit transaction",
		})
	}

	tokenPair, err := h.deps.JWT.GenerateTokenPair(strconv.FormatInt(foundUser.ID, 10))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token pair",
		})
	}

	response := OAuthCallbackHandlerResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}
	return c.JSON(response)
}

func (h *OAuthCallbackHandler) Identify() string {
	return "oauth-callback"
}

package oauth

import (
	"database/sql"
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

// @ID GoogleOAuth
// @Summary      Callback for Google OAuth
// @Description  Must not be called directly!
// @Tags         oauth
// @Accept       json
// @Produce      json
// @Param        code   query      string  true  "Google OAuth Code"
// @Success      200  {object}  AuthGoogleHandlerResponse
// @Failure      400  {object}  string
// @Failure      500  {object}  string
// @Router       /oauth/google/callback [get]
func (h *AuthGoogleHandler) Handle(c *fiber.Ctx) error {
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

	response := AuthGoogleHandlerResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}
	return c.JSON(response)
}

func (h *AuthGoogleHandler) Identify() string {
	return "auth-google"
}

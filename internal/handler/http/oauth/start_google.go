package oauth

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/solutionchallenge/ondaum-server/pkg/oauth"
	"github.com/solutionchallenge/ondaum-server/pkg/oauth/google"
	"go.uber.org/fx"
)

type StartGoogleHandlerDependencies struct {
	fx.In
	OAuth *oauth.Container
}

type StartGoogleHandlerResponse struct {
	AuthURL string `json:"auth_url"`
}

type StartGoogleHandler struct {
	deps StartGoogleHandlerDependencies
}

var _ http.Handler = &StartGoogleHandler{}

func NewStartGoogleHandler(deps StartGoogleHandlerDependencies) (*StartGoogleHandler, error) {
	return &StartGoogleHandler{deps: deps}, nil
}

// @ID StartGoogleOAuth
// @Summary      Get Google OAuth Authorization URL
// @Description  Returns the Google OAuth authorization URL, which includes the specified redirect URI (the URL where Google will send the authorization code after login).
// @Tags         oauth
// @Accept       json
// @Produce      json
// @Param        redirect query string true "Redirect URI (the client's callback URL where Google will redirect with the code)"
// @Success      200 {object} StartGoogleHandlerResponse
// @Failure      400 {object} http.Error
// @Router       /oauth/google/start [get]
func (h *StartGoogleHandler) Handle(c *fiber.Ctx) error {
	requestID := http.GetRequestID(c)
	redirectURI := c.Query("redirect")
	if redirectURI == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), nil, "Redirect URI is required"),
		)
	}
	if !strings.HasPrefix(redirectURI, "http://") && !strings.HasPrefix(redirectURI, "https://") {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), nil, "Redirect URI must be a valid URL starting with http:// or https://"),
		)
	}
	parsedURL, err := url.Parse(redirectURI)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), err, "Invalid redirect URI format"),
		)
	}

	if parsedURL.Host == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(c.UserContext(), nil, "Redirect URI must contain a valid host"),
		)
	}
	authURL, err := h.deps.OAuth.Use(google.Provider).GetAuthURL(requestID, redirectURI)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(c.UserContext(), err, "Unauthorized redirect URI"),
		)
	}

	enc := json.NewEncoder(c.Response().BodyWriter())
	enc.SetEscapeHTML(false)
	return enc.Encode(StartGoogleHandlerResponse{
		AuthURL: authURL,
	})
}

func (h *StartGoogleHandler) Identify() string {
	return "start-google"
}

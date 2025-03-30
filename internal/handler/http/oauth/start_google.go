package oauth

import (
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

type StartGoogleHandler struct {
	deps StartGoogleHandlerDependencies
}

var _ http.Handler = &StartGoogleHandler{}

func NewStartGoogleHandler(deps StartGoogleHandlerDependencies) (*StartGoogleHandler, error) {
	return &StartGoogleHandler{deps: deps}, nil
}

// @ID StartGoogleOAuth
// @Summary      Start Google OAuth
// @Description  This API redirects to Google OAuth, and finally redirects to the callback URL. (GoogleOAuthCallback)
// @Tags         oauth
// @Accept       json
// @Produce      json
// @Response     307  {string}  string
// @Router       /oauth/google/start [get]
func (h *StartGoogleHandler) Handle(c *fiber.Ctx) error {
	requestID := http.GetRequestID(c)
	authURL := h.deps.OAuth.Use(google.Provider).GetAuthURL(requestID)
	return c.Status(fiber.StatusTemporaryRedirect).Redirect(authURL)
}

func (h *StartGoogleHandler) Identify() string {
	return "start-google"
}

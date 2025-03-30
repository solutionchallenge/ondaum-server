package dependency

import (
	"github.com/solutionchallenge/ondaum-server/pkg/oauth"
	"github.com/solutionchallenge/ondaum-server/pkg/oauth/google"
	"go.uber.org/fx"
)

func NewOAuthModule(config oauth.Config) fx.Option {
	return fx.Module("oauth",
		fx.Provide(func() *oauth.Container {
			var clients []oauth.Client
			if config.Google.Enabled {
				clients = append(clients, google.NewClient(config))
			}
			return oauth.NewContainer(clients...)
		}),
	)
}

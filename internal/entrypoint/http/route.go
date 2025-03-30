package http

import (
	"github.com/solutionchallenge/ondaum-server/internal/dependency"
	"github.com/solutionchallenge/ondaum-server/internal/handler/http/oauth"
	"github.com/solutionchallenge/ondaum-server/internal/handler/http/sys"
	"github.com/solutionchallenge/ondaum-server/internal/handler/http/user"
	"go.uber.org/fx"
)

var PredefinedRoutes = []fx.Option{
	dependency.HttpRoute("GET", "/_sys/migrations", sys.NewGetMigrationsHandler),
	dependency.HttpRoute("GET", "/_sys/health", sys.NewGetHealthHandler),
	dependency.HttpRoute("GET", "/user/self", user.NewGetSelfHandler),
	dependency.HttpRoute("GET", "/oauth/google/start", oauth.NewStartGoogleHandler),
	dependency.HttpRoute("GET", "/oauth/google/callback", oauth.NewAuthGoogleHandler),
}

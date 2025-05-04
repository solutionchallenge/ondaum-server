package http

import (
	"github.com/solutionchallenge/ondaum-server/internal/dependency"
	"github.com/solutionchallenge/ondaum-server/internal/handler/http/debug"
	"github.com/solutionchallenge/ondaum-server/internal/handler/http/oauth"
	"github.com/solutionchallenge/ondaum-server/internal/handler/http/sys"
	"github.com/solutionchallenge/ondaum-server/internal/handler/http/user"
	"go.uber.org/fx"
)

var PredefinedRoutes = []fx.Option{
	dependency.HttpRoute("GET", "/_sys/migrations", sys.NewGetMigrationsHandler),
	dependency.HttpRoute("GET", "/_sys/health", sys.NewGetHealthHandler),
	dependency.HttpRoute("POST", "/_debug/user", debug.NewUpsertUserHandler),
	dependency.HttpRoute("POST", "/_debug/auth", debug.NewAuthUserHandler),
	dependency.HttpRoute("GET", "/_debug/users", debug.NewListUserHandler),
	dependency.HttpRoute("GET", "/user/self", user.NewGetSelfHandler),
	dependency.HttpRoute("PUT", "/user/privacy", user.NewUpsertPrivacyHandler),
	dependency.HttpRoute("PUT", "/user/addition", user.NewUpsertAdditionHandler),
	dependency.HttpRoute("GET", "/oauth/google/start", oauth.NewStartGoogleHandler),
	dependency.HttpRoute("POST", "/oauth/google/auth", oauth.NewAuthGoogleHandler),
}

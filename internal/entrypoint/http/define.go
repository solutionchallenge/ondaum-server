package http

import (
	"github.com/solutionchallenge/ondaum-server/internal/dependency"
	"github.com/solutionchallenge/ondaum-server/internal/handler/future"
	"github.com/solutionchallenge/ondaum-server/internal/handler/rest/chat"
	"github.com/solutionchallenge/ondaum-server/internal/handler/rest/debug"
	"github.com/solutionchallenge/ondaum-server/internal/handler/rest/oauth"
	"github.com/solutionchallenge/ondaum-server/internal/handler/rest/schema"
	"github.com/solutionchallenge/ondaum-server/internal/handler/rest/sys"
	"github.com/solutionchallenge/ondaum-server/internal/handler/rest/user"
	"github.com/solutionchallenge/ondaum-server/internal/handler/websocket"
	"go.uber.org/fx"
)

var PredefinedRoutes = []fx.Option{
	dependency.HttpRoute("GET", "/_sys/migrations", sys.NewGetMigrationsHandler),
	dependency.HttpRoute("GET", "/_sys/health", sys.NewGetHealthHandler),
	dependency.HttpRoute("GET", "/_sys/tokens", sys.NewGetTokensHandler),
	dependency.HttpRoute("POST", "/_debug/user", debug.NewUpsertUserHandler),
	dependency.HttpRoute("POST", "/_debug/auth", debug.NewAuthUserHandler),
	dependency.HttpRoute("GET", "/_debug/oauth", debug.NewOAuthCallbackHandler),
	dependency.HttpRoute("GET", "/_debug/users", debug.NewListUserHandler),
	dependency.HttpRoute("GET", "/_debug/chats", debug.NewListChatHandler),
	dependency.HttpRoute("GET", "/user/self", user.NewGetSelfHandler),
	dependency.HttpRoute("PUT", "/user/privacy", user.NewUpsertPrivacyHandler),
	dependency.HttpRoute("PUT", "/user/addition", user.NewUpsertAdditionHandler),
	dependency.HttpRoute("GET", "/oauth/google/start", oauth.NewStartGoogleHandler),
	dependency.HttpRoute("POST", "/oauth/google/auth", oauth.NewAuthGoogleHandler),
	dependency.HttpRoute("GET", "/chats", chat.NewListChatHandler),
	dependency.HttpRoute("GET", "/chat/:session_id", chat.NewGetChatHandler),
	dependency.HttpRoute("PUT", "/chat/:session_id/summary", chat.NewUpsertSummaryHandler),
	dependency.HttpRoute("GET", "/_schema/supported-emotions", schema.NewGetSupportedEmotionHandler),
}

var WebsocketRoutes = []fx.Option{
	dependency.WebsocketRoute("/chat/ws", websocket.NewChatHandler),
}

var FutureProcesses = []fx.Option{
	dependency.FutureProcess(future.ChatJobType, future.NewChatFutureHandler),
}

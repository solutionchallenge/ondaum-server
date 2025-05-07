package http

import (
	"github.com/solutionchallenge/ondaum-server/internal/dependency"
	"github.com/solutionchallenge/ondaum-server/internal/handler/future"
	"github.com/solutionchallenge/ondaum-server/internal/handler/rest/chat"
	"github.com/solutionchallenge/ondaum-server/internal/handler/rest/debug"
	"github.com/solutionchallenge/ondaum-server/internal/handler/rest/diagnosis"
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
	dependency.HttpRoute("PUT", "/user/privacy", user.NewUpsertUserPrivacyHandler),
	dependency.HttpRoute("PUT", "/user/addition", user.NewUpsertUserAdditionHandler),
	dependency.HttpRoute("GET", "/oauth/google/start", oauth.NewStartGoogleHandler),
	dependency.HttpRoute("POST", "/oauth/google/auth", oauth.NewAuthGoogleHandler),
	dependency.HttpRoute("GET", "/chats", chat.NewListChatHandler),
	dependency.HttpRoute("GET", "/chats/:session_id", chat.NewGetChatHandler),
	dependency.HttpRoute("GET", "/chats/:session_id/summary", chat.NewGetChatSummaryHandler),
	dependency.HttpRoute("PUT", "/chats/:session_id/summary", chat.NewUpsertChatSummaryHandler),
	dependency.HttpRoute("POST", "/chats/:session_id/archive", chat.NewArchiveChatHandler),
	dependency.HttpRoute("GET", "/diagnoses", diagnosis.NewListDiagnosisResultHandler),
	dependency.HttpRoute("POST", "/diagnoses/report", diagnosis.NewReportDiagnosisResultHandler),
	dependency.HttpRoute("GET", "/diagnoses/:diagnosis_id", diagnosis.NewGetDiagnosisResultHandler),
	dependency.HttpRoute("GET", "/diagnosis-papers", diagnosis.NewListDiagnosisPaperHandler),
	dependency.HttpRoute("GET", "/diagnosis-papers/:paper_id", diagnosis.NewGetDiagnosisPaperHandler),
	dependency.HttpRoute("GET", "/_schema/supported-emotions", schema.NewListSupportedEmotionHandler),
	dependency.HttpRoute("GET", "/_schema/supported-features", schema.NewListSupportedFeatureHandler),
	dependency.HttpRoute("GET", "/_schema/supported-diagnoses", schema.NewListSupportedDiagnosisHandler),
}

var WebsocketRoutes = []fx.Option{
	dependency.WebsocketRoute("/_ws/chat", websocket.NewChatHandler),
}

var FutureProcesses = []fx.Option{
	dependency.FutureProcess(future.ChatJobType, future.NewChatFutureHandler),
}

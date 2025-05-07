package http

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"github.com/benbjohnson/clock"
	"github.com/solutionchallenge/ondaum-server/internal/dependency"
	"github.com/solutionchallenge/ondaum-server/migration"
	"github.com/solutionchallenge/ondaum-server/pkg/database"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/solutionchallenge/ondaum-server/pkg/jwt"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"

	"go.uber.org/fx"
)

func Run(config AppConfig) {
	if config.Verbose {
		buf, err := json.MarshalIndent(config, "", "  ")
		encoded := base64.StdEncoding.EncodeToString(buf)
		if err == nil {
			value := fmt.Sprintf("Starting HTTP server with config: %s", string(encoded))
			utils.Log(utils.InfoLevel).BT().Send("%s", value)
		}
	}
	app := fx.New(
		fx.Provide(clock.New),
		fx.Supply(config.HttpConfig),
		fx.Supply(config.DatabaseConfig),
		fx.Supply(config.OAuthConfig),
		fx.Supply(config.JWTConfig),
		fx.Supply(config.FutureConfig),
		fx.Supply(config.LLMConfig),
		dependency.NewDatabaseModule(config.DatabaseConfig, utils.DebugLevel),
		dependency.ProvideMiddleware(http.NewJWTAuthMiddleware),
		dependency.NewHttpModule("/api/v1", PredefinedRoutes...),
		dependency.NewOAuthModule(config.OAuthConfig),
		dependency.NewWebsocketModule(WebsocketRoutes...),
		dependency.NewFutureModule(config.FutureConfig, FutureProcesses...),
		dependency.NewLLMModule(config.LLMConfig),
		fx.Provide(jwt.NewGenerator),
		fx.Invoke(func(db *sql.DB) {
			if config.Migration.Enabled {
				ctx := context.Background()
				if err := db.PingContext(ctx); err == nil {
					database.Migrate(db, migration.Collection...)
				} else {
					utils.Log(utils.ErrorLevel).Ctx(ctx).Err(err).BT().Send("Failed to migrate database")
				}
			}
		}),
	)
	utils.RunGracefully(
		[]os.Signal{syscall.SIGINT, syscall.SIGTERM},
		utils.Runner{
			RunningFunction: func() error {
				return app.Start(context.Background())
			},
			ShutdownHandler: func() error {
				return app.Stop(context.Background())
			},
		},
	)
}

package dependency

import (
	"context"
	"database/sql"

	"github.com/solutionchallenge/ondaum-server/pkg/database"
	"github.com/solutionchallenge/ondaum-server/pkg/database/memdb"
	"github.com/solutionchallenge/ondaum-server/pkg/database/mysql"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

func NewDatabaseModule(config database.Config, logLevel utils.LogLevel) fx.Option {
	return fx.Module("database",
		fx.Provide(func() (database.Connector, error) {
			return instantiateDatabaseConnector(config)
		}),
		fx.Provide(func(conn database.Connector) (*bun.DB, error) {
			hook := database.NewQueryLoggingHook(logLevel)
			return conn.ToBunDB(&hook)
		}),
		fx.Invoke(func(lc fx.Lifecycle, db *bun.DB) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return db.PingContext(ctx)
				},
				OnStop: func(ctx context.Context) error {
					return db.Close()
				},
			})
		}),
		fx.Provide(func(conn database.Connector) (*sql.DB, error) {
			hook := database.NewQueryLoggingHook(logLevel)
			return conn.ToSqlDB(&hook)
		}),
		fx.Invoke(func(lc fx.Lifecycle, db *sql.DB) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return db.PingContext(ctx)
				},
				OnStop: func(ctx context.Context) error {
					return db.Close()
				},
			})
		}),
	)
}

func instantiateDatabaseConnector(config database.Config) (database.Connector, error) {
	switch config.Kind {
	case memdb.Kind:
		return memdb.NewConnector(config)
	case mysql.Kind:
		return mysql.NewConnector(config)
	default:
		return database.Connector{}, utils.NewError("unsupported database driver: %s", config.Kind)
	}
}

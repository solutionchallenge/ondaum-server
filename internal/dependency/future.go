package dependency

import (
	"context"

	"github.com/benbjohnson/clock"
	"github.com/solutionchallenge/ondaum-server/pkg/future"
	dbfuture "github.com/solutionchallenge/ondaum-server/pkg/future/database"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

func NewFutureModule(config future.Config, process ...fx.Option) fx.Option {
	return fx.Module("future",
		fx.Provide(func(db *bun.DB, clk clock.Clock) future.Core {
			return dbfuture.NewCore(db, clk)
		}),
		fx.Provide(func(core future.Core) *future.Scheduler {
			return future.NewScheduler(config, core)
		}),
		fx.Options(process...),
		fx.Invoke(func(lc fx.Lifecycle, scheduler *future.Scheduler) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					scheduler.Start()
					return nil
				},
				OnStop: func(_ context.Context) error {
					scheduler.Stop()
					return nil
				},
			})
		}),
	)
}

func FutureProcess[DEP any, H future.Handler](actionType future.JobType, constructor func(dependencies DEP) (H, error)) fx.Option {
	return fx.Options(
		fx.Provide(fx.Private, constructor),
		fx.Invoke(func(scheduler *future.Scheduler, handler H) {
			scheduler.AddHandler(actionType, handler)
		}),
	)
}

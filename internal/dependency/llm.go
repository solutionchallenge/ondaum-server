package dependency

import (
	"context"

	"github.com/solutionchallenge/ondaum-server/pkg/llm"
	"github.com/solutionchallenge/ondaum-server/pkg/llm/gemini"
	"go.uber.org/fx"
)

func NewLLMModule(config llm.Config) fx.Option {
	return fx.Module("llm",
		fx.Provide(func() (llm.Client, error) {
			return gemini.NewClient(config)
		}),
		fx.Invoke(func(lc fx.Lifecycle, llm llm.Client) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					return nil
				},
				OnStop: func(_ context.Context) error {
					return llm.Close()
				},
			})
		}),
	)
}

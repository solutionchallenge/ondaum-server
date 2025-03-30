package dependency

import (
	"context"
	"fmt"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"go.uber.org/fx"

	"github.com/gofiber/swagger"
	_ "github.com/solutionchallenge/ondaum-server/docs"
)

func NewHttpModule(root string, routes ...fx.Option) fx.Option {
	return fx.Module("http",
		fx.Options(
			fx.Provide(
				fx.Private,
				fx.Annotate(
					func(middlewares []http.MiddlewareFunc) *fiber.App {
						app := fiber.New()
						app.Use(requestid.New())
						app.Use(cors.New(cors.Config{
							AllowOrigins: "*",
							AllowHeaders: "Origin, Content-Type, Accept, Authorization",
						}))
						app.Use(logger.New(logger.Config{
							Format: "[${ip}]:${port} ${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
						}))
						app.Use(pprof.New(pprof.Config{Prefix: path.Join(root, "_sys")}))
						for _, middleware := range middlewares {
							app.Use(middleware)
						}
						app.Get(path.Join(root, "_sys", "swagger/*"), swagger.HandlerDefault)
						return app
					},
					fx.ParamTags(`group:"middlewares"`),
				),
			),
			fx.Provide(func(app *fiber.App) fiber.Router {
				return http.AsRouter(app, root)
			}),
			fx.Options(routes...),
			fx.Invoke(func(lc fx.Lifecycle, server *fiber.App, config http.Config) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						return server.Listen(fmt.Sprintf("%s:%d", config.Host, config.Port))
					},
					OnStop: func(ctx context.Context) error {
						return server.ShutdownWithTimeout(config.Timeout.Shutdown)
					},
				})
			}),
		),
	)
}

func ProvideMiddleware(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			// fx.As(new(http.MiddlewareFunc)),
			fx.ResultTags(`group:"middlewares"`),
		),
	)
}

func HttpRoute[H http.Handler, DEP any](method string, path string, constructor func(dependencies DEP) (H, error)) fx.Option {
	return fx.Options(
		fx.Provide(constructor),
		fx.Invoke(func(router fiber.Router, instance H) {
			router.Add(method, path, instance.Handle).Name(instance.Identify())
		}),
	)
}

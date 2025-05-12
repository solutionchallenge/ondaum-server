package dependency

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/pkg/jwt"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	wspkg "github.com/solutionchallenge/ondaum-server/pkg/websocket"
	"go.uber.org/fx"
)

type WebsocketCoreManager struct {
	ConfiguredCoreMap map[string]struct{}
}

func NewWebsocketModule(routes ...fx.Option) fx.Option {
	return fx.Module("websocket",
		fx.Supply(
			&WebsocketCoreManager{
				ConfiguredCoreMap: make(map[string]struct{}),
			},
		),
		fx.Options(routes...),
	)
}

func WebsocketRoute[DEP any, H wspkg.Handler](path string, constructor func(dependencies DEP) (H, error)) fx.Option {
	return fx.Options(
		fx.Provide(fx.Private, constructor),
		fx.Invoke(
			func(manager *WebsocketCoreManager, router fiber.Router, jwt *jwt.Generator, handler H) {
				if _, ok := manager.ConfiguredCoreMap[path]; !ok {
					manager.ConfiguredCoreMap[path] = struct{}{}
					err := wspkg.EnableWebsocketCore(router, path, jwt)
					if err != nil {
						panic(err)
					}
					wspkg.Install(router, path, handler)
				} else {
					panic(utils.NewError("websocket core already configured for path %s", path))
				}
			},
		),
	)
}

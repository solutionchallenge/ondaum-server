package sys

import (
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type GetHealthHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type GetHealthHandler struct {
	deps GetHealthHandlerDependencies
}

func NewGetHealthHandler(deps GetHealthHandlerDependencies) (*GetHealthHandler, error) {
	return &GetHealthHandler{deps: deps}, nil
}

func (h *GetHealthHandler) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()
	if err := h.deps.DB.PingContext(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to ping database",
		})
	}
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

func (h *GetHealthHandler) Identify() string {
	return "get-health"
}

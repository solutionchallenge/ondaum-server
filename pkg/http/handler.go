package http

import (
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Identify() string
	Handle(c *fiber.Ctx) error
}

func Bind[H Handler](router fiber.Router, method string, path string, handler H) fiber.Router {
	return router.Add(method, path, handler.Handle).Name(handler.Identify())
}

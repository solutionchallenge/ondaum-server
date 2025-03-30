package http

import (
	"github.com/gofiber/fiber/v2"
)

func NewServer(config Config) *fiber.App {
	server := fiber.New(fiber.Config{
		BodyLimit:       config.Limit.Payload,
		Concurrency:     config.Limit.Concurrency,
		ReadTimeout:     config.Timeout.Read,
		WriteTimeout:    config.Timeout.Write,
		IdleTimeout:     config.Timeout.Idle,
		ReadBufferSize:  config.Buffer.Read,
		WriteBufferSize: config.Buffer.Write,
	})
	return server
}

func AsRouter(app *fiber.App, root string) fiber.Router {
	return app.Group(root)
}

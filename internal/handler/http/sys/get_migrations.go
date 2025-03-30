package sys

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/pkg/database"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type GetMigrationsHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type GetMigrationsHandlerResponse struct {
	MigrationName string    `json:"migration_name"`
	CreatedAt     time.Time `json:"created_at"`
}

type GetMigrationsHandler struct {
	deps GetMigrationsHandlerDependencies
}

func NewGetMigrationsHandler(deps GetMigrationsHandlerDependencies) (*GetMigrationsHandler, error) {
	return &GetMigrationsHandler{deps: deps}, nil
}

func (h *GetMigrationsHandler) Handle(c *fiber.Ctx) error {
	migrations, err := database.GetMigrationHistories(c.UserContext(), h.deps.DB.DB)
	if err != nil {
		utils.Log(utils.ErrorLevel).Ctx(c.UserContext()).Err(err).Send("Failed to get migrations")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get migrations",
		})
	}

	response := utils.Map(migrations, func(migration database.MigrationHistory) GetMigrationsHandlerResponse {
		return GetMigrationsHandlerResponse{
			MigrationName: migration.Name,
			CreatedAt:     migration.CreatedAt,
		}
	})

	return c.JSON(response)
}

func (h *GetMigrationsHandler) Identify() string {
	return "get-migrations"
}

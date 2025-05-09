package debug

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type UpsertUserHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type UpsertUserHandlerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UpsertUserHandlerResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UpsertUserHandler struct {
	deps UpsertUserHandlerDependencies
}

func NewUpsertUserHandler(deps UpsertUserHandlerDependencies) (*UpsertUserHandler, error) {
	return &UpsertUserHandler{deps: deps}, nil
}

// Must not be documented. (Debugging purpose only!)
func (h *UpsertUserHandler) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()
	if os.Getenv("FLAG_DEBUGGING_FEATURES_ENABLED") != "true" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	request := &UpsertUserHandlerRequest{}
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			http.NewError(ctx, err, "Invalid request"),
		)
	}

	user := &user.User{
		Email:    request.Email,
		Username: request.Username,
	}

	_, err := h.deps.DB.NewInsert().
		Model(user).
		On("DUPLICATE KEY UPDATE").
		Set("username = ?", request.Username).
		Set("updated_at = CURRENT_TIMESTAMP").
		Exec(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to upsert user"),
		)
	}

	response := UpsertUserHandlerResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
	return c.JSON(response)
}

func (h *UpsertUserHandler) Identify() string {
	return "upsert-user"
}

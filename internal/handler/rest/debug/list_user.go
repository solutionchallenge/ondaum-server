package debug

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type ListUserHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type ListUserHandlerResponse struct {
	Users []UserInfo `json:"users"`
}

type UserInfo struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ListUserHandler struct {
	deps ListUserHandlerDependencies
}

func NewListUserHandler(deps ListUserHandlerDependencies) (*ListUserHandler, error) {
	return &ListUserHandler{deps: deps}, nil
}

// Must not be documented. (Debugging purpose only!)
func (h *ListUserHandler) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()
	if os.Getenv("FLAG_DEBUGGING_FEATURES_ENABLED") != "true" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	var users []user.User
	err := h.deps.DB.NewSelect().
		Model(&users).
		Order("id ASC").
		Scan(ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(ctx, err, "Failed to list users"),
		)
	}

	userInfos := make([]UserInfo, len(users))
	for i, u := range users {
		userInfos[i] = UserInfo{
			ID:        u.ID,
			Email:     u.Email,
			Username:  u.Username,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	response := ListUserHandlerResponse{
		Users: userInfos,
	}
	return c.JSON(response)
}

func (h *ListUserHandler) Identify() string {
	return "list-user"
}

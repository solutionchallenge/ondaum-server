package user

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/http"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type GetSelfHandlerDependencies struct {
	fx.In
	DB *bun.DB
}

type GetSelfHandlerResponse struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type GetSelfHandler struct {
	deps GetSelfHandlerDependencies
}

var _ http.Handler = &GetSelfHandler{}

func NewGetSelfHandler(deps GetSelfHandlerDependencies) (*GetSelfHandler, error) {
	return &GetSelfHandler{
		deps: deps,
	}, nil
}

// @ID GetSelfUser
// @Summary      Get Self User Information
// @Description  This API returns the user's information.
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200 {object} GetSelfHandlerResponse
// @Failure      401 {object} http.Error
// @Failure      404 {object} http.Error
// @Failure      500 {object} http.Error
// @Router       /user/self [get]
// @Security     BearerAuth
func (h *GetSelfHandler) Handle(c *fiber.Ctx) error {
	rid := http.GetRequestID(c)
	id, err := http.GetUserID(c)
	if err != nil {
		utils.Log(utils.InfoLevel).Ctx(c.UserContext()).Err(err).RID(rid).Send("Unauthorized")
		return c.Status(fiber.StatusUnauthorized).JSON(
			http.NewError(err, "Unauthorized"),
		)
	}
	user := user.User{}
	err = h.deps.DB.NewSelect().
		Model(&user).
		Where("id = ?", id).
		Scan(c.UserContext())
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Log(utils.InfoLevel).Ctx(c.UserContext()).Err(err).RID(rid).Send("User not found for id: %v", id)
			return c.Status(fiber.StatusNotFound).JSON(
				http.NewError(err, "User not found"),
			)
		}
		utils.Log(utils.InfoLevel).Ctx(c.UserContext()).Err(err).RID(rid).Send("Failed to get user for id: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(
			http.NewError(err, "Failed to get user for id: "+strconv.FormatInt(id, 10)),
		)
	}

	response := GetSelfHandlerResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
	return c.JSON(response)
}

func (h *GetSelfHandler) Identify() string {
	return "get-self"
}

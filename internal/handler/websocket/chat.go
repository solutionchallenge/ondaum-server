package websocket

import (
	"github.com/benbjohnson/clock"
	fiberws "github.com/gofiber/websocket/v2"
	impl "github.com/solutionchallenge/ondaum-server/internal/handler/websocket/chat"
	ftpkg "github.com/solutionchallenge/ondaum-server/pkg/future"
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
	wspkg "github.com/solutionchallenge/ondaum-server/pkg/websocket"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type ChatHandlerDependencies struct {
	fx.In
	Future *ftpkg.Scheduler
	LLM    llm.Client
	DB     *bun.DB
	Clock  clock.Clock
}

type ChatHandler struct {
	deps ChatHandlerDependencies
}

func NewChatHandler(deps ChatHandlerDependencies) (*ChatHandler, error) {
	return &ChatHandler{deps: deps}, nil
}

// @ID ConnectChatWebsocket
// @Summary      Connect Chat Websocket
// @Description  Connect Chat Websocket. Reference the notion page for more information.
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        session_id query string false "Websocket Session ID (optional; if not provided, the server will use the most recent non-archived conversation or create a new one if none exists)"
// @Success      200 {object} wspkg.ResponseWrapper
// @Failure      426 {object} http.Error
// @Router       /_ws/chat [get]
// @Security     BearerAuth
func (h *ChatHandler) HandleMessage(c *fiberws.Conn, request wspkg.MessageWrapper) (wspkg.ResponseWrapper, bool, error) {
	return impl.HandleMessage(h.deps.DB, h.deps.Clock, h.deps.LLM, h.deps.Future, request)
}

func (h *ChatHandler) HandleConnect(c *fiberws.Conn, request wspkg.ConnectWrapper) (wspkg.ResponseWrapper, error) {
	return impl.HandleConnect(h.deps.DB, h.deps.Clock, request)
}

func (h *ChatHandler) HandleClose(_ *fiberws.Conn, _ wspkg.CloseWrapper) {}

func (h *ChatHandler) HandlePing(c *fiberws.Conn, request wspkg.PingWrapper) (wspkg.ResponseWrapper, bool, error) {
	return impl.HandlePing(h.deps.DB, request)
}

func (h *ChatHandler) Identify() string {
	return "ws-chat"
}

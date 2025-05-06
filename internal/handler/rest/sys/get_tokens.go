package sys

import (
	"github.com/gofiber/fiber/v2"
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
	"go.uber.org/fx"
)

type GetTokensDependencies struct {
	fx.In
	LLM llm.Client
}

type GetTokensResponse struct {
	TotalTokenCount      int `json:"total_token_count"`
	PromptTokenCount     int `json:"prompt_token_count"`
	CompletionTokenCount int `json:"completion_token_count"`
	CachedTokenCount     int `json:"cached_token_count"`
}

type GetTokensHandler struct {
	deps GetTokensDependencies
}

func NewGetTokensHandler(deps GetTokensDependencies) (*GetTokensHandler, error) {
	return &GetTokensHandler{deps: deps}, nil
}

func (h *GetTokensHandler) Handle(c *fiber.Ctx) error {
	response := GetTokensResponse{}
	statistics := h.deps.LLM.GetStatistics()
	response.TotalTokenCount = int(statistics.TotalTokens)
	response.PromptTokenCount = int(statistics.PromptTokens)
	response.CompletionTokenCount = int(statistics.CompletionTokens)
	response.CachedTokenCount = int(statistics.CachedTokens)
	return c.JSON(response)
}

func (h *GetTokensHandler) Identify() string {
	return "get-tokens"
}

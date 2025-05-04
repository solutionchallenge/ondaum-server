package gemini

import (
	"github.com/google/generative-ai-go/genai"
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
)

func AddStatistics(statistics *llm.Statistics, usage *genai.UsageMetadata) {
	statistics.TotalTokens += int64(usage.TotalTokenCount)
	statistics.PromptTokens += int64(usage.PromptTokenCount)
	statistics.CompletionTokens += int64(usage.CandidatesTokenCount)
	statistics.CachedTokens += int64(usage.CachedContentTokenCount)
}

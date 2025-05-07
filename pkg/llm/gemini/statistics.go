package gemini

import (
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
	"google.golang.org/genai"
)

func AddStatistics(statistics *llm.Statistics, usage *genai.GenerateContentResponseUsageMetadata) {
	statistics.TotalTokens += int64(usage.TotalTokenCount)
	statistics.PromptTokens += int64(usage.PromptTokenCount)
	statistics.CompletionTokens += int64(usage.CandidatesTokenCount)
	statistics.ThoughtsTokens += int64(usage.ThoughtsTokenCount)
	statistics.CachedTokens += int64(usage.CachedContentTokenCount)
}

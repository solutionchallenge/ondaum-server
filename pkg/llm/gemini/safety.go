package gemini

import (
	"fmt"

	"github.com/solutionchallenge/ondaum-server/pkg/llm"
	"google.golang.org/genai"
)

func ConfigToSafetySetting(config llm.Config) []*genai.SafetySetting {
	return []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: parseSafetyThreshold(config.Gemini.RedactionThreshold.Harrasement),
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: parseSafetyThreshold(config.Gemini.RedactionThreshold.HateSpeech),
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: parseSafetyThreshold(config.Gemini.RedactionThreshold.SexuallyExplicit),
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: parseSafetyThreshold(config.Gemini.RedactionThreshold.DangerousContent),
		},
		{
			Category:  genai.HarmCategoryCivicIntegrity,
			Threshold: parseSafetyThreshold(config.Gemini.RedactionThreshold.CivicIntegrity),
		},
	}
}

func parseSafetyThreshold(threshold string) genai.HarmBlockThreshold {
	switch threshold {
	case "none":
		return genai.HarmBlockThresholdBlockNone
	case "low":
		return genai.HarmBlockThresholdBlockLowAndAbove
	case "medium":
		return genai.HarmBlockThresholdBlockMediumAndAbove
	case "high":
		return genai.HarmBlockThresholdBlockOnlyHigh
	default:
		return genai.HarmBlockThresholdUnspecified
	}
}

func checkPromptBlocked(response *genai.GenerateContentResponse) error {
	if response.PromptFeedback != nil {
		switch response.PromptFeedback.BlockReason {
		case genai.BlockedReasonProhibitedContent, genai.BlockedReasonSafety:
			return fmt.Errorf(
				"blocked by gemini by inappropriate prompt: %v(%v)",
				response.PromptFeedback.BlockReasonMessage,
				response.PromptFeedback.SafetyRatings,
			)
		default:
			return fmt.Errorf(
				"blocked by gemini with unknown reason: %v(%v)",
				response.PromptFeedback.BlockReasonMessage,
				response.PromptFeedback.SafetyRatings,
			)
		}
	}
	return nil
}

func checkContentBlocked(feedbacks []map[string]any) error {
	for _, feedback := range feedbacks {
		if feedback["blocked"] == true {
			return fmt.Errorf("blocked by gemini by inappropriate content: %v(%v)", feedback["category"], feedback["probability"])
		}
	}
	return nil
}

func buildContentFeedbacks(response *genai.GenerateContentResponse) []map[string]any {
	feedback := make([]map[string]any, len(response.Candidates))
	for idx, candidate := range response.Candidates {
		if len(candidate.SafetyRatings) > 0 {
			for _, result := range candidate.SafetyRatings {
				feedback[idx] = map[string]any{
					"blocked":     result.Blocked,
					"category":    result.Category,
					"probability": result.Probability,
				}
			}
		}
	}
	return feedback
}

package gemini

import (
	"github.com/google/generative-ai-go/genai"
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
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
		// { // not supported yet
		// 	Category:  genai.HarmCategoryCivicIntegrity,
		// 	Threshold: parseSafetyThreshold(config.Gemini.RedactionThreshold.CivicIntegrity),
		// },
	}
}

func parseSafetyThreshold(threshold string) genai.HarmBlockThreshold {
	switch threshold {
	case "none":
		return genai.HarmBlockNone
	case "low":
		return genai.HarmBlockLowAndAbove
	case "medium":
		return genai.HarmBlockMediumAndAbove
	case "high":
		return genai.HarmBlockOnlyHigh
	default:
		return genai.HarmBlockUnspecified
	}
}

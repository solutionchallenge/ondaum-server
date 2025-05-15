package gemini

import (
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"google.golang.org/genai"
)

func BuildGenerativeConfig(client *Client, promptIdentifier string, rootpath ...string) (*genai.GenerateContentConfig, error) {
	var preparedPrompt *llm.PreparedPrompt = nil
	if promptIdentifier != "" {
		for _, prepared := range client.Config.Gemini.PreparedPrompts {
			if prepared.Identifier == promptIdentifier && prepared.PromptType == llm.PromptTypeSystemInstruction {
				preparedPrompt = &prepared
				break
			}
		}
	}
	systemInstruction := (*genai.Content)(nil)
	if preparedPrompt != nil {
		promptData, err := utils.ReadFileFrom(preparedPrompt.PromptFile, rootpath...)
		if err != nil {
			return nil, utils.WrapError(err, "failed to read prompt file")
		}
		systemInstruction = genai.NewContentFromText(promptData, genai.RoleUser)
	}
	safetySettings := ([]*genai.SafetySetting)(nil)
	if !preparedPrompt.DisableRedaction {
		safetySettings = ConfigToSafetySetting(client.Config)
	}
	chatConfig := &genai.GenerateContentConfig{
		ResponseMIMEType:  client.Config.Gemini.ResponseFormat,
		SafetySettings:    safetySettings,
		SystemInstruction: systemInstruction,
	}
	return chatConfig, nil
}

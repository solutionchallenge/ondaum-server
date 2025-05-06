package gemini

import (
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
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
		promptData, err := readPromptFile(preparedPrompt.PromptFile, rootpath...)
		if err != nil {
			return nil, err
		}
		systemInstruction = genai.NewContentFromText(promptData, genai.Role("system"))
	}
	chatConfig := &genai.GenerateContentConfig{
		ResponseMIMEType:  client.Config.Gemini.ResponseFormat,
		SafetySettings:    ConfigToSafetySetting(client.Config),
		SystemInstruction: systemInstruction,
	}
	return chatConfig, nil
}

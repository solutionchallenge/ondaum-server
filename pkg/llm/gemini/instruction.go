package gemini

import (
	"os"
	"path"

	"github.com/google/generative-ai-go/genai"
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
)

func ConfigToInstruction(config llm.Config, rootpath ...string) (*genai.Content, error) {
	var err error
	basepath := "./"
	if len(rootpath) > 0 && rootpath[0] != "" {
		basepath = rootpath[0]
	} else {
		basepath, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	filepath := path.Join(basepath, config.Gemini.SystemPrompt)
	instruction, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return &genai.Content{
		Parts: []genai.Part{genai.Text(instruction)},
		Role:  "system",
	}, nil
}

package gemini

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
	"google.golang.org/api/option"
)

var _ llm.Client = &Client{}

type Client struct {
	Config        llm.Config
	Core          *genai.Client
	Model         *genai.GenerativeModel
	Conversations map[string]llm.Conversation
	Statistics    llm.Statistics
	Mutex         sync.Mutex
}

func NewClient(config llm.Config) (*Client, error) {
	if !config.Gemini.Enabled {
		return nil, fmt.Errorf("gemini is not enabled")
	}
	core, err := genai.NewClient(context.Background(), option.WithAPIKey(config.Gemini.APIKey))
	if err != nil {
		return nil, err
	}
	model := core.GenerativeModel(config.Gemini.LLMModel)

	instruction, err := ConfigToInstruction(config)
	if err != nil {
		return nil, err
	}
	model.SystemInstruction = instruction
	model.ResponseMIMEType = config.Gemini.ResponseFormat
	model.SafetySettings = ConfigToSafetySetting(config)

	return &Client{
		Config:        config,
		Core:          core,
		Model:         model,
		Conversations: make(map[string]llm.Conversation),
	}, nil
}

func (client *Client) StartConversation(manager llm.HistoryManager, id ...string) llm.Conversation {
	client.Mutex.Lock()
	defer client.Mutex.Unlock()

	ConversationID := uuid.New().String()
	if len(id) > 0 && id[0] != "" {
		ConversationID = id[0]
	}

	conversation := NewConversation(ConversationID, client, manager)
	client.Conversations[conversation.ID] = conversation
	return conversation
}

func (client *Client) GetStatistics() llm.Statistics {
	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	return client.Statistics
}

func (client *Client) Close() error {
	client.Mutex.Lock()
	defer client.Mutex.Unlock()

	for _, conversation := range client.Conversations {
		conversation.End()
	}
	return nil
}

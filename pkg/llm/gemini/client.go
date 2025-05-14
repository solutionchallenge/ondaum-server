package gemini

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"google.golang.org/genai"
)

var _ llm.Client = &Client{}

type Client struct {
	Config        llm.Config
	Core          *genai.Client
	Conversations map[string]llm.Conversation
	Statistics    llm.Statistics
	Mutex         sync.Mutex
}

func NewClient(config llm.Config) (*Client, error) {
	if !config.Gemini.Enabled {
		return nil, utils.NewError("gemini is not enabled")
	}
	core, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  config.Gemini.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, utils.WrapError(err, "failed to create gemini client")
	}
	return &Client{
		Config:        config,
		Core:          core,
		Conversations: make(map[string]llm.Conversation),
	}, nil
}

func (client *Client) StartConversation(ctx context.Context, historyManager llm.HistoryManager, instructionIdentifier string, id ...string) (llm.Conversation, error) {
	ConversationID := uuid.New().String()
	if len(id) > 0 && id[0] != "" {
		ConversationID = id[0]
	}

	conversation, err := NewConversation(ctx, ConversationID, client, instructionIdentifier, historyManager)
	if err != nil {
		return nil, utils.WrapError(err, "failed to create gemini conversation")
	}

	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	client.Conversations[conversation.ID] = conversation
	return conversation, nil
}

func (client *Client) RunActionPrompt(ctx context.Context, instructionIdentifier string, promptIdentifier string, histories ...llm.Message) (llm.Message, error) {
	config, err := BuildGenerativeConfig(client, instructionIdentifier)
	if err != nil {
		return llm.Message{}, utils.WrapError(err, "failed to build generative config")
	}

	var prepared *llm.PreparedPrompt = nil
	for _, iterator := range client.Config.Gemini.PreparedPrompts {
		if iterator.Identifier == promptIdentifier && iterator.PromptType == llm.PromptTypeActionPrompt {
			prepared = &iterator
			break
		}
	}
	if prepared == nil {
		return llm.Message{}, utils.WrapError(err, "prepared prompt identifier '%s' not found", promptIdentifier)
	}

	prompt, err := utils.ReadFileFrom(prepared.PromptFile)
	if err != nil {
		return llm.Message{}, utils.WrapError(err, "ReadFileFrom failed for %s", prepared.PromptFile)
	}

	finalContents := []*genai.Content{}
	if len(histories) > 0 {
		historyContents := utils.Map(histories, func(message llm.Message) *genai.Content {
			return genai.NewContentFromText(message.Content, genai.Role(message.Role))
		})
		finalContents = append(finalContents, historyContents...)
	}

	currentUserTurnParts := []*genai.Part{genai.NewPartFromText(prompt)}
	if prepared.AttachmentFile != "" {
		reader, err := utils.OpenFileFrom(prepared.AttachmentFile)
		if err != nil {
			return llm.Message{}, utils.WrapError(err, "OpenFileFrom failed for %s", prepared.AttachmentFile)
		}
		defer reader.Close()

		uploadedFile, uploadErr := client.Core.Files.Upload(ctx, reader, &genai.UploadFileConfig{
			MIMEType:    prepared.AttachmentMime,
			DisplayName: prepared.AttachmentFile,
		})
		if uploadErr != nil {
			return llm.Message{}, utils.WrapError(uploadErr, "file upload failed for %s", prepared.AttachmentFile)
		}

		fileDataPart := genai.NewPartFromURI(uploadedFile.URI, prepared.AttachmentMime)
		currentUserTurnParts = append(currentUserTurnParts, fileDataPart)
	}

	currentUserTurnContent := genai.NewContentFromParts(currentUserTurnParts, genai.RoleUser)
	finalContents = append(finalContents, currentUserTurnContent)
	response, err := client.Core.Models.GenerateContent(
		ctx,
		client.Config.Gemini.LLMModel,
		finalContents,
		config,
	)
	if err != nil {
		return llm.Message{}, utils.WrapError(err, "GenerateContent failed")
	}

	AddStatistics(&client.Statistics, response.UsageMetadata)

	if err := checkPromptBlocked(response); err != nil {
		return llm.Message{}, utils.WrapError(err, "CheckPromptBlocked failed")
	}
	feedbacks := buildContentFeedbacks(response)
	if err := checkContentBlocked(feedbacks); err != nil {
		return llm.Message{}, utils.WrapError(err, "CheckContentBlocked failed")
	}

	return llm.Message{
		ID:      uuid.New().String(),
		Role:    llm.RoleAssistant,
		Content: response.Text(),
		Metadata: map[string]any{
			"feedbacks": feedbacks,
		},
	}, nil
}

func (client *Client) GetStatistics() llm.Statistics {
	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	return client.Statistics
}

func (client *Client) Close(ids ...string) error {
	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	if len(ids) <= 0 {
		client.Conversations = make(map[string]llm.Conversation)
	} else {
		for _, id := range ids {
			conversation, ok := client.Conversations[id]
			if ok && conversation != nil {
				delete(client.Conversations, id)
			}
		}
	}
	return nil
}

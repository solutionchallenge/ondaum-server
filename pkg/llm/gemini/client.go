package gemini

import (
	"context"
	"fmt"
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
		return nil, fmt.Errorf("gemini is not enabled")
	}
	core, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  config.Gemini.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
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
		return nil, err
	}

	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	client.Conversations[conversation.ID] = conversation
	return conversation, nil
}

func (client *Client) RunActionPrompt(ctx context.Context, instructionIdentifier string, promptIdentifier string, histories ...llm.Message) (llm.Message, error) {
	config, err := BuildGenerativeConfig(client, instructionIdentifier)
	if err != nil {
		return llm.Message{}, fmt.Errorf("BuildGenerativeConfig failed: %w", err)
	}

	var prepared *llm.PreparedPrompt = nil
	for _, iterator := range client.Config.Gemini.PreparedPrompts {
		if iterator.Identifier == promptIdentifier && iterator.PromptType == llm.PromptTypeActionPrompt {
			prepared = &iterator
			break
		}
	}
	if prepared == nil {
		return llm.Message{}, fmt.Errorf("prepared prompt identifier '%s' not found", promptIdentifier)
	}

	prompt, err := readPromptFile(prepared.PromptFile)
	if err != nil {
		return llm.Message{}, fmt.Errorf("readPromptFile failed for %s: %w", prepared.PromptFile, err)
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
		reader, err := openAttachmentFile(prepared.AttachmentFile)
		if err != nil {
			return llm.Message{}, fmt.Errorf("openAttachmentFile failed for %s: %w", prepared.AttachmentFile, err)
		}
		defer reader.Close()

		uploadedFile, uploadErr := client.Core.Files.Upload(ctx, reader, &genai.UploadFileConfig{
			MIMEType:    prepared.AttachmentMime,
			DisplayName: prepared.AttachmentFile,
		})
		if uploadErr != nil {
			return llm.Message{}, fmt.Errorf("file upload failed for %s: %w", prepared.AttachmentFile, uploadErr)
		}

		fileDataPart := genai.NewPartFromURI(uploadedFile.URI, prepared.AttachmentMime)

		contentForCaching := &genai.Content{
			Parts: []*genai.Part{fileDataPart},
		}

		cached, cacheErr := client.Core.Caches.Create(ctx, client.Config.Gemini.LLMModel, &genai.CreateCachedContentConfig{
			Contents:          []*genai.Content{contentForCaching},
			SystemInstruction: config.SystemInstruction,
			TTL:               prepared.AttachmentTTL,
			DisplayName:       fmt.Sprintf("cache_for_%s", prepared.AttachmentFile),
		})
		if cacheErr != nil {
			config.CachedContent = ""
			currentUserTurnParts = append(currentUserTurnParts, fileDataPart)
		} else {
			config.CachedContent = cached.Name
		}
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
		return llm.Message{}, fmt.Errorf("GenerateContent failed: %w", err)
	}

	if err := checkPromptBlocked(response); err != nil {
		return llm.Message{}, err
	}
	feedbacks := buildContentFeedbacks(response)
	if err := checkContentBlocked(feedbacks); err != nil {
		return llm.Message{}, err
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

func (client *Client) Close() error {
	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	for _, conversation := range client.Conversations {
		conversation.End()
	}
	client.Conversations = make(map[string]llm.Conversation)
	return nil
}

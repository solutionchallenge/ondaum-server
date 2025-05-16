package gemini

import (
	"context"

	"github.com/google/uuid"
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"google.golang.org/genai"
)

var (
	EmptyResponseErr = utils.NewError("empty response detected")
)

type Conversation struct {
	ID         string
	Client     *Client
	Session    *genai.Chat
	Statistics llm.Statistics
	Manager    llm.HistoryManager
}

func NewConversation(ctx context.Context, id string, client *Client, prompt string, manager llm.HistoryManager) (*Conversation, error) {
	histories := utils.Map(manager.Get(ctx, id), func(message llm.Message) *genai.Content {
		return genai.NewContentFromText(message.Content, genai.Role(message.Role))
	})
	config, err := BuildGenerativeConfig(client, prompt)
	if err != nil {
		return nil, utils.WrapError(err, "failed to build generative config")
	}
	session, err := client.Core.Chats.Create(
		ctx,
		client.Config.Gemini.LLMModel,
		config,
		histories,
	)
	if err != nil {
		return nil, utils.WrapError(err, "failed to create chatting session")
	}
	return &Conversation{
		ID:         id,
		Client:     client,
		Session:    session,
		Statistics: llm.Statistics{},
		Manager:    manager,
	}, nil
}

func (conversation *Conversation) Request(ctx context.Context, request llm.Message) (llm.Message, error) {
	conversation.Manager.Add(ctx, request)

	prompt := genai.NewPartFromText(request.Content)
	response, err := conversation.Session.SendMessage(ctx, *prompt)
	if err != nil {
		return llm.Message{}, utils.WrapError(err, "failed to send message")
	}

	AddStatistics(&conversation.Statistics, response.UsageMetadata)
	AddStatistics(&conversation.Client.Statistics, response.UsageMetadata)

	if err := checkPromptBlocked(response); err != nil {
		return llm.Message{}, utils.WrapError(err, "failed to check prompt blocked")
	}
	feedbacks := buildContentFeedbacks(response)
	if err := checkContentBlocked(feedbacks); err != nil {
		return llm.Message{}, utils.WrapError(err, "failed to check content blocked")
	}
	if response.Text() == "" {
		return llm.Message{}, utils.NewError("empty response detected")
	}
	message := llm.Message{
		ID:      uuid.New().String(),
		Role:    llm.RoleModel,
		Content: response.Text(),
		Metadata: map[string]any{
			"feedbacks": feedbacks,
		},
	}
	conversation.Manager.Add(ctx, message)
	return message, nil
}

func (conversation *Conversation) GetHistory(ctx context.Context) []llm.Message {
	return conversation.Manager.Get(ctx, conversation.ID)
}

func (conversation *Conversation) GetStatistics() llm.Statistics {
	return conversation.Statistics
}

func (conversation *Conversation) End() {
	conversation.Client.Close(conversation.ID)
}

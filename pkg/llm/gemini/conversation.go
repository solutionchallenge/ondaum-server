package gemini

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

type Conversation struct {
	ID         string
	Client     *Client
	Session    *genai.ChatSession
	Statistics llm.Statistics
	Manager    llm.HistoryManager
}

func NewConversation(id string, client *Client, manager llm.HistoryManager) *Conversation {
	histories := utils.Map(manager.Get(id), func(message llm.Message) *genai.Content {
		return &genai.Content{
			Role:  string(message.Role),
			Parts: []genai.Part{genai.Text(message.Content)},
		}
	})
	session := client.Model.StartChat()
	session.History = append(session.History, histories...)
	return &Conversation{
		ID:         id,
		Client:     client,
		Session:    session,
		Statistics: llm.Statistics{},
		Manager:    manager,
	}
}

func (conversation *Conversation) Request(request llm.Message) (llm.Message, error) {
	prompt := genai.Text(request.Content)
	response, err := conversation.Session.SendMessage(context.Background(), prompt)
	if err != nil {
		return llm.Message{}, err
	}

	AddStatistics(&conversation.Statistics, response.UsageMetadata)
	AddStatistics(&conversation.Client.Statistics, response.UsageMetadata)

	if len(response.Candidates) == 0 {
		switch response.PromptFeedback.BlockReason {
		case genai.BlockReasonSafety:
			return llm.Message{}, fmt.Errorf("blocked by gemini by inappropriate prompt: %v", response.PromptFeedback.BlockReason.String())
		default:
			return llm.Message{}, fmt.Errorf("blocked by gemini with unknown reason: %v", response.PromptFeedback.BlockReason.String())
		}
	}

	candidate := response.Candidates[0]
	metadata := map[string]any{}
	if len(candidate.SafetyRatings) > 0 {
		for _, result := range candidate.SafetyRatings {
			if result.Blocked {
				return llm.Message{}, fmt.Errorf("blocked by gemini by inappropriate content: %v(%v)", result.Category, result.Probability)
			}
			metadata["feedback"] = map[string]any{
				"category":    result.Category,
				"probability": result.Probability,
			}
		}
	}
	contents := ""
	for _, part := range candidate.Content.Parts {
		contents += string(part.(genai.Text))
	}
	message := llm.Message{
		ID:       uuid.New().String(),
		Role:     llm.RoleAssistant,
		Content:  contents,
		Metadata: metadata,
	}
	conversation.Manager.Add(request, message)
	return message, nil
}

func (conversation *Conversation) GetStatistics() llm.Statistics {
	return conversation.Statistics
}

func (conversation *Conversation) GetHistory() []llm.Message {
	return conversation.Manager.Get(conversation.ID)
}

func (conversation *Conversation) End() {}

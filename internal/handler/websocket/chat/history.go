package chat

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	domain "github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
)

type ChatHistoryManager struct {
	db             *bun.DB
	memoryCache    []llm.Message
	conversationID string
}

func NewChatHistoryManager(db *bun.DB, conversationID string) *ChatHistoryManager {
	return &ChatHistoryManager{
		db:             db,
		memoryCache:    []llm.Message{},
		conversationID: conversationID,
	}
}

func (h *ChatHistoryManager) Add(ctx context.Context, messages ...llm.Message) {
	h.memoryCache = append(h.memoryCache, messages...)
	h.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		chat := domain.Chat{}
		err := tx.NewSelect().Model(&chat).Where("session_id = ?", h.conversationID).Scan(ctx)
		if err != nil || chat.ID == 0 {
			utils.Log(utils.WarnLevel).CID(h.conversationID).Err(err).BT().Send("Failed to query chat")
			return utils.WrapError(err, "failed to query chat")
		}

		histories := make([]domain.History, 0, len(messages))
		for _, message := range messages {
			marshaled, err := json.Marshal(message.Metadata)
			if err != nil {
				utils.Log(utils.WarnLevel).CID(h.conversationID).Err(err).BT().Send("Failed to marshal metadata")
				continue
			}
			histories = append(histories, domain.History{
				ChatID:    chat.ID,
				MessageID: message.ID,
				Role:      string(message.Role),
				Content:   message.Content,
				Metadata:  marshaled,
			})
		}

		if len(histories) > 0 {
			_, err = tx.NewInsert().Model(&histories).Exec(ctx)
			if err != nil {
				utils.Log(utils.WarnLevel).CID(h.conversationID).Err(err).BT().Send("Failed to bulk insert chat history")
				return utils.WrapError(err, "failed to bulk insert chat history")
			}
		}
		return nil
	})
}

func (h *ChatHistoryManager) Get(ctx context.Context, conversationID string) []llm.Message {
	chat := &domain.Chat{}
	err := h.db.NewSelect().
		Model(chat).
		Relation("Histories", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("inserted_at ASC")
		}).
		Relation("Summary").
		Where("session_id = ?", conversationID).
		Order("created_at ASC").
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.Log(utils.WarnLevel).CID(conversationID).Err(err).BT().Send("Chat not found for session_id: %v", conversationID)
			return h.memoryCache
		}
		utils.Log(utils.ErrorLevel).CID(conversationID).Err(err).BT().Send("Failed to get chat for session_id: %v", conversationID)
		return h.memoryCache
	}
	utils.Log(utils.DebugLevel).CID(conversationID).BT().Send("Found %d histories", len(chat.Histories))
	return utils.Map(chat.Histories, func(history *domain.History) llm.Message {
		return llm.Message{
			ID:      history.MessageID,
			Role:    llm.Role(history.Role),
			Content: history.Content,
		}
	})
}

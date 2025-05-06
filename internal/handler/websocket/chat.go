package websocket

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/benbjohnson/clock"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	domain "github.com/solutionchallenge/ondaum-server/internal/domain/chat"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/internal/handler/future"
	ftpkg "github.com/solutionchallenge/ondaum-server/pkg/future"
	"github.com/solutionchallenge/ondaum-server/pkg/llm"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	wspkg "github.com/solutionchallenge/ondaum-server/pkg/websocket"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

const (
	ChatPayloadNotifyConversationFinished = "conversation_finished"
	ChatPayloadNotifyConversationArchived = "conversation_archived"
	ChatPayloadNotifyNewConversation      = "new_conversation"
	ChatPayloadNotifyExistingConversation = "existing_conversation"
)

const (
	ChatAutoFinishAfter = 30 * time.Minute
)

type ChatHandlerDependencies struct {
	fx.In
	Future *ftpkg.Scheduler
	LLM    llm.Client
	DB     *bun.DB
	Clock  clock.Clock
}

type ChatHandler struct {
	deps ChatHandlerDependencies
}

func NewChatHandler(deps ChatHandlerDependencies) (*ChatHandler, error) {
	return &ChatHandler{deps: deps}, nil
}

func (h *ChatHandler) Identify() string {
	return "ws-chat"
}

// @ID ConnectChatWebsocket
// @Summary      Connect Chat Websocket
// @Description  Connect Chat Websocket. Reference the notion page for more information.
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        session_id query string false "Websocket Session ID (if not provided, the server will create a new session)"
// @Success      200 {object} wspkg.ResponseWrapper
// @Failure      426 {object} http.Error
// @Router       /chat/ws [get]
// @Security     BearerAuth
func (h *ChatHandler) HandleMessage(c *fiberws.Conn, request wspkg.MessageWrapper) (wspkg.ResponseWrapper, bool, error) {
	if !request.Authorized || !h.checkAuthorization(request.UserID) {
		utils.Log(utils.InfoLevel).CID(request.SessionID).RID(request.MessageID).Send("Unauthorized")
		return wspkg.BuildRejectResponse(request), false, nil
	}

	var response wspkg.ResponseWrapper
	var shouldClose bool
	err := h.deps.DB.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		chat := &domain.Chat{}
		err := tx.NewSelect().
			Model(chat).
			Where("session_id = ? AND user_id = ?", request.SessionID, request.UserID).
			Scan(ctx)
		if err != nil {
			utils.Log(utils.ErrorLevel).CID(request.SessionID).RID(request.MessageID).Err(err).Send("Failed to query chat")
			return err
		}

		if !chat.ArchivedAt.IsZero() {
			utils.Log(utils.InfoLevel).CID(request.SessionID).RID(request.MessageID).Send("Cannot reactivate archived chat")
			response = wspkg.BuildRejectResponse(request)
			return nil
		}

		if !chat.FinishedAt.IsZero() {
			_, err = tx.NewUpdate().
				Model(chat).
				Set("finished_at = NULL").
				Where("id = ?", chat.ID).
				Exec(ctx)
			if err != nil {
				utils.Log(utils.ErrorLevel).CID(request.SessionID).RID(request.MessageID).Err(err).Send("Failed to reactivate chat")
				return err
			}
		}
		return nil
	})
	if err != nil {
		utils.Log(utils.ErrorLevel).CID(request.SessionID).RID(request.MessageID).Err(err).Send("Failed to run transaction")
		return wspkg.ResponseWrapper{}, false, err
	}

	_, err = h.upsertChattingEndFutureJob(context.Background(), request.SessionID, request.UserID, ChatAutoFinishAfter)
	if err != nil {
		utils.Log(utils.ErrorLevel).CID(request.SessionID).RID(request.MessageID).Err(err).Send("Failed to upsert future job")
		return wspkg.ResponseWrapper{}, false, err
	}

	manager := newChatHistoryManager(h.deps.DB, request.SessionID)
	conversation, err := h.deps.LLM.StartConversation(context.Background(), manager, "interactive_chat", request.SessionID)
	if err != nil {
		utils.Log(utils.ErrorLevel).CID(request.SessionID).RID(request.MessageID).Err(err).Send("Failed to start conversation")
		return wspkg.ResponseWrapper{}, false, err
	}
	llmResponse, err := conversation.Request(context.Background(), llm.Message{
		ConversationID: request.SessionID,
		ID:             request.MessageID,
		Role:           llm.RoleUser,
		Content:        request.Payload.(string),
	})
	if err != nil {
		utils.Log(utils.ErrorLevel).CID(request.SessionID).RID(request.MessageID).Err(err).Send("Failed to request to conversation")
		return wspkg.ResponseWrapper{}, false, err
	}
	marshaled, _ := json.Marshal(llmResponse.Metadata)
	utils.Log(utils.DebugLevel).CID(request.SessionID).RID(request.MessageID).Send("Message ID: %v", llmResponse.ID)
	utils.Log(utils.DebugLevel).CID(request.SessionID).RID(request.MessageID).Send("Message Metadata: %v", string(marshaled))

	response = wspkg.BuildResponseFrom(
		request, llmResponse.ID,
		wspkg.PredefinedActionData, llmResponse.Content,
	)

	return response, shouldClose, nil
}

func (h *ChatHandler) HandleConnect(c *fiberws.Conn, request wspkg.ConnectWrapper) (wspkg.ResponseWrapper, error) {
	if !request.Authorized || !h.checkAuthorization(request.UserID) {
		utils.Log(utils.InfoLevel).CID(request.SessionID).Send("Unauthorized")
		return wspkg.BuildRejectResponse(request), nil
	}

	chat := &domain.Chat{}
	err := h.deps.DB.NewSelect().
		Model(chat).
		Where("user_id = ? AND session_id = ?", request.UserID, request.SessionID).
		Scan(context.Background())
	if err != nil && err != sql.ErrNoRows {
		utils.Log(utils.ErrorLevel).CID(request.SessionID).Err(err).Send("Failed to query chat")
		return wspkg.ResponseWrapper{}, err
	}

	if err == sql.ErrNoRows {
		chat = &domain.Chat{
			UserID:       request.UserID,
			SessionID:    request.SessionID,
			StartedDate:  h.deps.Clock.Now().Truncate(24 * time.Hour),
			UserTimezone: h.deps.Clock.Now().Location().String(),
		}
		_, err = h.deps.DB.NewInsert().Model(chat).Exec(context.Background())
		if err != nil {
			utils.Log(utils.ErrorLevel).CID(request.SessionID).Err(err).Send("Failed to create chat")
			return wspkg.ResponseWrapper{}, err
		}
		return wspkg.BuildResponseFrom(
			request, uuid.New().String(),
			wspkg.PredefinedActionNotify, ChatPayloadNotifyNewConversation,
		), nil
	}

	if !chat.FinishedAt.IsZero() && chat.ArchivedAt.IsZero() {
		return wspkg.BuildResponseFrom(
			request, uuid.New().String(),
			wspkg.PredefinedActionNotify, ChatPayloadNotifyConversationFinished,
		), nil
	}

	if !chat.ArchivedAt.IsZero() {
		return wspkg.BuildCloseResponse(
			request, wspkg.PredefinedActionNotify, ChatPayloadNotifyConversationArchived,
		), nil
	}

	return wspkg.BuildResponseFrom(
		request, uuid.New().String(),
		wspkg.PredefinedActionNotify, ChatPayloadNotifyExistingConversation,
	), nil
}

func (h *ChatHandler) HandleClose(_ *fiberws.Conn, _ wspkg.CloseWrapper) {}

func (h *ChatHandler) HandlePing(c *fiberws.Conn, request wspkg.PingWrapper) (wspkg.ResponseWrapper, bool, error) {
	if !request.Authorized || !h.checkAuthorization(request.UserID) {
		utils.Log(utils.InfoLevel).CID(request.SessionID).Send("Unauthorized")
		return wspkg.BuildRejectResponse(request), false, nil
	}

	chat := &domain.Chat{}
	err := h.deps.DB.NewSelect().
		Model(chat).
		Where("user_id = ? AND session_id = ?", request.UserID, request.SessionID).
		Scan(context.Background())
	if err != nil {
		utils.Log(utils.ErrorLevel).CID(request.SessionID).Err(err).Send("Failed to query chat")
		return wspkg.ResponseWrapper{}, false, err
	}

	if !chat.FinishedAt.IsZero() && chat.ArchivedAt.IsZero() {
		return wspkg.BuildResponseFrom(
			request, uuid.New().String(),
			wspkg.PredefinedActionNotify, ChatPayloadNotifyConversationFinished,
		), false, nil
	}

	if !chat.ArchivedAt.IsZero() {
		return wspkg.BuildCloseResponse(
			request, wspkg.PredefinedActionNotify, ChatPayloadNotifyConversationArchived,
		), false, nil
	}

	return wspkg.BuildNoopResponse(request), false, nil
}

func (h *ChatHandler) checkAuthorization(userID int64) bool {
	user := &user.User{ID: userID}
	err := h.deps.DB.NewSelect().Model(user).Where("id = ?", userID).Scan(context.Background())
	return err == nil
}

func (h *ChatHandler) upsertChattingEndFutureJob(
	ctx context.Context, sessionID string, userID int64, triggerAfter time.Duration,
) (*ftpkg.Job, error) {
	job, err := h.deps.Future.FindBy(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		marshaled, err := json.Marshal(future.ChatFutureHandlerParams{
			UserID:         userID,
			ConversationID: sessionID,
		})
		if err != nil {
			return nil, err
		}
		job, err = h.deps.Future.Create(
			ctx, future.ChatJobType, string(marshaled),
			triggerAfter,
			sessionID,
		)
		if err != nil {
			return nil, err
		}
		return job, nil
	} else {
		err = h.deps.Future.Reschdule(
			ctx, job.ID, triggerAfter, true,
		)
		if err != nil {
			return nil, err
		}
		return job, nil
	}
}

type chatHistoryManager struct {
	db             *bun.DB
	memoryCache    []llm.Message
	conversationID string
}

func newChatHistoryManager(db *bun.DB, conversationID string) *chatHistoryManager {
	return &chatHistoryManager{
		db:             db,
		memoryCache:    []llm.Message{},
		conversationID: conversationID,
	}
}

func (h *chatHistoryManager) Add(ctx context.Context, messages ...llm.Message) {
	h.memoryCache = append(h.memoryCache, messages...)
	h.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		chat := domain.Chat{}
		err := tx.NewSelect().Model(&chat).Where("session_id = ?", h.conversationID).Scan(ctx)
		if err != nil || chat.ID == 0 {
			utils.Log(utils.WarnLevel).CID(h.conversationID).Err(err).Send("Failed to query chat")
			return err
		}

		histories := make([]domain.History, 0, len(messages))
		for _, message := range messages {
			marshaled, err := json.Marshal(message.Metadata)
			if err != nil {
				utils.Log(utils.WarnLevel).CID(h.conversationID).Err(err).Send("Failed to marshal metadata")
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
				utils.Log(utils.WarnLevel).CID(h.conversationID).Err(err).Send("Failed to bulk insert chat history")
				return err
			}
		}
		return nil
	})
}

func (h *chatHistoryManager) Get(ctx context.Context, conversationID string) []llm.Message {
	histories := []llm.Message{}
	err := h.db.NewSelect().Model(&domain.History{}).Where("session_id = ?", conversationID).Scan(ctx, &histories)
	if err != nil {
		utils.Log(utils.WarnLevel).CID(conversationID).Err(err).Send("Failed to get chat history")
		return h.memoryCache
	}
	return histories
}

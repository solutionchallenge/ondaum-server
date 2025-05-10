package chat

import (
	"time"

	"github.com/uptrace/bun"
)

type History struct {
	bun.BaseModel `bun:"table:chat_histories,alias:ch"`

	ID         int64     `json:"id" db:"id" bun:"id,pk,autoincrement"`
	Role       string    `json:"role" db:"role" bun:"role"`
	Content    string    `json:"content" db:"content" bun:"content"`
	Metadata   []byte    `json:"metadata" db:"metadata" bun:"metadata,type:json"`
	MessageID  string    `json:"message_id" db:"message_id" bun:"message_id"`
	InsertedAt time.Time `json:"inserted_at" db:"inserted_at" bun:"inserted_at,notnull,default:CURRENT_TIMESTAMP"`
	ChatID     int64     `json:"chat_id" db:"chat_id" bun:"chat_id,notnull"`

	Chat *Chat `json:"chat,omitempty" bun:"rel:belongs-to,join:chat_id=id"`
}

type HistoryDTO struct {
	ID       string    `json:"id"`
	When     time.Time `json:"when"`
	Role     string    `json:"role"`
	Content  string    `json:"content"`
	Metadata []byte    `json:"metadata"`
}

func (h *History) ToHistoryDTO() HistoryDTO {
	return HistoryDTO{
		ID:       h.MessageID,
		When:     h.InsertedAt,
		Role:     h.Role,
		Content:  h.Content,
		Metadata: h.Metadata,
	}
}

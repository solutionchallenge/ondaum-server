package chat

import (
	"time"

	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/uptrace/bun"
)

type Summary struct {
	bun.BaseModel `bun:"table:chat_summaries,alias:cs"`

	ID        int64                  `json:"id" db:"id" bun:"id,pk,autoincrement"`
	ChatID    int64                  `json:"chat_id" db:"chat_id" bun:"chat_id,notnull"`
	Title     string                 `json:"title" db:"title" bun:"title"`
	Text      string                 `json:"text" db:"text" bun:"text"`
	Keywords  []string               `json:"keywords" db:"keywords" bun:"keywords,type:json"`
	Emotions  common.EmotionRateList `json:"emotions" db:"emotions" bun:"emotions,type:json"`
	CreatedAt time.Time              `json:"created_at" db:"created_at" bun:"created_at,notnull,default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at" bun:"updated_at,notnull,default:CURRENT_TIMESTAMP"`

	Chat *Chat `json:"chat,omitempty" bun:"rel:belongs-to,join:chat_id=id"`
}

type SummaryDTO struct {
	Title    string                 `json:"title"`
	Text     string                 `json:"text"`
	Keywords []string               `json:"keywords"`
	Emotions common.EmotionRateList `json:"emotions"`
}

func (s *Summary) ToSummaryDTO() SummaryDTO {
	return SummaryDTO{
		Title:    s.Title,
		Text:     s.Text,
		Keywords: s.Keywords,
		Emotions: s.Emotions,
	}
}

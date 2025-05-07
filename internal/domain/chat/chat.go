package chat

import (
	"strconv"
	"time"

	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
)

type Chat struct {
	bun.BaseModel `bun:"table:chats,alias:c"`
	ID            int64     `json:"id" db:"id" bun:"id,pk,autoincrement"`
	UserID        int64     `json:"user_id" db:"user_id" bun:"user_id,notnull"`
	SessionID     string    `json:"session_id" db:"session_id" bun:"session_id,notnull"`
	StartedDate   time.Time `json:"started_date" db:"started_date" bun:"started_date,notnull"`
	UserTimezone  string    `json:"user_timezone" db:"user_timezone" bun:"user_timezone,notnull"`
	CreatedAt     time.Time `json:"created_at" db:"created_at" bun:"created_at,notnull,default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at" bun:"updated_at,notnull,default:CURRENT_TIMESTAMP"`
	FinishedAt    time.Time `json:"finished_at" db:"finished_at" bun:"finished_at"`
	ArchivedAt    time.Time `json:"archived_at" db:"archived_at" bun:"archived_at"`

	User      *user.User `json:"user,omitempty" bun:"rel:belongs-to,join:user_id=id"`
	Summary   *Summary   `json:"summary,omitempty" bun:"rel:has-one,join:id=chat_id"`
	Histories []*History `json:"histories,omitempty" bun:"rel:has-many,join:id=chat_id"`
}

type ChatWithSimplifiedSummaryDTO struct {
	ID           string                `json:"id"`
	UserID       string                `json:"user_id"`
	SessionID    string                `json:"session_id"`
	StartedDate  string                `json:"started_date"`
	UserTimezone string                `json:"user_timezone"`
	IsFinished   bool                  `json:"is_finished"`
	IsArchived   bool                  `json:"is_archived"`
	Summary      *SimplifiedSummaryDTO `json:"summary,omitempty"`
}

type ChatWithSimplifiedSummaryAndHistoriesDTO struct {
	ChatWithSimplifiedSummaryDTO
	Histories []HistoryDTO `json:"histories"`
}

func (c *Chat) ToChatWithSimplifiedSummaryDTO() ChatWithSimplifiedSummaryDTO {
	summary := (*SimplifiedSummaryDTO)(nil)
	if c.Summary != nil {
		dto := c.Summary.ToSimplifiedSummaryDTO()
		summary = &dto
	}
	return ChatWithSimplifiedSummaryDTO{
		ID:           strconv.FormatInt(c.ID, 10),
		UserID:       strconv.FormatInt(c.UserID, 10),
		SessionID:    c.SessionID,
		StartedDate:  c.StartedDate.Format(utils.TIME_FORMAT_DATE),
		UserTimezone: c.UserTimezone,
		IsFinished:   !c.FinishedAt.IsZero(),
		IsArchived:   c.ArchivedAt.IsZero(),
		Summary:      summary,
	}
}

func (c *Chat) ToChatWithSimplifiedSummaryAndHistoriesDTO() ChatWithSimplifiedSummaryAndHistoriesDTO {
	histories := make([]HistoryDTO, len(c.Histories))
	if len(c.Histories) > 0 {
		histories = utils.Map(c.Histories, func(h *History) HistoryDTO {
			return h.ToHistoryDTO()
		})
	}
	return ChatWithSimplifiedSummaryAndHistoriesDTO{
		ChatWithSimplifiedSummaryDTO: c.ToChatWithSimplifiedSummaryDTO(),
		Histories:                    histories,
	}
}

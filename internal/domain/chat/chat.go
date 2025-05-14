package chat

import (
	"strconv"
	"time"

	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
)

type Chat struct {
	bun.BaseModel `bun:"table:chats,alias:c"`
	ID            int64                   `json:"id" db:"id" bun:"id,pk,autoincrement"`
	UserID        int64                   `json:"user_id" db:"user_id" bun:"user_id,notnull"`
	SessionID     string                  `json:"session_id" db:"session_id" bun:"session_id,notnull"`
	StartedDate   time.Time               `json:"started_date" db:"started_date" bun:"started_date,notnull"`
	UserTimezone  string                  `json:"user_timezone" db:"user_timezone" bun:"user_timezone,notnull"`
	CreatedAt     time.Time               `json:"created_at" db:"created_at" bun:"created_at,notnull,default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time               `json:"updated_at" db:"updated_at" bun:"updated_at,notnull,default:CURRENT_TIMESTAMP"`
	FinishedAt    bun.NullTime            `json:"finished_at" db:"finished_at" bun:"finished_at"`
	ArchivedAt    bun.NullTime            `json:"archived_at" db:"archived_at" bun:"archived_at"`
	ChatDuration  common.NullableDuration `json:"chat_duration" db:"chat_duration" bun:"chat_duration"`

	User      *user.User `json:"user,omitempty" bun:"rel:belongs-to,join:user_id=id"`
	Summary   *Summary   `json:"summary,omitempty" bun:"rel:has-one,join:id=chat_id"`
	Histories []*History `json:"histories,omitempty" bun:"rel:has-many,join:id=chat_id"`
}

type ChatDTO struct {
	ID           string      `json:"id"`
	UserID       string      `json:"user_id"`
	SessionID    string      `json:"session_id"`
	StartedDate  string      `json:"started_date"`
	UserTimezone string      `json:"user_timezone"`
	ChatDuration string      `json:"chat_duration"`
	IsFinished   bool        `json:"is_finished"`
	IsArchived   bool        `json:"is_archived"`
	Summary      *SummaryDTO `json:"summary,omitempty"`
}

func (c *Chat) ToChatDTO() ChatDTO {
	summary := (*SummaryDTO)(nil)
	if c.Summary != nil {
		dto := c.Summary.ToSummaryDTO(c.Histories)
		summary = &dto
	}
	return ChatDTO{
		ID:           strconv.FormatInt(c.ID, 10),
		UserID:       strconv.FormatInt(c.UserID, 10),
		SessionID:    c.SessionID,
		StartedDate:  c.StartedDate.Format(utils.TIME_FORMAT_DATE),
		UserTimezone: c.UserTimezone,
		ChatDuration: c.ChatDuration.ToString(common.DurationFormatTime),
		IsFinished:   !c.FinishedAt.IsZero(),
		IsArchived:   !c.ArchivedAt.IsZero(),
		Summary:      summary,
	}
}

type ChatDTOWithHistory struct {
	ChatDTO
	Histories *[]HistoryDTO `json:"histories,omitempty"`
}

func (c *Chat) ToChatDTOWithHistory() ChatDTOWithHistory {
	base := c.ToChatDTO()
	histories := (*[]HistoryDTO)(nil)
	if len(c.Histories) > 0 {
		converted := utils.Map(c.Histories, func(h *History) HistoryDTO {
			return h.ToHistoryDTO()
		})
		histories = &converted
	}
	return ChatDTOWithHistory{
		ChatDTO:   base,
		Histories: histories,
	}
}

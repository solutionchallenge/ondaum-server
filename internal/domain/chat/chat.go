package chat

import (
	"time"

	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
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

	User      *user.User     `json:"user,omitempty" bun:"rel:belongs-to,join:user_id=id"`
	Histories []*ChatHistory `json:"histories,omitempty" bun:"rel:has-many,join:id=chat_id"`
}

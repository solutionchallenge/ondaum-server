package user

import (
	"time"

	"github.com/solutionchallenge/ondaum-server/pkg/oauth"
	"github.com/uptrace/bun"
)

type OAuth struct {
	bun.BaseModel `bun:"table:user_oauths,alias:uo"`

	ID           int64          `json:"id" db:"id" bun:"id,pk,autoincrement"`
	UserID       int64          `json:"user_id" db:"user_id" bun:"user_id,notnull"`
	Provider     oauth.Provider `json:"provider" db:"provider" bun:"provider,notnull"`
	ProviderCode string         `json:"provider_code" db:"provider_code" bun:"provider_code,notnull"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at" bun:"created_at,notnull,default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at" bun:"updated_at,notnull,default:CURRENT_TIMESTAMP"`

	User *User `json:"user,omitempty" bun:"rel:belongs-to,join:user_id=id"`
}

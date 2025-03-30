package user

import (
	"time"

	"github.com/solutionchallenge/ondaum-server/pkg/oauth"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64     `json:"id" db:"id" bun:"id,pk,autoincrement"`
	Email     string    `json:"email" db:"email" bun:"email,unique"`
	Username  string    `json:"username" db:"username" bun:"username"`
	CreatedAt time.Time `json:"created_at" db:"created_at" bun:"created_at,notnull,default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" bun:"updated_at,notnull,default:CURRENT_TIMESTAMP"`

	OAuths []*UserOAuth `json:"oauths,omitempty" bun:"rel:has-many,join:id=user_id"`
}

func (u *User) GetOAuth(provider oauth.Provider) *UserOAuth {
	for _, oauth := range u.OAuths {
		if oauth.Provider == provider {
			return oauth
		}
	}
	return nil
}

package user

import (
	"time"

	"github.com/uptrace/bun"
)

type UserGender string

const (
	UserGenderMale   UserGender = "male"
	UserGenderFemale UserGender = "female"
	UserGenderOther  UserGender = "other"
)

type UserPrivacy struct {
	bun.BaseModel `bun:"table:user_privacies,alias:up"`

	UserID    int64      `json:"user_id" db:"user_id" bun:"user_id,pk"`
	Gender    UserGender `json:"gender" db:"gender" bun:"gender,notnull"`
	Birthday  time.Time  `json:"birthday" db:"birthday" bun:"birthday,notnull"`
	CreatedAt time.Time  `json:"created_at" db:"created_at" bun:"created_at,notnull,default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at" bun:"updated_at,notnull,default:CURRENT_TIMESTAMP"`

	User *User `bun:"rel:belongs-to,join:user_id=id"`
}

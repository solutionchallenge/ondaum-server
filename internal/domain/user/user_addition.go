package user

import (
	"time"

	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/uptrace/bun"
)

type Addition struct {
	bun.BaseModel `bun:"table:user_additions,alias:ua"`

	UserID    int64              `json:"user_id" db:"user_id" bun:"user_id,pk"`
	Concerns  []string           `json:"concerns" db:"concerns" bun:"concerns,type:json,notnull"`
	Emotions  common.EmotionList `json:"emotions" db:"emotions" bun:"emotions,type:json,notnull"`
	CreatedAt time.Time          `json:"created_at" db:"created_at" bun:"created_at,notnull,default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time          `json:"updated_at" db:"updated_at" bun:"updated_at,notnull,default:CURRENT_TIMESTAMP"`

	User *User `json:"user,omitempty" bun:"rel:belongs-to,join:user_id=id"`
}

type SimplifiedAdditionDTO struct {
	Concerns []string           `json:"concerns"`
	Emotions common.EmotionList `json:"emotions"`
}

func (a *Addition) ToSimplifiedAdditionDTO() SimplifiedAdditionDTO {
	return SimplifiedAdditionDTO{}
}

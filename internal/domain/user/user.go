package user

import (
	"encoding/json"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/solutionchallenge/ondaum-server/pkg/oauth"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64     `json:"id" db:"id" bun:"id,pk,autoincrement"`
	Email     string    `json:"email" db:"email" bun:"email,unique"`
	Username  string    `json:"username" db:"username" bun:"username"`
	CreatedAt time.Time `json:"created_at" db:"created_at" bun:"created_at,notnull,default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" bun:"updated_at,notnull,default:CURRENT_TIMESTAMP"`

	OAuths   []*OAuth  `json:"oauths,omitempty" bun:"rel:has-many,join:id=user_id"`
	Privacy  *Privacy  `json:"privacy,omitempty" bun:"rel:has-one,join:id=user_id"`
	Addition *Addition `json:"addition,omitempty" bun:"rel:has-one,join:id=user_id"`
}

type UserDTO struct {
	ID       int64        `json:"id"`
	Email    string       `json:"email"`
	Username string       `json:"username"`
	Privacy  *PrivacyDTO  `json:"privacy,omitempty"`
	Addition *AdditionDTO `json:"addition,omitempty"`
}

func (u *User) ToUserDTO() UserDTO {
	privacy := (*PrivacyDTO)(nil)
	if u.Privacy != nil {
		dto := u.Privacy.ToPrivacyDTO()
		privacy = &dto
	}
	addition := (*AdditionDTO)(nil)
	if u.Addition != nil {
		dto := u.Addition.ToAdditionDTO()
		addition = &dto
	}
	return UserDTO{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
		Privacy:  privacy,
		Addition: addition,
	}
}

func (u *User) GetOAuth(provider oauth.Provider) *OAuth {
	for _, oauth := range u.OAuths {
		if oauth.Provider == provider {
			return oauth
		}
	}
	return nil
}

type UserMentalStateHint struct {
	Today    string              `json:"today"`
	Username *string             `json:"username,omitempty"`
	Gender   *UserGender         `json:"gender,omitempty"`
	Birthday *string             `json:"birthday,omitempty"`
	Concerns *[]string           `json:"concerns,omitempty"`
	Emotions *common.EmotionList `json:"emotions,omitempty"`
}

func (umsh *UserMentalStateHint) Marshal() string {
	marshaled, err := json.Marshal(umsh)
	if err != nil {
		return ""
	}
	return string(marshaled)
}

func (u *User) ToUserMentalStateHint(clk clock.Clock) *UserMentalStateHint {
	hint := &UserMentalStateHint{
		Today: clk.Now().Format(utils.TIME_FORMAT_DATE),
	}
	if u.Username != "" {
		hint.Username = &u.Username
	}
	if u.Privacy != nil {
		birthday := u.Privacy.Birthday.Format(utils.TIME_FORMAT_DATE)
		hint.Gender = &u.Privacy.Gender
		hint.Birthday = &birthday
	}
	if u.Addition != nil {
		hint.Concerns = &u.Addition.Concerns
		hint.Emotions = &u.Addition.Emotions
	}
	return hint
}

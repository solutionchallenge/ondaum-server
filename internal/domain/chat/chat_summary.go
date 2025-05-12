package chat

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
)

type MainTopic struct {
	BeginMessageID string `json:"begin_message_id" bun:"begin_message_id"`
	EndMessageID   string `json:"end_message_id" bun:"end_message_id"`
}

func (m *MainTopic) Scan(src interface{}) error {
	bytes, ok := src.([]byte)
	if !ok {
		return utils.NewError("failed to type assert MainTopic")
	}
	var result MainTopic
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}
	*m = result
	return nil
}

func (m *MainTopic) Value() (driver.Value, error) {
	buf, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return string(buf), nil
}

func (m *MainTopic) Validate() bool {
	return m.BeginMessageID != "" && m.EndMessageID != ""
}

func (m *MainTopic) ToString() string {
	result, _ := m.Value()
	return result.(string)
}

type Summary struct {
	bun.BaseModel `bun:"table:chat_summaries,alias:cs"`

	ID              int64                  `json:"id" db:"id" bun:"id,pk,autoincrement"`
	ChatID          int64                  `json:"chat_id" db:"chat_id" bun:"chat_id,notnull"`
	Title           string                 `json:"title" db:"title" bun:"title"`
	Text            string                 `json:"text" db:"text" bun:"text"`
	Keywords        []string               `json:"keywords" db:"keywords" bun:"keywords,type:json"`
	Emotions        common.EmotionRateList `json:"emotions" db:"emotions" bun:"emotions,type:json"`
	Recommendations []string               `json:"recommendations" db:"recommendations" bun:"recommendations,type:json"`
	PositiveScore   float64                `json:"positive_score" db:"positive_score" bun:"positive_score"`
	NegativeScore   float64                `json:"negative_score" db:"negative_score" bun:"negative_score"`
	NeutralScore    float64                `json:"neutral_score" db:"neutral_score" bun:"neutral_score"`
	MainTopic       MainTopic              `json:"main_topic" db:"main_topic" bun:"main_topic,type:json"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at" bun:"created_at,notnull,default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at" bun:"updated_at,notnull,default:CURRENT_TIMESTAMP"`

	Chat *Chat `json:"chat,omitempty" bun:"rel:belongs-to,join:chat_id=id"`
}

type SummaryDTO struct {
	Title           string                 `json:"title"`
	Text            string                 `json:"text"`
	Keywords        []string               `json:"keywords"`
	Emotions        common.EmotionRateList `json:"emotions"`
	Recommendations []string               `json:"recommendations"`
	PositiveScore   float64                `json:"positive_score"`
	NegativeScore   float64                `json:"negative_score"`
	NeutralScore    float64                `json:"neutral_score"`
	MainTopic       MainTopic              `json:"main_topic"`
}

func (s *Summary) ToSummaryDTO() SummaryDTO {
	return SummaryDTO{
		Title:           s.Title,
		Text:            s.Text,
		Keywords:        s.Keywords,
		Emotions:        s.Emotions,
		Recommendations: s.Recommendations,
		PositiveScore:   s.PositiveScore,
		NegativeScore:   s.NegativeScore,
		NeutralScore:    s.NeutralScore,
		MainTopic:       s.MainTopic,
	}
}

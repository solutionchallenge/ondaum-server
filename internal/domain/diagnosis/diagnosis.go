package diagnosis

import (
	"time"

	"github.com/solutionchallenge/ondaum-server/internal/domain/common"
	"github.com/solutionchallenge/ondaum-server/internal/domain/user"
	"github.com/uptrace/bun"
)

type Diagnosis struct {
	bun.BaseModel `bun:"table:diagnoses,alias:i"`

	ID                int64            `json:"id" db:"id" bun:"id,pk,autoincrement"`
	UserID            int64            `json:"user_id" db:"user_id" bun:"user_id,notnull"`
	Diagnosis         common.Diagnosis `json:"diagnosis" db:"diagnosis" bun:"diagnosis,notnull"`
	TotalScore        int64            `json:"total_score" db:"total_score" bun:"total_score,notnull"`
	ResultScore       int64            `json:"result_score" db:"result_score" bun:"result_score,notnull"`
	ResultName        string           `json:"result_name" db:"result_name" bun:"result_name,notnull"`
	ResultDescription string           `json:"result_description" db:"result_description" bun:"result_description,notnull"`
	ResultCritical    bool             `json:"result_critical" db:"result_critical" bun:"result_critical,notnull"`
	CreatedAt         time.Time        `json:"created_at" db:"created_at" bun:"created_at,notnull,default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time        `json:"updated_at" db:"updated_at" bun:"updated_at,notnull,default:CURRENT_TIMESTAMP"`

	User *user.User `json:"user,omitempty" bun:"rel:belongs-to,join:user_id=id"`
}

type DiagnosisDTO struct {
	ID                int64            `json:"id"`
	Diagnosis         common.Diagnosis `json:"diagnosis"`
	TotalScore        int64            `json:"total_score"`
	ResultScore       int64            `json:"result_score"`
	ResultName        string           `json:"result_name"`
	ResultDescription string           `json:"result_description"`
	ResultCritical    bool             `json:"result_critical"`
}

func (i *Diagnosis) ToDiagnosisDTO() DiagnosisDTO {
	return DiagnosisDTO{
		ID:                i.ID,
		Diagnosis:         i.Diagnosis,
		TotalScore:        i.TotalScore,
		ResultScore:       i.ResultScore,
		ResultName:        i.ResultName,
		ResultDescription: i.ResultDescription,
		ResultCritical:    i.ResultCritical,
	}
}

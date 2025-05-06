package database

import (
	"time"

	"github.com/solutionchallenge/ondaum-server/pkg/future"
	"github.com/uptrace/bun"
)

type FutureJob struct {
	bun.BaseModel    `bun:"future_jobs"`
	ID               string           `bun:"id,pk"`
	ActionIdentifier string           `bun:"action_identifier,notnull"`
	ActionType       future.JobType   `bun:"action_type,notnull"`
	ActionParams     string           `bun:"action_params,notnull"`
	TriggeredAt      time.Time        `bun:"triggered_at,notnull"`
	CompletedAt      time.Time        `bun:"completed_at"`
	Status           future.JobStatus `bun:"status,notnull"`
	Error            string           `bun:"error"`
	CreatedAt        time.Time        `bun:"created_at,notnull,default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time        `bun:"updated_at,notnull,default:CURRENT_TIMESTAMP"`
}

const FutureJobTableCreationSQL = `
CREATE TABLE IF NOT EXISTS future_jobs (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	action_identifier VARCHAR(255) NOT NULL,
	action_type VARCHAR(255) NOT NULL,
	action_params TEXT NOT NULL,
	triggered_at DATETIME NOT NULL,
	completed_at DATETIME,
	status VARCHAR(50) NOT NULL,
	error TEXT,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
`

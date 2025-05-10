package future

import (
	"time"
)

type JobType string
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCanceled  JobStatus = "canceled"
)

type Job struct {
	ID               string    `json:"id"`
	ActionIdentifier string    `json:"action_identifier"`
	ActionType       JobType   `json:"action_type"`
	ActionParams     string    `json:"action_params"`
	TriggeredAt      time.Time `json:"triggered_at"`
	CompletedAt      time.Time `json:"completed_at"`
	Status           JobStatus `json:"status"`
	Error            string    `json:"error"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

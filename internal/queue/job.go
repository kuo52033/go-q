package queue

import "time"

type JobId string

type JobStatus string

const (
	StatusQueued  JobStatus = "queued"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
	StatusDelayed    JobStatus = "delayed"
	StatusRetrying   JobStatus = "retrying"
)

type JobPayload map[string]interface{}

type Job struct {
	ID JobId `json:"id" redis:"id"`
	Type string `json:"type" redis:"type"`
	Payload JobPayload `json:"payload" redis:"payload"`
	Status JobStatus `json:"status" redis:"status"`
	CreatedAt time.Time `json:"created_at" redis:"created_at"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty" redis:"scheduled_at"`
	UpdatedAt time.Time `json:"updated_at" redis:"updated_at"`
	MaxAttempts int `json:"max_attempts" redis:"max_attempts"`
	AttemptCount int `json:"attempt_count" redis:"attempt_count"`
	LastError string `json:"last_error,omitempty" redis:"last_error"`
	ProcessedAt *time.Time `json:"processed_at,omitempty" redis:"processed_at"`
	Queue string `json:"queue" redis:"queue"`
}

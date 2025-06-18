package model

import (
	"encoding/json"
	"time"
)

type JobStatus string

const (
	StatusQueued  JobStatus = "queued"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
	StatusDelayed    JobStatus = "delayed"
	StatusRetrying   JobStatus = "retrying"
)

// 直接將其底層的 string 轉為 []byte 
func (js JobStatus) MarshalBinary() ([]byte, error) {
	return []byte(js), nil
}

type JobPayload map[string]interface{}

func (p JobPayload) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *JobPayload) ScanRedis(value string) error {
	return json.Unmarshal([]byte(value), p)
}
type Job struct {
	ID string `json:"id" redis:"id"`
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

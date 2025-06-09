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

type JobPayload map[string]interface{}

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

func (j *Job) ToMap() (map[string]interface{} ,error){
	jsonPayloadBytes, err := json.Marshal(j.Payload)
	if err != nil {
		return nil, err
	}

	jobMap := map[string]interface{}{
		"id": j.ID,
		"type": j.Type,
		"payload": jsonPayloadBytes,
		"status": string(j.Status),
		"created_at": j.CreatedAt.Format(time.RFC3339),
		"updated_at": j.UpdatedAt.Format(time.RFC3339),
		"max_attempts": j.MaxAttempts,
		"attempt_count": j.AttemptCount,
		"last_error": j.LastError,
		"queue": j.Queue,
	}

	if j.ScheduledAt != nil {
		jobMap["scheduled_at"] = j.ScheduledAt.Format(time.RFC3339)
	}

	if j.ProcessedAt != nil {
		jobMap["processed_at"] = j.ProcessedAt.Format(time.RFC3339)
	}


	return jobMap, nil
}

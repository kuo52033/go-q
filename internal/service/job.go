package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kuo-52033/go-q/internal/model"
	"github.com/kuo-52033/go-q/internal/myerror"
	"github.com/kuo-52033/go-q/internal/store"
)

type JobService interface {
	CreateJob(
		ctx context.Context, 
		jobType string, 
		payload model.JobPayload, 
		queueName string, 
		maxAttempts int,
	) (*model.Job, error)
}
type jobService struct {
	jobStore store.JobStore
}

func NewJobService(jobStore store.JobStore) JobService {
	return &jobService{jobStore: jobStore}
}

func (s *jobService) CreateJob(
	ctx context.Context, 
	jobType string, 
	payload model.JobPayload, 
	queueName string, 
	maxAttempts int,
) (*model.Job, error) {
	job := &model.Job{
		ID: uuid.New().String(),
		Type: jobType,
		Payload: payload,
		Status: model.StatusQueued,
		Queue: queueName,
		CreatedAt: time.Now(),
		AttemptCount: 0,
		MaxAttempts: maxAttempts,
	}	

	if err := s.jobStore.CreateJob(ctx, job); err != nil {
		return nil, myerror.InternalServerError(myerror.JOB_CREATE_FAILED, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if err := s.jobStore.PushJobToQueue(ctx, queueName, job.ID); err != nil {
		return nil, myerror.InternalServerError(myerror.JOB_CREATE_FAILED, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return job, nil	
}

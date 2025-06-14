package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kuo-52033/go-q/internal/model"
	"github.com/kuo-52033/go-q/internal/myerror"
)

type JobStore interface {
	SaveJob(ctx context.Context, job *model.Job) error
	EnqueueJobId(ctx context.Context, queueName string, jobID string) error
}

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
		jobStore JobStore
}

func NewJobService(jobStore JobStore) JobService {
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

	if err := s.jobStore.SaveJob(ctx, job); err != nil {
		return nil, myerror.InternalServerError(myerror.JOB_CREATE_FAILED, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if err := s.jobStore.EnqueueJobId(ctx, queueName, job.ID); err != nil {
		return nil, myerror.InternalServerError(myerror.JOB_CREATE_FAILED, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return job, nil	
}

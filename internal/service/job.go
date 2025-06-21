package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kuo-52033/go-q/internal/model"
	"github.com/kuo-52033/go-q/internal/utils/myerror"
)

type JobStore interface {
	SaveJob(ctx context.Context, job *model.Job) error
	EnqueueJobId(ctx context.Context, queueName string, jobID string) error
	DequeueJobId(ctx context.Context, queueName string, timeout time.Duration) (string, error)
	GetJobById(ctx context.Context, jobID string) (*model.Job, error)
	UpdateJobStatus(ctx context.Context, jobID string, status model.JobStatus) error
}

type JobService interface {
	CreateJob(
		ctx context.Context, 
		jobType string, 
		payload model.JobPayload, 
		queueName string, 
		maxAttempts int,
	) (*model.Job, error)
	GetJobById(ctx context.Context, jobID string) (*model.Job, error)
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
		return nil, myerror.InternalServerError(myerror.JOB_CREATE_FAILED, myerror.WithError(err))
	}

	if err := s.jobStore.EnqueueJobId(ctx, queueName, job.ID); err != nil {
		return nil, myerror.InternalServerError(myerror.JOB_CREATE_FAILED, myerror.WithError(err))
	}

	return job, nil	
}

func (s *jobService) GetJobById(ctx context.Context, jobID string) (*model.Job, error) {
	job, err := s.jobStore.GetJobById(ctx, jobID)
	if err != nil {
		return nil, myerror.InternalServerError(myerror.JOB_GET_FAILED, myerror.WithError(err))
	}
	return job, nil
}

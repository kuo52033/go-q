package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kuo-52033/go-q/internal/db"
	"github.com/kuo-52033/go-q/internal/error"
	"github.com/kuo-52033/go-q/internal/model"
)

type JobService interface {
	CreateJob(
		ctx context.Context, 
		jobType string, 
		payload model.JobPayload, 
		queueName string, 
		maxAttempts int,
	) (*model.Job, *error.Error)
}

type jobService struct {
	rdb db.RedisClient
}

func NewJobService(rdb db.RedisClient) JobService {
	return &jobService{rdb: rdb}
}

func (s *jobService) CreateJob(
	ctx context.Context, 
	jobType string, 
	payload model.JobPayload, 
	queueName string, 
	maxAttempts int,
) (*model.Job, *error.Error) {
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
	
	key := fmt.Sprintf("job:%s", job.ID)
	jobMap, err := job.ToMap()
	if err != nil {
		return nil, error.InternalServerError(error.JOB_CREATE_FAILED, map[string]interface{}{
			"error": err.Error(),
		})
	}
	err = s.rdb.HMSet(ctx, key, jobMap)
	if err != nil {
		return nil, error.InternalServerError(error.JOB_CREATE_FAILED, map[string]interface{}{
			"error": err.Error(),
		})
	}

	listKey := fmt.Sprintf("queue:%s", queueName)
	err = s.rdb.LPush(ctx, listKey, job.ID)
	if err != nil {
		return nil, error.InternalServerError(error.JOB_CREATE_FAILED, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return job, nil	
}

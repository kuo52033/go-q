package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kuo-52033/go-q/internal/db"
	"github.com/kuo-52033/go-q/internal/model"
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
	
	key := fmt.Sprintf("job:%s", job.ID)
	jobMap, err := job.ToMap()
	if err != nil {
		return nil, fmt.Errorf("failed to convert job to map: %w", err)
	}
	err = s.rdb.HMSet(ctx, key, jobMap)
	if err != nil {
		return nil, fmt.Errorf("failed to set job in hash: %w", err)
	}

	listKey := fmt.Sprintf("queue:%s", queueName)
	err = s.rdb.LPush(ctx, listKey, job.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to push job to list: %w", err)
	}

	return job, nil
}

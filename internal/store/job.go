package store

import (
	"context"
	"fmt"
	"github.com/kuo-52033/go-q/internal/model"
	"github.com/redis/go-redis/v9"
)

type JobStore interface {
	CreateJob(ctx context.Context, job *model.Job) error
	PushJobToQueue(ctx context.Context, queueName string, jobID string) error
}

type jobStore struct {
	rdb *redis.Client
}

func NewJobStore(rdb *redis.Client) JobStore {
	return &jobStore{rdb: rdb}
}

func (s *jobStore) CreateJob(ctx context.Context, job *model.Job) error {
	key := fmt.Sprintf("job:%s", job.ID)
	jobMap, err := job.ToMap()
	if err != nil {
		return err
	}
	return s.rdb.HMSet(ctx, key, jobMap).Err()
}

func (s *jobStore) PushJobToQueue(ctx context.Context, queueName string, jobID string) error {
	listKey := fmt.Sprintf("queue:%s", queueName)
	return s.rdb.LPush(ctx, listKey, jobID).Err()
}

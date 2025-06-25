package store

import (
	"context"
	"fmt"
	"time"

	"github.com/kuo-52033/go-q/internal/model"
	"github.com/redis/go-redis/v9"
)

type RedisJobStore struct {
	rdb *redis.Client
}

func NewRedisJobStore(rdb *redis.Client) *RedisJobStore {
	return &RedisJobStore{rdb: rdb}
}

func (s *RedisJobStore) SaveJob(ctx context.Context, job *model.Job) error {
	key := fmt.Sprintf("job:%s", job.ID)
	return s.rdb.HSet(ctx, key, job).Err()
}

func (s *RedisJobStore) EnqueueJobId(ctx context.Context, queueName string, jobID string) error {
	listKey := fmt.Sprintf("queue:%s", queueName)
	return s.rdb.LPush(ctx, listKey, jobID).Err()
}


func (s *RedisJobStore) DequeueJobId(ctx context.Context, queueName string, timeout time.Duration) (string, error) {
	listKey := fmt.Sprintf("queue:%s", queueName)
	jobID, err := s.rdb.BRPop(ctx, timeout, listKey).Result()
	if err != nil {
		return "", err
	}
	return jobID[1], nil
}

func (s *RedisJobStore) ReEnqueueJobId(ctx context.Context, queueName string, jobID string) error {
	listKey := fmt.Sprintf("queue:%s", queueName)
	return s.rdb.RPush(ctx, listKey, jobID).Err()
}

func (s *RedisJobStore) GetJobById(ctx context.Context, jobID string) (*model.Job, error) {
	var job model.Job
	key := fmt.Sprintf("job:%s", jobID)
	err := s.rdb.HGetAll(ctx, key).Scan(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (s *RedisJobStore) UpdateJobStatus(ctx context.Context, jobID string, status model.JobStatus) error {
	key := fmt.Sprintf("job:%s", jobID)
	return s.rdb.HSet(ctx, key, "status", status).Err()
}

package store

import (
	"context"
	"fmt"
	"github.com/kuo-52033/go-q/internal/model"
	"github.com/redis/go-redis/v9"
)

type RedisJobStore struct {
	rdb *redis.Client
}

func NewJobStore(rdb *redis.Client) *RedisJobStore {
	return &RedisJobStore{rdb: rdb}
}

func (s *RedisJobStore) SaveJobHash(ctx context.Context, job *model.Job) error {
	key := fmt.Sprintf("job:%s", job.ID)
	jobMap, err := job.ToMap()
	if err != nil {
		return err
	}
	return s.rdb.HMSet(ctx, key, jobMap).Err()
}

func (s *RedisJobStore) EnqueueJobId(ctx context.Context, queueName string, jobID string) error {
	listKey := fmt.Sprintf("queue:%s", queueName)
	return s.rdb.LPush(ctx, listKey, jobID).Err()
}

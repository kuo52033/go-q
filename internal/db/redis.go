package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	Close() error
	HMSet(ctx context.Context, key string, values ...interface{}) error
	LPush(ctx context.Context, key string, values ...interface{}) error
}

type redisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr string) (RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &redisClient{Client: rdb}, nil
}

func (r *redisClient) Close() error {
	return r.Client.Close()
}

func (r *redisClient) LPush(ctx context.Context, key string, values ...interface{}) error {
	return r.Client.LPush(ctx, key, values...).Err()
}

func (r *redisClient) HMSet(ctx context.Context, key string, value ...interface{}) error {
	return r.Client.HMSet(ctx, key, value...).Err()
}

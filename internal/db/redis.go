package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{
		client: client,
		ctx:    context.Background(),
	}
}

func (c *RedisClient) Close() error {
	return c.client.Close()
}

func (c *RedisClient) Ping() (string, error) {
	return c.client.Ping(c.ctx).Result()
}

func ConnectRedis(addr string) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	r := NewRedisClient(client)
	defer r.Close()

	if _, err := r.Ping(); err != nil {
		return nil, err
	}

	return r, nil
}

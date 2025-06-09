package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{
		Client: client,
		Ctx:    context.Background(),
	}
}

func (c *RedisClient) Close() error {
	return c.Client.Close()
}

func (c *RedisClient) Ping() (string, error) {
	return c.Client.Ping(c.Ctx).Result()
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

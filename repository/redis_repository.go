package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	Set(client *redis.Client, ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	Get(client *redis.Client, ctx context.Context, key string) (string, error)
	Del(client *redis.Client, ctx context.Context, key string) (int64, error)
	FlushDb(client *redis.Client, ctx context.Context) (string, error)
}

type RedisRepositoryImplementation struct {
}

func NewRedisRepository() RedisRepository {
	return &RedisRepositoryImplementation{}
}

func (repository *RedisRepositoryImplementation) Set(client *redis.Client, ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return client.Set(ctx, key, value, expiration).Result()
}

func (repository *RedisRepositoryImplementation) Get(client *redis.Client, ctx context.Context, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func (repository *RedisRepositoryImplementation) Del(client *redis.Client, ctx context.Context, key string) (int64, error) {
	return client.Del(ctx, key).Result()
}

func (repository *RedisRepositoryImplementation) FlushDb(client *redis.Client, ctx context.Context) (string, error) {
	return client.FlushDB(ctx).Result()
}

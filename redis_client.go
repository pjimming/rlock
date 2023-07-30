package rlock

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
	Ping(ctx context.Context) *redis.StatusCmd
	Context() context.Context
}

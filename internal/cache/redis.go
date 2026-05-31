package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(addr string) *Redis {
	return &Redis{
		client: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (r *Redis) SetSession(ctx context.Context, sessionID string, ttl time.Duration) error {
	return r.client.Set(ctx, "session:"+sessionID, 1, ttl).Err()
}

func (r *Redis) SessionExists(ctx context.Context, sessionID string) (bool, error) {
	n, err := r.client.Exists(ctx, "session:"+sessionID).Result()
	return n > 0, err
}

func (r *Redis) DeleteSession(ctx context.Context, sessionID string) error {
	return r.client.Del(ctx, "session:"+sessionID).Err()
}

func (r *Redis) IncrRateLimit(ctx context.Context, key string, window time.Duration) (int64, error) {
	pipe := r.client.Pipeline()
	incr := pipe.Incr(ctx, "rl:"+key)
	pipe.Expire(ctx, "rl:"+key, window)
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, err
	}
	return incr.Val(), nil
}

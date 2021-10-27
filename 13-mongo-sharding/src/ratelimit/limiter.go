package ratelimit

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const namespace = "rl"

func NewLimiter(client *redis.Client, action string, period time.Duration, limit int64) *Limiter {

	return &Limiter{client, action, period, limit}
}

type Limiter struct {
	client *redis.Client

	action string
	period time.Duration
	limit  int64
}

//go:embed incr_expirenx.lua
var incrExpireLua string
var incrExpireScript = redis.NewScript(incrExpireLua)

func (l *Limiter) CanDoAt(ctx context.Context, ts time.Time) (bool, error) {
	key := l.key(ts)
	ttlMs := l.period.Milliseconds()

	rawCount, err := incrExpireScript.Run(ctx, l.client, []string{key}, ttlMs).Result()
	if err != nil {
		return false, err
	}
	count := rawCount.(int64)

	return count <= l.limit, nil
}

func (l *Limiter) key(ts time.Time) string {
	interval := ts.UTC().UnixNano() / l.period.Nanoseconds()
	return fmt.Sprintf("%s:%s:%x", namespace, l.action, interval)
}

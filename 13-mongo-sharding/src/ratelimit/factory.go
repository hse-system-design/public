package ratelimit

import (
	"github.com/go-redis/redis/v8"
	"time"
)

func NewFactory(client *redis.Client) *Factory {
	return &Factory{client}
}

type Factory struct {
	client *redis.Client
}

func (f *Factory) NewLimiter(action string, period time.Duration, limit int64) *Limiter {
	return NewLimiter(f.client, action, period, limit)
}

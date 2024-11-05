package redis

import (
	"fmt"
	"time"

	"github.com/codemaestro64/cachey/store"
)

func WithAddress(address string) store.Option {
	return func(s store.Store) error {
		redisStore, ok := s.(*RedisStore)
		if !ok {
			return fmt.Errorf("invalid store type for redis options")
		}

		redisStore.config.address = address
		return nil
	}
}

func WithDB(db int) store.Option {
	return func(s store.Store) error {
		redisStore, ok := s.(*RedisStore)
		if !ok {
			return fmt.Errorf("invalid store type for redis options")
		}

		redisStore.config.db = db
		return nil
	}
}

func WithMaxRetries(maxRetries int) store.Option {
	return func(s store.Store) error {
		redisStore, ok := s.(*RedisStore)
		if !ok {
			return fmt.Errorf("invalid store type for redis options")
		}

		redisStore.config.maxRetries = maxRetries
		return nil
	}
}

func WithReadTimeout(timeout time.Duration) store.Option {
	return func(s store.Store) error {
		redisStore, ok := s.(*RedisStore)
		if !ok {
			return fmt.Errorf("invalid store type for redis options")
		}

		redisStore.config.readTimeout = timeout
		return nil
	}
}

func WithWriteTimeout(timeout time.Duration) store.Option {
	return func(s store.Store) error {
		redisStore, ok := s.(*RedisStore)
		if !ok {
			return fmt.Errorf("invalid store type for redis options")
		}

		redisStore.config.writeTimeout = timeout
		return nil
	}
}

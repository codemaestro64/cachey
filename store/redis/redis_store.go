package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/codemaestro64/cachey/store"
	"github.com/redis/go-redis/v9"
)

type config struct {
	address      string
	password     string
	db           int
	maxRetries   int
	readTimeout  time.Duration
	writeTimeout time.Duration
}

type RedisStore struct {
	config *config
	store  *redis.Client
}

func NewRedisStore() store.Store {
	defaultConfig := config{
		address:      "localhost:6379",
		password:     "",
		db:           0,
		maxRetries:   5,
		readTimeout:  10 * time.Second,
		writeTimeout: 10 * time.Second,
	}

	return &RedisStore{
		config: &defaultConfig,
	}
}

func (s *RedisStore) Init() error {
	if s.config == nil {
		return errors.New("redis store: configuration is missing")
	}

	s.store = redis.NewClient(&redis.Options{
		Addr:         s.config.address,
		Password:     s.config.password,
		DB:           s.config.db,
		MaxRetries:   s.config.maxRetries,
		ReadTimeout:  s.config.readTimeout,
		WriteTimeout: s.config.writeTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), s.config.readTimeout*time.Second)
	defer cancel()

	err := s.store.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("redis store: error pinging server: %v", err)
	}

	return nil
}

func (s *RedisStore) Has(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.readTimeout*time.Second)
	defer cancel()

	exists, err := s.store.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis store: error checking if key exists: %v", err)
	}

	return exists > 0, nil
}
func (s *RedisStore) Get(key string) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.readTimeout*time.Second)
	defer cancel()

	val, err := s.store.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("redis store: error getting cache data: %v", err)
	}

	return val, nil
}
func (s *RedisStore) Put(key string, data any, duration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.writeTimeout*time.Second)
	defer cancel()

	err := s.store.Set(ctx, key, data, duration).Err()
	if err != nil {
		return fmt.Errorf("redis store: error saving item to the store: %v", err)
	}

	return nil
}
func (s *RedisStore) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.writeTimeout*time.Second)
	defer cancel()

	err := s.store.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis store: error deleting key: %v", err)
	}

	return nil
}

func (s *RedisStore) Flush() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.readTimeout*time.Second)
	defer cancel()

	err := s.store.FlushDBAsync(ctx).Err()
	if err != nil {
		return fmt.Errorf("redis store: error flushing db: %v", err)
	}

	return nil
}

package redis

import (
	"context"
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
		store:  redis.NewClient(&redis.Options{}),
	}
}

func (s *RedisStore) Init() error {
	s.store.Options().Addr = s.config.address
	s.store.Options().Password = s.config.password
	s.store.Options().DB = s.config.db
	s.store.Options().MaxRetries = s.config.maxRetries
	s.store.Options().ReadTimeout = s.config.readTimeout
	s.store.Options().WriteTimeout = s.config.writeTimeout

	err := s.store.Ping(context.Background()).Err()
	if err != nil {
		return fmt.Errorf("redis store: error pinging server: %v", err)
	}

	return nil
}

func (s *RedisStore) Has(key string) (bool, error) {
	exists, err := s.store.Exists(context.Background(), key).Result()
	if err != nil {
		return false, fmt.Errorf("redis store: error checking if key exists: %v", err)
	}

	return exists > 0, nil
}
func (s *RedisStore) Get(key string) (any, error) {
	val, err := s.store.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("redis store: error getting cache data: %v", err)
	}

	return val, nil
}
func (s *RedisStore) Put(key string, data any, duration time.Duration) error {
	err := s.store.Set(context.Background(), key, data, duration).Err()
	if err != nil {
		return fmt.Errorf("redis store: error saving item to the store: %v", err)
	}

	return nil
}
func (s *RedisStore) Delete(key string) error {
	err := s.store.Del(context.Background(), key).Err()
	if err != nil {
		return fmt.Errorf("redis store: error deleting key: %v", err)
	}

	return nil
}

func (s *RedisStore) Flush() error {
	err := s.store.FlushDB(context.Background()).Err()
	if err != nil {
		return fmt.Errorf("redis store: error flushing db: %v", err)
	}

	return nil
}

func (s *RedisStore) FlushExpired() {
}

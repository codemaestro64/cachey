package memory

import (
	"time"

	"github.com/codemaestro64/cachey/store"
	"github.com/jellydator/ttlcache/v3"
)

type MemoryStore struct {
	store *ttlcache.Cache[string, any]
}

func NewMemoryStore() store.Store {
	return &MemoryStore{
		store: ttlcache.New[string, any](),
	}
}

func (s *MemoryStore) Init() error {
	return nil
}

func (s *MemoryStore) Has(key string) (bool, error) {
	return s.store.Has(key), nil
}

func (s *MemoryStore) Get(key string) (any, error) {
	item := s.store.Get(key)
	if item == nil || item.IsExpired() {
		return nil, nil
	}

	return item.Value(), nil
}

func (s *MemoryStore) Put(key string, data any, duration time.Duration) error {
	s.store.Set(key, data, duration)

	return nil
}
func (s *MemoryStore) Delete(key string) error {
	s.store.Delete(key)

	return nil
}

func (s *MemoryStore) Flush() error {
	s.store.DeleteAll()

	return nil
}

func (s *MemoryStore) FlushExpired() {
	s.store.DeleteExpired()
}

package store

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type MemoryStore struct {
	store *ttlcache.Cache[string, any]
}

func NewMemoryStore() (Store, error) {
	return &MemoryStore{
		store: ttlcache.New[string, any](),
	}, nil
}

func (s *MemoryStore) Has(key string) bool {
	return s.store.Has(key)
}
func (s *MemoryStore) Get(key string) any {
	item := s.store.Get(key)
	if item == nil || item.IsExpired() {
		return nil
	}

	return item.Value()
}
func (s *MemoryStore) Put(key string, data any, duration time.Duration) {
	s.store.Set(key, data, duration)
}
func (s *MemoryStore) Delete(key string) {
	s.store.Delete(key)
}

func (s *MemoryStore) Flush() {
	s.store.DeleteAll()
}

func (s *MemoryStore) FlushExpired() {
	s.store.DeleteExpired()
}

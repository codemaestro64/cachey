package store

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type MemoryStore struct {
	store *ttlcache.Cache[string, any]
}

func NewMemoryStore() Store {
	return &MemoryStore{
		store: ttlcache.New[string, any](),
	}
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
func (s *MemoryStore) Put(key string, data any, duration int) {
	var ttl time.Duration
	if duration > 0 {
		ttl = time.Second * time.Duration(duration)
	} else {
		ttl = ttlcache.NoTTL
	}

	s.store.Set(key, data, ttl)
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

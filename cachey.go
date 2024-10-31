package cachey

import (
	"time"

	"github.com/codemaestro64/cachey/store"
)

type Cache struct {
	store store.Store
}

func New(providerName string) (*Cache, error) {
	return nil, nil
}

func RegisterProvider(providerName string) error {
	return nil
}

func (c *Cache) Has(key string) bool {
	return false
}

func (c *Cache) Get(key string) any {
	return nil
}

func (c *Cache) GetOrDefault(key string, defaultFunc func() any) any {
	return nil
}

func (c *Cache) Remember(key string, duration time.Duration, rememberFunc func() any) any {
	return nil
}

func (c *Cache) RememberForever(key string, rememberFunc func() any) any {
	return nil
}

func (c *Cache) Pull(key string) any {
	return nil
}

func (c *Cache) PullOrDefault(key string, defaultFunc func() any) any {
	return nil
}

func (c *Cache) Put(key string, data any, duration time.Duration) {

}

func (c *Cache) Forever(key string, data any) {

}

func (c *Cache) Add(key string, data any, duration time.Duration) {

}

func (c *Cache) Forget(key string) {

}

func (c *Cache) Flush() {

}

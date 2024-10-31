package cachey

import (
	"fmt"
	"time"

	"github.com/codemaestro64/cachey/store"
)

// Cache represents a caching mechanism that wraps a store implementation.
type Cache struct {
	store store.Store // The underlying store for caching data.
}

// Supported cache provider constants.
const (
	MemoryStore = "memory" // Name of the memory store provider.

	ForeverDuration = -1 // Duration to store data indefinitely.
)

// stores maps provider names to their corresponding store constructors.
var stores = map[string]func() store.Store{
	MemoryStore: store.NewMemoryStore, // Registering the memory store.
}

// New initializes a new Cache instance using the specified provider name.
// It returns an error if the provider is not registered.
func New(providerName string) (*Cache, error) {
	storeConstructor, ok := stores[providerName]
	if !ok {
		return nil, fmt.Errorf("cache provider `%s` is not registered", providerName)
	}
	return &Cache{store: storeConstructor()}, nil
}

// RegisterProvider registers a new cache provider with the given name and constructor function.
// Returns an error if the provider is already registered.
func RegisterProvider(providerName string, providerFunc func() store.Store) error {
	if _, exists := stores[providerName]; exists {
		return fmt.Errorf("cache provider `%s` is already registered", providerName)
	}
	stores[providerName] = providerFunc
	return nil
}

// Has checks if a value exists in the cache for the given key.
// Returns true if the key exists, false otherwise.
func (c *Cache) Has(key string) bool {
	return c.store.Has(key)
}

// Get retrieves the value associated with the given key from the cache.
// Returns nil if the key does not exist.
func (c *Cache) Get(key string) any {
	return c.store.Get(key)
}

// GetOrDefault retrieves the value associated with the given key.
// If the key does not exist, it calls the provided defaultFunc to get a default value.
func (c *Cache) GetOrDefault(key string, defaultFunc func() any) any {
	if data := c.Get(key); data != nil {
		return data
	}
	return defaultFunc()
}

// Remember retrieves the value for the specified key from the cache.
// If it does not exist, it calls rememberFunc to generate the value,
// stores it in the cache with the specified duration, and returns it.
func (c *Cache) Remember(key string, duration time.Duration, rememberFunc func() any) any {
	if data := c.Get(key); data != nil {
		return data
	}
	data := rememberFunc()
	c.Put(key, data, duration)
	return data
}

// RememberForever retrieves the value for the specified key from the cache.
// If it does not exist, it calls rememberFunc to generate the value,
// and stores it indefinitely in the cache.
func (c *Cache) RememberForever(key string, rememberFunc func() any) any {
	return c.Remember(key, ForeverDuration, rememberFunc)
}

// Pull retrieves the value for the specified key from the cache and
// removes it from the cache. Returns the value or nil if it doesn't exist.
func (c *Cache) Pull(key string) any {
	data := c.Get(key)
	c.Forget(key)
	return data
}

// PullOrDefault retrieves the value for the specified key from the cache.
// If it does not exist, it calls defaultFunc to get a default value,
// removes the key from the cache, and returns the value.
func (c *Cache) PullOrDefault(key string, defaultFunc func() any) any {
	data := c.GetOrDefault(key, defaultFunc)
	c.Forget(key)
	return data
}

// Put stores the given data in the cache under the specified key
// with the provided duration. If the duration is zero, the data is
// stored indefinitely.
func (c *Cache) Put(key string, data any, duration time.Duration) {
	c.store.Put(key, data, duration)
}

// Forever stores the given data in the cache under the specified key
// indefinitely, ignoring the duration.
func (c *Cache) Forever(key string, data any) {
	c.Put(key, data, ForeverDuration)
}

// Add stores the given data in the cache under the specified key
// only if the key does not already exist. If the key exists, no action is taken.
func (c *Cache) Add(key string, data any, duration time.Duration) {
	if !c.Has(key) {
		c.store.Put(key, data, duration)
	}
}

// Forget removes the value associated with the specified key from the cache.
func (c *Cache) Forget(key string) {
	c.store.Delete(key)
}

// Flush removes all values from the cache.
func (c *Cache) Flush() {
	c.store.Flush()
}

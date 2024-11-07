package cachey

import (
	"fmt"
	"time"

	"github.com/codemaestro64/cachey/store"
	"github.com/codemaestro64/cachey/store/memory"
	"github.com/codemaestro64/cachey/store/redis"
)

// Cache represents a caching mechanism that wraps a store implementation.
type Cache struct {
	store store.Store // The underlying store for caching data.
}

// Supported cache store constants.
const (
	MemoryStore = "memory" // Name of the memory store.
	RedisStore  = "redis"  // Name of the redis store

	ForeverDuration = -1 // Duration to store data indefinitely.
)

type StoreConstructorFunc func() store.Store

// stores maps store names to their corresponding store constructors.
var stores = map[string]StoreConstructorFunc{
	MemoryStore: memory.NewMemoryStore,
	RedisStore:  redis.NewRedisStore,
}

// New initializes a new Cache instance using the specified store name.
// It returns an error if the store is not registered.
func New(storeName string, options ...store.Option) (*Cache, error) {
	storeConstructor, ok := stores[storeName]
	if !ok {
		return nil, fmt.Errorf("cache store `%s` is not registered", storeName)
	}

	store := storeConstructor()

	// apply options to the store
	for _, option := range options {
		err := option(store)
		if err != nil {
			return nil, err
		}
	}

	// initialize store with applied config
	err := store.Init()
	if err != nil {
		return nil, err
	}

	return &Cache{store: store}, nil
}

// Registerstore registers a new cache store with the given name and constructor function.
// Returns an error if the store is already registered.
func RegisterStore(storeName string, constructorFunc StoreConstructorFunc) error {
	if _, exists := stores[storeName]; exists {
		return fmt.Errorf("cache store `%s` is already registered", storeName)
	}
	stores[storeName] = constructorFunc
	return nil
}

// Has checks if a value exists in the cache for the given key.
// Returns true if the key exists, false otherwise.
func (c *Cache) Has(key string) (bool, error) {
	return c.store.Has(key)
}

// Get retrieves the value associated with the given key from the cache.
// Returns nil if the key does not exist.
func (c *Cache) Get(key string) (any, error) {
	return c.store.Get(key)
}

// GetOrDefault retrieves the value associated with the given key.
// If the key does not exist, it calls the provided defaultFunc to get a default value.
func (c *Cache) GetOrDefault(key string, defaultFunc func() any) (any, error) {
	data, err := c.Get(key)
	if err != nil {
		return nil, err
	}

	if data != nil {
		return data, nil
	}
	return defaultFunc(), nil
}

// Remember retrieves the value for the specified key from the cache.
// If it does not exist, it calls rememberFunc to generate the value,
// stores it in the cache with the specified duration, and returns it.
func (c *Cache) Remember(key string, duration time.Duration, rememberFunc func() any) (any, error) {
	data, err := c.Get(key)
	if err != nil {
		return nil, err
	}

	if data != nil {
		return data, nil
	}

	data = rememberFunc()
	c.Put(key, data, duration)
	return data, nil
}

// RememberForever retrieves the value for the specified key from the cache.
// If it does not exist, it calls rememberFunc to generate the value,
// and stores it indefinitely in the cache.
func (c *Cache) RememberForever(key string, rememberFunc func() any) (any, error) {
	return c.Remember(key, ForeverDuration, rememberFunc)
}

// Pull retrieves the value for the specified key from the cache and
// removes it from the cache. Returns the value or nil if it doesn't exist.
func (c *Cache) Pull(key string) (any, error) {
	data, err := c.Get(key)
	if err != nil {
		return nil, err
	}

	c.Forget(key)
	return data, nil
}

// PullOrDefault retrieves the value for the specified key from the cache.
// If it does not exist, it calls defaultFunc to get a default value,
// removes the key from the cache, and returns the value.
func (c *Cache) PullOrDefault(key string, defaultFunc func() any) (any, error) {
	data, err := c.GetOrDefault(key, defaultFunc)
	if err != nil {
		return nil, err
	}

	c.Forget(key)
	return data, nil
}

// Put stores the given data in the cache under the specified key
// with the provided duration. If the duration is zero, the data is
// stored indefinitely.
func (c *Cache) Put(key string, data any, duration time.Duration) error {
	return c.store.Put(key, data, duration)
}

// Forever stores the given data in the cache under the specified key
// indefinitely, ignoring the duration.
func (c *Cache) Forever(key string, data any) {
	c.Put(key, data, ForeverDuration)
}

// Add stores the given data in the cache under the specified key
// only if the key does not already exist. If the key exists, no action is taken.
func (c *Cache) Add(key string, data any, duration time.Duration) error {
	has, err := c.Has(key)
	if err != nil {
		return err
	}

	if !has {
		return c.store.Put(key, data, duration)
	}

	return nil
}

// Forget removes the value associated with the specified key from the cache.
func (c *Cache) Forget(key string) error {
	return c.store.Delete(key)
}

// Flush removes all values from the cache.
func (c *Cache) Flush() error {
	return c.store.Flush()
}

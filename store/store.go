package store

import (
	"time"
)

// Store defines the methods required for a caching store.
type Store interface {
	// Init initializes the store, readying it for use
	Init() error

	// Has checks if a value exists in the store for the given key.
	Has(key string) (bool, error)

	// Get retrieves the value associated with the given key.
	// Returns nil if the key does not exist.
	Get(key string) (any, error)

	// Put stores the value under the specified key with a duration.
	Put(key string, data any, duration time.Duration) error

	// Delete removes the value associated with the specified key.
	Delete(key string) error

	// Flush removes all values from the store.
	Flush() error
}

type Option func(store Store) error

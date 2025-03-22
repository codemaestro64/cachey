package redis

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

// TestRedisStore runs exhaustive unit tests for RedisStore
func TestRedisStore(t *testing.T) {
	// Start a mock Redis server
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	// Create a new RedisStore with the mock Redis address
	store := &RedisStore{
		config: &config{
			address:      mr.Addr(),
			password:     "",
			db:           0,
			maxRetries:   3,
			readTimeout:  5 * time.Second,
			writeTimeout: 5 * time.Second,
		},
	}

	// Initialize the store
	err = store.Init()
	assert.NoError(t, err, "Failed to initialize Redis store")

	// Test Put method
	t.Run("Put", func(t *testing.T) {
		err := store.Put("test_key", "test_value", 10*time.Second)
		assert.NoError(t, err, "Failed to set value in Redis")
	})

	// Test Get method
	t.Run("Get", func(t *testing.T) {
		val, err := store.Get("test_key")
		assert.NoError(t, err, "Failed to get value from Redis")
		assert.Equal(t, "test_value", val, "Stored value does not match expected value")
	})

	// Test Has method (should return true)
	t.Run("Has - Key Exists", func(t *testing.T) {
		exists, err := store.Has("test_key")
		assert.NoError(t, err, "Failed to check if key exists")
		assert.True(t, exists, "Key should exist but does not")
	})

	// Test Has method (should return false)
	t.Run("Has - Key Does Not Exist", func(t *testing.T) {
		exists, err := store.Has("non_existent_key")
		assert.NoError(t, err, "Failed to check if key exists")
		assert.False(t, exists, "Key should not exist but does")
	})

	// Test Delete method
	t.Run("Delete", func(t *testing.T) {
		err := store.Delete("test_key")
		assert.NoError(t, err, "Failed to delete key from Redis")

		exists, err := store.Has("test_key")
		assert.NoError(t, err, "Failed to check if deleted key exists")
		assert.False(t, exists, "Deleted key should not exist but does")
	})

	// Test Flush method
	t.Run("Flush", func(t *testing.T) {
		_ = store.Put("key1", "value1", 10*time.Second)
		_ = store.Put("key2", "value2", 10*time.Second)

		err := store.Flush()
		assert.NoError(t, err, "Failed to flush Redis DB")

		exists, _ := store.Has("key1")
		assert.False(t, exists, "Key1 should not exist after flush")

		exists, _ = store.Has("key2")
		assert.False(t, exists, "Key2 should not exist after flush")
	})

	// Test expiration
	err = store.Put("expiring_key", "value", 2*time.Second)
	assert.NoError(t, err, "Failed to set expiring key")

	// Simulate time passing in miniredis
	mr.FastForward(3 * time.Second)

	// Try accessing the key
	exists, _ := store.Has("expiring_key")
	assert.False(t, exists, "Expiring key should be removed by Redis after expiration")
}

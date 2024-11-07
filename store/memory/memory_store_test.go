package memory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMemoryStore(t *testing.T) {
	store := NewMemoryStore()
	assert.NotEqual(t, nil, store)
}

func TestMemoryStore_PutAndGet(t *testing.T) {
	store := NewMemoryStore()
	key := "testKey"
	value := "testValue"

	store.Put(key, value, time.Second)

	cachedValue, _ := store.Get(key)
	assert.Equal(t, value, cachedValue)

	// Wait for expiration
	time.Sleep(2 * time.Second)
	exists, _ := store.Has(key)
	assert.Equal(t, false, exists)
}

func TestMemoryStore_Has(t *testing.T) {
	store := NewMemoryStore()
	key := "testKey"
	value := "testValue"

	store.Put(key, value, time.Second)
	exists, _ := store.Has(key)
	assert.Equal(t, true, exists)

	// Wait for expiration
	time.Sleep(2 * time.Second)
	exists, _ = store.Has(key)
	assert.Equal(t, false, exists)
}

func TestMemoryStore_Delete(t *testing.T) {
	store := NewMemoryStore()
	key := "testKey"
	value := "testValue"
	store.Put(key, value, time.Minute)

	store.Delete(key)
	exists, _ := store.Has(key)
	assert.Equal(t, false, exists)
}

func TestMemoryStore_Flush(t *testing.T) {
	store := NewMemoryStore()
	store.Put("key1", "value1", time.Minute)
	store.Put("key2", "value2", time.Minute)

	exists1, _ := store.Has("key1")
	exists2, _ := store.Has("key2")

	assert.Equal(t, true, exists1)
	assert.Equal(t, true, exists2)

	store.Flush()

	has1, _ := store.Has("key1")
	has2, _ := store.Has("key2")
	assert.Equal(t, false, has1)
	assert.Equal(t, false, has2)
}

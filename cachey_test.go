package cachey

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testCacheGetOrDefault(t *testing.T, cache *Cache) {
	key := "key"
	defaultVal := "defaultVal"

	// 1. Try to get the value from cache, which should not exist initially.
	gotVal, err := cache.GetOrDefault(key, func() any {
		return defaultVal
	})
	assert.NoError(t, err)
	assert.Equal(t, defaultVal, gotVal) // Should return default value, since key doesn't exist

	// 2. Verify that the default value is not stored in the cache
	hasKey, err := cache.Has(key)
	assert.NoError(t, err)
	assert.False(t, hasKey) // Key should not exist in the cache

	// 3. Now set the value in the cache and test GetOrDefault again
	err = cache.Put(key, "storedValue", ForeverDuration)
	assert.NoError(t, err)

	// 4. Call GetOrDefault again, and it should return the cached value now
	gotValAgain, err := cache.GetOrDefault(key, func() any {
		return defaultVal
	})
	assert.NoError(t, err)
	assert.Equal(t, "storedValue", gotValAgain) // Should return the cached value

	// 5. Verify the cache still holds the stored value (and not the default)
	cachedVal, err := cache.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, "storedValue", cachedVal)

}

func testCachePull(t *testing.T, cache *Cache) {
	key := "key"
	val := "val"

	// add item to the cache
	err := cache.Put(key, val, ForeverDuration)
	assert.NoError(t, err)

	// pull the item
	cachedVal, err := cache.Pull(key)
	assert.NoError(t, err)

	// assert value is not nil
	assert.NotNil(t, cachedVal)

	// assert value is same as stored value
	assert.Equal(t, val, cachedVal)

	// ensure cache doesn't have pulled key
	has, err := cache.Has(key)
	assert.NoError(t, err)

	assert.Equal(t, false, has)
}

func testCachePullOrDefault(t *testing.T, cache *Cache) {
	key := "key"
	defaultVal := "defaultVal"

	// 1. First, call PullOrDefault when the key does not exist in the cache.
	pulledVal, err := cache.PullOrDefault(key, func() any {
		return defaultVal
	})
	assert.NoError(t, err)
	assert.Equal(t, defaultVal, pulledVal)

	// 2. Verify the value is not stored in the cache (key should not exist).
	hasKey, err := cache.Has(key)
	assert.NoError(t, err)
	assert.False(t, hasKey)

	// 3. Now add the key to the cache, and test PullOrDefault again
	err = cache.Put(key, "storedValue", ForeverDuration)
	assert.NoError(t, err)

	// 4. Call PullOrDefault again, it should return the cached value and remove the key from the cache
	pulledValAgain, err := cache.PullOrDefault(key, func() any {
		return "newDefault"
	})
	assert.NoError(t, err)
	assert.Equal(t, "storedValue", pulledValAgain)

	// 5. Verify the key has been removed from the cache
	hasKey, err = cache.Has(key)
	assert.NoError(t, err)
	assert.False(t, hasKey)
}

func testCacheRemember(t *testing.T, cache *Cache) {
	key := "key"
	rememberedValue := "value"

	val, err := cache.Remember(key, time.Second, func() any {
		return rememberedValue
	})
	assert.NoError(t, err)

	// ensure returned value and remembered values are equal
	assert.Equal(t, rememberedValue, val)

	// get item from cache
	cachedVal, err := cache.Get(key)
	assert.NoError(t, err)

	// ensure cachedVal equals remembderedValue
	assert.Equal(t, rememberedValue, cachedVal)

	// wait 2 seconds to ensure expiration
	time.Sleep(2 * time.Second)

	has, err := cache.Has(key)
	assert.NoError(t, err)

	// ensure key has been deleted
	assert.Equal(t, false, has)
}

func testCacheAdd(t *testing.T, cache *Cache) {
	key := "key"
	val := "val"

	err := cache.Add(key, val, ForeverDuration)
	assert.NoError(t, err)

	// ensure item was added to cache
	cachedVal, err := cache.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, val, cachedVal)

	// ensure it doesn't add an already present key
	// add item to cache
	key1 := "key1"
	val1 := "val1"

	err = cache.Put(key1, val1, ForeverDuration)
	assert.NoError(t, err)

	// add new val with same key
	val2 := "val2"
	err = cache.Add(key1, val2, ForeverDuration)
	assert.NoError(t, err)

	// get val
	cachedVal1, err := cache.Get(key1)
	assert.NoError(t, err)

	// ensure original value did not change
	assert.Equal(t, val1, cachedVal1)
}

func runAllTests(t *testing.T, cache *Cache) {
	t.Run("Test GetOrDefault", func(t *testing.T) {
		testCacheGetOrDefault(t, cache)
	})

	t.Run("Test Pull", func(t *testing.T) {
		testCachePull(t, cache)
	})

	t.Run("Test PullOrDefault", func(t *testing.T) {
		testCachePullOrDefault(t, cache)
	})

	t.Run("Test Remember", func(t *testing.T) {
		testCacheRemember(t, cache)
	})

	t.Run("Test Add", func(t *testing.T) {
		testCacheAdd(t, cache)
	})
}

func TestMemoryCache(t *testing.T) {
	memoryCache, _ := New(MemoryStore)
	runAllTests(t, memoryCache)
}

func TestRedisCache(t *testing.T) {

}

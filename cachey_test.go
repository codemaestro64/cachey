package cachey

import (
	"testing"
	"time"

	"github.com/codemaestro64/cachey/store"
)

func TestNewCache(t *testing.T) {
	cache, err := New("memory")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cache == nil {
		t.Fatal("expected a Cache instance, got nil")
	}
}

func TestCache_PutAndGet(t *testing.T) {
	memStore := store.NewMemoryStore()
	cache := &Cache{store: memStore}

	cache.Put("key1", "value1", 0)
	value := cache.Get("key1")
	if value != "value1" {
		t.Fatalf("expected value1, got %v", value)
	}
}

func TestCache_Has(t *testing.T) {
	memStore := store.NewMemoryStore()
	cache := &Cache{store: memStore}

	cache.Put("key1", "value1", 0)
	if !cache.Has("key1") {
		t.Fatal("expected key1 to exist")
	}
	if cache.Has("nonexistent") {
		t.Fatal("expected nonexistent to not exist")
	}
}

func TestCache_Forget(t *testing.T) {
	memStore := store.NewMemoryStore()
	cache := &Cache{store: memStore}

	cache.Put("key1", "value1", 0)
	cache.Forget("key1")
	if cache.Has("key1") {
		t.Fatal("expected key1 to be deleted")
	}
}

func TestCache_Flush(t *testing.T) {
	memStore := store.NewMemoryStore()
	cache := &Cache{store: memStore}

	cache.Put("key1", "value1", 0)
	cache.Flush()
	if cache.Has("key1") {
		t.Fatal("expected cache to be empty after flush")
	}
}

func TestCache_GetOrDefault(t *testing.T) {
	memStore := store.NewMemoryStore()
	cache := &Cache{store: memStore}
	defaultValue := "default"

	value := cache.GetOrDefault("nonexistent", func() any { return defaultValue })
	if value != defaultValue {
		t.Fatalf("expected %v, got %v", defaultValue, value)
	}
}

func TestCache_Remember(t *testing.T) {
	memStore := store.NewMemoryStore()
	cache := &Cache{store: memStore}

	value := cache.Remember("key1", time.Second, func() any { return "rememberedValue" })
	if value != "rememberedValue" {
		t.Fatalf("expected rememberedValue, got %v", value)
	}

	if !cache.Has("key1") {
		t.Fatal("expected key1 to exist in the cache")
	}

	// Wait for 2 seconds to confirm expiration
	time.Sleep(2 * time.Second)
	if cache.Has("key1") {
		t.Fatal("expected key1 to be expired after duration")
	}
}

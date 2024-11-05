package memory

import (
	"testing"
	"time"
)

func TestNewMemoryStore(t *testing.T) {
	store, _ := NewMemoryStore()
	if store == nil {
		t.Fatal("Expected new MemoryStore, got nil")
	}
}

func TestMemoryStore_PutAndGet(t *testing.T) {
	store, _ := NewMemoryStore()
	key := "testKey"
	value := "testValue"
	duration := 1 * time.Second

	store.Put(key, value, duration)

	retrievedValue, _ := store.Get(key)
	if retrievedValue != value {
		t.Errorf("Expected %v, got %v", value, retrievedValue)
	}

	// Wait for expiration
	time.Sleep(duration + 2*time.Second)
	retrievedValue, _ = store.Get(key)
	if retrievedValue != nil {
		t.Errorf("Expected nil after expiration, got %v", retrievedValue)
	}
}

func TestMemoryStore_Has(t *testing.T) {
	store, _ := NewMemoryStore()
	key := "testKey"
	value := "testValue"
	duration := 1 * time.Second

	store.Put(key, value, duration)
	has, _ := store.Has(key)
	if !has {
		t.Errorf("Expected true, got false")
	}

	// Wait for expiration
	time.Sleep(duration + 2*time.Second)
	has, _ = store.Has(key)
	if has {
		t.Errorf("Expected false after expiration, got true")
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	store, _ := NewMemoryStore()
	key := "testKey"
	value := "testValue"
	store.Put(key, value, time.Minute)

	store.Delete(key)
	has, _ := store.Has(key)
	if has {
		t.Errorf("Expected false after deletion, got true")
	}
}

func TestMemoryStore_Flush(t *testing.T) {
	store, _ := NewMemoryStore()
	store.Put("key1", "value1", time.Minute)
	store.Put("key2", "value2", time.Minute)

	store.Flush()
	has1, _ := store.Has("key1")
	has2, _ := store.Has("key2")
	if has1 || has2 {
		t.Errorf("Expected both keys to be deleted after flush")
	}
}

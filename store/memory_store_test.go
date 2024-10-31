package store

import (
	"testing"
	"time"
)

func TestNewMemoryStore(t *testing.T) {
	store := NewMemoryStore()
	if store == nil {
		t.Fatal("Expected new MemoryStore, got nil")
	}
}

func TestMemoryStore_PutAndGet(t *testing.T) {
	store := NewMemoryStore()
	key := "testKey"
	value := "testValue"
	duration := 1 * time.Second

	store.Put(key, value, duration)

	retrievedValue := store.Get(key)
	if retrievedValue != value {
		t.Errorf("Expected %v, got %v", value, retrievedValue)
	}

	// Wait for expiration
	time.Sleep(duration + 2*time.Second)
	retrievedValue = store.Get(key)
	if retrievedValue != nil {
		t.Errorf("Expected nil after expiration, got %v", retrievedValue)
	}
}

func TestMemoryStore_Has(t *testing.T) {
	store := NewMemoryStore()
	key := "testKey"
	value := "testValue"
	duration := 1 * time.Second

	store.Put(key, value, duration)
	if !store.Has(key) {
		t.Errorf("Expected true, got false")
	}

	// Wait for expiration
	time.Sleep(duration + 2*time.Second)
	if store.Has(key) {
		t.Errorf("Expected false after expiration, got true")
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	store := NewMemoryStore()
	key := "testKey"
	value := "testValue"
	store.Put(key, value, time.Minute)

	store.Delete(key)
	if store.Has(key) {
		t.Errorf("Expected false after deletion, got true")
	}
}

func TestMemoryStore_Flush(t *testing.T) {
	store := NewMemoryStore()
	store.Put("key1", "value1", time.Minute)
	store.Put("key2", "value2", time.Minute)

	store.Flush()
	if store.Has("key1") || store.Has("key2") {
		t.Errorf("Expected both keys to be deleted after flush")
	}
}

package internal

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cache := NewCache(5 * time.Second)

	if cache == nil {
		t.Error("expected cache to be created, got nil")
	}

	if cache.entries == nil {
		t.Error("expected entries map to be initialized, got nil")
	}
}

func TestAddAndGet(t *testing.T) {
	cache := NewCache(5 * time.Second)

	key := "test-key"
	val := []byte("test-value")

	cache.Add(key, val)
	result, found := cache.Get(key)

	if !found {
		t.Error("expected to find key in cache")
	}

	if string(result) != string(val) {
		t.Errorf("expected %s, got %s", string(val), string(result))
	}
}

func TestGetNonExistentKey(t *testing.T) {
	cache := NewCache(5 * time.Second)

	_, found := cache.Get("does-not-exist")

	if found {
		t.Error("expected found to be false for non-existent key")
	}
}

func TestAddOverwritesExistingKey(t *testing.T) {
	cache := NewCache(5 * time.Second)

	key := "test-key"
	val1 := []byte("first-value")
	val2 := []byte("second-value")

	cache.Add(key, val1)
	cache.Add(key, val2)

	result, found := cache.Get(key)

	if !found {
		t.Error("expected to find key in cache")
	}

	if string(result) != string(val2) {
		t.Errorf("expected %s, got %s", string(val2), string(result))
	}
}

func TestReapLoopRemovesOldEntries(t *testing.T) {
	interval := 50 * time.Millisecond
	cache := NewCache(interval)

	cache.Add("old-key", []byte("old-value"))

	// Verify entry exists
	_, found := cache.Get("old-key")
	if !found {
		t.Error("expected key to exist immediately after adding")
	}

	// Wait for reapLoop to run (interval + buffer)
	time.Sleep(interval*2 + 10*time.Millisecond)

	// Entry should be removed
	_, found = cache.Get("old-key")
	if found {
		t.Error("expected old entry to be removed by reapLoop")
	}
}

func TestReapLoopKeepsNewEntries(t *testing.T) {
	interval := 100 * time.Millisecond
	cache := NewCache(interval)

	cache.Add("new-key", []byte("new-value"))

	// Wait less than interval
	time.Sleep(interval / 2)

	// Entry should still exist
	_, found := cache.Get("new-key")
	if !found {
		t.Error("expected new entry to still exist before interval passes")
	}
}

func TestMultipleEntries(t *testing.T) {
	cache := NewCache(5 * time.Second)

	entries := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
		"key3": []byte("value3"),
	}

	for k, v := range entries {
		cache.Add(k, v)
	}

	for k, expected := range entries {
		result, found := cache.Get(k)
		if !found {
			t.Errorf("expected to find key %s", k)
		}
		if string(result) != string(expected) {
			t.Errorf("for key %s: expected %s, got %s", k, string(expected), string(result))
		}
	}
}

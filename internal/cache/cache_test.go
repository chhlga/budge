package cache

import (
	"sync"
	"testing"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := New(3)

	cache.Set("key1", "value1")

	val, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1")
	}
	if val != "value1" {
		t.Errorf("Expected 'value1', got '%v'", val)
	}
}

func TestCache_GetMissing(t *testing.T) {
	cache := New(3)

	_, found := cache.Get("nonexistent")
	if found {
		t.Error("Expected not to find nonexistent key")
	}
}

func TestCache_LRUEviction(t *testing.T) {
	cache := New(3)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	if _, found := cache.Get("key1"); !found {
		t.Error("Expected to find key1")
	}
	if _, found := cache.Get("key2"); !found {
		t.Error("Expected to find key2")
	}
	if _, found := cache.Get("key3"); !found {
		t.Error("Expected to find key3")
	}

	cache.Set("key4", "value4")

	if _, found := cache.Get("key1"); found {
		t.Error("Expected key1 to be evicted (least recently used)")
	}
	if _, found := cache.Get("key2"); !found {
		t.Error("Expected to find key2")
	}
	if _, found := cache.Get("key3"); !found {
		t.Error("Expected to find key3")
	}
	if _, found := cache.Get("key4"); !found {
		t.Error("Expected to find key4")
	}
}

func TestCache_LRUUpdateOnGet(t *testing.T) {
	cache := New(3)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	cache.Get("key1")

	cache.Set("key4", "value4")

	if _, found := cache.Get("key1"); !found {
		t.Error("Expected key1 to still exist (was accessed recently)")
	}
	if _, found := cache.Get("key2"); found {
		t.Error("Expected key2 to be evicted (least recently used)")
	}
}

func TestCache_UpdateValue(t *testing.T) {
	cache := New(3)

	cache.Set("key1", "value1")
	cache.Set("key1", "value2")

	val, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1")
	}
	if val != "value2" {
		t.Errorf("Expected 'value2', got '%v'", val)
	}
}

func TestCache_Clear(t *testing.T) {
	cache := New(3)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	cache.Clear()

	if _, found := cache.Get("key1"); found {
		t.Error("Expected cache to be empty after Clear()")
	}
	if _, found := cache.Get("key2"); found {
		t.Error("Expected cache to be empty after Clear()")
	}
}

func TestCache_Delete(t *testing.T) {
	cache := New(3)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	cache.Delete("key1")

	if _, found := cache.Get("key1"); found {
		t.Error("Expected key1 to be deleted")
	}
	if _, found := cache.Get("key2"); !found {
		t.Error("Expected key2 to still exist")
	}
}

func TestCache_ThreadSafety(t *testing.T) {
	cache := New(100)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := string(rune('a' + id))
				cache.Set(key, id)
				cache.Get(key)
			}
		}(i)
	}

	wg.Wait()
}

func TestCache_Size(t *testing.T) {
	cache := New(5)

	if cache.Len() != 0 {
		t.Errorf("Expected empty cache to have size 0, got %d", cache.Len())
	}

	cache.Set("key1", "value1")
	if cache.Len() != 1 {
		t.Errorf("Expected size 1, got %d", cache.Len())
	}

	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	if cache.Len() != 3 {
		t.Errorf("Expected size 3, got %d", cache.Len())
	}

	cache.Delete("key1")
	if cache.Len() != 2 {
		t.Errorf("Expected size 2 after delete, got %d", cache.Len())
	}

	cache.Clear()
	if cache.Len() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", cache.Len())
	}
}

func TestCache_ZeroCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for zero capacity cache")
		}
	}()

	New(0)
}

func TestCache_NegativeCapacity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for negative capacity cache")
		}
	}()

	New(-1)
}

// Package ratelimit provides a map-based implementation of a TTL (Time to Live) rate limiter.
// It allows you to track usage count for each key and automatically evict stale entries from the map.
//
// Usage:
//
//	m := ratelimit.New(60) // Creates a new TTLMap with max TTL of 60 seconds
//	m.IncrementUsage("key") // Increments usage count for "key"
//	count := m.GetUsageCount("key") // Gets the current usage count for "key"
//
// The TTLMap instance created with New() will automatically evict entries that have not been accessed for more than
// the specified TTL. This eviction process is done asynchronously by a goroutine.
package ratelimit

import (
	"sync"
	"time"
)

// item is a struct that represents a value in the map, along with its creation timestamp
type item struct {
	value     int   // The usage count for this item
	createdTs int64 // The UNIX timestamp when this item was created
}

// TTLMap is a thread-safe map-based implementation of a TTL (Time to Live) rate limiter
type TTLMap struct {
	m map[string]*item // The underlying map that holds the key-value pairs
	l sync.Mutex       // The mutex used to synchronize access to the map
}

// New returns a new TTLMap instance with a maximum TTL of maxTTL seconds.
// The returned instance automatically evicts stale entries every second.
func New(maxTTL int) (m *TTLMap) {
	m = &TTLMap{m: make(map[string]*item, 0)}

	go func() {
		for now := range time.Tick(time.Second) {
			m.l.Lock()
			for k, v := range m.m {
				if now.Unix()-v.createdTs > int64(maxTTL) {
					delete(m.m, k)
				}
			}
			m.l.Unlock()
		}
	}()

	return
}

// IncrementUsage increments the usage count for the specified key.
// If the key does not exist in the map, it is added with a usage count of 1.
func (m *TTLMap) IncrementUsage(k string) {
	m.l.Lock()

	_, ok := m.m[k]
	if !ok {
		m.m[k] = &item{
			value:     1,
			createdTs: time.Now().Unix(),
		}
	} else {
		m.m[k].value++
	}

	m.l.Unlock()
}

// GetUsageCount returns the current usage count for the specified key.
// If the key does not exist in the map, it returns 0.
func (m *TTLMap) GetUsageCount(k string) (v int) {
	m.l.Lock()

	if it, ok := m.m[k]; ok {
		v = it.value
	}

	m.l.Unlock()

	return v
}

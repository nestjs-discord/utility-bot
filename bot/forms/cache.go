package forms

import (
	"fmt"
	"github.com/dgraph-io/ristretto"
)

func (f *Forms) initCacheInstance() error {
	cache, err := ristretto.NewCache(&ristretto.Config[string, bool]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return fmt.Errorf("failed to create cache instance for mod actions: %w", err)
	}
	f.modActionsCache = cache

	return nil
}

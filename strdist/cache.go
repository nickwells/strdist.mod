package strdist

import "fmt"

const DfltMaxCacheSize = 3

// cacheEntry holds the details of a cache entry.
type cacheEntry[T any] struct {
	e        T
	useCount int
}

// cache represents a cache of values of type T referenced by a name
type cache[T any] struct {
	cache        map[string]cacheEntry[T]
	maxCacheSize int
}

// newCache returns a newly instantiated cache and an error which will be
// non-nil if the maxCacheSize is < 0. Note that a cache with a zero
// maxCacheSize will always be empty.
func newCache[T any](maxCacheSize int) (*cache[T], error) {
	if maxCacheSize < 0 {
		return nil,
			fmt.Errorf("the maxCacheSize (%d) must be >= 0", maxCacheSize)
	}

	return &cache[T]{
		cache:        map[string]cacheEntry[T]{},
		maxCacheSize: maxCacheSize,
	}, nil
}

// Desc returns a string describing the cache configuration
func (c cache[T]) Desc() string {
	return fmt.Sprintf("cache sz: %3d", c.maxCacheSize)
}

// clearLeastUsedCacheEntry removes the cache entry which has been used least
// frequently
func (c *cache[T]) clearLeastUsedCacheEntry() {
	if c.maxCacheSize == 0 { // there is no cache
		return
	}

	if len(c.cache) == 0 { // the cache has no entries
		return
	}

	if len(c.cache) < c.maxCacheSize { // the cache is not full
		return
	}

	var leastUsedEntry string

	leastUsedCount := -1

	for k, ce := range c.cache {
		if leastUsedCount != -1 {
			if ce.useCount >= leastUsedCount {
				continue
			}
		}

		leastUsedCount = ce.useCount
		leastUsedEntry = k
	}

	delete(c.cache, leastUsedEntry)
}

// setCachedEntry sets the value of the cached entry in the cache. It first
// removes the least heavily used entry (if necessary).
func (c *cache[T]) setCachedEntry(key string, val T) {
	c.clearLeastUsedCacheEntry()

	c.cache[key] = cacheEntry[T]{
		e:        val,
		useCount: 1,
	}
}

// getCachedEntry gets the value of the cached entry from the cache,
// returning the entry value and a bool indicating whether the entry was
// found in the map
func (c *cache[T]) getCachedEntry(key string) (T, bool) {
	ce, ok := c.cache[key]
	if ok {
		ce.useCount++
		c.cache[key] = ce
	}

	return ce.e, ok
}

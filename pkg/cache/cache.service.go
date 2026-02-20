package cache

import (
	"runtime"
	"sync"
	"time"
)

type ItemCache struct {
	Object     any
	Expiration int64
}

func (ic ItemCache) Expired(now int64) bool {
	return ic.Expiration > 0 && now > ic.Expiration
}

type Cache struct {
	*cache
}

type cache struct {
	mutex             sync.RWMutex
	items             map[string]ItemCache
	janitor           *janitor
	limit             int
	defaultExpiration time.Duration
}

func New(defaultExpiration time.Duration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]ItemCache)
	return newCacheWithJanitor(defaultExpiration, cleanupInterval, items)
}

func newCacheWithJanitor(defaultExpiration time.Duration, cleanupInterval time.Duration, items map[string]ItemCache) *Cache {
	c := newCache(defaultExpiration, items)
	C := &Cache{c}
	if cleanupInterval > 0 {
		runJanitor(c, cleanupInterval)
		runtime.SetFinalizer(C, stopJanitor)
	}
	return C
}

func newCache(defaultExpiration time.Duration, items map[string]ItemCache) *cache {
	if defaultExpiration == 0 {
		defaultExpiration = -1
	}
	c := &cache{
		defaultExpiration: defaultExpiration,
		items:             items,
	}
	return c
}

func (c *cache) Get(k string) (any, bool) {
	c.mutex.RLock()
	item, found := c.items[k]
	if !found {
		c.mutex.RUnlock()
		return nil, false
	}
	if item.Expired(time.Now().UnixNano()) {
		c.mutex.RUnlock()
		return nil, false
	}
	c.mutex.RUnlock()
	return item.Object, true
}

func (c *cache) GetWithExpiration(k string) (any, time.Time, bool) {
	var expirationTime time.Time
	c.mutex.RLock()
	item, found := c.items[k]
	if !found {
		c.mutex.RUnlock()
		return nil, expirationTime, false
	}

	if item.Expiration > 0 {
		expirationTime = time.Unix(0, item.Expiration)
	}

	if item.Expired(time.Now().UnixNano()) {
		c.mutex.RUnlock()
		return nil, expirationTime, false
	}

	c.mutex.RUnlock()
	return item.Object, expirationTime, true
}

func (c *cache) get(k string) (any, bool) {
	item, found := c.items[k]
	if !found {
		return nil, false
	}
	if item.Expired(time.Now().UnixNano()) {
		return nil, false
	}
	return item.Object, true
}

func (c *cache) Set(key string, value any) {
	c.set(key, value, 0)
}

func (c *cache) SetWithDuration(key string, value any, duration time.Duration) {
	var e int64
	if duration == 0 {
		duration = c.defaultExpiration
	}
	if duration > 0 {
		e = time.Now().Add(duration).UnixNano()
	}
	c.mutex.Lock()
	c.items[key] = ItemCache{Object: value, Expiration: e}
	c.mutex.Unlock()
}

func (c *cache) set(key string, value any, duration time.Duration) {
	var e int64
	if duration == 0 {
		duration = c.defaultExpiration
	}
	if duration > 0 {
		e = time.Now().Add(duration).UnixNano()
	}
	c.mutex.Lock()
	c.items[key] = ItemCache{Object: value, Expiration: e}
	c.mutex.Unlock()
}

func (c *cache) DeleteExpired() {
	var expiredKeys []string
	now := time.Now().UnixNano()
	c.mutex.Lock()
	for key, item := range c.items {
		if item.Expired(now) {
			expiredKeys = append(expiredKeys, key)
		}
	}
	c.mutex.Unlock()
	for _, key := range expiredKeys {
		c.Delete(key)
	}
}

func (c *cache) Delete(key string) {
	c.mutex.Lock()
	delete(c.items, key)
	c.mutex.Unlock()
}

func (c *cache) Items() map[string]ItemCache {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	m := make(map[string]ItemCache, len(c.items))
	now := time.Now().UnixNano()
	for k, v := range c.items {
		if v.Expired(now) {
			continue
		}
		m[k] = v
	}
	return m
}

func (c *cache) ItemCount() int {
	c.mutex.RLock()
	n := len(c.items)
	c.mutex.RUnlock()
	return n
}

func (c *cache) Flush() {
	c.mutex.Lock()
	c.items = map[string]ItemCache{}
	c.mutex.Unlock()
}

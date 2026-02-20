package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetGet(t *testing.T) {
	durationZero := time.Duration(0)
	c := New(durationZero, durationZero)

	wg := sync.WaitGroup{}
	for i := range [1000]int{} {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", i)
			c.Set(key, i)
		}(i)
	}

	wg.Wait()

	for i := range [1000]int{} {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", i)
			val, ok := c.Get(key)
			if !ok {
				t.Errorf("key %s must be will exist", key)
			}

			valInt, ok := val.(int)
			if !ok {
				t.Errorf("key %s value %v must be will int", key, val)
			}

			if valInt != i {
				t.Errorf("key %s value %d must be will equal to %d", key, valInt, i)
			}
		}(i)
	}

	wg.Wait()
}
func TestItemCache_Expired(t *testing.T) {
	tableTest := []struct {
		Name       string
		Expiration int64
		Now        int64
		Expected   bool
	}{
		{
			Name:       "should false when expiration is less or equal than zero",
			Expiration: 0,
			Expected:   false,
		},
		{
			Name:       "should false when now is less than expiration",
			Expiration: 2,
			Now:        1,
			Expected:   false,
		},
		{
			Name:       "should true when now is greater than expiration",
			Expiration: 1,
			Now:        2,
			Expected:   true,
		},
	}

	for _, test := range tableTest {
		t.Run(test.Name, func(t *testing.T) {
			item := &ItemCache{
				Expiration: test.Expiration,
			}

			assert.Equal(t, test.Expected, item.Expired(test.Now))
		})
	}
}

func Test_New(t *testing.T) {
	tests := []struct {
		name                      string
		defaultExpiration         time.Duration
		cleanupInterval           time.Duration
		defaultExpirationExpected time.Duration
		items                     map[string]ItemCache
	}{
		{
			name:                      "Test with zero duration",
			defaultExpiration:         time.Duration(0),
			cleanupInterval:           time.Duration(0),
			defaultExpirationExpected: time.Duration(-1),
			items:                     map[string]ItemCache{},
		},
		{
			name:                      "Test with non-zero duration",
			defaultExpiration:         time.Duration(1),
			cleanupInterval:           time.Duration(1),
			defaultExpirationExpected: time.Duration(1),
			items:                     map[string]ItemCache{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheActual := New(tt.defaultExpiration, tt.cleanupInterval)
			assert.NotNil(t, cacheActual)

			cacheExpected := &cache{
				defaultExpiration: tt.defaultExpirationExpected,
				items:             tt.items,
			}

			assert.Equal(t, cacheExpected.items, cacheActual.items)
			assert.NotNil(t, &cacheExpected.mutex)
			assert.Zero(t, cacheActual.limit)
			assert.Equal(t, cacheExpected.defaultExpiration, cacheActual.defaultExpiration)
		})
	}
}

func Test_newCache(t *testing.T) {
	tests := []struct {
		name                      string
		defaultExpiration         time.Duration
		defaultExpirationExpected time.Duration
		items                     map[string]ItemCache
	}{
		{
			name:                      "Test with zero duration",
			defaultExpiration:         time.Duration(0),
			defaultExpirationExpected: time.Duration(-1),
			items:                     map[string]ItemCache{},
		},
		{
			name:                      "Test with non-zero duration",
			defaultExpiration:         time.Duration(1),
			defaultExpirationExpected: time.Duration(1),
			items:                     map[string]ItemCache{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheActual := newCache(tt.defaultExpiration, tt.items)
			assert.NotNil(t, cacheActual)

			cacheExpected := &cache{
				defaultExpiration: tt.defaultExpirationExpected,
				items:             tt.items,
			}

			assert.Equal(t, cacheExpected.items, cacheActual.items)
			assert.NotNil(t, &cacheExpected.mutex)
			assert.Zero(t, cacheActual.limit)
			assert.Equal(t, cacheExpected.defaultExpiration, cacheActual.defaultExpiration)
		})
	}
}

func TestCache_Get(t *testing.T) {
	testTable := []struct {
		Name          string
		isSetEnable   bool
		setValue      any
		expectedValue any
		expectedOk    bool
	}{
		{
			Name:          "key is not exist",
			isSetEnable:   false,
			setValue:      "",
			expectedValue: nil,
			expectedOk:    false,
		},
		{
			Name:          "key is exist",
			isSetEnable:   true,
			setValue:      "123456",
			expectedValue: "123456",
			expectedOk:    true,
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			c := New(time.Duration(0), time.Duration(0))
			if test.isSetEnable {
				c.Set("key", test.setValue)
			}

			value, ok := c.Get("key")
			assert.Equal(t, test.expectedOk, ok)
			assert.Equal(t, test.expectedValue, value)
		})
	}
}

func TestCache_get(t *testing.T) {
	testTable := []struct {
		Name          string
		isSetEnable   bool
		setValue      any
		expectedValue any
		expectedOk    bool
	}{
		{
			Name:          "key is not exist",
			isSetEnable:   false,
			setValue:      "",
			expectedValue: nil,
			expectedOk:    false,
		},
		{
			Name:          "key is exist",
			isSetEnable:   true,
			setValue:      "123456",
			expectedValue: "123456",
			expectedOk:    true,
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			c := New(time.Duration(0), time.Duration(0))
			if test.isSetEnable {
				c.Set("key", test.setValue)
			}

			value, ok := c.get("key")
			assert.Equal(t, test.expectedOk, ok)
			assert.Equal(t, test.expectedValue, value)
		})
	}
}

func TestCache_GetWithExpiration(t *testing.T) {
	testTable := []struct {
		Name          string
		isSetEnable   bool
		setValue      any
		expectedValue any
		exepectedTime time.Time
		expectedOk    bool
	}{
		{
			Name:          "key is not exist",
			isSetEnable:   false,
			setValue:      "",
			expectedValue: nil,
			exepectedTime: time.Time{},
			expectedOk:    false,
		},
		{
			Name:          "key is exist",
			isSetEnable:   true,
			setValue:      "123456",
			expectedValue: "123456",
			exepectedTime: time.Time{},
			expectedOk:    true,
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			c := New(time.Duration(0), time.Duration(0))
			if test.isSetEnable {
				c.Set("key", test.setValue)
			}

			value, expirationTime, ok := c.GetWithExpiration("key")
			assert.Equal(t, test.expectedOk, ok)
			assert.Equal(t, test.exepectedTime, expirationTime)
			assert.Equal(t, test.expectedValue, value)
		})
	}
}

func TestCache_SetWithDuration(t *testing.T) {
	testTable := []struct {
		name              string
		value             any
		duration          time.Duration
		defaultExpiration time.Duration
		itemCacheExpected ItemCache
		inDeltaExpiration time.Duration
	}{
		{
			name:              "Test with zero duration",
			value:             "test",
			duration:          time.Duration(0) * time.Second,
			defaultExpiration: time.Duration(1) * time.Second,
			itemCacheExpected: ItemCache{
				Object:     "test",
				Expiration: time.Now().UnixNano(),
			},
			inDeltaExpiration: time.Duration(2) * time.Second,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := New(test.defaultExpiration, time.Duration(0))
			c.SetWithDuration("key", test.value, test.duration)
			itemCacheActual, ok := c.items["key"]
			assert.True(t, ok)
			assert.Equal(t, test.itemCacheExpected.Object, itemCacheActual.Object)
			assert.InDelta(t, test.itemCacheExpected.Expiration, itemCacheActual.Expiration, float64(test.inDeltaExpiration))
		})
	}
}

func TestCache_DeleteExpired(t *testing.T) {}

func TestCache_Delete(t *testing.T) {
	c := New(time.Duration(0), time.Duration(0))
	c.set("key", "value", time.Duration(0))
	c.Delete("key")
	_, ok := c.items["key"]
	assert.False(t, ok)
}

func TestCache_GetItems(t *testing.T) {
	c := New(time.Duration(0), time.Duration(0))
	c.set("key", "value", time.Duration(0))
	expectedMap := map[string]ItemCache{
		"key": {
			Object:     "value",
			Expiration: 0,
		},
	}

	currentMap := c.Items()
	assert.Equal(t, expectedMap, currentMap)
}

func TestCache_ItemCount(t *testing.T) {
	c := New(time.Duration(0), time.Duration(0))
	c.set("key", "value", time.Duration(0))
	assert.Equal(t, 1, c.ItemCount())
}

func TestCache_Items(t *testing.T) {
	c := New(time.Duration(0), time.Duration(0))
	c.set("key", "value", time.Duration(0))
	assert.NotEmpty(t, c.items)
	c.Flush()
	assert.Empty(t, c.items)
}

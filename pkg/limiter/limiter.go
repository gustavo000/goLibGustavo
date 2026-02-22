package limiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.RWMutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func New(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Clean old requests
	var valid []time.Time
	for _, t := range rl.requests[clientIP] {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.limit {
		rl.requests[clientIP] = valid
		return false
	}

	rl.requests[clientIP] = append(valid, now)
	return true
}

func (rl *RateLimiter) GetCount(clientIP string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return len(rl.requests[clientIP])
}

func (rl *RateLimiter) Reset(clientIP string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.requests, clientIP)
}

package auth_service

import (
	"sync"
	"time"
)

type rateLimiter struct {
	mu          sync.Mutex
	limit       int
	window      time.Duration
	bucket      map[string]*rateBucket
	lastCleanup time.Time
}

type rateBucket struct {
	count int
	reset time.Time
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		limit:  limit,
		window: window,
		bucket: make(map[string]*rateBucket),
	}
}

func NewRateLimiter(limit int, window time.Duration) *rateLimiter {
	return newRateLimiter(limit, window)
}

func (l *rateLimiter) Allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.lastCleanup.IsZero() {
		l.lastCleanup = now
	}
	if now.Sub(l.lastCleanup) >= l.window {
		for k, b := range l.bucket {
			if now.After(b.reset) {
				delete(l.bucket, k)
			}
		}
		l.lastCleanup = now
	}

	b, ok := l.bucket[key]
	if !ok || now.After(b.reset) {
		l.bucket[key] = &rateBucket{count: 1, reset: now.Add(l.window)}
		return true
	}
	if b.count >= l.limit {
		return false
	}
	b.count++
	return true
}

package initialize

import (
	"sync"

	"golang.org/x/time/rate"
)

type AutomateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
}

var alimit *AutomateLimiter

func NewAutomateLimiter() *AutomateLimiter {
	if alimit == nil {
		alimit = &AutomateLimiter{
			limiters: make(map[string]*rate.Limiter),
		}
	}
	return alimit
}

func (rl *AutomateLimiter) GetLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, ok := rl.limiters[key]
	if !ok {
		limiter = rate.NewLimiter(rate.Limit(1.0/3.0), 10) // 1 request per 3 seconds = 20 requests per minute, 10 requests được thực thi ngay lập tức sau đó sẽ bị giới hạn theo rate 1/3 requests/second
		rl.limiters[key] = limiter
	}
	return limiter
}

func (rl *AutomateLimiter) Allow(key string) bool {
	limiter := rl.GetLimiter(key)
	return limiter.Allow()
}

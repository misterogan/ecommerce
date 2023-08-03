package ratelimiter

import (
	"net/http"
	"sync"
	"time"
)

type Limiter struct {
	limit    int
	duration time.Duration
	ips      map[string]int
	mutex    sync.Mutex
}

func NewLimiter(limit int, duration time.Duration) *Limiter {
	return &Limiter{
		limit:    limit,
		duration: duration,
		ips:      make(map[string]int),
	}
}

func (lim *Limiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		lim.mutex.Lock()
		defer lim.mutex.Unlock()

		count, exists := lim.ips[ip]
		if !exists {
			lim.ips[ip] = 1
		} else {
			count++
			if count > lim.limit {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			lim.ips[ip] = count
		}

		go func() {
			<-time.After(lim.duration)
			lim.mutex.Lock()
			defer lim.mutex.Unlock()
			delete(lim.ips, ip)
		}()

		next(w, r)
	}
}

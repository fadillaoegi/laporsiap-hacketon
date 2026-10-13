package httpapi

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type rateBucket struct {
	started time.Time
	count   int
}
type localLimiter struct {
	mu      sync.Mutex
	buckets map[string]rateBucket
	limit   int
}

func (l *localLimiter) allow(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	for key, b := range l.buckets {
		if now.Sub(b.started) >= time.Minute {
			delete(l.buckets, key)
		}
	}
	b := l.buckets[host]
	if b.started.IsZero() {
		b.started = now
	}
	if b.count >= l.limit || (len(l.buckets) >= 1024 && b.count == 0) {
		return false
	}
	b.count++
	l.buckets[host] = b
	return true
}
func tooMany(w http.ResponseWriter, err error) {
	w.Header().Set("Retry-After", "60")
	ErrorResponse(w, http.StatusTooManyRequests, err)
}

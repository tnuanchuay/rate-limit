package rate_limit

import (
	"net/http"
	"sync"
	"time"
)

type RateLimit struct {
	lock sync.Mutex
	tick <-chan time.Time
}

func New(request_per_sec float64) RateLimit {
	p := 1 / request_per_sec
	p = p / 1e-9
	intP := int64(p)

	return RateLimit{
		tick: time.Tick(time.Duration(intP)),
	}
}

func (rl RateLimit) PromiseTick() {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	<-rl.tick
}

func (rl RateLimit) Available() bool {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	if len(rl.tick) > 0 {
		<-rl.tick
		return true
	}

	return false
}

func (rl RateLimit) PromiseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.PromiseTick()
		next.ServeHTTP(w, r)
	})
}

func (rl RateLimit) TimeoutMiddleware(next http.Handler, onTimeout http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.Available() {
			onTimeout.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

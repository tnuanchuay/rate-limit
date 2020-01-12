package rate_limit

import (
	"net/http"
	"time"
)

type RateLimit struct {
	tick   <- chan time.Time
}

func New(request_per_sec float64) RateLimit {
	p := 1 / request_per_sec
	p = p / 1e-9
	intP := int64(p)

	return RateLimit{
		tick:        time.Tick(time.Duration(intP)),
	}
}

func (rl RateLimit) Tick() {
	<-rl.tick
}

func (rl RateLimit) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.Tick()
		next.ServeHTTP(w, r)
	})
}

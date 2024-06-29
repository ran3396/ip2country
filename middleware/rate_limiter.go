package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// This function returns a middleware function that limits the number of requests per second.
// It uses a mutex to synchronize access to the request count and the last request time.
// If the time since the last request is greater than 1 second, the request count is reset to 0.
// If the request count exceeds the limit, a 429 Too Many Requests response is sent.
func RateLimiter(rps int) mux.MiddlewareFunc {
	var mu sync.Mutex
	var lastRequestTime time.Time
	var requestCount int

	return func(next http.Handler) http.Handler {
		// This is the actual middleware function that enforces the rate limit.
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			// Defer the unlock to ensure it's called even if the handler panics.
			defer mu.Unlock()

			now := time.Now()
			if now.Sub(lastRequestTime).Seconds() > 1 {
				lastRequestTime = now
				requestCount = 0
			}

			if requestCount >= rps {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			requestCount++
			next.ServeHTTP(w, r)
		})
	}
}

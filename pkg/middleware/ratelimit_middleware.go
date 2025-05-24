// pkg/middleware/ratelimit_middleware.go
package middleware

import (
	"net/http"
	"sync" // Untuk map yang aman dari concurrent access

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Untuk rate limiting per IP, kita butuh map untuk menyimpan limiter per IP
var (
	mu       sync.Mutex
	limiters = make(map[string]*rate.Limiter)
)

// RateLimitMiddlewarePerIP menerapkan rate limiting berdasarkan IP client.
// r adalah rate (misal, 1 request per detik)
// b adalah burst size (misal, 5 request burst)
func RateLimitMiddlewarePerIP(r rate.Limit, b int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(r, b)
			limiters[ip] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many requests",
				"message": "You have exceeded the request limit. Please try again later.",
			})
			return
		}
		c.Next()
	}
}

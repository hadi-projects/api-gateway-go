// pkg/middleware/logging_middleware.go
package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next() // Proses request berikutnya dalam chain

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			path = path + "?" + c.Request.URL.RawQuery
		}

		log.Printf("[GATEWAY] | %3d | %13v | %15s | %-7s | %s",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				log.Printf("[ERROR] %s", e)
			}
		}
	}
}

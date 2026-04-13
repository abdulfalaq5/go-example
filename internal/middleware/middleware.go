package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a simple request-logging middleware that prints method, path,
// status code, latency, and client IP for every request.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		fmt.Printf("[GIN] %s | %d | %v | %s | %s %s\n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			status,
			latency,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
		)
	}
}

// Recovery catches panics and returns 500 instead of crashing the server.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("[PANIC RECOVERED] %v\n", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "internal server error",
				})
			}
		}()
		c.Next()
	}
}

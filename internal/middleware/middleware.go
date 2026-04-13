package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a structured request-logging middleware.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		args := []any{
			slog.Int("status", status),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("ip", c.ClientIP()),
			slog.String("latency", latency.String()),
		}

		if len(c.Errors) > 0 {
			args = append(args, slog.String("errors", c.Errors.String()))
			slog.Error("Request failed", args...)
		} else if status >= 400 && status < 500 {
			slog.Warn("Client error", args...)
		} else if status >= 500 {
			slog.Error("Server error", args...)
		} else {
			slog.Info("Request handled", args...)
		}
	}
}

// Recovery catches panics, logs them via slog, and returns 500 instead of crashing the server.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Panic recovered", 
					slog.Any("error", err), 
					slog.String("path", c.Request.URL.Path),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "internal server error",
				})
			}
		}()
		c.Next()
	}
}

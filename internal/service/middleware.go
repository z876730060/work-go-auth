package service

import (
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func BaseMiddleware() gin.HandlerFunc {
	total := atomic.Int64{}
	count := atomic.Int64{}
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("recover from panic", "err", err)
			}
		}()
		start := time.Now()
		defer func() {
			slog.Info("request duration", "path", c.Request.URL.Path, "method", c.Request.Method, "status", c.Writer.Status(), "duration", time.Since(start))
		}()
		total.Add(1)
		count.Add(1)
		defer count.Add(-1)
		slog.Info("request count", "total", total.Load(), "current", count.Load())
		c.Next()
	}
}

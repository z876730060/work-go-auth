package service

import (
	"compress/gzip"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/auth/internal/service/common"
)

type HealthService struct {
	log *slog.Logger
}

func NewHealthService(log *slog.Logger) *HealthService {
	return &HealthService{
		log: log.With(common.HANDLER, "healthHandler"),
	}
}

func (h *HealthService) Health(c *gin.Context) {
	c.Header("Connection", "keep-alive")
	c.Header("Content-Encoding", "gzip")
	gzipWriter := gzip.NewWriter(c.Writer)
	defer gzipWriter.Close()
	gzipWriter.Write([]byte("ok"))
	gzipWriter.Flush()
}

func (h *HealthService) Register(e *gin.Engine) {
	e.GET("/health", h.Health)
}

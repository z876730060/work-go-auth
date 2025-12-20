package service

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/auth/internal/service/common"
	"github.com/z876730060/auth/pkg/feign"
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
	news, err := feign.DataServiceInstance.GetNew(1)
	if err != nil {
		h.log.Error("get news error", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	h.log.Info("get news", "news", news)
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (h *HealthService) Register(e *gin.Engine) {
	e.GET("/health", h.Health)
}

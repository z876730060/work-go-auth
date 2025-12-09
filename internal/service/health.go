package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (h *HealthService) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (h *HealthService) Register(e *gin.Engine) {
	e.GET("/health", h.Health)
}

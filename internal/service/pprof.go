package service

import (
	"log/slog"
	"net/http/pprof"
	"strings"

	"github.com/gin-gonic/gin"
)

type pprofHandler struct {
	l *slog.Logger
}

func NewPprofHandler(l *slog.Logger) *pprofHandler {
	return &pprofHandler{l: l}
}

func (h *pprofHandler) Register(e *gin.Engine) {
	e.GET("/debug/pprof/*urlPath", h.ServeHTTP)
}

func (h *pprofHandler) ServeHTTP(c *gin.Context) {
	urlPath := c.Param("urlPath")
	switch {
	case strings.HasPrefix(urlPath, "/cmdline"):
		pprof.Cmdline(c.Writer, c.Request)
	case strings.HasPrefix(urlPath, "/profile"):
		pprof.Profile(c.Writer, c.Request)
	case strings.HasPrefix(urlPath, "/symbol"):
		pprof.Symbol(c.Writer, c.Request)
	case strings.HasPrefix(urlPath, "/trace"):
		pprof.Trace(c.Writer, c.Request)
	default:
		pprof.Index(c.Writer, c.Request)
	}
}

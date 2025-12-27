package service

import (
	"log/slog"
	"net/http/pprof"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/auth/internal/service/common"
)

type pprofHandler struct {
	l *slog.Logger
}

func NewPprofHandler(l *slog.Logger) *pprofHandler {
	return &pprofHandler{l: l.With(common.HANDLER, "pprofHandler")}
}

func (h *pprofHandler) Register(e *gin.Engine) {
	if !Cfg.Application.Debug.Enable {
		return
	}
	h.l.Info("Register pprof handler", "username", Cfg.Application.Debug.Username)
	// 只有管理员的角色才可以访问
	e.GET("/debug/pprof/*urlPath", gin.BasicAuth(gin.Accounts{Cfg.Application.Debug.Username: Cfg.Application.Debug.Password}), h.ServeHTTP)
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

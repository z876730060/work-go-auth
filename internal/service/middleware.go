package service

import (
	"log/slog"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/auth/internal/service/common"
	"github.com/z876730060/auth/internal/service/user"
)

func AuthMiddleware(l *slog.Logger) gin.HandlerFunc {
	l = l.With(common.HANDLER, "authMiddleware")
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		l.Info("Authorization", "header", header)

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := common.ValidateJavaJWT(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		l.Info("Authorization", "claims", claims)

		var userRole []user.UserRole
		db.Model(&user.UserRole{}).Where("user_id = ?", claims.UserID).Find(&userRole)
		roles := make([]uint, 0)
		for _, role := range userRole {
			roles = append(roles, role.RoleID)
		}

		c.Set("userId", claims.UserID)
		c.Set("role", roles)
		c.Set("username", claims.Username)
		// TODO: 验证 Authorization header 是否有效
		// 例如，检查是否包含有效的 token 或其他验证逻辑
		// 如果无效，返回 401 Unauthorized 错误
		// 如果有效，继续处理请求
		c.Next()
	}
}

func BaseMiddleware(l *slog.Logger) gin.HandlerFunc {
	l = l.With(common.HANDLER, "baseMiddleware")
	total := atomic.Int64{}
	count := atomic.Int64{}
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				l.Error("recover from panic", "err", err)
			}
		}()
		start := time.Now()
		defer func() {
			l.Info("request duration", "path", c.Request.URL.Path, "method", c.Request.Method, "status", c.Writer.Status(), "duration", time.Since(start))
		}()
		total.Add(1)
		count.Add(1)
		defer count.Add(-1)
		l.Info("request count", "total", total.Load(), "current", count.Load())
		c.Next()
	}
}

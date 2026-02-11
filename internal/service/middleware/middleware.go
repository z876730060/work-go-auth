package middleware

import (
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/auth/internal/service/common"
	"gorm.io/gorm"
)

func AuthMiddleware(l *slog.Logger, db *gorm.DB) gin.HandlerFunc {
	l = l.With("interceptor", "authMiddleware")
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		l.Debug("Authorization", "header", header)

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := common.ValidateJavaJWT(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		l.Debug("Authorization", "claims", claims)

		c.Set("userId", claims.UserID)
		c.Set("role", claims.RoleIds)
		c.Set("roleKeys", claims.RoleKeys)
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
			l.Info("request duration", "remote", c.Request.RemoteAddr, "path", c.Request.URL.Path, "method", c.Request.Method, "status", c.Writer.Status(), "duration", time.Since(start))
		}()
		total.Add(1)
		count.Add(1)
		defer count.Add(-1)
		l.Info("request count", "total", total.Load(), "current", count.Load())
		c.Next()
	}
}

func RoleMiddleware(l *slog.Logger, roleKey ...string) gin.HandlerFunc {
	l = l.With("interceptor", "roleMiddleware")
	return func(c *gin.Context) {
		roleKeys := c.GetStringSlice("roleKeys")
		if len(roleKeys) == 0 {
			l.Error("roleKeys is empty")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if !sliceContains(roleKey, roleKeys) {
			l.Error("roleKeys not contains roleKey", "roleKey", roleKey, "roleKeys", roleKeys)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}

func sliceContains(slice1 []string, slice2 []string) bool {
	for _, str := range slice2 {
		if slices.Contains(slice1, str) {
			return true
		}
	}
	return false
}

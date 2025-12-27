package login

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/steambap/captcha"
	"github.com/z876730060/auth/internal/service/common"
	"github.com/z876730060/auth/internal/service/role"
	"github.com/z876730060/auth/internal/service/user"
	"gorm.io/gorm"
)

// Handler login处理器
type Handler struct {
	db          *gorm.DB
	l           *slog.Logger
	redisClient *redis.Client
	info        common.Info
}

func NewHandler(baseHandler common.BaseHandler) *Handler {
	return &Handler{
		db:          baseHandler.DB,
		l:           baseHandler.Log.With(slog.String(common.HANDLER, "loginHandler")),
		redisClient: baseHandler.RedisClient,
		info:        baseHandler.Info,
	}
}

func (h *Handler) Register(e *gin.Engine) {
	e.POST("/login", h.Login)
	e.GET("/captcha", h.Captcha)
}

// Login 登录
func (h *Handler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr("invalid request body", h.info))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, common.RespErr(err.Error(), h.info))
		return
	}

	// 校验用户名是否存在
	var u user.User
	if err := h.db.Where(user.User{Username: req.Username}).First(&u).Error; err != nil {
		c.JSON(http.StatusUnauthorized, common.RespErr("username or password is incorrect", h.info))
		return
	}

	// 校验密码是否正确
	if u.Password != req.Password {
		c.JSON(http.StatusUnauthorized, common.RespErr("username or password is incorrect", h.info))
		return
	}

	var roles []role.Role
	h.db.Model(&role.Role{}).Where("ID in (?)", h.db.Model(&user.UserRole{}).Select("role_id").Where("user_id = ?", u.ID)).Find(&roles)
	roleIds := make([]uint, 0)
	roleKeys := make([]string, 0)
	for _, role := range roles {
		roleIds = append(roleIds, role.ID)
		roleKeys = append(roleKeys, role.Key)
	}

	h.l.Info("user login", "username", u.Username, "roleIds", roleIds, "roleKeys", roleKeys)

	// 生成JWT token
	token, err := common.GenerateCompatibleToken(u.ID, u.Username, roleIds, roleKeys)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	// 缓存token
	if err := h.redisClient.Set(c, fmt.Sprintf("jwt:user:%s", u.Username), token, time.Hour).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}

	c.JSON(http.StatusOK, common.RespOk("login success", token, h.info))
}

// Captcha 获取验证码
func (h *Handler) Captcha(c *gin.Context) {
	data, err := captcha.New(150, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.RespErr(err.Error(), h.info))
		return
	}
	data.WriteImage(c.Writer)
}

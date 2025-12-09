package service

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	cap "github.com/steambap/captcha"
)

const (
	RedisKeyCaptcha = "captcha:"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// Captcha 图片验证码
func (a *AuthService) Captcha(ctx *gin.Context) {
	// create a captcha of 150x50px
	data, _ := cap.New(150, 50)
	text := strings.ToLower(data.Text)
	slog.Info("captcha", "text", text)

	requestId := requestid.Get(ctx)
	redisKey := RedisKeyCaptcha + requestId
	// set captcha to redis
	if err := redisClient.Set(ctx.Request.Context(), redisKey, text, time.Minute).Err(); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.SetCookie("captcha", requestId, 60, "/", "", false, true)

	// send image data to client
	data.WriteImage(ctx.Writer)
}

// Login 登录
func (a *AuthService) Login(ctx *gin.Context) {
	type LoginReq struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		RequestID string `json:"request_id"`
		Captcha   string `json:"captcha"`
		Platform  string `json:"platform"`
	}
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	requestId, err := ctx.Cookie("captcha")
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	redisKey := RedisKeyCaptcha + requestId
	// get captcha from redis
	captcha, err := redisClient.Get(ctx.Request.Context(), redisKey).Result()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if req.Captcha != captcha {
		ctx.JSON(400, gin.H{"error": "captcha not match"})
		return
	}
	defer redisClient.Del(ctx.Request.Context(), redisKey)
	// check username and password

	ctx.JSON(200, gin.H{"message": "login success"})
}

// Logout 退出登录
func (a *AuthService) Logout(ctx *gin.Context) {

}

func (a *AuthService) Register(e *gin.Engine) {
	e.GET("/captcha", a.Captcha)
	e.POST("/login", a.Login)
	e.POST("/logout", a.Logout)
}

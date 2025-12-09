package service

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Devtool struct {
}

func NewDevtool() *Devtool {
	return &Devtool{}
}

// Register 注册路由
func (d *Devtool) Register(e *gin.Engine) {
	if Cfg.Application.Env != "dev" {
		return
	}
	e.GET("/.well-known/appspecific/com.chrome.devtools.json", d.GetDevtool)
}

func (d *Devtool) GetDevtool(ctx *gin.Context) {
	type WorkSpace struct {
		Root string `json:"root"`
		Name string `json:"name"`
	}

	wd, err := os.Getwd()
	slog.Info("get current working directory", "wd", wd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "get current working directory failed: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"workspaces": WorkSpace{
			Root: wd,
			Name: Cfg.Application.Name,
		},
	})
}

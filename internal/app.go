package internal

import (
	"log/slog"
	"os"
	"os/signal"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/z876730060/auth/internal/cloud"
	"github.com/z876730060/auth/internal/service"
	"github.com/z876730060/auth/internal/utils"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

// Run 运行应用
func (a *App) Run() {
	service.InitConfig()
	service.InitRedis()
	service.InitDB()
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.Use(gin.Recovery())
	service.InitRoute(e)
	addr := getAddress()
	go e.Run(addr)

	cloud.RegisterManagerInstance.Register(service.Cfg)
	defer cloud.RegisterManagerInstance.Unregister(service.Cfg)

	slog.Info("app start, listen on " + addr)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	<-exit
	slog.Info("receive interrupt signal, exit")
}

// getAddress 获取监听地址
func getAddress() string {
	ip := utils.GetEnv("IP", service.Cfg.Application.IP)
	port := utils.GetEnv("PORT", strconv.Itoa(service.Cfg.Application.Port))

	return ip + ":" + port
}

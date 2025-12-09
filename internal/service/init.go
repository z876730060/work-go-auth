package service

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strconv"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/z876730060/auth/internal/service/common"
	"github.com/z876730060/auth/internal/service/login"
	"github.com/z876730060/auth/internal/service/menu"
	"github.com/z876730060/auth/internal/service/role"
	"github.com/z876730060/auth/internal/service/user"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db          *gorm.DB
	Cfg         Config
	redisClient *redis.Client
)

func InitDB() {
	if !Cfg.DB.Enable {
		return
	}
	var gormDialector = getDialector()

	var err error
	db, err = gorm.Open(gormDialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("db connect failed: " + err.Error())
	}

	menu.InitMenuTable(db)
	role.InitRoleTable(db)
	user.InitUserTable(db)
	slog.Info("db connect success")
}

// InitRoute 初始化路由
func InitRoute(e *gin.Engine) {

	info := common.Info{
		Version:   Cfg.Application.Version,
		GoVersion: runtime.Version(),
	}

	e.Use(BaseMiddleware(), requestid.New())
	NewHealthService().Register(e)
	login.NewHandler(slog.Default(), db, info, redisClient).Register(e)
	e.Use(AuthMiddleware())

	role.NewHandler(slog.Default(), db, info).Register(e)
	user.NewHandler(slog.Default(), db, redisClient, info).Register(e)
	menu.NewHandler(slog.Default(), db, info).Register(e)
	menu.NewMicroAppHandler(slog.Default(), info, db).Register(e)
	slog.Info("route register success")
}

func InitConfig() {
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		panic("read config file failed: " + err.Error())
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		panic("unmarshal config file failed: " + err.Error())
	}
	slog.Info("config load success", "path", viper.ConfigFileUsed())
}

func InitRedis() {
	if !Cfg.Redis.Enable {
		return
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr: Cfg.Redis.Ip + ":" + strconv.Itoa(int(Cfg.Redis.Port)),
		DB:   Cfg.Redis.DB,
	})
	pong := redisClient.Ping(context.Background())
	if err := pong.Err(); err != nil {
		slog.Error("redis ping failed", "err", err)
		return
	}
	slog.Info("redis connect success")
}

func getDialector() gorm.Dialector {
	switch Cfg.DB.Type {
	case "mysql":
		return mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			Cfg.DB.Username, Cfg.DB.Password, Cfg.DB.Ip, Cfg.DB.Port, Cfg.DB.DBName))
	case "postgres":
		return postgres.Open(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			Cfg.DB.Ip, Cfg.DB.Port, Cfg.DB.Username, Cfg.DB.Password, Cfg.DB.DBName))
	default:
		panic("db type not support")
	}
}

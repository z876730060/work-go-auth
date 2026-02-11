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
	"github.com/z876730060/auth/internal/service/dictionary"
	"github.com/z876730060/auth/internal/service/login"
	"github.com/z876730060/auth/internal/service/menu"
	"github.com/z876730060/auth/internal/service/middleware"
	"github.com/z876730060/auth/internal/service/role"
	"github.com/z876730060/auth/internal/service/user"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const ()

var (
	db          *gorm.DB
	Cfg         Config
	redisClient *redis.Client
)

// InitDB 初始化数据库
func InitDB() {
	if !Cfg.DB.Enable {
		return
	}
	var gormDialector = getDialector()

	var err error
	db, err = gorm.Open(gormDialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		panic("db connect failed: " + err.Error())
	}

	menu.InitMenuTable(db)
	role.InitRoleTable(db)
	user.InitUserTable(db)
	dictionary.InitDictionaryTable(db)
	slog.Info("db connect success")
}

// InitRoute 初始化路由
func InitRoute(e *gin.Engine) {

	info := common.Info{
		Version:   Cfg.Application.Version,
		GoVersion: runtime.Version(),
	}

	l := slog.Default().With("service", "auth")
	baseHandler := common.NewBaseHandler(l, db, redisClient, info)

	e.Use(middleware.BaseMiddleware(l), requestid.New())
	NewHealthService(l).Register(e)
	NewPprofHandler(l).Register(e)
	login.NewHandler(baseHandler).Register(e)
	e.Use(middleware.AuthMiddleware(l, db))

	role.NewHandler(baseHandler).Register(e)
	user.NewHandler(baseHandler).Register(e)
	menu.NewHandler(baseHandler).Register(e)
	menu.NewMicroAppHandler(baseHandler).Register(e)
	dictionary.NewHandler(baseHandler).Register(e)
	slog.Info("route register success")
}

// InitConfig 初始化配置
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

// InitRedis 初始化缓存
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

// getDialector 获取数据源
func getDialector() gorm.Dialector {
	switch Cfg.DB.Type {
	case "mysql":
		return mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			Cfg.DB.Username, Cfg.DB.Password, Cfg.DB.Ip, Cfg.DB.Port, Cfg.DB.DBName))
	case "postgres":
		return postgres.Open(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			Cfg.DB.Ip, Cfg.DB.Port, Cfg.DB.Username, Cfg.DB.Password, Cfg.DB.DBName))
	case "sqlite":
		return sqlite.Open(Cfg.DB.DBName)
	default:
		panic("db type not support")
	}
}

package common

import (
	"log/slog"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	HANDLER = "handler"
)

type BaseHandler struct {
	Log         *slog.Logger
	DB          *gorm.DB
	RedisClient *redis.Client
	Info        Info
}

func NewBaseHandler(log *slog.Logger, db *gorm.DB, redisClient *redis.Client, info Info) BaseHandler {
	return BaseHandler{
		Log:         log,
		DB:          db,
		RedisClient: redisClient,
		Info:        info,
	}
}

package global

import (
	"ginProject/config"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	GvaDb     *gorm.DB
	GvaRedis  *redis.Client
	GvaConfig config.Server
	GvaVp    *viper.Viper
	GvaLog   *zap.Logger
)
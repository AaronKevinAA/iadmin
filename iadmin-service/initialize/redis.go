package initialize

import (
	"ginProject/global"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func Redis()  {
	redisCfg := global.GvaConfig.Redis
	client := redis.NewClient(&redis.Options{
		Addr: redisCfg.Addr,
		Password: redisCfg.Password,
		DB: redisCfg.DB,
	})
	_,err := client.Ping().Result()
	if err!= nil{
		global.GvaLog.Error("redis 连接失败, err:", zap.Any("err", err))
	}else{
		global.GvaLog.Info("redis 连接成功:", zap.Any("redisCfg", redisCfg))
		global.GvaRedis = client
	}
}
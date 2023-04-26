package initialize

import (
	"chatgpt-web-go/global"
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func InitRedis() {
	Setup()
}

func Setup() error {
	client := redis.NewClient(&redis.Options{
		Addr:            global.Cfg.Redis.Host,
		Password:        global.Cfg.Redis.Password,
		ConnMaxIdleTime: global.Cfg.Redis.IdleTimeout,
		MaxIdleConns:    global.Cfg.Redis.MaxIdle,
		PoolSize:        global.Cfg.Redis.MaxActive,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.Gzap.Error("redis connect ping failed, err:", zap.Any("err", err))
		return err
	} else {
		global.Gredis = client
		return nil
	}
}

package redis

import (
	"chatgpt-web-go/src/global"
	"context"
	"encoding/json"
	"time"
)

var ctx = context.Background()

func Set(key string, data interface{}, exTime time.Duration) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = global.Gredis.SetNX(ctx, key, value, exTime).Err()
	if err != nil {
		return err
	}
	return nil
}

func Exist(key string) bool {
	val, err := global.Gredis.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	if val == 0 {
		return false
	}
	return true
}

func Get(key string) (string, error) {
	value, err := global.Gredis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func Delete(key string) (bool, error) {
	err := global.Gredis.Del(ctx, key).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func LikeDeletes(key string) error {
	keys, err := global.Gredis.Keys(ctx, "*"+key+"*").Result()
	if err != nil {
		return err
	}
	err = global.Gredis.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}
	return nil
}

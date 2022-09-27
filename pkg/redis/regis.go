package redisService

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/trucktrace/pkg/logger"
)

var ctx = context.Background()

type RedisService struct {
	redis *redis.Client
}

func (redisService *RedisService) SetKeyValue(key string, value string) error {
	_, err := redisService.redis.Set(ctx, key, value, 0).Result()
	if err != nil {
		logger.ErrorLogger("SetKeyValue", "Can't set key:value").Error("Error - " + err.Error())
		return err
	}
	logger.ErrorLogger("SetKeyValue", "Set key:value successful").Info("Set key:value successful")

	return nil
}

func (redisService *RedisService) DeleteKeyValue(key string) error {
	err := redisService.redis.Del(ctx, key).Err()

	if err != nil {
		logger.ErrorLogger("DeleteKeyValue", "Can't delete key:value").Error("Error - " + err.Error())
		return err
	}

	logger.ErrorLogger("DeleteKeyValue", "Delete key:value successful").Info("Delete key:value successful")

	return nil
}

func (redisService *RedisService) GetValueFromRedisStore(key string) (value string, getError error) {
	token, err := redisService.redis.Get(ctx, key).Result()

	if err != nil {
		logger.ErrorLogger("GetValueFromRedisStore", "Can't get value from redis store").Error("Error - " + err.Error())

		return "", err
	}

	logger.ErrorLogger("GetValueFromRedisStore", "Value from redis store was gotten").Info("Value from redis store was gotten")

	return token, nil
}

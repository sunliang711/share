package database

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Ctx context.Context
	Rdb *redis.Client
)

func init() {
	Ctx = context.Background()
	if viper.GetBool("redis.enable") {
		Rdb = redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis.address"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.DB"),
			PoolSize: viper.GetInt("redis.poolSize"),
		})
		_, err := Rdb.Ping(Ctx).Result()
		if err != nil {
			logrus.Fatalf("Connect redis error: %v", err)
		}
		logrus.Info("Redis connected")

	}
}

func closeRedis() {
	if viper.GetBool("redis.enable") {
		Rdb.Close()
	}
}

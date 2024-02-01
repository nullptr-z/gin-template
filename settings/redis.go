package settings

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rds *redis.Client

func InitRedis() error {

	rds = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.addr"), viper.GetInt("redis.port")),
		Password: viper.GetString("reds.password"),
		DB:       viper.GetInt("reds.db"),
		PoolSize: viper.GetInt("reds.pool_size"),
	})

	if _, err := rds.Ping().Result(); err != nil {
		zap.L().Error("Redis ping", zap.Error(err))
		return err
	}
	fmt.Println("Redis initialized........")
	return nil
}

func CloseRedis() {
	rds.Close()
}

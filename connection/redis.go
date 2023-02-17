package connection

import (
	"context"
	"fmt"
	"kedai/backend/be-kedai/config"

	"github.com/redis/go-redis/v9"
)

var (
	cacheConfig = config.Cache
	rdc         *redis.Client
)

var ctx = context.Background()

func ConnectCache() error {
	rdc = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cacheConfig.Host, cacheConfig.Port),
		Password: cacheConfig.Password,
		DB:       0,
	})

	err := rdc.Ping(ctx).Err()

	return err
}

func GetCache() *redis.Client {
	return rdc
}

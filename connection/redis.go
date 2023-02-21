package connection

import (
	"context"
	"crypto/tls"
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
		Addr:      fmt.Sprintf("%s:%s", cacheConfig.Host, cacheConfig.Port),
		Username:  cacheConfig.Username,
		Password:  cacheConfig.Password,
		DB:        0,
		TLSConfig: &tls.Config{},
	})

	err := rdc.Ping(ctx).Err()

	return err
}

func GetCache() *redis.Client {
	return rdc
}

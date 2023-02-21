package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCache interface {
	StoreToken(userId int, accessToken string, refreshToken string) (error)
}

type userCacheImpl struct {
	rdc *redis.Client
}

type UserCConfig struct {
	RDC *redis.Client
}

func NewUserCache(cfg *UserCConfig) UserCache {
	return &userCacheImpl{
		rdc: cfg.RDC,
	}
}

func (r *userCacheImpl) StoreToken(userId int, accessToken string, refreshToken string) (error) {
	refreshKey := fmt.Sprintf("%d:%s", userId, refreshToken)
	refreshTime := time.Duration(24) * time.Hour

	accessKey := fmt.Sprintf("%d:%s", userId, accessToken)
	accessTime := time.Duration(5) * time.Minute

	errRefresh := r.rdc.Set(context.Background(), refreshKey, 0, refreshTime).Err()
	if errRefresh != nil {
		return errRefresh
	}

	errAccess := r.rdc.Set(context.Background(), accessKey, 0, accessTime).Err()
	if errAccess != nil {
		return errAccess
	}

	return nil
}
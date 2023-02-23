package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCache interface {
	StoreToken(userId int, accessToken string, refreshToken string) error
	FindToken(userId int, token string) error
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

func (r *userCacheImpl) StoreToken(userId int, accessToken string, refreshToken string) error {
	refreshKey := fmt.Sprintf("user_%d:%s", userId, refreshToken)
	refreshTime := time.Duration(24) * time.Hour

	accessKey := fmt.Sprintf("user_%d:%s", userId, accessToken)
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

func (r *userCacheImpl) FindToken(userId int, token string) error {
	key := fmt.Sprintf("user_%d:%s", userId, token)
	err := r.rdc.Get(context.Background(), key).Err()

	if err != nil {
		return err
	}

	return nil
}

package cache

import (
	"context"
	"fmt"
	"kedai/backend/be-kedai/config"
	jwttoken "kedai/backend/be-kedai/internal/utils/jwtToken"

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
	refreshTime := jwttoken.ParseTokenAgeFromENV(config.GetEnv("REFRESH_TOKEN_AGE", ""), "refresh")

	accessKey := fmt.Sprintf("user_%d:%s", userId, accessToken)
	accessTime := jwttoken.ParseTokenAgeFromENV(config.GetEnv("ACCESS_TOKEN_AGE", ""), "access")

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

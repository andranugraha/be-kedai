package cache

import (
	"context"
	"fmt"
	"kedai/backend/be-kedai/config"
	jwttoken "kedai/backend/be-kedai/internal/utils/jwtToken"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCache interface {
	StoreToken(userId int, accessToken string, refreshToken string) error
	DeleteToken(key string) error
	FindToken(userId int, token string) error
	DeleteAllByID(userId int) error
	DeleteRefreshTokenAndAccessToken(userId int, refreshToken string, accessToken string) error
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

func (r *userCacheImpl) StoreUserPasswordAndVerificationCode(userId int, newPassword string, verificationCode string) error {
	expireTime := time.Minute * 10
	key := fmt.Sprintf("user_%d-updatePassword", userId)

	err := r.rdc.HSet(context.Background(), key, "newPassword", newPassword, "verificationCode", verificationCode).Err()
	if err != nil {
		return err
	}

	err = r.rdc.Expire(context.Background(), key, expireTime).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *userCacheImpl) FindUserPasswordAndVerificationCode(userId int) (string, string, error) {
	key := fmt.Sprintf("user_%d-updatePassword", userId)

	newPassword, err := r.rdc.HGet(context.Background(), key, "newPassword").Result()
	if err != nil {
		return "", "", err
	}

	verificationCode, err := r.rdc.HGet(context.Background(), key, "verificationCode").Result()
	if err != nil {
		return "", "", err
	}

	return newPassword, verificationCode, nil
}

func (r *userCacheImpl) FindToken(userId int, token string) error {
	key := fmt.Sprintf("user_%d:%s", userId, token)
	err := r.rdc.Get(context.Background(), key).Err()

	if err != nil {
		return err
	}

	return nil
}

func (r *userCacheImpl) DeleteToken(key string) error {
	return r.rdc.Del(context.Background(), key).Err()
}

func (r *userCacheImpl) DeleteAllByID(userId int) error {
	ctx := context.Background()

	iter := r.rdc.Scan(ctx, 0, fmt.Sprintf("user_%d:*", userId), 0).Iterator()
	for iter.Next(ctx) {
		if err := r.rdc.Unlink(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (r *userCacheImpl) DeleteRefreshTokenAndAccessToken(userId int, refreshToken string, accessToken string) error {
	ctx := context.Background()

	refreshKey := fmt.Sprintf("user_%d:%s", userId, refreshToken)
	accessKey := fmt.Sprintf("user_%d:%s", userId, accessToken)

	if err := r.rdc.Unlink(ctx, refreshKey, accessKey).Err(); err != nil {
		return err
	}

	return nil
}

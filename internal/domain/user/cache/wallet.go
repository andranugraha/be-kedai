package cache

import (
	"context"
	"fmt"
	errs "kedai/backend/be-kedai/internal/common/error"
	"time"

	"github.com/redis/go-redis/v9"
)

type WalletCache interface {
	StorePinAndVerificationCode(userID int, pin string, verificationCode string) error
	FindPinAndVerificationCode(userID int) (string, string, error)
	DeletePinAndVerificationCode(userID int) error
	StoreResetPinToken(userID int, token string) error
	FindResetPinToken(token string) error
	DeleteResetPinToken(token string) error
}

type walletCacheImpl struct {
	rdc *redis.Client
}

type WalletCConfig struct {
	RDC *redis.Client
}

func NewWalletCache(cfg *WalletCConfig) WalletCache {
	return &walletCacheImpl{
		rdc: cfg.RDC,
	}
}

func (c *walletCacheImpl) StorePinAndVerificationCode(userID int, pin string, verificationCode string) error {
	expireTime := time.Minute * 5
	key := fmt.Sprintf("user_%d-updatePin", userID)

	err := c.rdc.HSet(context.Background(), key, "newPin", pin, "verificationCode", verificationCode).Err()
	if err != nil {
		return err
	}

	return c.rdc.Expire(context.Background(), key, expireTime).Err()
}

func (c *walletCacheImpl) FindPinAndVerificationCode(userID int) (string, string, error) {
	key := fmt.Sprintf("user_%d-updatePin", userID)

	pin, err := c.rdc.HGet(context.Background(), key, "newPin").Result()
	if err != nil {
		if err == redis.Nil {
			err = errs.ErrVerificationCodeNotFound
		}
		return "", "", err
	}

	verificationCode, err := c.rdc.HGet(context.Background(), key, "verificationCode").Result()
	if err != nil {
		if err == redis.Nil {
			err = errs.ErrVerificationCodeNotFound
		}
		return "", "", err
	}

	return pin, verificationCode, nil
}

func (c *walletCacheImpl) DeletePinAndVerificationCode(userID int) error {
	key := fmt.Sprintf("user_%d-updatePin", userID)

	err := c.rdc.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *walletCacheImpl) StoreResetPinToken(userID int, token string) error {
	expireTime := time.Minute * 5
	key := fmt.Sprintf("resetPinToken:%s", token)

	err := c.rdc.SetNX(context.Background(), key, userID, expireTime).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *walletCacheImpl) FindResetPinToken(token string) error {
	key := fmt.Sprintf("resetPinToken:%s", token)

	_, err := c.rdc.Get(context.Background(), key).Int()
	if err != nil {
		if err == redis.Nil {
			err = errs.ErrResetPasswordTokenNotFound
		}

		return err
	}

	return nil
}

func (c *walletCacheImpl) DeleteResetPinToken(token string) error {
	key := fmt.Sprintf("resetPinToken:%s", token)
	err := c.rdc.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}

	return nil
}

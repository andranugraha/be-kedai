package cache

import (
	"context"
	"fmt"
	"kedai/backend/be-kedai/config"
	errs "kedai/backend/be-kedai/internal/common/error"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type WalletCache interface {
	FindWalletStepUpErrorCount(walletId int) (*int, error)
	StoreOrIncrementWalletStepUpErrorCount(walletId int) error
	DeleteErrorCount(walletId int) error
	BlockWallet(walletId int) error
	CheckIsWalletBlocked(walletId int) error
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

func (r *walletCacheImpl) FindWalletStepUpErrorCount(walletId int) (*int, error) {
	key := fmt.Sprintf("wallet_%d:step_up_error_count", walletId)
	walletErrorCount, err := r.rdc.Get(context.Background(), key).Int()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	return &walletErrorCount, nil
}

func (r *walletCacheImpl) StoreOrIncrementWalletStepUpErrorCount(walletId int) error {
	key := fmt.Sprintf("wallet_%d:step_up_error_count", walletId)
	err := r.rdc.Incr(context.Background(), key).Err()
	if err != nil {
		if err == redis.Nil {
			errorCount := 1
			expiryDuration, _ := strconv.Atoi(config.GetEnv("BLOCKED_WALLET_AGE", ""))
			err = r.rdc.Set(context.Background(), key, errorCount, time.Minute*time.Duration(expiryDuration)).Err()
			return err
		}

		return err
	}

	return nil
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
			err = errs.ErrResetPinTokenNotFound
		}

		return err
	}

	return nil
}

func (r *walletCacheImpl) DeleteErrorCount(walletId int) error {
	key := fmt.Sprintf("wallet_%d:step_up_error_count", walletId)
	err := r.rdc.Del(context.Background(), key).Err()
	if err != nil {
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

func (r *walletCacheImpl) BlockWallet(walletId int) error {
	key := fmt.Sprintf("wallet_%d:blocked", walletId)
	expiryDuration, _ := strconv.Atoi(config.GetEnv("BLOCKED_WALLET_AGE", ""))
	err := r.rdc.Set(context.Background(), key, 0, time.Minute*time.Duration(expiryDuration)).Err()
	if err != nil {
		return err
	}

	return errs.ErrWalletTemporarilyBlocked
}

func (r *walletCacheImpl) CheckIsWalletBlocked(walletId int) error {
	key := fmt.Sprintf("wallet_%d:blocked", walletId)
	wallet, err := r.rdc.Get(context.Background(), key).Int()
	if err != nil {
		if err == redis.Nil {
			return nil
		}

		return err
	}

	if wallet == 0 {
		return nil
	}

	return errs.ErrWalletTemporarilyBlocked
}

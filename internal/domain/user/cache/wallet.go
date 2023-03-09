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

func (r *walletCacheImpl) DeleteErrorCount(walletId int) error {
	key := fmt.Sprintf("wallet_%d:step_up_error_count", walletId)
	err := r.rdc.Del(context.Background(), key).Err()
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

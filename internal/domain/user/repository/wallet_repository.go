package repository

import (
	errRes "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type WalletRepository interface {
	Create(wallet *model.Wallet) (*model.Wallet, error)
}

type walletRepositoryImpl struct {
	db *gorm.DB
}

type WalletRConfig struct {
	DB *gorm.DB
}

func NewWalletRepository(cfg *WalletRConfig) WalletRepository {
	return &walletRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *walletRepositoryImpl) Create(wallet *model.Wallet) (*model.Wallet, error) {
	err := r.db.Create(wallet).Error
	if err != nil {
		if errRes.IsDuplicateKeyError(err) {
			return nil, errRes.ErrWalletAlreadyExist
		}

		return nil, err
	}

	return wallet, nil
}

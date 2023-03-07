package repository

import (
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type WalletHistoryRepository interface {
	Create(*gorm.DB, *model.WalletHistory) error
}

type walletHistoryRepoImpl struct {
	db *gorm.DB
}

type WalletHistoryRConfig struct {
	DB *gorm.DB
}

func NewWalletHistoryRepository(cfg *WalletHistoryRConfig) WalletHistoryRepository {
	return &walletHistoryRepoImpl{
		db: cfg.DB,
	}
}

func (r *walletHistoryRepoImpl) Create(tx *gorm.DB, history *model.WalletHistory) error {
	err := tx.Create(&history).Error
	return err
}

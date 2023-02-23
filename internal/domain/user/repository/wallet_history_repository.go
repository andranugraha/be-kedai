package repository

import (
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type WalletHistoryRepository interface {
	Create(*gorm.DB, *model.WalletHistory) (error)
}

type walletHistoryImpl struct {
	db *gorm.DB
}

type WalletHConfig struct {
	DB *gorm.DB
}

func NewWalletHistoryRepository(cfg *WalletHConfig) WalletHistoryRepository {
	return &walletHistoryImpl{
		db: cfg.DB,
	}
}

func (r *walletHistoryImpl) Create(tx *gorm.DB ,history *model.WalletHistory) (error) {
	err := tx.Create(&history).Error
	return err
}

package repository

import (
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type WalletHistoryRepository interface {
	Create(*gorm.DB, *model.WalletHistory) error
	GetWalletHistoryById(id int) ([]*model.WalletHistory, error)
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

func (r *walletHistoryRepoImpl) GetWalletHistoryById(id int) ([]*model.WalletHistory, error) {
	var histories []*model.WalletHistory

	err := r.db.Where("wallet_id = ?", id).Order("created_at desc").Find(&histories).Error
	if err != nil {
		return nil, err
	}

	return histories, nil
}

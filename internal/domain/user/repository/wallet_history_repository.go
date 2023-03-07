package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"math"

	"gorm.io/gorm"
)

type WalletHistoryRepository interface {
	Create(*gorm.DB, *model.WalletHistory) error
	GetHistoryDetailById(ref string, wallet *model.Wallet) (*model.WalletHistory, error)
	GetWalletHistoryById(req dto.WalletHistoryRequest, id int) ([]*model.WalletHistory, int64, int, error)
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

func (r *walletHistoryRepoImpl) GetHistoryDetailById(ref string, wallet *model.Wallet) (*model.WalletHistory, error) {
	var history *model.WalletHistory

	err := r.db.Where("reference = ? and wallet_id = ?", ref, wallet.ID).First(&history).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrWalletHistoryDoesNotExist
		}

		return nil, err
	}

	return history, nil
}

func (r *walletHistoryRepoImpl) GetWalletHistoryById(req dto.WalletHistoryRequest, id int) ([]*model.WalletHistory, int64, int, error) {
	var (
		histories []*model.WalletHistory
		totalRows int64
		totalPage int
	)

	err := r.db.Where("wallet_id = ?", id).Order("created_at desc").Limit(req.Limit).Offset(req.Offset()).Find(&histories).Error
	if err != nil {
		return nil, 0, 0, err
	}

	r.db.Model(&model.WalletHistory{}).Where("wallet_id = ?", id).Count(&totalRows)
	totalPage = int(math.Ceil(float64(totalRows) / float64(req.Limit)))

	return histories, totalRows, totalPage, nil
}

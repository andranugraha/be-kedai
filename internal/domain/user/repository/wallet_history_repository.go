package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	shopDto "kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"math"

	"gorm.io/gorm"
)

type WalletHistoryRepository interface {
	Create(*gorm.DB, *model.WalletHistory) error
	CreateMultiple(*gorm.DB, []*model.WalletHistory) error
	GetHistoryDetailById(ref string, wallet *model.Wallet) (*model.WalletHistory, error)
	GetWalletHistoryById(req dto.WalletHistoryRequest, id int) ([]*model.WalletHistory, int64, int, error)
	GetShopFinanceReleased(shopId int) (*shopDto.ShopFinanceReleased, error)
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

func (r *walletHistoryRepoImpl) GetShopFinanceReleased(shopId int) (*shopDto.ShopFinanceReleased, error) {
	var (
		shopFinanceReleased = &shopDto.ShopFinanceReleased{}
	)

	query := r.db.
		Model(&model.WalletHistory{}).
		Select(`
	SUM(amount) AS "total",
	SUM(CASE WHEN date >= date_trunc('month', now()) THEN amount ELSE 0 END) AS "month",
	SUM(CASE WHEN date >= date_trunc('week', now()) THEN amount ELSE 0 END) AS "week"
`).
		Joins("JOIN wallets ON wallets.id = wallet_histories.wallet_id").
		Joins("JOIN users ON users.id = wallets.user_id").
		Joins("JOIN shops ON shops.user_id = users.id").
		Where("type = 'Withdrawal'").Where("shops.id = ?", shopId)

	err := query.Find(&shopFinanceReleased).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shopFinanceReleased, nil
		}
		return nil, err
	}

	return shopFinanceReleased, nil
}

func (r *walletHistoryRepoImpl) CreateMultiple(tx *gorm.DB, histories []*model.WalletHistory) error {
	err := tx.Create(&histories).Error
	return err
}

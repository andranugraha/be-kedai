package repository

import (
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type UserVoucherRepository interface {
	GetUsedMarketplaceByUserID(userID int) ([]*model.UserVoucher, error)
	GetUsedShopByUserID(userID int) ([]*model.UserVoucher, error)
}

type userVoucherRepositoryImpl struct {
	db *gorm.DB
}

type UserVoucherRConfig struct {
	DB *gorm.DB
}

func NewUserVoucherRepository(cfg *UserVoucherRConfig) UserVoucherRepository {
	return &userVoucherRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *userVoucherRepositoryImpl) GetUsedMarketplaceByUserID(userID int) ([]*model.UserVoucher, error) {
	var userVouchers []*model.UserVoucher

	err := r.db.Where("user_id = ?", userID).Where("is_used = ?", true).Not("marketplace_voucher_id IS NULL").Find(&userVouchers).Error
	if err != nil {
		return nil, err
	}

	return userVouchers, nil
}

func (r *userVoucherRepositoryImpl) GetUsedShopByUserID(userID int) ([]*model.UserVoucher, error) {
	var userVouchers []*model.UserVoucher

	err := r.db.Where("user_id = ?", userID).Where("is_used = ?", true).Not("shop_voucher_id IS NULL").Find(&userVouchers).Error
	if err != nil {
		return nil, err
	}

	return userVouchers, nil
}

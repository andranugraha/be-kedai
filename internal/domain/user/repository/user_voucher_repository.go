package repository

import (
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type UserVoucherRepository interface {
	GetUsedMarketplaceByUserID(userID int) ([]*model.UserVoucher, error)
	GetUsedShopByUserID(userID int) ([]*model.UserVoucher, error)
	UpdateMarketplaceVoucherToUnused(tx *gorm.DB, userID int, marketplaceVoucherId int) error
	UpdateShopVoucherToUnused(tx *gorm.DB, userID int, shopVoucherId int) error
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

func (r *userVoucherRepositoryImpl) UpdateMarketplaceVoucherToUnused(tx *gorm.DB, userID int, marketplaceVoucherId int) error {
	res := tx.Model(&model.UserVoucher{}).
		Where("user_id = ?", userID).
		Where("marketplace_voucher_id = ?", marketplaceVoucherId).
		Update("is_used", false)
	if err := res.Error; err != nil {
		return err
	}

	if res.RowsAffected == 0 {
		return commonErr.ErrVoucherNotFound
	}

	return nil
}

func (r *userVoucherRepositoryImpl) UpdateShopVoucherToUnused(tx *gorm.DB, userID int, shopVoucherId int) error {
	res := tx.Model(&model.UserVoucher{}).
		Where("user_id = ?", userID).
		Where("shop_voucher_id = ?", shopVoucherId).
		Update("is_used", false)
	if err := res.Error; err != nil {
		return err
	}

	if res.RowsAffected == 0 {
		return commonErr.ErrVoucherNotFound
	}

	return nil
}

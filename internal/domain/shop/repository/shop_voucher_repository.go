package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"time"

	userRepo "kedai/backend/be-kedai/internal/domain/user/repository"

	"gorm.io/gorm"
)

type ShopVoucherRepository interface {
	GetShopVoucher(shopId int) ([]*model.ShopVoucher, error)
	GetValidByIdAndUserId(id, userId int) (*model.ShopVoucher, error)
	GetValidByUserIDAndShopID(userID int, shopID int) ([]*model.ShopVoucher, error)
}

type shopVoucherRepositoryImpl struct {
	db                    *gorm.DB
	userVoucherRepository userRepo.UserVoucherRepository
}

type ShopVoucherRConfig struct {
	DB                    *gorm.DB
	UserVoucherRepository userRepo.UserVoucherRepository
}

func NewShopVoucherRepository(cfg *ShopVoucherRConfig) ShopVoucherRepository {
	return &shopVoucherRepositoryImpl{
		db:                    cfg.DB,
		userVoucherRepository: cfg.UserVoucherRepository,
	}
}

func (r *shopVoucherRepositoryImpl) GetShopVoucher(shopId int) ([]*model.ShopVoucher, error) {
	var shopVoucher []*model.ShopVoucher
	publicVoucher := true
	now := time.Now()
	err := r.db.Where("shop_id = ?", shopId).Where("is_hidden != ?", publicVoucher).Where("? < expired_at", now).Find(&shopVoucher).Error
	if err != nil {
		return nil, err
	}

	return shopVoucher, nil
}

func (r *shopVoucherRepositoryImpl) GetValidByIdAndUserId(id, userId int) (*model.ShopVoucher, error) {
	var (
		shopVoucher      model.ShopVoucher
		invalidVoucherID []int
	)

	userVoucher, err := r.userVoucherRepository.GetUsedShopByUserID(userId)
	if err != nil {
		return nil, err
	}

	for _, voucher := range userVoucher {
		if voucher.ShopVoucherId != nil {
			invalidVoucherID = append(invalidVoucherID, *voucher.ShopVoucherId)
		}
	}

	db := r.db

	if len(invalidVoucherID) > 0 {
		db = db.Where("id NOT IN (?)", invalidVoucherID)
	}

	err = db.Where("expired_at > now()").First(&shopVoucher, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrInvalidVoucher
		}
		return nil, err
	}

	return &shopVoucher, nil
}

func (r *shopVoucherRepositoryImpl) GetValidByUserIDAndShopID(userID int, shopID int) ([]*model.ShopVoucher, error) {
	var shopVouchers []*model.ShopVoucher
	var invalidVoucherID []int

	userVoucher, err := r.userVoucherRepository.GetUsedShopByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, voucher := range userVoucher {
		if voucher.ShopVoucherId != nil {
			invalidVoucherID = append(invalidVoucherID, *voucher.ShopVoucherId)
		}
	}

	db := r.db
	if len(invalidVoucherID) > 0 {
		db = db.Not("id IN (?)", invalidVoucherID)
	}

	publicVoucher := true
	err = db.Where("shop_id = ?", shopID).
		Where("is_hidden != ?", publicVoucher).
		Where("? < expired_at", time.Now()).
		Find(&shopVouchers).Error
	if err != nil {
		return nil, err
	}

	return shopVouchers, nil
}

package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/model"

	"gorm.io/gorm"
)

type ShopRepository interface {
	FindShopById(id int) (*model.Shop, error)
	FindShopByUserId(userId int) (*model.Shop, error)
	FindShopBySlug(slug string) (*model.Shop, error)
}

type shopRepositoryImpl struct {
	db          *gorm.DB
	voucherRepo ShopVoucherRepository
}

type ShopRConfig struct {
	DB          *gorm.DB
	VoucherRepo ShopVoucherRepository
}

func NewShopRepository(cfg *ShopRConfig) ShopRepository {
	return &shopRepositoryImpl{
		db:          cfg.DB,
		voucherRepo: cfg.VoucherRepo,
	}
}

func (r *shopRepositoryImpl) FindShopById(id int) (*model.Shop, error) {
	var shop model.Shop

	err := r.db.Where("id = ?", id).First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrShopNotFound
		}

		return nil, err
	}

	return &shop, err
}

func (r *shopRepositoryImpl) FindShopByUserId(userId int) (*model.Shop, error) {
	var shop model.Shop

	err := r.db.Where("user_id = ?", userId).First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrShopNotFound
		}

		return nil, err
	}

	return &shop, err
}

func (r *shopRepositoryImpl) FindShopBySlug(slug string) (*model.Shop, error) {
	var shop model.Shop

	err := r.db.Where("slug = ?", slug).Preload("ShopCategory").First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrShopNotFound
		}
	}

	voucher, errVoucher := r.voucherRepo.GetShopVoucher(shop.ID)
	if errVoucher != nil {
		return nil, errVoucher
	}

	shop.ShopVoucher = voucher

	return &shop, nil
}

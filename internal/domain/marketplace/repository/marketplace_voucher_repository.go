package repository

import (
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"time"

	userRepo "kedai/backend/be-kedai/internal/domain/user/repository"

	"gorm.io/gorm"
)

type MarketplaceVoucherRepository interface {
	GetMarketplaceVoucher(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error)
	GetValidByUserID(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error)
}

type marketplaceVoucherRepositoryImpl struct {
	db                    *gorm.DB
	userVoucherRepository userRepo.UserVoucherRepository
}

type MarketplaceVoucherRConfig struct {
	DB                    *gorm.DB
	UserVoucherRepository userRepo.UserVoucherRepository
}

func NewMarketplaceVoucherRepository(cfg *MarketplaceVoucherRConfig) MarketplaceVoucherRepository {
	return &marketplaceVoucherRepositoryImpl{
		db:                    cfg.DB,
		userVoucherRepository: cfg.UserVoucherRepository,
	}
}

func (r *marketplaceVoucherRepositoryImpl) GetMarketplaceVoucher(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error) {
	var marketplaceVoucher []*model.MarketplaceVoucher

	db := r.db

	if req.CategoryId != 0 {
		db = db.Where("category_id = ?", req.CategoryId)
	}
	if req.PaymentMethodId != 0 {
		db = db.Where("payment_method_id = ?", req.PaymentMethodId)
	}

	publicVoucher := true
	err := db.Where("expired_at > ?", time.Now()).Where("is_hidden != ?", publicVoucher).Find(&marketplaceVoucher).Error
	if err != nil {
		return nil, err
	}

	return marketplaceVoucher, nil
}

func (r *marketplaceVoucherRepositoryImpl) GetValidByUserID(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error) {
	var marketplaceVoucher []*model.MarketplaceVoucher
	var invalidVoucherID []int

	userVoucher, err := r.userVoucherRepository.GetUsedMarketplaceByUserID(req.UserId)
	if err != nil {
		return nil, err
	}

	for _, voucher := range userVoucher {
		invalidVoucherID = append(invalidVoucherID, voucher.MarketplaceVoucherId)
	}

	db := r.db

	if req.CategoryId != 0 {
		db = db.Where("category_id = ?", req.CategoryId)
	}
	if req.PaymentMethodId != 0 {
		db = db.Where("payment_method_id = ?", req.PaymentMethodId)
	}

	if len(invalidVoucherID) > 0 {
		db = db.Not("id IN (?)", invalidVoucherID)
	}

	publicVoucher := true
	err = db.Where("expired_at > ?", time.Now()).
		Where("is_hidden != ?", publicVoucher).
		Find(&marketplaceVoucher).Error
	if err != nil {
		return nil, err
	}

	return marketplaceVoucher, nil
}

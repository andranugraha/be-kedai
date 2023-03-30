package repository

import (
	"errors"
	"fmt"
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"math"
	"time"

	"kedai/backend/be-kedai/internal/common/constant"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	userRepo "kedai/backend/be-kedai/internal/domain/user/repository"

	"gorm.io/gorm"
)

type MarketplaceVoucherRepository interface {
	GetMarketplaceVoucher(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error)
	GetMarketplaceVoucherAdminByCode(voucherCode string) (*dto.AdminMarketplaceVoucher, error)
	GetMarketplaceVoucherAdmin(request *dto.AdminVoucherFilterRequest) ([]*dto.AdminMarketplaceVoucher, int64, int, error)
	GetValidByUserID(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error)
	GetValid(id, userID, PaymentMethodID int) (*model.MarketplaceVoucher, error)
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

func (r *marketplaceVoucherRepositoryImpl) GetMarketplaceVoucherAdminByCode(voucherCode string) (*dto.AdminMarketplaceVoucher, error) {
	var voucher dto.AdminMarketplaceVoucher

	now := time.Now()
	query := r.db.Where("code = ?", voucherCode)

	query = query.Select("marketplace_vouchers.*, "+
		"CASE WHEN expired_at >= ? THEN ? "+
		"ELSE ? "+
		"END as status", now, constant.VoucherPromotionStatusOngoing, constant.VoucherPromotionStatusExpired)

	err := query.First(&voucher).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commonErr.ErrVoucherNotFound
		}
		return nil, err
	}

	return &voucher, nil
}

func (r *marketplaceVoucherRepositoryImpl) GetMarketplaceVoucherAdmin(request *dto.AdminVoucherFilterRequest) ([]*dto.AdminMarketplaceVoucher, int64, int, error) {
	var (
		marketplaceVoucher []*dto.AdminMarketplaceVoucher
		totalRows          int64
		totalPages         int
	)

	now := time.Now()
	query := r.db

	if request.Name != "" {
		query = query.Where("marketplace_vouchers.name ILIKE ?", fmt.Sprintf("%%%s%%", request.Name))
	}
	if request.Code != "" {
		query = query.Where("marketplace_vouchers.code ILIKE ?", fmt.Sprintf("%%%s%%", request.Code))
	}

	switch request.Status {
	case constant.VoucherPromotionStatusOngoing:
		query = query.Where("? <= marketplace_vouchers.expired_at", now)
	case constant.VoucherPromotionStatusExpired:
		query = query.Where("? > marketplace_vouchers.expired_at", now)
	}

	query = query.Select("marketplace_vouchers.*, "+
		"CASE WHEN expired_at >= ? THEN ? "+
		"ELSE ? "+
		"END as status", now, constant.VoucherPromotionStatusOngoing, constant.VoucherPromotionStatusExpired)

	query = query.Session(&gorm.Session{})

	err := query.Model(&model.MarketplaceVoucher{}).Distinct("marketplace_vouchers.id").Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages = int(math.Ceil(float64(totalRows) / float64(request.Limit)))

	err = query.
		Order("marketplace_vouchers.created_at desc").Limit(request.Limit).Offset(request.Offset()).Find(&marketplaceVoucher).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return marketplaceVoucher, totalRows, totalPages, nil
}

func (r *marketplaceVoucherRepositoryImpl) GetValidByUserID(req *dto.GetMarketplaceVoucherRequest) ([]*model.MarketplaceVoucher, error) {
	var marketplaceVoucher []*model.MarketplaceVoucher
	var invalidVoucherID []int

	userVoucher, err := r.userVoucherRepository.GetUsedMarketplaceByUserID(req.UserId)
	if err != nil {
		return nil, err
	}

	for _, voucher := range userVoucher {
		invalidVoucherID = append(invalidVoucherID, *voucher.MarketplaceVoucherId)
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

	if req.Code != "" {
		db = db.Where("code = ?", req.Code)
	} else {
		publicVoucher := true
		db = db.Where("is_hidden != ?", publicVoucher)
	}

	err = db.Where("expired_at > ?", time.Now()).
		Find(&marketplaceVoucher).Error
	if err != nil {
		return nil, err
	}

	return marketplaceVoucher, nil
}

func (r *marketplaceVoucherRepositoryImpl) GetValid(id, userID, PaymentMethodID int) (*model.MarketplaceVoucher, error) {
	var (
		marketplaceVoucher model.MarketplaceVoucher
		invalidVoucherID   []int
	)

	userVoucher, err := r.userVoucherRepository.GetUsedMarketplaceByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, voucher := range userVoucher {
		invalidVoucherID = append(invalidVoucherID, *voucher.MarketplaceVoucherId)
	}

	db := r.db

	if PaymentMethodID != 0 {
		db = db.Where("payment_method_id = ?", PaymentMethodID).Or("payment_method_id is null")
	}

	if len(invalidVoucherID) > 0 {
		db = db.Not("id IN (?)", invalidVoucherID)
	}

	err = db.Where("id = ?", id).
		Where("expired_at > ?", time.Now()).
		First(&marketplaceVoucher).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commonErr.ErrInvalidVoucher
		}

		return nil, err
	}

	return &marketplaceVoucher, nil
}

package repository

import (
	"errors"
	"fmt"
	"kedai/backend/be-kedai/internal/common/constant"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"math"
	"time"

	userRepo "kedai/backend/be-kedai/internal/domain/user/repository"

	"gorm.io/gorm"
)

type ShopVoucherRepository interface {
	GetShopVoucher(shopId int) ([]*model.ShopVoucher, error)
	GetSellerVoucher(shopId int, request *dto.SellerVoucherFilterRequest) ([]*dto.SellerVoucher, int64, int, error)
	GetValidByIdAndUserId(id, userId int) (*model.ShopVoucher, error)
	GetValidByUserIDAndShopID(dto.GetValidShopVoucherRequest, int) ([]*model.ShopVoucher, error)
	Create(shopId int, request *dto.CreateVoucherRequest) (*model.ShopVoucher, error)
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

func (r *shopVoucherRepositoryImpl) GetSellerVoucher(shopId int, request *dto.SellerVoucherFilterRequest) ([]*dto.SellerVoucher, int64, int, error) {
	var (
		vouchers   []*dto.SellerVoucher
		totalRows  int64
		totalPages int
	)

	now := time.Now()
	query := r.db.Where("shop_id = ?", shopId)

	if request.Name != "" {
		query = query.Where("shop_vouchers.name ILIKE ?", fmt.Sprintf("%%%s%%", request.Name))
	}
	if request.Code != "" {
		query = query.Where("shop_vouchers.code ILIKE ?", fmt.Sprintf("%%%s%%", request.Code))
	}

	switch request.Status {
	case constant.VoucherPromotionStatusOngoing:
		query = query.Where("shop_vouchers.start_from <= ? AND ? < shop_vouchers.expired_at", now, now)
	case constant.VoucherPromotionStatusUpcoming:
		query = query.Where("? < shop_vouchers.start_from", now)
	case constant.VoucherPromotionStatusExpired:
		query = query.Where("shop_vouchers.expired_at <= ?", now)
	}

	query = query.Select("shop_vouchers.*, "+
		"CASE WHEN start_from <= ? AND expired_at >= ? THEN ? "+
		"WHEN start_from > ? THEN ? "+
		"ELSE ? "+
		"END as status", now, now, constant.VoucherPromotionStatusOngoing, now, constant.VoucherPromotionStatusUpcoming, constant.VoucherPromotionStatusExpired)

	query = query.Session(&gorm.Session{})

	err := query.Model(&model.ShopVoucher{}).Distinct("shop_vouchers.id").Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages = int(math.Ceil(float64(totalRows) / float64(request.Limit)))

	err = query.
		Order("shop_vouchers.created_at desc").Limit(request.Limit).Offset(request.Offset()).Find(&vouchers).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return vouchers, totalRows, totalPages, nil
}

func (r *shopVoucherRepositoryImpl) Create(shopId int, request *dto.CreateVoucherRequest) (*model.ShopVoucher, error) {
	tx := r.db.Begin()
	defer tx.Commit()

	voucher := &model.ShopVoucher{
		Name:         request.Name,
		Code:         request.Code,
		Amount:       request.Amount,
		Type:         request.Type,
		IsHidden:     request.IsHidden,
		Description:  request.Description,
		MinimumSpend: request.MinimumSpend,
		UsedQuota:    0,
		TotalQuota:   request.TotalQuota,
		StartFrom:    request.StartFrom,
		ExpiredAt:    request.ExpiredAt,
		ShopId:       shopId,
	}

	err := tx.Create(voucher).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return voucher, nil
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

func (r *shopVoucherRepositoryImpl) GetValidByUserIDAndShopID(req dto.GetValidShopVoucherRequest, shopID int) ([]*model.ShopVoucher, error) {
	var shopVouchers []*model.ShopVoucher
	var invalidVoucherID []int

	userVoucher, err := r.userVoucherRepository.GetUsedShopByUserID(req.UserID)
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

	if req.Code != "" {
		db = db.Where("code = ?", req.Code)
	} else {
		publicVoucher := true
		db = db.Where("is_hidden != ?", publicVoucher)
	}
	err = db.Where("shop_id = ?", shopID).
		Where("? < expired_at", time.Now()).
		Find(&shopVouchers).Error
	if err != nil {
		return nil, err
	}

	return shopVouchers, nil
}

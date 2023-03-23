package repository

import (
	"fmt"
	"kedai/backend/be-kedai/internal/common/constant"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"math"
	"time"

	"gorm.io/gorm"
)

type ShopPromotionRepository interface {
}

type shopPromotionRepositoryImpl struct {
	db *gorm.DB
}

type ShopPromotionRConfig struct {
	DB *gorm.DB
}

func NewShopPromotionRepository(cfg *ShopPromotionRConfig) ShopPromotionRepository {
	return &shopPromotionRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *shopPromotionRepositoryImpl) GetSellerPromotions(shopId int, request *dto.SellerPromotionFilterRequest) ([]*dto.SellerPromotion, int64, int, error) {
	var (
		promotions []*dto.SellerPromotion
		totalRows  int64
		totalPages int
	)

	now := time.Now()
	query := r.db.Where("shop_id = ?", shopId)

	if request.Name != "" {
		query = query.Where("shop_promotions.name ILIKE ?", fmt.Sprintf("%%%s%%", request.Name))
	}

	switch request.Status {
	case constant.VoucherPromotionStatusOngoing:
		query = query.Where("shop_promotions.start_from <= ? AND ? < shop_promotions.expired_at", now, now)
	case constant.VoucherPromotionStatusUpcoming:
		query = query.Where("? < shop_promotions.start_from", now)
	case constant.VoucherPromotionStatusExpired:
		query = query.Where("shop_promotions.expired_at <= ?", now)
	}

	query = query.Select("shop_promotions.*, "+
		"CASE WHEN start_period <= ? AND end_period >= ? THEN ? "+
		"WHEN start_from > ? THEN ? "+
		"ELSE ? "+
		"END as status", now, now, constant.VoucherPromotionStatusOngoing, now, constant.VoucherPromotionStatusUpcoming, constant.VoucherPromotionStatusExpired)

	query = query.Session(&gorm.Session{})

	err := query.Model(&model.ShopPromotion{}).Distinct("shop_promotions.id").Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages = int(math.Ceil(float64(totalRows) / float64(request.Limit)))

	err = query.
		Order("shop_promotions.created_at desc").Limit(request.Limit).Offset(request.Offset()).Find(&promotions).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return promotions, totalRows, totalPages, nil
}

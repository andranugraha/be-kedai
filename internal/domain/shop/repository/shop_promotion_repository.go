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
	GetSellerPromotions(shopId int, request *dto.SellerPromotionFilterRequest) ([]*dto.SellerPromotion, int64, int, error)
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
	query := r.db.
		Select("shop_promotions.*,"+
			"CASE WHEN start_period <= ? AND end_period >= ? THEN ? "+
			"WHEN start_from > ? THEN ? "+
			"ELSE ? "+
			"END as status, "+
			"products.id as product_id, products.name as product_name, products.code as product_code, "+
			"(SELECT url FROM product_medias pm WHERE pm.product_id = products.id LIMIT 1) AS image_url").
		Joins("JOIN shops ON shops.id = shop_promotions.shop_id").
		Joins("JOIN skus ON skus.id = shop_promotions.sku_id").
		Joins("JOIN products ON products.id = skus.product_id").
		Where("shop_promotions.shop_id = ?", shopId)

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
		Order("shop_promotions.created_at desc").Preload("SKUs.Variants").Preload("SKUs.Promotion").Limit(request.Limit).Offset(request.Offset()).Find(&promotions).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return promotions, totalRows, totalPages, nil
}

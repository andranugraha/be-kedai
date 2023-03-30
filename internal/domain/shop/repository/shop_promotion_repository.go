package repository

import (
	"errors"
	"fmt"
	"kedai/backend/be-kedai/internal/common/constant"
	errs "kedai/backend/be-kedai/internal/common/error"
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"math"
	"time"

	productRepo "kedai/backend/be-kedai/internal/domain/product/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShopPromotionRepository interface {
	GetSellerPromotions(shopId int, request *dto.SellerPromotionFilterRequest) ([]*dto.SellerPromotion, int64, int, error)
	GetSellerPromotionById(shopId int, promotionId int) (*dto.SellerPromotion, error)
	Update(shopPromotion *model.ShopPromotion, productPromotion []*productModel.ProductPromotion) error
	Create(shopID int, request *dto.CreateShopPromotionRequest) (*dto.CreateShopPromotionResponse, error)
}

type shopPromotionRepositoryImpl struct {
	db                *gorm.DB
	productRepository productRepo.ProductRepository
}

type ShopPromotionRConfig struct {
	DB                *gorm.DB
	ProductRepository productRepo.ProductRepository
}

func NewShopPromotionRepository(cfg *ShopPromotionRConfig) ShopPromotionRepository {
	return &shopPromotionRepositoryImpl{
		db:                cfg.DB,
		productRepository: cfg.ProductRepository,
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
		Model(&model.ShopPromotion{}).
		Where("shop_promotions.shop_id = ?", shopId)

	if request.Name != "" {
		query = query.Where("shop_promotions.name ILIKE ?", fmt.Sprintf("%%%s%%", request.Name))
	}

	switch request.Status {
	case constant.VoucherPromotionStatusOngoing:
		query = query.Where("shop_promotions.start_period <= ? AND ? <= shop_promotions.end_period", now, now)
	case constant.VoucherPromotionStatusUpcoming:
		query = query.Where("? < shop_promotions.start_period", now)
	case constant.VoucherPromotionStatusExpired:
		query = query.Where("shop_promotions.end_period < ?", now)
	}

	query = query.Select("shop_promotions.*, "+
		"CASE WHEN start_period <= ? AND end_period >= ? THEN ? "+
		"WHEN start_period > ? THEN ? "+
		"ELSE ? "+
		"END as status", now, now, constant.VoucherPromotionStatusOngoing, now, constant.VoucherPromotionStatusUpcoming, constant.VoucherPromotionStatusExpired)

	query = query.Session(&gorm.Session{})

	err := query.Model(&model.ShopPromotion{}).Distinct("shop_promotions.id").Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages = int(math.Ceil(float64(totalRows) / float64(request.Limit)))

	err = query.
		Order("shop_promotions.created_at desc").
		Limit(request.Limit).
		Offset(request.Offset()).Find(&promotions).Error
	if err != nil {
		return nil, 0, 0, err
	}

	for _, promotions := range promotions {
		products, err := r.productRepository.GetWithPromotions(shopId, promotions.ID)
		if err != nil {
			return nil, 0, 0, err
		}
		promotions.Product = products
	}

	return promotions, totalRows, totalPages, nil
}

func (r *shopPromotionRepositoryImpl) GetSellerPromotionById(shopId int, promotionId int) (*dto.SellerPromotion, error) {
	var promotion *dto.SellerPromotion

	now := time.Now()
	query := r.db.
		Model(&model.ShopPromotion{}).
		Where("shop_promotions.shop_id = ? AND shop_promotions.id = ?", shopId, promotionId)

	query = query.Select("shop_promotions.*, "+
		"CASE WHEN start_period <= ? AND end_period >= ? THEN ? "+
		"WHEN start_period > ? THEN ? "+
		"ELSE ? "+
		"END as status", now, now, constant.VoucherPromotionStatusOngoing, now, constant.VoucherPromotionStatusUpcoming, constant.VoucherPromotionStatusExpired)

	err := query.First(&promotion).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrPromotionNotFound
		}
		return nil, err
	}

	products, err := r.productRepository.GetWithPromotions(shopId, promotion.ID)
	if err != nil {
		return nil, err
	}
	promotion.Product = products

	return promotion, nil
}

func (r *shopPromotionRepositoryImpl) Update(shopPromotion *model.ShopPromotion, productPromotions []*productModel.ProductPromotion) error {
	tx := r.db.Begin()
	defer tx.Commit()

	res := tx.Clauses(clause.Returning{}).Updates(shopPromotion)
	if err := res.Error; err != nil {
		tx.Rollback()
		return err
	}

	if res.RowsAffected < 1 {
		return errs.ErrPromotionNotFound
	}

	for _, productPromotion := range productPromotions {
		res := tx.Save(productPromotion)
		if err := res.Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func (r *shopPromotionRepositoryImpl) Create(shopID int, request *dto.CreateShopPromotionRequest) (*dto.CreateShopPromotionResponse, error) {
	if err := request.ValidateDateRange(); err != nil {
		return nil, err
	}

	tx := r.db.Begin()
	defer tx.Commit()

	shopPromotion := request.GenerateShopPromotion()
	shopPromotion.ShopId = shopID

	err := tx.Create(shopPromotion).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	productPromotions := []*productModel.ProductPromotion{}

	for _, pp := range request.ProductPromotions {
		productPromotions = append(productPromotions, &productModel.ProductPromotion{
			Type:          pp.Type,
			Amount:        pp.Amount,
			Stock:         pp.Stock,
			IsActive:      *pp.IsActive,
			PurchaseLimit: pp.PurchaseLimit,
			SkuId:         pp.SkuId,
			PromotionId:   shopPromotion.ID,
		})
	}

	err = tx.Create(productPromotions).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	response := &dto.CreateShopPromotionResponse{
		ShopPromotion:     *shopPromotion,
		ProductPromotions: productPromotions,
	}

	return response, nil
}

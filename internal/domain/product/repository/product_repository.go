package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	model "kedai/backend/be-kedai/internal/domain/product/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetByID(ID int) (*model.Product, error)
	GetByCode(Code string) (*model.Product, error)
	GetRecommendation(categoryId int) ([]*model.Product, error)
}

type productRepositoryImpl struct {
	db *gorm.DB
}

type ProductRConfig struct {
	DB *gorm.DB
}

func NewProductRepository(cfg *ProductRConfig) ProductRepository {
	return &productRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *productRepositoryImpl) GetByID(ID int) (*model.Product, error) {
	var product model.Product

	err := r.db.Where("id = ?", ID).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrProductDoesNotExist
		}

		return nil, err
	}

	return &product, err
}

func (r *productRepositoryImpl) GetByCode(Code string) (*model.Product, error) {
	var product model.Product

	err := r.db.Where("code = ?", Code).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrProductDoesNotExist
		}

		return nil, err
	}

	return &product, nil
}

func (r *productRepositoryImpl) GetRecommendation(categoryId int) ([]*model.Product, error) {
	var (
		products []*model.Product
		limit    = 5
		isActive = true
	)

	db := r.db.Select(`products.*, min(s.price) as min_price, max(s.price) as max_price,
	max(case when pp.type = 'nominal' then ROUND(cast(pp.amount / s.price * 100 as numeric), 2) else pp.amount end) as promotion_percent,
	(select url from product_medias pm where pm.product_id = products.id limit 1) as image_url`).
		Joins("join skus s on s.product_id = products.id").
		Joins("join shops sh on sh.id = products.shop_id").
		Joins("left join product_promotions pp on pp.sku_id = s.id and (select count(id) from shop_promotions sp where pp.promotion_id = sp.id and now() between sp.start_period and sp.end_period) > 0").
		Group("products.id")

	err := db.Where("products.category_id = ? and products.is_active = ?", categoryId, isActive).Limit(limit).Order("products.sold desc, products.rating desc").Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

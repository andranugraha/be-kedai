package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"math"

	"gorm.io/gorm"
)

type ShopRepository interface {
	FindShopById(id int) (*model.Shop, error)
	FindShopByUserId(userId int) (*model.Shop, error)
	FindShopBySlug(slug string) (*model.Shop, error)
	FindShopByKeyword(req dto.FindShopRequest) ([]*dto.FindShopResponse, int64, int, error)
}

type shopRepositoryImpl struct {
	db *gorm.DB
}

type ShopRConfig struct {
	DB *gorm.DB
}

func NewShopRepository(cfg *ShopRConfig) ShopRepository {
	return &shopRepositoryImpl{
		db: cfg.DB,
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

	return &shop, nil
}

func (r *shopRepositoryImpl) FindShopByKeyword(req dto.FindShopRequest) ([]*dto.FindShopResponse, int64, int, error) {
	var (
		shopList  []*dto.FindShopResponse
		totalRows int64
		totalPage int
		isActive  = true
	)

	db := r.db.Select(`shops.*, count(p.id) as product_count`).
		Joins("left join products p on shops.id = p.shop_id and p.is_active = ?", isActive).
		Group("shops.id").Where("shops.name ILIKE ?", "%"+req.Keyword+"%")

	db.Model(&model.Shop{}).Count(&totalRows)
	totalPage = int(math.Ceil(float64(totalRows) / float64(req.Limit)))

	err := db.Order("rating desc").Limit(req.Limit).Offset(req.Offset()).Find(&shopList).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return shopList, totalRows, totalPage, nil
}

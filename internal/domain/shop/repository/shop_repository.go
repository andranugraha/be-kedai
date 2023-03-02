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
	FindShopByKeyword(keyword string) ([]*model.Shop, error)
	FindMostRatedShopByKeyword(keyword string) (*model.Shop, error)
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

func (r *shopRepositoryImpl) FindShopByKeyword(keyword string) ([]*model.Shop, error) {
	var shopList []*model.Shop

	err := r.db.Where("name ILIKE ?", "%"+keyword+"%").Order("rating desc").Find(&shopList).Error
	if err != nil {
		return nil, err
	}

	return shopList, nil	
}

func (r *shopRepositoryImpl) FindMostRatedShopByKeyword(keyword string) (*model.Shop, error) {
	var topShop *model.Shop

	err := r.db.Where("name ILIKE ?", "%"+keyword+"%").Order("rating desc").First(&topShop).Error
	if err != nil {
		return nil, err
	}

	return topShop, nil
}

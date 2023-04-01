package repository

import (
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"math"

	"gorm.io/gorm"
)

type ShopCategoryRepository interface {
	GetByShopID(shopID int, req dto.GetSellerCategoriesRequest) ([]*dto.ShopCategory, int64, int, error)
}

type shopCategoryRepositoryImpl struct {
	db *gorm.DB
}

type ShopCategoryRConfig struct {
	DB *gorm.DB
}

func NewShopCategoryRepository(cfg *ShopCategoryRConfig) ShopCategoryRepository {
	return &shopCategoryRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *shopCategoryRepositoryImpl) GetByShopID(shopID int, req dto.GetSellerCategoriesRequest) (res []*dto.ShopCategory, totalRows int64, totalPages int, err error) {
	db := r.db.Model(&model.ShopCategory{}).Where("shop_id = ?", shopID).
		Joins("left join shop_category_products scp on scp.shop_category_id = shop_categories.id").
		Select("shop_categories.id, shop_categories.name, shop_categories.shop_id, shop_categories.is_active, count(scp.id) as total_product").
		Group("shop_categories.id")

	if req.Status != "" {
		db = db.Where("is_active = ?", req.Status == "enabled")
	}

	err = db.Count(&totalRows).Error
	if err != nil {
		return
	}

	totalPages = int(math.Ceil(float64(totalRows) / float64(req.Limit)))

	err = db.Offset(req.Offset()).Limit(req.Limit).Find(&res).Error
	if err != nil {
		return
	}

	return
}

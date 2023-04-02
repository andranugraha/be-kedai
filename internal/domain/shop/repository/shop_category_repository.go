package repository

import (
	"errors"
	"fmt"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	productRepo "kedai/backend/be-kedai/internal/domain/product/repository"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"math"

	"gorm.io/gorm"
)

type ShopCategoryRepository interface {
	GetByShopID(shopID int, req dto.GetSellerCategoriesRequest) ([]*dto.ShopCategory, int64, int, error)
	GetByIDAndShopID(id, shopID int) (*dto.ShopCategory, error)
	GetCategoryByIDAndShopID(id, shopID int) (*model.ShopCategory, error)
	Create(shopCategory *model.ShopCategory) error
	Update(shopCategory *model.ShopCategory) error
	Delete(id, shopId int) error
}

type shopCategoryRepositoryImpl struct {
	db          *gorm.DB
	productRepo productRepo.ProductRepository
}

type ShopCategoryRConfig struct {
	DB          *gorm.DB
	ProductRepo productRepo.ProductRepository
}

func NewShopCategoryRepository(cfg *ShopCategoryRConfig) ShopCategoryRepository {
	return &shopCategoryRepositoryImpl{
		db:          cfg.DB,
		productRepo: cfg.ProductRepo,
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

	if req.Search != "" {
		db = db.Where("name ilike ?", fmt.Sprintf("%%%s%%", req.Search))
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

func (r *shopCategoryRepositoryImpl) GetByIDAndShopID(id, shopID int) (res *dto.ShopCategory, err error) {
	err = r.db.Model(&model.ShopCategory{}).Where("shop_categories.id = ?", id).Where("shop_categories.shop_id = ?", shopID).
		Joins("left join shop_category_products scp on scp.shop_category_id = shop_categories.id").
		Select("shop_categories.id, shop_categories.name, shop_categories.shop_id, shop_categories.is_active, count(scp.id) as total_product").
		Group("shop_categories.id").
		First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = commonErr.ErrCategoryNotFound
			return
		}
	}

	err = r.db.Table("shop_category_products").Where("shop_category_id = ?", id).
		Select("products.id, products.code, products.name, (select url from product_medias pm where pm.product_id = products.id and deleted_at is null limit 1) as image_url, min(skus.price) as min_price, max(skus.price) as max_price, sum(skus.stock) as stock").
		Joins("join products on products.id = shop_category_products.product_id").
		Joins("join skus on skus.product_id = products.id").
		Group("products.id").
		Find(&res.Products).Error

	return
}

func (r *shopCategoryRepositoryImpl) GetCategoryByIDAndShopID(id, shopID int) (res *model.ShopCategory, err error) {
	err = r.db.Model(&model.ShopCategory{}).Where("shop_categories.id = ?", id).Where("shop_categories.shop_id = ?", shopID).
		First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = commonErr.ErrCategoryNotFound
		}

		return
	}

	return
}

func (r *shopCategoryRepositoryImpl) Create(shopCategory *model.ShopCategory) error {
	for _, categoryProduct := range shopCategory.Products {
		product, err := r.productRepo.GetByID(categoryProduct.ProductId)
		if err != nil {
			return err
		}

		if product.ShopID != shopCategory.ShopId {
			return commonErr.ErrProductDoesNotExist
		}
	}

	tx := r.db.Begin()
	defer tx.Commit()

	if err := tx.Create(&shopCategory).Error; err != nil {
		tx.Rollback()
		if commonErr.IsDuplicateKeyError(err) {
			return commonErr.ErrCategoryAlreadyExist
		}

		if commonErr.IsForeignKeyError(err) {
			return commonErr.ErrProductDoesNotExist
		}

		return err
	}

	return nil
}

func (r *shopCategoryRepositoryImpl) Update(shopCategory *model.ShopCategory) error {
	tx := r.db.Begin()
	defer tx.Commit()

	if err := tx.Where("shop_category_id = ?", shopCategory.ID).Unscoped().Delete(&model.ShopCategoryProduct{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, categoryProduct := range shopCategory.Products {
		product, err := r.productRepo.GetByID(categoryProduct.ProductId)
		if err != nil {
			return err
		}

		if product.ShopID != shopCategory.ShopId {
			return commonErr.ErrProductDoesNotExist
		}
	}

	if err := tx.Save(&shopCategory).Error; err != nil {
		tx.Rollback()
		if commonErr.IsDuplicateKeyError(err) {
			return commonErr.ErrCategoryAlreadyExist
		}

		if commonErr.IsForeignKeyError(err) {
			return commonErr.ErrProductDoesNotExist
		}

		return err
	}

	return nil
}

func (r *shopCategoryRepositoryImpl) Delete(id, shopId int) error {
	tx := r.db.Begin()
	defer tx.Commit()

	res := tx.Where("id = ?", id).Where("shop_id = ?", shopId).Delete(&model.ShopCategory{})
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	if res.RowsAffected == 0 {
		return commonErr.ErrCategoryNotFound
	}

	if err := tx.Where("shop_category_id = ?", id).Unscoped().Delete(&model.ShopCategoryProduct{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

package repository

import (
	"errors"
	"fmt"
	"kedai/backend/be-kedai/internal/common/constant"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	model "kedai/backend/be-kedai/internal/domain/product/model"
	"math"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetByID(ID int) (*model.Product, error)
	GetByCode(code string) (*dto.ProductDetail, error)
	GetByShopID(shopID int, request *dto.ShopProductFilterRequest) ([]*dto.ProductDetail, int64, int, error)
	GetRecommendationByCategory(productId int, categoryId int) ([]*dto.ProductResponse, error)
	ProductSearchFiltering(req dto.ProductSearchFilterRequest, shopId int) ([]*dto.ProductResponse, int64, int, error)
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

func (r *productRepositoryImpl) GetByCode(code string) (*dto.ProductDetail, error) {
	var product dto.ProductDetail

	query := r.db.Select(`products.*, min(s.price) as min_price, max(s.price) as max_price, sum(s.stock) as total_stock,
	max(case when pp.type = 'nominal' then pp.amount / s.price else pp.amount end) as promotion_percent
	`).
		Joins("join skus s on s.product_id = products.id").
		Joins("left join product_promotions pp on pp.sku_id = s.id and (select count(id) from shop_promotions sp where pp.promotion_id = sp.id and now() between sp.start_period and sp.end_period) > 0").
		Group("products.id")

	err := query.Where("code = ?", code).Preload("SKU").Preload("VariantGroup.Variant").Preload("Media").Preload("Bulk").Preload("Shop.Address.Subdistrict").First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrProductDoesNotExist
		}

		return nil, err
	}

	if len(product.VariantGroup) > 0 {
		product.SKU = nil
	}

	return &product, nil
}

func (r *productRepositoryImpl) GetRecommendationByCategory(productId int, categoryId int) ([]*dto.ProductResponse, error) {
	var (
		products []*dto.ProductResponse
		limit    = 5
		isActive = true
	)

	db := r.db.Select(`products.*, min(s.price) as min_price, max(s.price) as max_price,
	max(case when pp.type = 'nominal' then pp.amount / s.price else pp.amount end) as promotion_percent,
	(select url from product_medias pm where pm.product_id = products.id limit 1) as image_url`).
		Joins("join skus s on s.product_id = products.id").
		Joins("left join product_promotions pp on pp.sku_id = s.id and (select count(id) from shop_promotions sp where pp.promotion_id = sp.id and now() between sp.start_period and sp.end_period) > 0").
		Group("products.id")

	err := db.Where("products.category_id = ? and products.is_active = ? and products.id != ?", categoryId, isActive, productId).Limit(limit).Order("products.sold desc, products.rating desc").Find(&products).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCategoryDoesNotExist
		}
		return nil, err
	}

	return products, nil
}

func (r *productRepositoryImpl) GetByShopID(shopID int, request *dto.ShopProductFilterRequest) ([]*dto.ProductDetail, int64, int, error) {
	var (
		products   []*dto.ProductDetail
		totalRows  int64
		totalPages int
		active     = true
	)

	query := r.db.Select(`products.*, min(s.price) as min_price, max(s.price) as max_price, sum(s.stock) as total_stock,
	max(case when pp.type = 'nominal' then pp.amount / s.price else pp.amount end) as promotion_percent,
	(select url from product_medias pm where pm.product_id = products.id limit 1) as image_url
	`).
		Joins("join skus s on s.product_id = products.id").
		Joins("left join product_promotions pp on pp.sku_id = s.id and (select count(id) from shop_promotions sp where pp.promotion_id = sp.id and now() between sp.start_period and sp.end_period) > 0")

	if request.ShopProductCategoryID > 0 {
		query.Joins("left join shop_category_products scp on products.id = scp.product_id").Where("scp.id = ?", request.ShopProductCategoryID)
	}

	query = query.Where("is_active = ?", active)
	query = query.Group("products.id")

	var priceSort string
	if request.PriceSort == constant.SortByPriceHigh {
		priceSort = "desc"
	} else {
		priceSort = "asc"
	}

	query.Order(fmt.Sprintf("min(s.price) %s", priceSort))

	switch request.Sort {
	case constant.SortByRecommended:
		query = query.Order("products.rating desc, products.sold desc")
	case constant.SortByLatest:
		query = query.Order("products.created_at desc")
	case constant.SortByTopSales:
		query = query.Order("products.sold desc")
	default:
		query = query.Order("products.created_at desc")
	}

	err := query.Model(&model.Product{}).Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages = int(math.Ceil(float64(totalRows) / float64(request.Limit)))

	query = query.Where("products.shop_id = ?", shopID)
	if request.ExceptionID > 0 {
		query = query.Where("products.id != ?", request.ExceptionID)
	}
	err = query.Limit(request.Limit).Offset(request.Offset()).Find(&products).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return products, totalRows, totalPages, nil
}

func (r *productRepositoryImpl) ProductSearchFiltering(req dto.ProductSearchFilterRequest, shopId int) ([]*dto.ProductResponse, int64, int, error) {
	var (
		productList []*dto.ProductResponse
		totalRows   int64
		totalPages  int
		active      = true
	)

	db := r.db.Select(`products.*, min(s.price) as min_price, max(s.price) as max_price, 
	concat(c.name, ', ', p.name) as address, 
	max(case when pp.type = 'nominal' then pp.amount / s.price else pp.amount end) as promotion_percent, 
	(select url from product_medias pm where products.id = pm.product_id limit 1) as image_url`).
		Joins("join skus s ON s.product_id = products.id").
		Joins("join shops sh ON sh.id = products.shop_id").
		Joins("join user_addresses ua ON ua.id = sh.address_id").
		Joins("join cities c ON c.id = ua.city_id").
		Joins("join provinces p ON p.id = c.province_id").
		Joins("left join product_promotions pp ON pp.sku_id = s.id and (select count(id) from shop_promotions sp where pp.promotion_id = sp.id and now() between sp.start_period and sp.end_period) > 0").
		Group("products.id, c.name, p.name")

	db = db.Where("products.is_active = ?", active).Where("products.name ILIKE ?", "%"+req.Keyword+"%")

	if req.CategoryId > 0 {
		db = db.Where("products.category_id = ?", req.CategoryId)
	}

	if req.MinRating > 0 {
		db = db.Where("products.rating >= ?", req.MinRating)
	}

	if len(req.CityIds) > 0 {
		db = db.Where("ua.city_id in (?)", req.CityIds)
	}

	if req.MinPrice > 0 || req.MaxPrice > 0 {
		if req.MinPrice > 0 {
			db = db.Where("s.id = (select id from skus where product_id = products.id and skus.price >= ? limit 1)", req.MinPrice)
		}

		if req.MaxPrice > 0 {
			db = db.Where("s.id = (select id from skus where product_id = products.id and skus.price <= ? limit 1)", req.MaxPrice)
		}
	}

	if shopId != 0 {
		db = db.Where("products.shop_id = ?", shopId)
	}

	switch req.Sort {
	case constant.SortByRecommended:
		db = db.Order("products.rating desc, products.sold desc")
	case constant.SortByLatest:
		db = db.Order("products.created_at desc")
	case constant.SortByTopSales:
		db = db.Order("products.sold desc")
	case constant.SortByPriceLow:
		db = db.Where("s.id = (select id from skus where product_id = products.id order by price asc limit 1)").Group("s.id").Order("s.price asc")
	case constant.SortByPriceHigh:
		db = db.Where("s.id = (select id from skus where product_id = products.id order by price asc limit 1)").Group("s.id").Order("s.price desc")
	default:
		db = db.Order("products.created_at desc")
	}

	db.Model(&model.Product{}).Count(&totalRows)
	totalPages = int(math.Ceil(float64(totalRows) / float64(req.Limit)))

	err := db.Model(&model.Product{}).Limit(req.Limit).Offset(req.Offset()).Find(&productList).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return productList, totalRows, totalPages, nil
}

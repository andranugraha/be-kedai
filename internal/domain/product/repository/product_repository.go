package repository

import (
	"errors"
	"fmt"
	"kedai/backend/be-kedai/internal/common/constant"
	errs "kedai/backend/be-kedai/internal/common/error"
	orderDto "kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	model "kedai/backend/be-kedai/internal/domain/product/model"
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	"math"
	"time"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetByID(ID int) (*model.Product, error)
	GetActiveByID(ID int) (*model.Product, error)
	GetByCode(code string) (*dto.ProductDetail, error)
	GetByShopID(shopID int, request *dto.ShopProductFilterRequest) ([]*dto.ProductDetail, int64, int, error)
	GetRecommendationByCategory(productId int, categoryId int) ([]*dto.ProductResponse, error)
	ProductSearchFiltering(req dto.ProductSearchFilterRequest, shopId int, categoryIDs []int) ([]*dto.ProductResponse, int64, int, error)
	GetBySellerID(shopID int, request *dto.SellerProductFilterRequest) ([]*dto.SellerProduct, int64, int, error)
	GetWithPromotions(shopID int, promotionID int) ([]*dto.SellerProductPromotionResponse, error)
	SearchAutocomplete(req dto.ProductSearchAutocomplete) ([]*dto.ProductResponse, error)
	GetSellerProductByCode(shopID int, productCode string) (*model.Product, error)
	AddViewCount(productID int) error
	UpdateActivation(shopID int, code string, isActive bool) error
	Create(shopID int, request *dto.CreateProductRequest, courierServices []*shopModel.CourierService) (*model.Product, error)
	GetRecommended(req *dto.GetRecommendedProductRequest) ([]*dto.ProductResponse, int64, int, error)
	Update(shopID int, code string, payload *dto.CreateProductRequest, courierServices []*shopModel.CourierService) (*model.Product, error)
	AddSoldCount(tx *gorm.DB, items []*orderDto.TransactionItem) error
}

type productRepositoryImpl struct {
	db                       *gorm.DB
	variantGroupRepo         VariantGroupRepository
	skuRepository            SkuRepository
	productVariantRepository ProductVariantRepository
	discussionRepository     DiscussionRepository
	productMediaRepository   ProductMediaRepository
}

type ProductRConfig struct {
	DB                       *gorm.DB
	VariantGroupRepo         VariantGroupRepository
	SkuRepository            SkuRepository
	ProductVariantRepository ProductVariantRepository
	DiscussionRepository     DiscussionRepository
	ProductMediaRepository   ProductMediaRepository
}

func NewProductRepository(cfg *ProductRConfig) ProductRepository {
	return &productRepositoryImpl{
		db:                       cfg.DB,
		variantGroupRepo:         cfg.VariantGroupRepo,
		skuRepository:            cfg.SkuRepository,
		productVariantRepository: cfg.ProductVariantRepository,
		discussionRepository:     cfg.DiscussionRepository,
		productMediaRepository:   cfg.ProductMediaRepository,
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

func (r *productRepositoryImpl) GetActiveByID(ID int) (*model.Product, error) {
	var (
		product model.Product
		active  = true
	)
	err := r.db.Where("is_active = ?", active).First(&product, ID).Error
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
		Joins("join skus s on s.product_id = products.id AND s.deleted_at IS NULL").
		Joins("left join product_promotions pp on pp.sku_id = s.id and (select count(id) from shop_promotions sp where pp.promotion_id = sp.id and now() between sp.start_period and sp.end_period) > 0").
		Group("products.id")

	err := query.Where("code = ?", code).Preload("SKU", func(query *gorm.DB) *gorm.DB {
		return query.Select("skus.id, skus.price, skus.stock, skus.product_id")
	}).Preload("VariantGroup.Variant").Preload("Media").Preload("Bulk").Preload("Shop.Address.Subdistrict").First(&product).Error
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

func (r *productRepositoryImpl) AddSoldCount(tx *gorm.DB, items []*orderDto.TransactionItem) error {
	for _, ti := range items {
		err := tx.Model(&model.Product{}).Where("id = ?", ti.Sku.ProductId).Update("sold", gorm.Expr("sold + ?", ti.Quantity)).Error
		if err != nil {
			continue
		}
	}

	return nil
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
		Joins("join skus s on s.product_id = products.id AND s.deleted_at IS NULL").
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
		Joins("join skus s on s.product_id = products.id AND s.deleted_at IS NULL").
		Joins("left join product_promotions pp on pp.sku_id = s.id and (select count(id) from shop_promotions sp where pp.promotion_id = sp.id and now() between sp.start_period and sp.end_period) > 0")

	if request.ShopProductCategoryID > 0 {
		query = query.Joins("left join shop_category_products scp on products.id = scp.product_id").
			Joins("join shop_categories sc on scp.shop_category_id = sc.id").
			Where("sc.id = ?", request.ShopProductCategoryID).Where("sc.is_active = true")
	}

	query = query.Where("products.is_active = ?", active)
	query = query.Group("products.id")

	query = query.Where("products.shop_id = ?", shopID)
	if request.ExceptionID > 0 {
		query = query.Where("products.id != ?", request.ExceptionID)
	}

	query = query.Session(&gorm.Session{})

	err := query.Model(&model.Product{}).Distinct("products.id").Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages = int(math.Ceil(float64(totalRows) / float64(request.Limit)))

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

	var priceSort string
	if request.PriceSort == constant.SortByPriceHigh {
		priceSort = "desc"
	} else {
		priceSort = "asc"
	}

	query.Order(fmt.Sprintf("min(s.price) %s", priceSort))

	err = query.Limit(request.Limit).Offset(request.Offset()).Find(&products).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return products, totalRows, totalPages, nil
}

func (r *productRepositoryImpl) ProductSearchFiltering(req dto.ProductSearchFilterRequest, shopId int, categoryIDs []int) ([]*dto.ProductResponse, int64, int, error) {
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
		Joins("join skus s ON s.product_id = products.id AND s.deleted_at IS NULL").
		Joins("join shops sh ON sh.id = products.shop_id").
		Joins("join user_addresses ua ON ua.id = sh.address_id").
		Joins("join cities c ON c.id = ua.city_id").
		Joins("join provinces p ON p.id = c.province_id").
		Joins("left join product_promotions pp ON pp.sku_id = s.id and (select count(id) from shop_promotions sp where pp.promotion_id = sp.id and now() between sp.start_period and sp.end_period) > 0").
		Group("products.id, c.name, p.name")

	db = db.Where("products.is_active = ?", active).Where("products.name ILIKE ?", "%"+req.Keyword+"%")

	if req.CategoryId > 0 {
		db = db.Where("products.category_id IN ?", categoryIDs)
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

	countQuery := db.Session(&gorm.Session{})
	countQuery.Model(&model.Product{}).Distinct("products.id").Count(&totalRows)
	totalPages = int(math.Ceil(float64(totalRows) / float64(req.Limit)))

	err := db.Model(&model.Product{}).Limit(req.Limit).Offset(req.Offset()).Find(&productList).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return productList, totalRows, totalPages, nil
}

func (r *productRepositoryImpl) GetBySellerID(shopID int, request *dto.SellerProductFilterRequest) ([]*dto.SellerProduct, int64, int, error) {
	var (
		products   []*dto.SellerProduct
		totalRows  int64
		totalPages int
	)

	query := r.db.
		Select(`products.*,
		MIN(skus.price) AS min_price,
		SUM(skus.stock) AS total_stock,
		(SELECT url FROM product_medias pm WHERE pm.product_id = products.id LIMIT 1) AS image_url
	`).
		Joins("JOIN skus ON skus.product_id = products.id AND skus.deleted_at IS NULL").
		Group("products.id")

	query = query.Where("products.shop_id = ?", shopID)
	if request.Name != "" {
		query = query.Where("products.name ILIKE ?", fmt.Sprintf("%%%s%%", request.Name))
	}
	if request.Sku != "" {
		query = query.Where("skus.sku ILIKE ?", fmt.Sprintf("%%%s%%", request.Sku))
	}
	if request.Sales > 0 {
		query = query.Where("products.sold >= ?", request.Sales)
	}
	if request.Stock > 0 {
		query = query.Having("SUM(skus.stock) >= ?", request.Stock)
	}

	switch request.Status {
	case constant.ProductStatusLive:
		query = query.Where("products.is_active")
	case constant.ProductStatusDelisted:
		query = query.Not("products.is_active")
	case constant.ProductStatusSoldOut:
		query = query.Where("skus.stock = 0")
	}

	if request.IsPromoted != nil && !*request.IsPromoted {
		now := time.Now()
		startPeriod := request.StartPeriod
		endPeriod := request.EndPeriod

		if startPeriod.IsZero() {
			startPeriod = now
		}

		if endPeriod.IsZero() {
			endPeriod = now
		}

		query = query.Where(`products.id NOT IN
			(SELECT skus.product_id
			FROM skus
			JOIN product_promotions ON product_promotions.sku_id = skus.id
			JOIN shop_promotions ON shop_promotions.id = product_promotions.promotion_id
			WHERE shop_promotions.start_period <= ? AND shop_promotions.end_period >= ?)`,
			startPeriod, endPeriod)
	}

	query = query.Session(&gorm.Session{})

	err := query.Model(&model.Product{}).Distinct("products.id").Count(&totalRows).Error
	if err != nil {
		return nil, 0, 0, err
	}
	totalPages = int(math.Ceil(float64(totalRows) / float64(request.Limit)))

	query = query.Order("products.is_active DESC")
	switch request.Sort {
	case constant.SortByLowSales:
		query = query.Order("products.sold ASC")
	case constant.SortByTopSales:
		query = query.Order("products.sold DESC")
	case constant.SortByPriceLow:
		query = query.Order("min_price ASC")
	case constant.SortByPriceHigh:
		query = query.Order("min_price DESC")
	case constant.SortByStockLow:
		query = query.Order("total_stock ASC")
	case constant.SortByStockHigh:
		query = query.Order("total_stock DESC")
	default:
		query = query.Order("products.created_at DESC")
	}

	err = query.Preload("Bulk").Preload("SKUs.Variants").Limit(request.Limit).Offset((request.Page - 1) * request.Limit).Find(&products).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return products, totalRows, totalPages, nil
}

func (r *productRepositoryImpl) GetWithPromotions(shopID int, promotionID int) ([]*dto.SellerProductPromotionResponse, error) {
	var (
		products []*dto.SellerProductPromotion
	)

	query := r.db.
		Select(`products.id,
			products.name,
			products.code,
			(SELECT url FROM product_medias pm WHERE pm.product_id = products.id LIMIT 1) AS image_url
		`).
		Joins(`
		JOIN skus ON skus.product_id = products.id
		JOIN product_promotions ON product_promotions.sku_id = skus.id
	`).
		Group("products.id").
		Where("products.shop_id = ? AND product_promotions.promotion_id = ?", shopID, promotionID)

	err := query.Preload("SKUs.Variants").Preload("SKUs.Promotion").Find(&products).Error
	if err != nil {
		return nil, err
	}

	convertedProducts := dto.ConvertSellerProductPromotions(products)

	return convertedProducts, nil
}

func (r *productRepositoryImpl) SearchAutocomplete(req dto.ProductSearchAutocomplete) ([]*dto.ProductResponse, error) {
	var (
		products []*dto.ProductResponse
		active   = true
	)

	db := r.db.Select(`products.*, (select url from product_medias pm where products.id = pm.product_id limit 1) as image_url`)

	db = db.Where("products.is_active = ?", active).Where("products.name ILIKE ?", "%"+req.Keyword+"%").Order("products.rating desc")

	err := db.Limit(req.Limit).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepositoryImpl) GetSellerProductByCode(shopID int, productCode string) (*model.Product, error) {
	var product model.Product

	err := r.db.
		Where("shop_id = ?", shopID).Where("code = ?", productCode).
		Preload("Bulk").
		Preload("Media").
		Preload("VariantGroup.Variant").
		Preload("SKUs.Variants").
		First(&product).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrProductDoesNotExist
		}

		return nil, err
	}
	return &product, nil
}

func (r *productRepositoryImpl) AddViewCount(productID int) error {
	res := r.db.
		Model(&model.Product{}).
		Where("id = ?", productID).
		Where("is_active = ?", true).
		Update("view", gorm.Expr("view + ?", 1))

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errs.ErrProductDoesNotExist
	}

	return nil
}

func (r *productRepositoryImpl) UpdateActivation(shopID int, code string, isActive bool) error {
	res := r.db.Model(&model.Product{}).Where("code = ?", code).Where("shop_id = ?", shopID).Update("is_active", isActive)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errs.ErrProductDoesNotExist
	}

	return nil
}

func (r *productRepositoryImpl) Create(shopID int, request *dto.CreateProductRequest, courierServices []*shopModel.CourierService) (*model.Product, error) {
	tx := r.db.Begin()
	defer tx.Commit()

	product := request.GenerateProduct()
	product.ShopID = shopID
	product.CourierService = courierServices

	err := tx.Create(product).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	variantGroups := request.GenerateVariantGroups()
	for _, vg := range variantGroups {
		vg.ProductID = product.ID
	}

	if variantGroups != nil {
		err = r.variantGroupRepo.Create(tx, variantGroups)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	skus := request.GenerateSKU(variantGroups)

	for _, s := range skus {
		s.ProductId = product.ID
	}

	err = r.skuRepository.Create(tx, skus)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	productVariants := []*model.ProductVariant{}
	for _, sku := range skus {
		for _, v := range sku.Variants {
			productVariants = append(productVariants, &model.ProductVariant{
				SkuId:     sku.ID,
				VariantId: v.ID,
			})
		}
	}

	if len(productVariants) > 0 {
		err = r.productVariantRepository.Create(tx, productVariants)
		if err != nil {
			return nil, err
		}
	}

	product.VariantGroup = variantGroups
	product.SKUs = skus

	return product, nil
}

func (r *productRepositoryImpl) GetRecommended(req *dto.GetRecommendedProductRequest) (recommendedProducts []*dto.ProductResponse, totalRows int64, totalPages int, err error) {

	var (
		isActive = true
	)

	db := r.db.Select(`products.*, min(s.price) as min_price, max(s.price) as max_price,
		concat(c.name, ', ', p.name) as address, 
		max(case when pp.type = 'nominal' then pp.amount / s.price else pp.amount end) as promotion_percent,
		(select url from product_medias pm where pm.product_id = products.id limit 1) as image_url,
		(select id from skus s where products.id = s.product_id limit 1) as default_sku_id`).
		Joins("join skus s on s.product_id = products.id AND s.deleted_at IS NULL").
		Joins("join shops sh ON sh.id = products.shop_id").
		Joins("join user_addresses ua ON ua.id = sh.address_id").
		Joins("join cities c ON c.id = ua.city_id").
		Joins("join provinces p ON p.id = c.province_id").
		Joins("left join product_promotions pp on pp.sku_id = s.id and (select count(id) from shop_promotions sp where pp.promotion_id = sp.id and now() between sp.start_period and sp.end_period) > 0").
		Where("products.is_active = ?", isActive).
		Group("products.id,c.name,p.name")

	errCount := db.Model(&model.Product{}).Count(&totalRows).Error
	if errCount != nil {
		return nil, 0, 0, errCount
	}

	totalPages = int(math.Ceil(float64(totalRows) / float64(req.Limit)))

	err = db.
		Limit(req.Limit).
		Offset(req.Offset()).
		Order("products.sold desc, products.rating desc").
		Find(&recommendedProducts).Error

	if err != nil {
		return nil, 0, 0, err
	}

	return recommendedProducts, totalRows, totalPages, nil
}

func (r *productRepositoryImpl) Update(shopID int, code string, payload *dto.CreateProductRequest, courierServices []*shopModel.CourierService) (*model.Product, error) {

	tx := r.db.Begin()
	defer tx.Commit()

	product, err := r.GetSellerProductByCode(shopID, code)
	if err != nil {
		return nil, err
	}

	updatedProduct := payload.GenerateProduct()
	updatedProduct.Code = product.Code
	updatedProduct.ID = product.ID
	updatedProduct.ShopID = shopID
	updatedProduct.CourierService = courierServices
	updatedProduct.View = product.View
	updatedProduct.CreatedAt = product.CreatedAt
	updatedProduct.Rating = product.Rating
	updatedProduct.Sold = product.Sold

	if payload.BulkPrice == nil {
		errBulk := tx.Where("product_id=?", product.ID).Delete(&model.ProductBulkPrice{}).Error
		if errBulk != nil {
			tx.Rollback()
			return nil, errBulk
		}
	} else {
		if product.Bulk != nil {
			updatedProduct.Bulk.ID = product.Bulk.ID
		}
	}

	var media []*model.ProductMedia

	errDelete := r.productMediaRepository.Delete(tx, product.ID)
	if errDelete != nil {
		tx.Rollback()
		return nil, errDelete
	}

	for _, value := range product.Media {

		for _, reqMedia := range updatedProduct.Media {
			if value.Url == reqMedia.Url {
				media = append(media, value)
				break
			}
		}
	}

	for _, reqMedia := range updatedProduct.Media {
		found := false
		for _, value := range media {
			if value.Url == reqMedia.Url {
				found = true
				break
			}
		}

		if !found {
			media = append(media, &model.ProductMedia{
				Url: reqMedia.Url,
			})
		}
	}

	variantGroups := payload.GenerateVariantGroups()
	for _, vg := range variantGroups {
		vg.ProductID = product.ID

		for _, value := range product.VariantGroup {
			if vg.Name == value.Name {
				vg.CreatedAt = value.CreatedAt
				vg.UpdatedAt = value.UpdatedAt
				break
			}
		}
	}

	var newVarGroup []*model.VariantGroup

	newVarGroup, err = r.variantGroupRepo.Update(tx, product.ID, variantGroups)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	skus := payload.GenerateSKU(variantGroups)

	for _, s := range skus {
		s.ProductId = product.ID
		for idx, sVariant := range s.Variants {
			for _, varGroup := range newVarGroup {
				for _, variant := range varGroup.Variant {
					if variant.ID == sVariant.ID {
						s.Variants[idx].ID = variant.ID
					}
				}
			}
		}
	}

	err = r.skuRepository.Update(tx, product.ID, skus)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Where("id = ?", product.ID).Save(&updatedProduct).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return updatedProduct, nil
}

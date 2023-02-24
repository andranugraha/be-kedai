package repository

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/constant"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"math"

	"gorm.io/gorm"
)

type UserWishlistRepository interface {
	GetUserWishlists(req dto.GetUserWishlistsRequest) ([]*model.UserWishlist, int64, int, error)
	GetUserWishlist(userWishlist *model.UserWishlist) (*model.UserWishlist, error)
	AddUserWishlist(userWishlist *model.UserWishlist) (*model.UserWishlist, error)
	RemoveUserWishlist(userWishlist *model.UserWishlist) error
}

type userWishlistRepositoryImpl struct {
	db *gorm.DB
}

type UserWishlistRConfig struct {
	DB *gorm.DB
}

func NewUserWishlistRepository(cfg *UserWishlistRConfig) UserWishlistRepository {
	return &userWishlistRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *userWishlistRepositoryImpl) GetUserWishlists(req dto.GetUserWishlistsRequest) (userWishlists []*model.UserWishlist, totalRows int64, totalPages int, err error) {
	var (
		active = true
	)

	db := r.db.Preload("Product", func(db *gorm.DB) *gorm.DB {
		return db.Select(`products.*, min(s.price) as min_price, max(s.price) as max_price, 
			concat(c.name, ', ', p.name) as address, 
			max(case when pp.type = 'nominal' then ROUND(cast(pp.amount / s.price * 100 as numeric), 2) else pp.amount end) as promotion_percent, 
			(select url from product_medias pm where products.id = pm.product_id limit 1) as image_url`).
			Joins("join skus s ON s.product_id = products.id").
			Joins("join shops sh ON sh.id = products.shop_id").
			Joins("join user_addresses ua ON ua.id = sh.address_id").
			Joins("join cities c ON c.id = ua.city_id").
			Joins("join provinces p ON p.id = c.province_id").
			Joins("left join product_promotions pp ON pp.sku_id = s.id and (select count(id) from shop_promotions sp where pp.promotion_id = sp.id and now() between sp.start_period and sp.end_period) > 0").
			Group("products.id, c.name, p.name")
	}).Where("user_wishlists.user_id = ?", req.UserId)

	db = db.Joins("join products p ON p.id = user_wishlists.product_id and p.is_active = ?", active).
		Group("user_wishlists.id, p.id")

	if req.CategoryID > 0 {
		db = db.Where("p.category_id = ?", req.CategoryID)
	}

	if req.MinRating > 0 {
		db = db.Where("p.rating >= ?", req.MinRating)
	}

	if len(req.CityIds) > 0 {
		db = db.Joins("join shops sh ON sh.id = p.shop_id").
			Joins("join user_addresses ua ON ua.id = sh.address_id").
			Where("ua.city_id in (?)", req.CityIds)
	}

	if req.MinPrice > 0 || req.MaxPrice > 0 {
		db = db.Joins("join skus s ON s.product_id = p.id").Group("s.id")

		if req.MinPrice > 0 {
			db = db.Where("s.id = (select id from skus where product_id = p.id and skus.price >= ? limit 1)", req.MinPrice)
		}

		if req.MaxPrice > 0 {
			db = db.Where("s.id = (select id from skus where product_id = p.id and skus.price <= ? limit 1)", req.MaxPrice)
		}
	}

	switch req.Sort {
	case constant.SortByRecommended:
		db = db.Order("p.rating desc, p.sold desc")
	case constant.SortByLatest:
		db = db.Order("user_wishlists.created_at desc")
	case constant.SortByTopSales:
		db = db.Order("p.sold desc")
	case constant.SortByPriceLow:
		if req.MinPrice == 0 && req.MaxPrice == 0 {
			db = db.Joins("join skus s ON s.product_id = p.id and s.id = (select id from skus where product_id = p.id order by price asc limit 1)").Group("s.id")
		}

		db = db.Where("s.id = (select id from skus where product_id = p.id order by price asc limit 1)").Order("s.price asc")
	case constant.SortByPriceHigh:
		if req.MinPrice == 0 && req.MaxPrice == 0 {
			db = db.Joins("join skus s ON s.product_id = p.id and s.id = (select id from skus where product_id = p.id order by price asc limit 1)").Group("s.id")
		}

		db = db.Where("s.id = (select id from skus where product_id = p.id order by price asc limit 1)").Order("s.price desc")
	default:
		db = db.Order("user_wishlists.created_at desc")
	}

	db.Model(&userWishlists).Count(&totalRows)
	totalPages = int(math.Ceil(float64(totalRows) / float64(req.Limit)))

	err = db.Find(&userWishlists).Error
	if err != nil {
		return
	}

	return
}

func (r *userWishlistRepositoryImpl) GetUserWishlist(userWishlist *model.UserWishlist) (*model.UserWishlist, error) {
	var res model.UserWishlist

	err := r.db.Where("user_id = ? AND product_id = ?", userWishlist.UserID, userWishlist.ProductID).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrProductNotInWishlist
		}

		return nil, err
	}

	return &res, nil
}

func (r *userWishlistRepositoryImpl) AddUserWishlist(userWishlist *model.UserWishlist) (*model.UserWishlist, error) {
	err := r.db.Create(userWishlist).Error
	if err != nil {
		if errs.IsDuplicateKeyError(err) {
			return nil, errs.ErrProductInWishlist
		}
		return nil, err
	}

	return userWishlist, nil
}

func (r *userWishlistRepositoryImpl) RemoveUserWishlist(userWishlist *model.UserWishlist) error {
	// hard delete
	res := r.db.Unscoped().Where("user_id = ? AND product_id = ?", userWishlist.UserID, userWishlist.ProductID).Delete(&model.UserWishlist{})
	if err := res.Error; err != nil {
		return err
	}

	if res.RowsAffected < 1 {
		return errs.ErrProductNotInWishlist
	}

	return nil
}

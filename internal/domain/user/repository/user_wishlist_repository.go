package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type UserWishlistRepository interface {
	GetUserWishlists(req dto.GetUserWishlistsRequest) ([]*model.UserWishlist, error)
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

func (r *userWishlistRepositoryImpl) GetUserWishlists(req dto.GetUserWishlistsRequest) ([]*model.UserWishlist, error) {
	var (
		userWishlists []*model.UserWishlist
		active        = true
	)

	err := r.db.Preload("Product", func(db *gorm.DB) *gorm.DB {
		return db.Select(`products.*, min(s.price) as min_price, max(s.price) as max_price, 
			concat(c.name, ', ', p.name) as address, 
			max(case when pp.type = 'nominal' then ROUND(cast(pp.amount / s.price * 100 as numeric), 2) else pp.amount end) as promotion_percent,
			count(t.id) as total_sold`).
			Joins("join skus s ON s.product_id = products.id").
			Joins("join shops sh ON sh.id = products.shop_id").
			Joins("join user_addresses ua ON ua.id = sh.address_id").
			Joins("join cities c ON c.id = ua.city_id").
			Joins("join provinces p ON p.id = c.province_id").
			Joins("left join product_promotions pp ON pp.sku_id = s.id and (select count(id) from shop_promotions sp where pp.promotion_id = sp.id and now() between sp.start_period and sp.end_period) > 0").
			Joins("left join transactions t on s.id = t.sku_id and t.id in (select id from invoice_per_shops ips where t.invoice_id = ips.id and ips.status = 'COMPLETED')").
			Where("products.is_active = ?", active).
			Group("products.id, c.name, p.name")
	}).Where("user_id = ?", req.UserId).Find(&userWishlists).Error
	if err != nil {
		return nil, err
	}

	return userWishlists, nil
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

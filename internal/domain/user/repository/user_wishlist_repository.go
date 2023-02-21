package repository

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type UserWishlistRepository interface {
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

func (r *userWishlistRepositoryImpl) RemoveUserWishlist(userWishlist *model.UserWishlist) error {
	res := r.db.Where("user_id = ? AND product_id = ?", userWishlist.UserID, userWishlist.ProductID).Delete(&model.UserWishlist{})
	if err := res.Error; err != nil {
		return err
	}

	if res.RowsAffected < 1 {
		return errs.ErrProductNotInWishlist
	}

	return nil
}

package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type UserWishlistRepository interface {
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

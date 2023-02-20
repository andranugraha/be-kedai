package repository

import (
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type UserWishlistRepository interface {
	// GetByUserIDAndProductID(userID, productID int) (*model.UserWishlist, error)
	AddUserWishlist(userWishlist *model.UserWishlist) (*model.UserWishlist, error)
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

// func (r *userWishlistRepositoryImpl) GetByUserIDAndProductID(userID, productID int) (*model.UserWishlist, error) {
// 	var userWishlist model.UserWishlist

// 	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&userWishlist).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errs.ErrUserWishlistNotExist
// 		}

// 		return nil, err
// 	}

// 	return &userWishlist, nil
// }

func (r *userWishlistRepositoryImpl) AddUserWishlist(userWishlist *model.UserWishlist) (*model.UserWishlist, error) {
	err := r.db.Create(userWishlist).Error
	if err != nil {
		return userWishlist, err
	}

	return userWishlist, nil
}

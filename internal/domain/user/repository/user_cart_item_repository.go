package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type UserCartItemRepository interface {
	CreateCartItem(cartItem *model.UserCartItem) (*model.UserCartItem, error)
	GetCartItemByUserIdAndSkuId(userId int, skuId int) (*model.UserCartItem, error)
	UpdateCartItem(cartItem *model.UserCartItem) (*model.UserCartItem, error)
}

type userCartItemRepository struct {
	db *gorm.DB
}

type UserCartItemRConfig struct {
	DB *gorm.DB
}

func NewUserCartItemRepository(cfg *UserCartItemRConfig) UserCartItemRepository {
	return &userCartItemRepository{
		db: cfg.DB,
	}
}

func (r *userCartItemRepository) CreateCartItem(cartItem *model.UserCartItem) (*model.UserCartItem, error) {

	err := r.db.Create(cartItem).Error
	if err != nil {
		return nil, err
	}

	return cartItem, nil

}

func (r *userCartItemRepository) GetCartItemByUserIdAndSkuId(userId int, skuId int) (*model.UserCartItem, error) {
	var cartItem model.UserCartItem

	err := r.db.Where("user_id = ? AND sku_id = ?", userId, skuId).First(&cartItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCartItemNotFound
		}
		return nil, err
	}

	return &cartItem, nil
}

func (r *userCartItemRepository) UpdateCartItem(cartItem *model.UserCartItem) (*model.UserCartItem, error) {
	res := r.db.Model(&cartItem).Updates(cartItem)
	if err := res.Error; err != nil {
		return nil, err
	}

	if res.RowsAffected < 1 {
		return nil, errs.ErrCartItemNotFound
	}

	return cartItem, nil
}

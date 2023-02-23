package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
)

type UserCartItemRepository interface {
	CreateCartItem(cartItem *model.CartItem) (*model.CartItem, error)
	GetCartItemByUserIdAndSkuId(userId int, skuId int) (*model.CartItem, error)
	UpdateCartItem(cartItem *model.CartItem) (*model.CartItem, error)
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

func (r *userCartItemRepository) CreateCartItem(cartItem *model.CartItem) (*model.CartItem, error) {

	err := r.db.Create(cartItem).Preload("Sku.Product.Shop").Error
	if err != nil {
		return nil, err
	}

	return cartItem, nil

}

func (r *userCartItemRepository) GetCartItemByUserIdAndSkuId(userId int, skuId int) (*model.CartItem, error) {
	var cartItem model.CartItem

	err := r.db.Where("user_id = ? AND sku_id = ?", userId, skuId).Preload("Sku.Product.Shop").First(&cartItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCartItemNotFound
		}
		return nil, err
	}

	return &cartItem, nil
}

func (r *userCartItemRepository) UpdateCartItem(cartItem *model.CartItem) (*model.CartItem, error) {
	res := r.db.Model(&cartItem).Updates(cartItem)
	if err := res.Error; err != nil {
		return nil, err
	}

	if res.RowsAffected < 1 {
		return nil, errs.ErrCartItemNotFound
	}

	return cartItem, nil
}

package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"math"

	"gorm.io/gorm"
)

type UserCartItemRepository interface {
	CreateCartItem(cartItem *model.CartItem) (*model.CartItem, error)
	GetAllCartItem(*dto.GetCartItemsRequest) (cartItems []*model.CartItem, totalRows int64, totalPages int, err error)
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

func (r *userCartItemRepository) GetAllCartItem(req *dto.GetCartItemsRequest) (cartItems []*model.CartItem, totalRows int64, totalPages int, err error) {

	db := r.db.Where("cart_items.user_id = ?", req.UserId).
		Joins("left join skus s on cart_items.sku_id = s.id").
		Joins("left join products p on s.product_id = p.id").
		Joins("left join shops sh on p.shop_id = sh.id").
		Group("sh.id, cart_items.id").
		Order("cart_items.created_at").
		Preload("Sku.Product.Shop.Address.City").
		Preload("Sku.Product.Shop.Address.Province").
		Preload("Sku.Variants.Group").
		Preload("Sku.Promotion")

	db.Model(&model.CartItem{}).Count(&totalRows)

	totalPages = 1
	if req.Limit > 0 {
		totalPages = int(math.Ceil(float64(totalRows) / float64(req.Limit)))
	}

	err = db.Limit(req.Limit).Offset(req.Offset()).Find(&cartItems).Error
	if err != nil {
		return
	}

	return
}

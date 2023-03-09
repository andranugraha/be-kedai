package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"math"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserCartItemRepository interface {
	CreateCartItem(cartItem *model.CartItem) (*model.CartItem, error)
	GetAllCartItem(*dto.GetCartItemsRequest) (cartItems []*model.CartItem, totalRows int64, totalPages int, err error)
	GetCartItemByUserIdAndSkuId(userId int, skuId int) (*model.CartItem, error)
	UpdateCartItem(cartItem *model.CartItem) (*model.CartItem, error)
	GetCartItemByIdAndUserId(id, userId int) (*model.CartItem, error)
	DeleteCartItemBySkuIdsAndUserId(tx *gorm.DB, skuIds []int, userIds int) error
	DeleteCartItem(*dto.DeleteCartItemRequest) error
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
	var totalCartItem int64
	maxCartItem := 200
	err := r.db.Model(&model.CartItem{}).Where("user_id = ?", cartItem.UserId).Count(&totalCartItem).Error

	if err != nil {
		return nil, err
	}

	if totalCartItem >= int64(maxCartItem) {
		return nil, errs.ErrCartItemLimitExceeded
	}

	err = r.db.Create(cartItem).Preload("Sku.Product.Shop").Error
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
	res := r.db.Where("id = ?", cartItem.ID).Where("user_id = ?", cartItem.UserId).Clauses(clause.Returning{}).Updates(cartItem)
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
		Where(`sh.id IN (SELECT DISTINCT sh.id from shops sh 
			JOIN products p on p.shop_id = sh.id 
			JOIN skus s on s.product_id = p.id 
			JOIN cart_items ci on ci.sku_id = s.id 
			WHERE ci.user_id = ?
			ORDER BY sh.id LIMIT ? OFFSET ?)`, req.UserId, req.Limit, req.Offset()).
		Order("cart_items.created_at").
		Preload("Sku.Product.Shop.Address.City").
		Preload("Sku.Product.Shop.Address.Province").
		Preload("Sku.Product.Shop.Address.Subdistrict").
		Preload("Sku.Variants.Group").
		Preload("Sku.Promotion")

	db.Model(&model.CartItem{}).Count(&totalRows)

	totalPages = 1
	if req.Limit > 0 {
		totalPages = int(math.Ceil(float64(totalRows) / float64(req.Limit)))
	}

	err = db.Find(&cartItems).Error
	if err != nil {
		return
	}

	return
}

func (r *userCartItemRepository) GetCartItemByIdAndUserId(id, userId int) (*model.CartItem, error) {
	var cartItem model.CartItem

	err := r.db.Where("cart_items.id = ? AND cart_items.user_id = ?", id, userId).
		Joins("join skus s on cart_items.sku_id = s.id").
		Joins("join products p on s.product_id = p.id").
		Where("p.is_active = ?", true).
		Preload("Sku.Product.Bulk").Preload("Sku.Promotion").First(&cartItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrCartItemNotFound
		}
		return nil, err
	}

	return &cartItem, nil
}

func (r *userCartItemRepository) DeleteCartItemBySkuIdsAndUserId(tx *gorm.DB, skuIds []int, userId int) error {
	err := tx.Unscoped().Where("sku_id IN ?", skuIds).Where("user_id = ?", userId).Delete(&model.CartItem{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *userCartItemRepository) DeleteCartItem(req *dto.DeleteCartItemRequest) error {
	res := r.db.Unscoped().Where("id in (?) and user_id = ?", req.CartItemIds, req.UserId).Delete(&model.CartItem{})
	if err := res.Error; err != nil {
		return err
	}

	if res.RowsAffected < 1 {
		return errs.ErrCartItemNotFound
	}

	return nil
}

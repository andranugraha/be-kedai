package service

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	productService "kedai/backend/be-kedai/internal/domain/product/service"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
)

type UserCartItemService interface {
	CreateCartItem(*dto.UserCartItemRequest) (*model.UserCartItem, error)
}

type userCartItemServiceImpl struct {
	cartItemRepository repository.UserCartItemRepository
	skuService         productService.SkuService
	userService        UserService
}

type UserCartItemSConfig struct {
	CartItemRepository repository.UserCartItemRepository
	SkuService         productService.SkuService
	UserService        UserService
}

func NewUserCartItemService(cfg *UserCartItemSConfig) UserCartItemService {
	return &userCartItemServiceImpl{
		cartItemRepository: cfg.CartItemRepository,
		skuService:         cfg.SkuService,
		userService:        cfg.UserService,
	}
}

func (s *userCartItemServiceImpl) CreateCartItem(cartItemReq *dto.UserCartItemRequest) (*model.UserCartItem, error) {
	var result *model.UserCartItem
	totalQuantity := cartItemReq.Quantity

	// Check user
	_, err := s.userService.GetByID(cartItemReq.UserId)
	if err != nil {
		return nil, err
	}

	// Check sku (product)
	sku, err := s.skuService.GetByID(cartItemReq.SkuId)
	if err != nil {
		return nil, err
	}

	// Check if cart item already exists
	sameCartItem, err := s.cartItemRepository.GetCartItemByUserIdAndSkuId(cartItemReq.UserId, cartItemReq.SkuId)

	if errors.Is(err, errs.ErrCartItemNotFound) {
		cartItem := cartItemReq.ToUserCartItem()
		result, err = s.cartItemRepository.CreateCartItem(cartItem)
		if err != nil {
			return nil, err
		}

		return result, err
	}

	if sameCartItem != nil {
		totalQuantity += sameCartItem.Quantity
		// Update Cart Quantity
		return nil, nil
	}

	// Check sku quantity
	if sku.Stock < totalQuantity {
		return nil, errs.ErrSkuQuantityNotEnough
	}

	return result, nil
}

package service

import (
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
)

type UserCartItemService interface {
	CreateCartItem(*dto.UserCartItemRequest) (*model.UserCartItem, error)
}

type userCartItemServiceImpl struct {
	cartItemRepository repository.UserCartItemRepository
}

type UserCartItemSConfig struct {
	CartItemRepository repository.UserCartItemRepository
}

func NewUserCartItemService(cfg *UserCartItemSConfig) UserCartItemService {
	return &userCartItemServiceImpl{
		cartItemRepository: cfg.CartItemRepository,
	}
}

func (s *userCartItemServiceImpl) CreateCartItem(cartItemReq *dto.UserCartItemRequest) (*model.UserCartItem, error) {
	cartItem := cartItemReq.ToUserCartItem()

	result, err := s.cartItemRepository.CreateCartItem(cartItem)
	if err != nil {
		return nil, err
	}

	return result, nil
}

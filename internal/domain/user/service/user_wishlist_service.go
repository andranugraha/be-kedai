package service

import (
	productService "kedai/backend/be-kedai/internal/domain/product/service"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
)

type UserWishlistService interface {
	AddUserWishlist(req *dto.UserWishlistRequest) (*model.UserWishlist, error)
}

type userWishlistServiceImpl struct {
	userWishlistRepository repository.UserWishlistRepository
	userService            UserService
	productService         productService.ProductService
}

type UserWishlistSConfig struct {
	UserWishlistRepository repository.UserWishlistRepository
	UserService            UserService
	ProductService         productService.ProductService
}

func NewUserWishlistService(cfg *UserWishlistSConfig) UserWishlistService {
	return &userWishlistServiceImpl{
		userWishlistRepository: cfg.UserWishlistRepository,
		userService:            cfg.UserService,
		productService:         cfg.ProductService,
	}
}

func (s *userWishlistServiceImpl) AddUserWishlist(req *dto.UserWishlistRequest) (*model.UserWishlist, error) {
	var userWishlist model.UserWishlist

	user, err := s.userService.GetByID(req.UserID)
	if err != nil {
		return nil, err
	}

	product, err := s.productService.GetByCode(req.ProductCode)
	if err != nil {
		return nil, err
	}

	userWishlist.UserID = user.ID
	userWishlist.ProductID = product.ID

	return s.userWishlistRepository.AddUserWishlist(&userWishlist)
}

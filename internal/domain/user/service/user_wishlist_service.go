package service

import (
	productService "kedai/backend/be-kedai/internal/domain/product/service"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
)

type UserWishlistService interface {
	GetUserWishlist(req *dto.UserWishlistRequest) (*model.UserWishlist, error)
	RemoveUserWishlist(req *dto.UserWishlistRequest) error
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

func (s *userWishlistServiceImpl) GetUserWishlist(req *dto.UserWishlistRequest) (*model.UserWishlist, error) {
	var userWishlist model.UserWishlist

	user, err := s.userService.GetByID(req.UserId)
	if err != nil {
		return nil, err
	}

	product, err := s.productService.GetByID(req.ProductId)
	if err != nil {
		return nil, err
	}

	userWishlist.UserID = user.ID
	userWishlist.ProductID = product.ID

	return s.userWishlistRepository.GetUserWishlist(&userWishlist)
}

func (s *userWishlistServiceImpl) AddUserWishlist(req *dto.UserWishlistRequest) (*model.UserWishlist, error) {
	var userWishlist model.UserWishlist

	user, err := s.userService.GetByID(req.UserId)
	if err != nil {
		return nil, err
	}

	product, err := s.productService.GetByID(req.ProductId)
	if err != nil {
		return nil, err
	}

	userWishlist.UserID = user.ID
	userWishlist.ProductID = product.ID

	return s.userWishlistRepository.AddUserWishlist(&userWishlist)
}

func (s *userWishlistServiceImpl) RemoveUserWishlist(req *dto.UserWishlistRequest) error {
	var userWishlist model.UserWishlist

	user, err := s.userService.GetByID(req.UserId)
	if err != nil {
		return err
	}

	product, err := s.productService.GetByID(req.ProductId)
	if err != nil {
		return err
	}

	userWishlist.UserID = user.ID
	userWishlist.ProductID = product.ID

	return s.userWishlistRepository.RemoveUserWishlist(&userWishlist)
}

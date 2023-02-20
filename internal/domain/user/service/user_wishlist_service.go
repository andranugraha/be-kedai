package service

import (
	productRepo "kedai/backend/be-kedai/internal/domain/product/repository"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
)

type UserWishlistService interface {
	AddUserWishlist(req *dto.UserWishlistRequest) (*model.UserWishlist, error)
}

type userWishlistServiceImpl struct {
	userWishlistRepository repository.UserWishlistRepository
	userRepository         repository.UserRepository
	productRepository      productRepo.ProductRepository
}

type UserWishlistSConfig struct {
	UserWishlistRepository repository.UserWishlistRepository
	UserRepository         repository.UserRepository
	ProductRepository      productRepo.ProductRepository
}

func NewUserWishlistService(cfg *UserWishlistSConfig) UserWishlistService {
	return &userWishlistServiceImpl{
		userWishlistRepository: cfg.UserWishlistRepository,
		userRepository:         cfg.UserRepository,
		productRepository:      cfg.ProductRepository,
	}
}

func (s *userWishlistServiceImpl) AddUserWishlist(req *dto.UserWishlistRequest) (*model.UserWishlist, error) {
	var userWishlist model.UserWishlist

	user, err := s.userRepository.GetByID(req.UserID)
	if err != nil {
		return nil, err
	}

	product, err := s.productRepository.GetByCode(req.ProductCode)
	if err != nil {
		return nil, err
	}

	userWishlist.UserID = user.ID
	userWishlist.ProductID = product.ID

	return s.userWishlistRepository.AddUserWishlist(&userWishlist)
}

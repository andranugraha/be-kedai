package service

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/repository"
	"strings"
)

type ShopService interface {
	FindShopById(id int) (*model.Shop, error)
	FindShopByUserId(userId int) (*model.Shop, error)
	FindShopBySlug(slug string) (*model.Shop, error)
	FindShopByKeyword(req *dto.FindShopRequest) (*commonDto.PaginationResponse, error)
}

type shopServiceImpl struct {
	shopRepository repository.ShopRepository
}

type ShopSConfig struct {
	ShopRepository repository.ShopRepository
}

func NewShopService(cfg *ShopSConfig) ShopService {
	return &shopServiceImpl{
		shopRepository: cfg.ShopRepository,
	}
}

func (s *shopServiceImpl) FindShopById(id int) (*model.Shop, error) {
	return s.shopRepository.FindShopById(id)
}

func (s *shopServiceImpl) FindShopByUserId(userId int) (*model.Shop, error) {
	return s.shopRepository.FindShopByUserId(userId)
}

func (s *shopServiceImpl) FindShopBySlug(slug string) (*model.Shop, error) {
	shop, err := s.shopRepository.FindShopBySlug(slug)
	if err != nil {
		return nil, err
	}

	return shop, nil
}

func (s *shopServiceImpl) FindShopByKeyword(req *dto.FindShopRequest) (*commonDto.PaginationResponse, error) {
	validateKeyword := strings.Trim(req.Keyword, " ")
	if validateKeyword == "" {
		return &commonDto.PaginationResponse{
			Data:       nil,
			TotalRows:  0,
			TotalPages: 0,
			Limit:      req.Limit,
			Page:       req.Page,
		}, nil
	}

	res, rows, pages, err := s.shopRepository.FindShopByKeyword(req)
	if err != nil {
		return nil, err
	}

	return &commonDto.PaginationResponse{
		Data:       res,
		TotalRows:  rows,
		TotalPages: pages,
		Limit:      req.Limit,
		Page:       req.Page,
	}, nil
}

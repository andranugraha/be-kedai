package service

import (
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
	"kedai/backend/be-kedai/internal/domain/marketplace/repository"
)

type CategotyService interface {
	AddCategory(req *dto.CategoryDTO) error
}

type categotyServiceImpl struct {
	categoryRepository repository.CategoryRepository
}

type CategorySConfig struct {
	CategoryRepository repository.CategoryRepository
}

func NewCategoryService(cfg *CategorySConfig) CategotyService {
	return &categotyServiceImpl{
		categoryRepository: cfg.CategoryRepository,
	}
}

func (s *categotyServiceImpl) AddCategory(req *dto.CategoryDTO) error {
	return s.categoryRepository.AddCategory(req)
}

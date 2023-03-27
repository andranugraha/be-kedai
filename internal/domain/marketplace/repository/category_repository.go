package repository

import (
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	AddCategory(req *dto.CategoryDTO) error
}

type categoryRepositoryImpl struct {
	db *gorm.DB
}

type CategoryRConfig struct {
	DB *gorm.DB
}

func NewCategoryRepository(cfg *CategoryRConfig) CategoryRepository {
	return &categoryRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *categoryRepositoryImpl) AddCategory(req *dto.CategoryDTO) error {
	return nil
}

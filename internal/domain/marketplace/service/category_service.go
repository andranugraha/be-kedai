package service

import (
	"kedai/backend/be-kedai/internal/domain/marketplace/repository"
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
)

type CategotyService interface {
	AddCategory(req *dto.CategoryDTO) error
}

type categotyServiceImpl struct {
	categoryRepository repository
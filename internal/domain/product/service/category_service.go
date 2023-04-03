package service

import (
	"kedai/backend/be-kedai/internal/common/dto"
	categoryDto "kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/product/repository"
)

type CategoryService interface {
	GetCategories(categoryDto.GetCategoriesRequest) (*dto.PaginationResponse, error)
	GetCategoryLineAgesFromBottom(categoryID int) ([]*model.Category, error)
	GetCategoryIDLineAgesFromTop(categoryID int) ([]int, error)
	AddCategory(category *model.Category) error
}

type categoryServiceImpl struct {
	categoryRepo repository.CategoryRepository
}

type CategorySConfig struct {
	CategoryRepo repository.CategoryRepository
}

func NewCategoryService(cfg *CategorySConfig) CategoryService {
	return &categoryServiceImpl{
		categoryRepo: cfg.CategoryRepo,
	}
}

func (c *categoryServiceImpl) GetCategories(query categoryDto.GetCategoriesRequest) (res *dto.PaginationResponse, err error) {
	categories, totalRows, totalPages, err := c.categoryRepo.GetAll(query)
	if err != nil {
		return
	}

	if query.WithPrice {
		for _, c := range categories {
			getCategoryMinPrice(c)
		}
	}

	for _, c := range categories {
		removeChildren(c, query.Depth)
	}

	res = &dto.PaginationResponse{
		Data:       categories,
		Limit:      query.Limit,
		Page:       query.Page,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}

	return
}

func (c *categoryServiceImpl) GetCategoryLineAgesFromBottom(categoryID int) ([]*model.Category, error) {
	return c.categoryRepo.GetLineageFromBottom(categoryID)
}

func (c *categoryServiceImpl) GetCategoryIDLineAgesFromTop(categoryID int) ([]int, error) {
	return c.categoryRepo.GetLineageFromTop(categoryID)
}

func removeChildren(category *model.Category, depth int) {
	if depth == 0 {
		category.Children = []*model.Category{}
		return
	}

	for _, c := range category.Children {
		removeChildren(c, depth-1)
	}
}

func getCategoryMinPrice(category *model.Category) float64 {
	var minPrice float64
	for _, c := range category.Children {
		if c.MinPrice != nil && (*c.MinPrice < minPrice || minPrice == 0) {
			minPrice = *c.MinPrice
			continue
		}

		if len(c.Children) > 0 {
			minPrice = getCategoryMinPrice(c)
		}

	}

	category.MinPrice = &minPrice
	return minPrice
}

func (s *categoryServiceImpl) AddCategory(category *model.Category) error {
	return s.categoryRepo.AddCategory(category)
}

package repository

import (
	"kedai/backend/be-kedai/internal/domain/product/dto"
	model "kedai/backend/be-kedai/internal/domain/product/model"
	"math"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetAll(dto.GetCategoriesRequest) ([]*model.Category, int64, int, error)
	GetLineageFromBottom(categoryID int) ([]*model.Category, error)
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

func (c *categoryRepositoryImpl) GetAll(query dto.GetCategoriesRequest) (categories []*model.Category, totalRows int64, totalPages int, err error) {
	db := nestedPreload(c.db, query)

	if query.ParentID != 0 {
		db = db.Where("categories.id = ?", query.ParentID)
	} else {
		db = db.Where("categories.parent_id is null")
	}

	db.Model(&categories).Count(&totalRows)

	totalPages = 1
	if query.Limit > 0 {
		totalPages = int(math.Ceil(float64(totalRows) / float64(query.Limit)))
	}

	err = db.Scopes(scope(query)).Limit(query.Limit).Offset(query.Offset()).Find(&categories).Error
	if err != nil {
		return
	}

	return
}

func (r *categoryRepositoryImpl) GetLineageFromBottom(categoryID int) ([]*model.Category, error) {
	var categories []*model.Category

	query := `
		WITH RECURSIVE category_lineages AS (
			SELECT *
			FROM categories
			WHERE id = ?
			UNION
			SELECT c.*
			FROM categories AS c
			JOIN category_lineages AS cl ON cl.parent_id = c.id
		)
		SELECT * FROM category_lineages ORDER BY id ASC
	`

	err := r.db.Raw(query, categoryID).Scan(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func nestedPreload(db *gorm.DB, query dto.GetCategoriesRequest) *gorm.DB {
	return db.Preload("Children", func(db *gorm.DB) *gorm.DB {
		return nestedPreload(db, dto.GetCategoriesRequest{
			WithPrice: query.WithPrice,
		}).Scopes(scope(query))
	})
}

func scope(query dto.GetCategoriesRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if query.WithPrice {
			db = db.Joins("LEFT JOIN products p ON p.category_id = categories.id").Joins("LEFT JOIN skus s on p.id = s.product_id").
				Select("categories.*, MIN(s.price) as min_price").
				Group("categories.id")
		}

		return db
	}
}

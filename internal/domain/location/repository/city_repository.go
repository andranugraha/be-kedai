package repository

import (
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"
	"math"

	"gorm.io/gorm"
)

type CityRepository interface {
	GetAll(dto.GetCitiesRequest) ([]*model.City, int64, int, error)
}

type cityRepositoryImpl struct {
	db *gorm.DB
}

type CityRConfig struct {
	DB *gorm.DB
}

func NewCityRepository(cfg *CityRConfig) CityRepository {
	return &cityRepositoryImpl{
		db: cfg.DB,
	}
}

func (c *cityRepositoryImpl) GetAll(req dto.GetCitiesRequest) (cities []*model.City, totalRows int64, totalPages int, err error) {
	db := c.db.Scopes(req.Scope())
	db.Model(&cities).Count(&totalRows)

	totalPages = 1
	if req.Limit > 0 {
		totalPages = int(math.Ceil(float64(totalRows) / float64(req.Limit)))
	}

	err = db.Limit(int(req.Limit)).Offset(req.Offset()).Find(&cities).Error
	if err != nil {
		return
	}

	return
}

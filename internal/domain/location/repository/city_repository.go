package repository

import (
	"kedai/backend/be-kedai/internal/domain/location/model"

	"gorm.io/gorm"
)

type CityRepository interface {
	GetCities() ([]*model.City, error)
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

func (c *cityRepositoryImpl) GetCities() ([]*model.City, error) {
	var cities []*model.City
	err := c.db.Find(&cities).Error
	if err != nil {
		return nil, err
	}

	return cities, nil
}

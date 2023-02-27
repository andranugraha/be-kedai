package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/model"

	"gorm.io/gorm"
)

type DistrictRepository interface {
	GetByID(cityID int) (*model.District, error)
}

type districtRepositoryImpl struct {
	db *gorm.DB
}

type DistrictRConfig struct {
	DB *gorm.DB
}

func NewDistrictRepository(cfg *DistrictRConfig) DistrictRepository {
	return &districtRepositoryImpl{
		db: cfg.DB,
	}
}

func (c *districtRepositoryImpl) GetByID(districtID int) (district *model.District, err error) {
	err = c.db.First(&district, districtID).Error
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrDistrictNotFound
		}
		return
	}

	return
}

package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"

	"gorm.io/gorm"
)

type DistrictRepository interface {
	GetByID(cityID int) (*model.District, error)
	GetAll(req dto.GetDistrictsRequest) (districts []*model.District, err error)
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

func (c *districtRepositoryImpl) GetAll(req dto.GetDistrictsRequest) (districts []*model.District, err error) {
	db := c.db
	if req.CityID != 0 {
		db = db.Where("districts.city_id = ?", req.CityID)
	}

	err = db.Find(&districts).Error
	if err != nil {
		return
	}

	return
}

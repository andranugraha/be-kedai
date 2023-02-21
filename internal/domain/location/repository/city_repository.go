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
	db := c.db
	if req.ProvinceID != 0 {
		db = db.Where("cities.province_id = ?", req.ProvinceID)
	}
	if req.Sort == "most_shops" {
		db = db.Joins("left join user_addresses ua2 on cities.id = ua2.city_id and (select count(ua.id) from user_addresses ua inner join shops s on ua.id = s.address_id) > 0").
			Group("cities.id").
			Order("count(ua2.id) desc, cities.name asc")
	}

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

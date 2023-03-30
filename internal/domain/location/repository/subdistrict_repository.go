package repository

import (
	"errors"
	"fmt"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"
	"strings"

	"gorm.io/gorm"
)

type SubdistrictRepository interface {
	GetByID(subdistrictID int) (*model.Subdistrict, error)
	GetAll(req dto.GetSubdistrictsRequest) (subdistricts []*model.Subdistrict, err error)
	GetDetailByNameAndPostalCode(subdistrictName string, postalCode string) (*model.Subdistrict, error)
	GetDetailByNameAndDistrictCityName(subdistrictName, districtName, cityName string) (*model.Subdistrict, error)
}

type subdistrictRepositoryImpl struct {
	db *gorm.DB
}

type SubdistrictRConfig struct {
	DB *gorm.DB
}

func NewSubdistrictRepository(cfg *SubdistrictRConfig) SubdistrictRepository {
	return &subdistrictRepositoryImpl{
		db: cfg.DB,
	}
}

func (c *subdistrictRepositoryImpl) GetByID(subdistrictID int) (subdistrict *model.Subdistrict, err error) {
	err = c.db.First(&subdistrict, subdistrictID).Error
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrSubdistrictNotFound
		}
		return
	}

	return
}

func (c *subdistrictRepositoryImpl) GetAll(req dto.GetSubdistrictsRequest) (subdistricts []*model.Subdistrict, err error) {
	db := c.db
	if req.DistrictID != 0 {
		db = db.Where("subdistricts.district_id = ?", req.DistrictID)
	}

	err = db.Find(&subdistricts).Error
	if err != nil {
		return
	}

	return
}

func (c *subdistrictRepositoryImpl) GetDetailByNameAndPostalCode(subdistrictName string, postalCode string) (subdistrict *model.Subdistrict, err error) {
	err = c.db.Where("lower(name) ilike ? AND levenshtein(postal_code, ?) < 2", fmt.Sprintf("%%%s%%", strings.ToLower(subdistrictName)), postalCode).Preload("District.City.Province").First(&subdistrict).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrSubdistrictNotFound
		}
		return
	}

	return
}

func (c *subdistrictRepositoryImpl) GetDetailByNameAndDistrictCityName(subdistrictName, districtName, cityName string) (subdistrict *model.Subdistrict, err error) {
	err = c.db.Where("lower(subdistricts.name) ilike ? AND lower(districts.name) ilike ? AND lower(cities.name) ilike ?",
		fmt.Sprintf("%%%s%%", strings.ToLower(subdistrictName)), fmt.Sprintf("%%%s%%", strings.ToLower(districtName)), fmt.Sprintf("%%%s%%", strings.ToLower(cityName))).
		Joins("JOIN districts ON districts.id = subdistricts.district_id").
		Joins("JOIN cities ON cities.id = districts.city_id").
		Preload("District.City.Province").
		First(&subdistrict).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrSubdistrictNotFound
		}
		return
	}

	return
}

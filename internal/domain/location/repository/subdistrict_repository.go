package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"

	"gorm.io/gorm"
)

type SubdistrictRepository interface {
	GetByID(subdistrictID int) (*model.Subdistrict, error)
	GetAll(req dto.GetSubdistrictsRequest) (subdistricts []*model.Subdistrict, err error)
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

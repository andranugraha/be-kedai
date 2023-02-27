package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/model"

	"gorm.io/gorm"
)

type SubdistrictRepository interface {
	GetByID(subdistrictID int) (*model.Subdistrict, error)
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

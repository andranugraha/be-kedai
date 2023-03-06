package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/model"

	"gorm.io/gorm"
)

type ProvinceRepository interface {
	GetAll() (provinces []*model.Province, err error)
	GetByID(provinceID int) (*model.Province, error)
}

type provinceRepositoryImpl struct {
	db *gorm.DB
}

type ProvinceRConfig struct {
	DB *gorm.DB
}

func NewProvinceRepository(cfg *ProvinceRConfig) ProvinceRepository {
	return &provinceRepositoryImpl{
		db: cfg.DB,
	}
}

func (c *provinceRepositoryImpl) GetAll() (provinces []*model.Province, err error) {
	err = c.db.Find(&provinces).Error
	if err != nil {
		return
	}

	return
}

func (c *provinceRepositoryImpl) GetByID(provinceID int) (province *model.Province, err error) {
	err = c.db.First(&province, provinceID).Error
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrProvinceNotFound
		}
		return
	}

	return
}

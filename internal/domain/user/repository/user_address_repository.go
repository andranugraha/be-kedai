package repository

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/model"

	"gorm.io/gorm"
)

type UserAddressRepository interface {
	AddUserAddress(*model.UserAddress) (*model.UserAddress, error)
}

type userAddressRepository struct {
	db              *gorm.DB
	userProfileRepo UserProfileRepository
}

type UserAddressRConfig struct {
	DB              *gorm.DB
	UserProfileRepo UserProfileRepository
}

func NewUserAddressRepository(cfg *UserAddressRConfig) UserAddressRepository {
	return &userAddressRepository{
		db:              cfg.DB,
		userProfileRepo: cfg.UserProfileRepo,
	}
}

func (r *userAddressRepository) AddUserAddress(newAddress *model.UserAddress) (*model.UserAddress, error) {
	var totalRows int64 = 0
	var maxAddress int64 = 10

	err := r.db.Model(&model.UserAddress{}).Where("user_id = ?", newAddress.UserID).Count(&totalRows).Error
	if err != nil {
		return nil, err
	}

	if totalRows >= maxAddress {
		return nil, errs.ErrMaxAddress
	}

	if totalRows == 0 {
		newAddress.IsDefault = true
	}

	tx := r.db.Begin()
	defer tx.Commit()

	err = r.db.Create(newAddress).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if newAddress.IsDefault {
		err = r.userProfileRepo.UpdateDefaultAddressId(tx, newAddress.UserID, newAddress.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return newAddress, nil
}

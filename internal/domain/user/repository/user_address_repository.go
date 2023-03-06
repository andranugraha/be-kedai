package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/model"

	"gorm.io/gorm"
)

type UserAddressRepository interface {
	GetUserAddressByIdAndUserId(addressId int, userId int) (*model.UserAddress, error)
	AddUserAddress(*model.UserAddress) (*model.UserAddress, error)
	GetAllUserAddress(userId int) ([]*model.UserAddress, error)
	DefaultAddressTransaction(tx *gorm.DB, userId int, addressId int) error
	UpdateUserAddress(*model.UserAddress) (*model.UserAddress, error)
	DeleteUserAddress(addressId int, userId int) error
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
		err = r.DefaultAddressTransaction(tx, newAddress.UserID, newAddress.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return newAddress, nil
}

func (r *userAddressRepository) GetAllUserAddress(userId int) ([]*model.UserAddress, error) {
	var addresses []*model.UserAddress

	err := r.db.Where("user_id = ?", userId).
		Preload("Subdistrict").
		Preload("District").
		Preload("City").
		Preload("Province").
		Find(&addresses).Error
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *userAddressRepository) DefaultAddressTransaction(tx *gorm.DB, userId int, addressId int) error {
	err := r.userProfileRepo.UpdateDefaultAddressId(tx, userId, addressId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *userAddressRepository) UpdateUserAddress(address *model.UserAddress) (*model.UserAddress, error) {
	tx := r.db.Begin()
	defer tx.Commit()

	res := tx.Model(&model.UserAddress{}).Where("id = ?", address.ID).Updates(address)
	if err := res.Error; err != nil {
		return nil, err
	}

	if res.RowsAffected == 0 {
		return nil, errs.ErrAddressNotFound
	}

	if address.IsDefault {
		err := r.DefaultAddressTransaction(tx, address.UserID, address.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return address, nil
}

func (r *userAddressRepository) GetUserAddressByIdAndUserId(addressId int, userId int) (*model.UserAddress, error) {
	var address model.UserAddress

	err := r.db.Where("id = ? AND user_id = ?", addressId, userId).First(&address).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrAddressNotFound
		}
		return nil, err
	}

	return &address, nil
}

func (r *userAddressRepository) DeleteUserAddress(addressId int, userId int) error {

	res := r.db.Delete(&model.UserAddress{}, "id = ? AND user_id = ?", addressId, userId)
	if err := res.Error; err != nil {
		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrAddressNotFound
	}

	return nil
}

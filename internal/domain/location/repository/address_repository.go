package repository

import (
	"context"
	"errors"
	"kedai/backend/be-kedai/internal/common/constant"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"

	errs "kedai/backend/be-kedai/internal/common/error"
	shopRepo "kedai/backend/be-kedai/internal/domain/shop/repository"
	userRepo "kedai/backend/be-kedai/internal/domain/user/repository"

	"googlemaps.github.io/maps"
	"gorm.io/gorm"
)

type AddressRepository interface {
	SearchAddress(req *dto.SearchAddressRequest) ([]*dto.SearchAddressResponse, error)
	GetUserAddressByIdAndUserId(addressId int, userId int) (*model.UserAddress, error)
	AddUserAddress(*model.UserAddress) (*model.UserAddress, error)
	GetAllUserAddress(userId int) ([]*model.UserAddress, error)
	DefaultAddressTransaction(tx *gorm.DB, userId int, addressId int) error
	PickupAddressTransaction(tx *gorm.DB, userId int, addressId int) error
	UpdateUserAddress(*model.UserAddress) (*model.UserAddress, error)
	DeleteUserAddress(addressId int, userId int) error
}

type addressRepositoryImpl struct {
	db              *gorm.DB
	googleMaps      *maps.Client
	userProfileRepo userRepo.UserProfileRepository
	shopRepo        shopRepo.ShopRepository
}

type AddressRConfig struct {
	DB              *gorm.DB
	GoogleMaps      *maps.Client
	UserProfileRepo userRepo.UserProfileRepository
	ShopRepo        shopRepo.ShopRepository
}

func NewAddressRepository(cfg *AddressRConfig) AddressRepository {
	return &addressRepositoryImpl{
		db:              cfg.DB,
		googleMaps:      cfg.GoogleMaps,
		userProfileRepo: cfg.UserProfileRepo,
		shopRepo:        cfg.ShopRepo,
	}
}

func (c *addressRepositoryImpl) SearchAddress(req *dto.SearchAddressRequest) (addresses []*dto.SearchAddressResponse, err error) {
	ctx := context.Background()
	defer ctx.Done()

	autoCompleteRequest := &maps.PlaceAutocompleteRequest{
		Input: req.Keyword,
		Components: map[maps.Component][]string{
			maps.ComponentCountry: {"id"},
		},
		Language: "id",
	}

	autocomplete, err := c.googleMaps.PlaceAutocomplete(ctx, autoCompleteRequest)
	if err != nil {
		return
	}

	for _, place := range autocomplete.Predictions {
		addresses = append(addresses, &dto.SearchAddressResponse{
			PlaceID:     place.PlaceID,
			Description: place.Description,
		})
	}

	return
}

func (r *addressRepositoryImpl) AddUserAddress(newAddress *model.UserAddress) (*model.UserAddress, error) {
	var totalRows int64 = 0
	var trueValue = true

	err := r.db.Model(&model.UserAddress{}).Where("user_id = ?", newAddress.UserID).Count(&totalRows).Error
	if err != nil {
		return nil, err
	}

	if totalRows >= constant.MaxAddressLimit {
		return nil, errs.ErrMaxAddress
	}

	if totalRows == 0 {
		newAddress.IsDefault = &trueValue
	}

	tx := r.db.Begin()
	defer tx.Commit()

	err = r.db.Create(newAddress).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if *newAddress.IsDefault {
		err = r.DefaultAddressTransaction(tx, newAddress.UserID, newAddress.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if *newAddress.IsPickup {
		err = r.PickupAddressTransaction(tx, newAddress.UserID, newAddress.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return newAddress, nil
}

func (r *addressRepositoryImpl) PickupAddressTransaction(tx *gorm.DB, userId int, addressId int) error {
	err := r.shopRepo.UpdateShopAddressIdByUserId(tx, userId, addressId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *addressRepositoryImpl) GetAllUserAddress(userId int) ([]*model.UserAddress, error) {
	var addresses []*model.UserAddress

	err := r.db.Where("user_id = ?", userId).
		Preload("Subdistrict").
		Preload("District").
		Preload("City").
		Preload("Province").
		Order("created_at desc").
		Find(&addresses).Error
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *addressRepositoryImpl) DefaultAddressTransaction(tx *gorm.DB, userId int, addressId int) error {
	err := r.userProfileRepo.UpdateDefaultAddressId(tx, userId, addressId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *addressRepositoryImpl) UpdateUserAddress(address *model.UserAddress) (*model.UserAddress, error) {
	tx := r.db.Begin()
	defer tx.Commit()

	res := tx.Model(&model.UserAddress{}).Where("id = ?", address.ID).Updates(address)
	if err := res.Error; err != nil {
		return nil, err
	}

	if res.RowsAffected == 0 {
		return nil, errs.ErrAddressNotFound
	}

	if *address.IsDefault {
		err := r.DefaultAddressTransaction(tx, address.UserID, address.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if *address.IsPickup {
		err := r.PickupAddressTransaction(tx, address.UserID, address.ID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return address, nil
}

func (r *addressRepositoryImpl) GetUserAddressByIdAndUserId(addressId int, userId int) (*model.UserAddress, error) {
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

func (r *addressRepositoryImpl) DeleteUserAddress(addressId int, userId int) error {

	res := r.db.Delete(&model.UserAddress{}, "id = ? AND user_id = ?", addressId, userId)
	if err := res.Error; err != nil {
		return err
	}

	if res.RowsAffected == 0 {
		return errs.ErrAddressNotFound
	}

	return nil
}

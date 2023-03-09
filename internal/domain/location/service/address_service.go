package service

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/location/repository"
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	shopService "kedai/backend/be-kedai/internal/domain/shop/service"
	userService "kedai/backend/be-kedai/internal/domain/user/service"
)

type AddressService interface {
	AddUserAddress(*dto.AddressRequest) (*model.UserAddress, error)
	UpdateUserAddress(*dto.AddressRequest) (*model.UserAddress, error)
	GetAllUserAddress(userId int) ([]*model.UserAddress, error)
	PreCheckAddress(*dto.AddressRequest) (*model.UserAddress, error)
	DeleteUserAddress(addressId int, userId int) error
	GetUserAddressByIdAndUserId(addressId int, userId int) (*model.UserAddress, error)
}

type addressService struct {
	addressRepo        repository.AddressRepository
	provinceService    ProvinceService
	districtService    DistrictService
	subdistrictService SubdistrictService
	cityService        CityService
	userProfileService userService.UserProfileService
	shopService        shopService.ShopService
}

type AddressSConfig struct {
	AddressRepo        repository.AddressRepository
	ProvinceService    ProvinceService
	DistrictService    DistrictService
	SubdistrictService SubdistrictService
	CityService        CityService
	UserProfileService userService.UserProfileService
	ShopService        shopService.ShopService
}

func NewAddressService(cfg *AddressSConfig) AddressService {
	return &addressService{
		addressRepo:        cfg.AddressRepo,
		provinceService:    cfg.ProvinceService,
		districtService:    cfg.DistrictService,
		subdistrictService: cfg.SubdistrictService,
		cityService:        cfg.CityService,
		userProfileService: cfg.UserProfileService,
		shopService:        cfg.ShopService,
	}
}

func (s *addressService) AddUserAddress(newAddress *dto.AddressRequest) (*model.UserAddress, error) {
	address, err := s.PreCheckAddress(newAddress)
	if err != nil {
		return nil, err
	}

	address, err = s.addressRepo.AddUserAddress(address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (s *addressService) GetAllUserAddress(userId int) ([]*model.UserAddress, error) {
	profile, err := s.userProfileService.GetProfile(userId)
	if err != nil {
		return nil, err
	}

	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil && !errors.Is(err, errs.ErrShopNotFound) {
		return nil, err
	}

	if shop == nil {
		shop = &shopModel.Shop{}
	}

	addresses, err := s.addressRepo.GetAllUserAddress(userId)
	if err != nil {
		return nil, err
	}

	return dto.ToAddressList(addresses, profile.DefaultAddressID, &shop.AddressID), nil
}

func (s *addressService) PreCheckAddress(newAddress *dto.AddressRequest) (*model.UserAddress, error) {
	var address *model.UserAddress

	subdistrict, err := s.subdistrictService.GetSubdistrictByID(newAddress.SubdistrictID)
	if err != nil {
		return nil, err
	}

	district, err := s.districtService.GetDistrictByID(subdistrict.DistrictID)
	if err != nil {
		return nil, err
	}

	city, err := s.cityService.GetCityByID(district.CityID)
	if err != nil {
		return nil, err
	}

	province, err := s.provinceService.GetProvinceByID(city.ProvinceID)
	if err != nil {
		return nil, err
	}

	address = newAddress.ToUserAddress()
	address.ProvinceID = province.ID
	address.CityID = city.ID
	address.DistrictID = district.ID

	return address, nil
}

func (s *addressService) UpdateUserAddress(updatedAddress *dto.AddressRequest) (*model.UserAddress, error) {
	address, err := s.addressRepo.GetUserAddressByIdAndUserId(updatedAddress.ID, updatedAddress.UserID)
	if err != nil {
		return nil, err
	}

	if !*(updatedAddress.IsDefault) {
		profile, err := s.userProfileService.GetProfile(address.UserID)
		if err != nil {
			return nil, err
		}

		if profile.DefaultAddressID != nil && *profile.DefaultAddressID == address.ID {
			return nil, errs.ErrMustHaveAtLeastOneDefaultAddress
		}
	}

	if !*(updatedAddress.IsPickup) {
		shop, err := s.shopService.FindShopByUserId(address.UserID)
		if err != nil {
			return nil, err
		}

		if shop.AddressID == updatedAddress.ID {
			return nil, errs.ErrMustHaveAtLeastOnePickupAddress
		}
	}

	address, err = s.PreCheckAddress(updatedAddress)
	if err != nil {
		return nil, err
	}

	address, err = s.addressRepo.UpdateUserAddress(address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (s *addressService) DeleteUserAddress(addressId int, userId int) error {
	address, err := s.addressRepo.GetUserAddressByIdAndUserId(addressId, userId)
	if err != nil {
		return err
	}

	profile, err := s.userProfileService.GetProfile(userId)
	if err != nil {
		return err
	}

	if profile.DefaultAddressID != nil && *profile.DefaultAddressID == address.ID {
		return errs.ErrMustHaveAtLeastOneDefaultAddress
	}

	shop, err := s.shopService.FindShopByUserId(userId)
	if err != nil && !errors.Is(err, errs.ErrShopNotFound) {
		return err
	}

	if shop != nil && shop.AddressID == address.ID {
		return errs.ErrMustHaveAtLeastOnePickupAddress
	}

	err = s.addressRepo.DeleteUserAddress(addressId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *addressService) GetUserAddressByIdAndUserId(addressId int, userId int) (*model.UserAddress, error) {
	address, err := s.addressRepo.GetUserAddressByIdAndUserId(addressId, userId)
	if err != nil {
		return nil, err
	}

	return address, nil
}

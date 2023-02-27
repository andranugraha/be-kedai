package service

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"
	locationService "kedai/backend/be-kedai/internal/domain/location/service"
	"kedai/backend/be-kedai/internal/domain/user/repository"
)

type UserAddressService interface {
	AddUserAddress(*dto.AddressRequest) (*model.UserAddress, error)
	UpdateUserAddress(*dto.AddressRequest) (*model.UserAddress, error)
	GetAllUserAddress(userId int) ([]*model.UserAddress, error)
	PreCheckAddress(*dto.AddressRequest) (*model.UserAddress, error)
}

type userAddressService struct {
	userAddressRepo    repository.UserAddressRepository
	provinceService    locationService.ProvinceService
	districtService    locationService.DistrictService
	subdistrictService locationService.SubdistrictService
	cityService        locationService.CityService
	userProfileService UserProfileService
}

type UserAddressSConfig struct {
	UserAddressRepo    repository.UserAddressRepository
	ProvinceService    locationService.ProvinceService
	DistrictService    locationService.DistrictService
	SubdistrictService locationService.SubdistrictService
	CityService        locationService.CityService
	UserProfileService UserProfileService
}

func NewUserAddressService(cfg *UserAddressSConfig) UserAddressService {
	return &userAddressService{
		userAddressRepo:    cfg.UserAddressRepo,
		provinceService:    cfg.ProvinceService,
		districtService:    cfg.DistrictService,
		subdistrictService: cfg.SubdistrictService,
		cityService:        cfg.CityService,
		userProfileService: cfg.UserProfileService,
	}
}

func (s *userAddressService) AddUserAddress(newAddress *dto.AddressRequest) (*model.UserAddress, error) {
	address, err := s.PreCheckAddress(newAddress)
	if err != nil {
		return nil, err
	}

	address, err = s.userAddressRepo.AddUserAddress(address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (s *userAddressService) GetAllUserAddress(userId int) ([]*model.UserAddress, error) {
	profile, err := s.userProfileService.GetProfile(userId)
	if err != nil {
		return nil, err
	}

	addresses, err := s.userAddressRepo.GetAllUserAddress(userId)
	if err != nil {
		return nil, err
	}

	return dto.ToAddressList(addresses, profile.DefaultAddressID), nil
}

func (s *userAddressService) PreCheckAddress(newAddress *dto.AddressRequest) (*model.UserAddress, error) {
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

func (s *userAddressService) UpdateUserAddress(newAddress *dto.AddressRequest) (*model.UserAddress, error) {
	address, err := s.userAddressRepo.GetUserAddressByIdAndUserId(newAddress.ID, newAddress.UserID)
	if err != nil {
		return nil, err
	}

	profile, err := s.userProfileService.GetProfile(newAddress.UserID)
	if err != nil {
		return nil, err
	}

	if profile.DefaultAddressID != nil && *profile.DefaultAddressID == address.ID && !*(newAddress.IsDefault) {
		return nil, errs.ErrMustHaveAtLeastOneDefaultAddress
	}

	address, err = s.PreCheckAddress(newAddress)
	if err != nil {
		return nil, err
	}

	address, err = s.userAddressRepo.UpdateUserAddress(address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

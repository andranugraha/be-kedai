package service_test

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/location/service"
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddUserAddress(t *testing.T) {
	var (
		falseValue = false
	)
	type input struct {
		data        *dto.AddressRequest
		err         error
		beforeTests func(mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockAddressRepo *mocks.AddressRepository)
	}
	type expected struct {
		data *model.UserAddress
		err  error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when GetSubdistrictByID return error",
			input: input{
				data: &dto.AddressRequest{
					SubdistrictID: 1,
				},
				err: nil,
				beforeTests: func(mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockAddressRepo *mocks.AddressRepository) {
					mockSubdistrictService.On("GetSubdistrictByID", 1).Return(nil, errs.ErrSubdistrictNotFound)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrSubdistrictNotFound,
			},
		},
		{
			description: "should return error when GetDistrictByID return error",
			input: input{
				data: &dto.AddressRequest{
					SubdistrictID: 1,
				},
				err: errs.ErrDistrictNotFound,
				beforeTests: func(mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockAddressRepo *mocks.AddressRepository) {
					mockSubdistrictService.On("GetSubdistrictByID", 1).Return(&model.Subdistrict{
						ID:         1,
						DistrictID: 1,
					}, nil)
					mockDistrictService.On("GetDistrictByID", 1).Return(nil, errs.ErrDistrictNotFound)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrDistrictNotFound,
			},
		},
		{
			description: "should return error when GetCityByID return error",
			input: input{
				data: &dto.AddressRequest{
					SubdistrictID: 1,
				},
				err: errs.ErrCityNotFound,
				beforeTests: func(mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockAddressRepo *mocks.AddressRepository) {
					mockSubdistrictService.On("GetSubdistrictByID", 1).Return(&model.Subdistrict{
						ID:         1,
						DistrictID: 1,
					}, nil)
					mockDistrictService.On("GetDistrictByID", 1).Return(&model.District{
						ID:     1,
						CityID: 1,
					}, nil)
					mockCityService.On("GetCityByID", 1).Return(nil, errs.ErrCityNotFound)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrCityNotFound,
			},
		},
		{
			description: "should return error when GetProvinceByID return error",
			input: input{
				data: &dto.AddressRequest{
					SubdistrictID: 1,
				},
				err: errs.ErrProvinceNotFound,
				beforeTests: func(mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockAddressRepo *mocks.AddressRepository) {
					mockSubdistrictService.On("GetSubdistrictByID", 1).Return(&model.Subdistrict{
						ID:         1,
						DistrictID: 1,
					}, nil)
					mockDistrictService.On("GetDistrictByID", 1).Return(&model.District{
						ID:     1,
						CityID: 1,
					}, nil)
					mockCityService.On("GetCityByID", 1).Return(&model.City{
						ID:         1,
						ProvinceID: 1,
					}, nil)
					mockProvinceService.On("GetProvinceByID", 1).Return(nil, errs.ErrProvinceNotFound)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrProvinceNotFound,
			},
		},
		{
			description: "should return error when AddUserAddress return error",
			input: input{
				data: &dto.AddressRequest{
					SubdistrictID: 1,
				},
				err: errs.ErrInternalServerError,
				beforeTests: func(mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockAddressRepo *mocks.AddressRepository) {
					mockSubdistrictService.On("GetSubdistrictByID", 1).Return(&model.Subdistrict{
						ID:         1,
						DistrictID: 1,
					}, nil)
					mockDistrictService.On("GetDistrictByID", 1).Return(&model.District{
						ID:     1,
						CityID: 1,
					}, nil)
					mockCityService.On("GetCityByID", 1).Return(&model.City{
						ID:         1,
						ProvinceID: 1,
					}, nil)
					mockProvinceService.On("GetProvinceByID", 1).Return(&model.Province{
						ID: 1,
					}, nil)
					mockAddressRepo.On("AddUserAddress", &model.UserAddress{
						SubdistrictID: 1,
						DistrictID:    1,
						CityID:        1,
						ProvinceID:    1,
						IsDefault:     &falseValue,
						IsPickup:      &falseValue,
					}).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInternalServerError,
			},
		},
		{
			description: "should return address and nil error when success",
			input: input{
				data: &dto.AddressRequest{
					SubdistrictID: 1,
				},
				err: nil,
				beforeTests: func(mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockAddressRepo *mocks.AddressRepository) {
					mockSubdistrictService.On("GetSubdistrictByID", 1).Return(&model.Subdistrict{
						ID:         1,
						DistrictID: 1,
					}, nil)
					mockDistrictService.On("GetDistrictByID", 1).Return(&model.District{
						ID:     1,
						CityID: 1,
					}, nil)
					mockCityService.On("GetCityByID", 1).Return(&model.City{
						ID:         1,
						ProvinceID: 1,
					}, nil)
					mockProvinceService.On("GetProvinceByID", 1).Return(&model.Province{
						ID: 1,
					}, nil)
					mockAddressRepo.On("AddUserAddress", &model.UserAddress{
						SubdistrictID: 1,
						DistrictID:    1,
						CityID:        1,
						ProvinceID:    1,
						IsDefault:     &falseValue,
						IsPickup:      &falseValue,
					}).Return(&model.UserAddress{
						SubdistrictID: 1,
						DistrictID:    1,
						CityID:        1,
						ProvinceID:    1,
						IsPickup:      &falseValue,
						IsDefault:     &falseValue,
					}, nil)
				},
			},
			expected: expected{
				data: &model.UserAddress{
					SubdistrictID: 1,
					DistrictID:    1,
					CityID:        1,
					ProvinceID:    1,
					IsDefault:     &falseValue,
					IsPickup:      &falseValue,
				},
				err: nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockSubdistrictService := mocks.NewSubdistrictService(t)
			mockDistrictService := mocks.NewDistrictService(t)
			mockCityService := mocks.NewCityService(t)
			mockProvinceService := mocks.NewProvinceService(t)
			mockAddressRepo := mocks.NewAddressRepository(t)
			c.beforeTests(mockSubdistrictService, mockDistrictService, mockCityService, mockProvinceService, mockAddressRepo)

			userAddressService := service.NewAddressService(&service.AddressSConfig{
				SubdistrictService: mockSubdistrictService,
				DistrictService:    mockDistrictService,
				CityService:        mockCityService,
				ProvinceService:    mockProvinceService,
				AddressRepo:        mockAddressRepo,
			})

			got, err := userAddressService.AddUserAddress(c.input.data)

			assert.Equal(t, c.expected.data, got)
			assert.ErrorIs(t, c.expected.err, err)
		})
	}

}

func TestGetAllUserAddress(t *testing.T) {
	var defaultAddressId int = 1
	var falseValue bool = false
	var trueValue bool = true

	type input struct {
		userId      int
		err         error
		beforeTests func(mockAddressRepo *mocks.AddressRepository, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService)
	}
	type expected struct {
		data []*model.UserAddress
		err  error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when GetProfile return error",
			input: input{
				userId: 1,
				err:    errs.ErrInternalServerError,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockUserProfileService.On("GetProfile", 1).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInternalServerError,
			},
		},
		{
			description: "should return error when FindShopByUserId return error",
			input: input{
				userId: 1,
				err:    errs.ErrInternalServerError,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						UserID: 1,
					}, nil)
					mockShopService.On("FindShopByUserId", 1).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInternalServerError,
			},
		},
		{
			description: "should return error when GetAllUserAddress return error",
			input: input{
				userId: 1,
				err:    errs.ErrInternalServerError,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						UserID: 1,
					}, nil)
					mockShopService.On("FindShopByUserId", 1).Return(nil, nil)
					mockAddressRepo.On("GetAllUserAddress", 1).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInternalServerError,
			},
		},
		{
			description: "should return address and nil error when success",
			input: input{
				userId: 1,
				err:    nil,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						UserID:           1,
						DefaultAddressID: &defaultAddressId,
					}, nil)
					mockShopService.On("FindShopByUserId", 1).Return(nil, nil)
					mockAddressRepo.On("GetAllUserAddress", 1).Return([]*model.UserAddress{
						{
							ID:            1,
							SubdistrictID: 1,
							DistrictID:    1,
							CityID:        1,
							ProvinceID:    1,
						},
					}, nil)
				},
			},
			expected: expected{
				data: []*model.UserAddress{
					{
						ID:            1,
						SubdistrictID: 1,
						DistrictID:    1,
						CityID:        1,
						ProvinceID:    1,
						IsDefault:     &trueValue,
						IsPickup:      &falseValue,
					},
				},
				err: nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockAddressRepo := mocks.NewAddressRepository(t)
			mockUserProfileService := mocks.NewUserProfileService(t)
			mockShopService := mocks.NewShopService(t)
			c.beforeTests(mockAddressRepo, mockUserProfileService, mockShopService)

			userAddressService := service.NewAddressService(&service.AddressSConfig{
				AddressRepo:        mockAddressRepo,
				UserProfileService: mockUserProfileService,
				ShopService:        mockShopService,
			})

			got, err := userAddressService.GetAllUserAddress(c.input.userId)

			assert.Equal(t, c.expected.data, got)
			assert.ErrorIs(t, c.expected.err, err)
		})
	}

}

func TestUpdateUserAddress(t *testing.T) {
	var defaultAddressId = 1
	var falseValue = false
	type input struct {
		addressId   int
		data        *dto.AddressRequest
		err         error
		beforeTests func(mockAddressRepo *mocks.AddressRepository, mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService)
	}
	type expected struct {
		data *model.UserAddress
		err  error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when GetUserAddressByIdAndUserId return error",
			input: input{
				addressId: 1,
				data: &dto.AddressRequest{
					SubdistrictID: 1,
					ID:            1,
					UserID:        1,
				},
				err: errs.ErrInternalServerError,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInternalServerError,
			},
		},
		{
			description: "should return error when GetProfile return error",
			input: input{
				addressId: 1,
				data: &dto.AddressRequest{
					SubdistrictID: 1,
					ID:            1,
					UserID:        1,
					IsDefault:     &falseValue,
				},
				err: errs.ErrInternalServerError,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInternalServerError,
			},
		},
		{
			description: "should return error ErrMustHaveAtLeastOneDefaultAddress when profile default address is not nil and profile default address is equal to address id",
			input: input{
				addressId: 1,
				data: &dto.AddressRequest{
					SubdistrictID: 1,
					ID:            1,
					UserID:        1,
					IsDefault:     &falseValue,
				},
				err: errs.ErrMustHaveAtLeastOneDefaultAddress,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						DefaultAddressID: &defaultAddressId,
					}, nil)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrMustHaveAtLeastOneDefaultAddress,
			},
		},
		{
			description: "should return error when FindShopByUserId return error",
			input: input{
				addressId: 1,
				data: &dto.AddressRequest{
					SubdistrictID: 1,
					ID:            1,
					UserID:        1,
					IsDefault:     &falseValue,
					IsPickup:      &falseValue,
				},
				err: errs.ErrInternalServerError,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						DefaultAddressID: nil,
					}, nil)
					mockShopService.On("FindShopByUserId", 1).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInternalServerError,
			},
		},
		{
			description: "should return error ErrMustHaveAtLeastOnePickupAddress when shop pickup address is not nil and shop pickup address is equal to address id",
			input: input{
				addressId: 1,
				data: &dto.AddressRequest{
					SubdistrictID: 1,
					ID:            1,
					UserID:        1,
					IsDefault:     &falseValue,
					IsPickup:      &falseValue,
				},
				err: errs.ErrMustHaveAtLeastOnePickupAddress,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						DefaultAddressID: nil,
					}, nil)
					mockShopService.On("FindShopByUserId", 1).Return(&shopModel.Shop{
						AddressID: 1,
					}, nil)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrMustHaveAtLeastOnePickupAddress,
			},
		},
		{
			description: "should return error when GetSubdistrictById return error",
			input: input{
				addressId: 1,
				data: &dto.AddressRequest{
					SubdistrictID: 1,
					ID:            1,
					UserID:        1,
					IsDefault:     &falseValue,
					IsPickup:      &falseValue,
				},
				err: errs.ErrInternalServerError,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						DefaultAddressID: nil,
					}, nil)
					mockShopService.On("FindShopByUserId", 1).Return(&shopModel.Shop{
						AddressID: 2,
					}, nil)
					mockSubdistrictService.On("GetSubdistrictByID", 1).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInternalServerError,
			},
		},
		{
			description: "should return error when GetDistrictByID return error",
			input: input{
				addressId: 1,
				data: &dto.AddressRequest{
					SubdistrictID: 1,
					ID:            1,
					UserID:        1,
					IsDefault:     &falseValue,
					IsPickup:      &falseValue,
				},
				err: errs.ErrInternalServerError,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						DefaultAddressID: nil,
					}, nil)
					mockShopService.On("FindShopByUserId", 1).Return(&shopModel.Shop{
						AddressID: 2,
					}, nil)
					mockSubdistrictService.On("GetSubdistrictByID", 1).Return(&model.Subdistrict{
						ID:         1,
						DistrictID: 1,
					}, nil)
					mockDistrictService.On("GetDistrictByID", 1).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInternalServerError,
			},
		},
		{
			description: "should return error when GetCityByID return error",
			input: input{
				addressId: 1,
				data: &dto.AddressRequest{
					SubdistrictID: 1,
					ID:            1,
					UserID:        1,
					IsDefault:     &falseValue,
					IsPickup:      &falseValue,
				},
				err: errs.ErrInternalServerError,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						DefaultAddressID: nil,
					}, nil)
					mockShopService.On("FindShopByUserId", 1).Return(&shopModel.Shop{
						AddressID: 2,
					}, nil)
					mockSubdistrictService.On("GetSubdistrictByID", 1).Return(&model.Subdistrict{
						ID:         1,
						DistrictID: 1,
					}, nil)
					mockDistrictService.On("GetDistrictByID", 1).Return(&model.District{
						ID:     1,
						CityID: 1,
					}, nil)
					mockCityService.On("GetCityByID", 1).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInternalServerError,
			},
		},
		{
			description: "should return error when GetProvinceByID return error",
			input: input{
				addressId: 1,
				data: &dto.AddressRequest{
					SubdistrictID: 1,
					ID:            1,
					UserID:        1,
					IsDefault:     &falseValue,
					IsPickup:      &falseValue,
				},
				err: errs.ErrInternalServerError,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						DefaultAddressID: nil,
					}, nil)
					mockShopService.On("FindShopByUserId", 1).Return(&shopModel.Shop{
						AddressID: 2,
					}, nil)
					mockSubdistrictService.On("GetSubdistrictByID", 1).Return(&model.Subdistrict{
						ID:         1,
						DistrictID: 1,
					}, nil)
					mockDistrictService.On("GetDistrictByID", 1).Return(&model.District{
						ID:     1,
						CityID: 1,
					}, nil)
					mockCityService.On("GetCityByID", 1).Return(&model.City{
						ID:         1,
						ProvinceID: 1,
					}, nil)
					mockProvinceService.On("GetProvinceByID", 1).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInternalServerError,
			},
		},
		{
			description: "should return error when UpdateUserAddress return error",
			input: input{
				addressId: 1,
				data: &dto.AddressRequest{
					SubdistrictID: 1,
					ID:            1,
					UserID:        1,
					IsDefault:     &falseValue,
					IsPickup:      &falseValue,
				},
				err: errs.ErrInternalServerError,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						DefaultAddressID: nil,
					}, nil)
					mockShopService.On("FindShopByUserId", 1).Return(&shopModel.Shop{
						AddressID: 2,
					}, nil)
					mockSubdistrictService.On("GetSubdistrictByID", 1).Return(&model.Subdistrict{
						ID:         1,
						DistrictID: 1,
					}, nil)
					mockDistrictService.On("GetDistrictByID", 1).Return(&model.District{
						ID:     1,
						CityID: 1,
					}, nil)
					mockCityService.On("GetCityByID", 1).Return(&model.City{
						ID:         1,
						ProvinceID: 1,
					}, nil)
					mockProvinceService.On("GetProvinceByID", 1).Return(&model.Province{
						ID: 1,
					}, nil)
					mockAddressRepo.On("UpdateUserAddress", &model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
						ProvinceID:    1,
						CityID:        1,
						DistrictID:    1,
						IsDefault:     &falseValue,
						IsPickup:      &falseValue,
					}).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInternalServerError,
			},
		},
		{
			description: "should return updated address when success",
			input: input{
				addressId: 1,
				data: &dto.AddressRequest{
					SubdistrictID: 1,
					ID:            1,
					UserID:        1,
					IsDefault:     &falseValue,
					IsPickup:      &falseValue,
				},
				err: nil,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						DefaultAddressID: nil,
					}, nil)
					mockShopService.On("FindShopByUserId", 1).Return(&shopModel.Shop{
						AddressID: 2,
					}, nil)
					mockSubdistrictService.On("GetSubdistrictByID", 1).Return(&model.Subdistrict{
						ID:         1,
						DistrictID: 1,
					}, nil)
					mockDistrictService.On("GetDistrictByID", 1).Return(&model.District{
						ID:     1,
						CityID: 1,
					}, nil)
					mockCityService.On("GetCityByID", 1).Return(&model.City{
						ID:         1,
						ProvinceID: 1,
					}, nil)
					mockProvinceService.On("GetProvinceByID", 1).Return(&model.Province{
						ID: 1,
					}, nil)
					mockAddressRepo.On("UpdateUserAddress", &model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
						ProvinceID:    1,
						CityID:        1,
						DistrictID:    1,
						IsDefault:     &falseValue,
						IsPickup:      &falseValue,
					}).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
						ProvinceID:    1,
						CityID:        1,
						DistrictID:    1,
					}, nil)
				},
			},
			expected: expected{
				data: &model.UserAddress{
					ID:            1,
					UserID:        1,
					SubdistrictID: 1,
					ProvinceID:    1,
					CityID:        1,
					DistrictID:    1,
				},
				err: nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockAddressRepo := mocks.NewAddressRepository(t)
			mockSubdistrictService := mocks.NewSubdistrictService(t)
			mockDistrictService := mocks.NewDistrictService(t)
			mockCityService := mocks.NewCityService(t)
			mockProvinceService := mocks.NewProvinceService(t)
			mockUserProfileService := mocks.NewUserProfileService(t)
			mockShopService := mocks.NewShopService(t)
			c.beforeTests(mockAddressRepo, mockSubdistrictService, mockDistrictService, mockCityService, mockProvinceService, mockUserProfileService, mockShopService)

			userAddressService := service.NewAddressService(&service.AddressSConfig{
				AddressRepo:        mockAddressRepo,
				SubdistrictService: mockSubdistrictService,
				DistrictService:    mockDistrictService,
				CityService:        mockCityService,
				ProvinceService:    mockProvinceService,
				UserProfileService: mockUserProfileService,
				ShopService:        mockShopService,
			})

			got, err := userAddressService.UpdateUserAddress(c.input.data)

			assert.Equal(t, c.expected.data, got)
			assert.ErrorIs(t, c.expected.err, err)
		})
	}

}

func TestDeleteUserAddress(t *testing.T) {
	var defaultAddressId = 1
	type input struct {
		addressId   int
		userId      int
		err         error
		beforeTests func(mockAddressRepo *mocks.AddressRepository, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService)
	}
	type expected struct {
		err error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when GetUserAddressByIdAndUserId return error",
			input: input{
				addressId: 1,
				userId:    1,
				err:       errs.ErrInternalServerError,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				err: errs.ErrInternalServerError,
			},
		},
		{
			description: "should return error when GetProfile return error",
			input: input{
				addressId: 1,
				userId:    1,
				err:       errs.ErrInternalServerError,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				err: errs.ErrInternalServerError,
			},
		},
		{
			description: "should return error when profile default address is equal to address id",
			input: input{
				addressId: 1,
				userId:    1,
				err:       errs.ErrMustHaveAtLeastOneDefaultAddress,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						DefaultAddressID: &defaultAddressId,
					}, nil)
				},
			},
			expected: expected{
				err: errs.ErrMustHaveAtLeastOneDefaultAddress,
			},
		},
		{
			description: "should return error when FindShopByUserId return error",
			input: input{
				addressId: 1,
				userId:    1,
				err:       errs.ErrInternalServerError,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						DefaultAddressID: nil,
					}, nil)
					mockShopService.On("FindShopByUserId", 1).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				err: errs.ErrInternalServerError,
			},
		},
		{
			description: "should return error ErrMustHaveAtLeastOnePickupAddress when shop pickup address is equal to address id",
			input: input{
				addressId: 1,
				userId:    1,
				err:       errs.ErrMustHaveAtLeastOnePickupAddress,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						DefaultAddressID: nil,
					}, nil)
					mockShopService.On("FindShopByUserId", 1).Return(&shopModel.Shop{
						AddressID: 1,
					}, nil)
				},
			},
			expected: expected{
				err: errs.ErrMustHaveAtLeastOnePickupAddress,
			},
		},
		{
			description: "should return error when DeleteUserAddress return error",
			input: input{
				addressId: 1,
				userId:    1,
				err:       errs.ErrInternalServerError,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						DefaultAddressID: nil,
					}, nil)
					mockShopService.On("FindShopByUserId", 1).Return(nil, nil)
					mockAddressRepo.On("DeleteUserAddress", 1, 1).Return(errs.ErrInternalServerError)
				},
			},
			expected: expected{
				err: errs.ErrInternalServerError,
			},
		},
		{
			description: "should return nil when success delete user address",
			input: input{
				addressId: 1,
				userId:    1,
				err:       nil,
				beforeTests: func(mockAddressRepo *mocks.AddressRepository, mockUserProfileService *mocks.UserProfileService, mockShopService *mocks.ShopService) {
					mockAddressRepo.On("GetUserAddressByIdAndUserId", 1, 1).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						SubdistrictID: 1,
					}, nil)
					mockUserProfileService.On("GetProfile", 1).Return(&userModel.UserProfile{
						DefaultAddressID: nil,
					}, nil)
					mockShopService.On("FindShopByUserId", 1).Return(nil, nil)
					mockAddressRepo.On("DeleteUserAddress", 1, 1).Return(nil)
				},
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockAddressRepo := mocks.NewAddressRepository(t)
			mockUserProfileService := mocks.NewUserProfileService(t)
			mockShopService := mocks.NewShopService(t)
			c.beforeTests(mockAddressRepo, mockUserProfileService, mockShopService)

			userAddressService := service.NewAddressService(&service.AddressSConfig{
				AddressRepo:        mockAddressRepo,
				UserProfileService: mockUserProfileService,
				ShopService:        mockShopService,
			})

			err := userAddressService.DeleteUserAddress(c.input.addressId, c.input.userId)

			assert.ErrorIs(t, c.expected.err, err)
		})
	}
}

package service_test

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/user/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddUserAddress(t *testing.T) {
	type input struct {
		data        *dto.AddAddressRequest
		err         error
		beforeTests func(mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserAddressRepo *mocks.UserAddressRepository)
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
				data: &dto.AddAddressRequest{
					SubdistrictID: 1,
				},
				err: nil,
				beforeTests: func(mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserAddressRepo *mocks.UserAddressRepository) {
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
				data: &dto.AddAddressRequest{
					SubdistrictID: 1,
				},
				err: errs.ErrDistrictNotFound,
				beforeTests: func(mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserAddressRepo *mocks.UserAddressRepository) {
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
				data: &dto.AddAddressRequest{
					SubdistrictID: 1,
				},
				err: errs.ErrCityNotFound,
				beforeTests: func(mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserAddressRepo *mocks.UserAddressRepository) {
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
				data: &dto.AddAddressRequest{
					SubdistrictID: 1,
				},
				err: errs.ErrProvinceNotFound,
				beforeTests: func(mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserAddressRepo *mocks.UserAddressRepository) {
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
				data: &dto.AddAddressRequest{
					SubdistrictID: 1,
				},
				err: errs.ErrInternalServerError,
				beforeTests: func(mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserAddressRepo *mocks.UserAddressRepository) {
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
					mockUserAddressRepo.On("AddUserAddress", &model.UserAddress{
						SubdistrictID: 1,
						DistrictID:    1,
						CityID:        1,
						ProvinceID:    1,
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
				data: &dto.AddAddressRequest{
					SubdistrictID: 1,
				},
				err: nil,
				beforeTests: func(mockSubdistrictService *mocks.SubdistrictService, mockDistrictService *mocks.DistrictService, mockCityService *mocks.CityService, mockProvinceService *mocks.ProvinceService, mockUserAddressRepo *mocks.UserAddressRepository) {
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
					mockUserAddressRepo.On("AddUserAddress", &model.UserAddress{
						SubdistrictID: 1,
						DistrictID:    1,
						CityID:        1,
						ProvinceID:    1,
					}).Return(&model.UserAddress{
						SubdistrictID: 1,
						DistrictID:    1,
						CityID:        1,
						ProvinceID:    1,
					}, nil)
				},
			},
			expected: expected{
				data: &model.UserAddress{
					SubdistrictID: 1,
					DistrictID:    1,
					CityID:        1,
					ProvinceID:    1,
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
			mockUserAddressRepo := mocks.NewUserAddressRepository(t)
			c.beforeTests(mockSubdistrictService, mockDistrictService, mockCityService, mockProvinceService, mockUserAddressRepo)

			userAddressService := service.NewUserAddressService(&service.UserAddressSConfig{
				SubdistrictService: mockSubdistrictService,
				DistrictService:    mockDistrictService,
				CityService:        mockCityService,
				ProvinceService:    mockProvinceService,
				UserAddressRepo:    mockUserAddressRepo,
			})

			got, err := userAddressService.AddUserAddress(c.input.data)

			assert.Equal(t, c.expected.data, got)
			assert.ErrorIs(t, c.expected.err, err)
		})
	}

}

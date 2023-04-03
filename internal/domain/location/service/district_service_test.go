package service_test

import (
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/location/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDistricts(t *testing.T) {
	type input struct {
		data        dto.GetDistrictsRequest
		err         error
		beforeTests func(mockDistrictRepo *mocks.DistrictRepository, mockCache *mocks.LocationCache)
	}
	type expected struct {
		data []*model.District
		err  error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return districts when cache is empty",
			input: input{
				data: dto.GetDistrictsRequest{
					CityID: 1,
				},
				err: nil,
				beforeTests: func(mockDistrictRepo *mocks.DistrictRepository, mockCache *mocks.LocationCache) {
					mockCache.On("GetDistricts", dto.GetDistrictsRequest{
						CityID: 1,
					}).Return(nil)

					mockDistrictRepo.On("GetAll", dto.GetDistrictsRequest{
						CityID: 1,
					}).Return([]*model.District{
						{
							ID:   1,
							Name: "test district",
						},
					}, nil)

					mockCache.On("StoreDistricts", dto.GetDistrictsRequest{
						CityID: 1,
					}, []*model.District{
						{
							ID:   1,
							Name: "test district",
						},
					})
				},
			},
			expected: expected{
				data: []*model.District{
					{
						ID:   1,
						Name: "test district",
					},
				},
				err: nil,
			},
		},
		{
			description: "should return districts when cache hit",
			input: input{
				data: dto.GetDistrictsRequest{
					CityID: 1,
				},
				err: nil,
				beforeTests: func(mockDistrictRepo *mocks.DistrictRepository, mockCache *mocks.LocationCache) {
					mockCache.On("GetDistricts", dto.GetDistrictsRequest{
						CityID: 1,
					}).Return([]*model.District{
						{
							ID:   1,
							Name: "test district",
						},
					})
				},
			},
			expected: expected{
				data: []*model.District{
					{
						ID:   1,
						Name: "test district",
					},
				},
				err: nil,
			},
		},
		{
			description: "should return error when district repo return error",
			input: input{
				data: dto.GetDistrictsRequest{
					CityID: 1,
				},
				err: nil,
				beforeTests: func(mockDistrictRepo *mocks.DistrictRepository, mockCache *mocks.LocationCache) {
					mockCache.On("GetDistricts", dto.GetDistrictsRequest{
						CityID: 1,
					}).Return(nil)

					mockDistrictRepo.On("GetAll", dto.GetDistrictsRequest{
						CityID: 1,
					}).Return(nil, assert.AnError)
				},
			},
			expected: expected{
				data: nil,
				err:  assert.AnError,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockDistrictRepo := mocks.NewDistrictRepository(t)
			mockCache := mocks.NewLocationCache(t)
			c.beforeTests(mockDistrictRepo, mockCache)
			districtService := service.NewDistrictService(&service.DistrictSConfig{
				DistrictRepo: mockDistrictRepo,
				Cache:        mockCache,
			})

			data, err := districtService.GetDistricts(c.input.data)

			assert.Equal(t, c.expected.err, err)
			assert.Equal(t, c.expected.data, data)
		})
	}

}

func TestGetDistrictByID(t *testing.T) {
	type input struct {
		data        int
		err         error
		beforeTests func(mockDistrictRepo *mocks.DistrictRepository)
	}
	type expected struct {
		data *model.District
		err  error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return district and error",
			input: input{
				data: 1,
				err:  nil,
				beforeTests: func(mockDistrictRepo *mocks.DistrictRepository) {
					mockDistrictRepo.On("GetByID", 1).Return(&model.District{
						ID: 1,
					}, nil)
				},
			},
			expected: expected{
				data: &model.District{
					ID: 1,
				},
				err: nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockDistrictRepo := mocks.NewDistrictRepository(t)
			c.beforeTests(mockDistrictRepo)

			districtService := service.NewDistrictService(&service.DistrictSConfig{
				DistrictRepo: mockDistrictRepo,
			})

			got, err := districtService.GetDistrictByID(c.input.data)

			assert.Equal(t, c.expected.data, got)
			assert.ErrorIs(t, c.expected.err, err)
		})
	}
}

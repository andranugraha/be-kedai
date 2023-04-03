package service_test

import (
	"kedai/backend/be-kedai/internal/common/dto"
	errorResponse "kedai/backend/be-kedai/internal/common/error"
	locationDto "kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/location/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCities(t *testing.T) {
	var (
		req = locationDto.GetCitiesRequest{
			Limit: 10,
			Page:  1,
		}
		res = &dto.PaginationResponse{
			Data: []*model.City{
				{
					ID:         1,
					ProvinceID: 1,
					Name:       "Kota Jakarta Pusat",
				},
			},
			Limit:      10,
			Page:       1,
			TotalRows:  1,
			TotalPages: 1,
		}
	)
	tests := []struct {
		name       string
		want       *dto.PaginationResponse
		wantErr    error
		beforeTest func(mockCityRepo *mocks.CityRepository, mockLocationCache *mocks.LocationCache)
	}{
		{
			name:    "should return cities with pagination when get all success",
			want:    res,
			wantErr: nil,
			beforeTest: func(mockCityRepo *mocks.CityRepository, mockLocationCache *mocks.LocationCache) {
				mockLocationCache.On("GetCities", req).Return(nil)
				mockCityRepo.On("GetAll", req).Return(res.Data, res.TotalRows, res.TotalPages, nil)
				mockLocationCache.On("StoreCities", req, res)
			},
		},
		{
			name:    "should return cities with pagination when cache hit",
			want:    res,
			wantErr: nil,
			beforeTest: func(mockCityRepo *mocks.CityRepository, mockLocationCache *mocks.LocationCache) {
				mockLocationCache.On("GetCities", req).Return(res)
			},
		},
		{
			name:    "should return error when get all failed",
			want:    nil,
			wantErr: errorResponse.ErrInternalServerError,
			beforeTest: func(mockCityRepo *mocks.CityRepository, mockLocationCache *mocks.LocationCache) {
				mockLocationCache.On("GetCities", req).Return(nil)
				mockCityRepo.On("GetAll", req).Return(nil, int64(0), 0, errorResponse.ErrInternalServerError)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := mocks.NewCityRepository(t)
			mockCache := mocks.NewLocationCache(t)
			test.beforeTest(mockRepo, mockCache)
			cityService := service.NewCityService(&service.CitySConfig{
				CityRepo: mockRepo,
				Cache:    mockCache,
			})

			got, err := cityService.GetCities(req)

			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, test.wantErr, err)
		})
	}
}

func TestGetCityByID(t *testing.T) {
	type input struct {
		data        int
		err         error
		beforeTests func(mockCityRepo *mocks.CityRepository)
	}
	type expected struct {
		data *model.City
		err  error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return city and error",
			input: input{
				data: 1,
				err:  nil,
				beforeTests: func(mockCityRepo *mocks.CityRepository) {
					mockCityRepo.On("GetByID", 1).Return(&model.City{
						ID: 1,
					}, nil)
				},
			},
			expected: expected{
				data: &model.City{
					ID: 1,
				},
				err: nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockCityRepo := mocks.NewCityRepository(t)
			c.beforeTests(mockCityRepo)

			cityService := service.NewCityService(&service.CitySConfig{
				CityRepo: mockCityRepo,
			})

			got, err := cityService.GetCityByID(c.input.data)

			assert.Equal(t, c.expected.data, got)
			assert.ErrorIs(t, c.expected.err, err)
		})
	}
}

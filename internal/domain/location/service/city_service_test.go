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
	tests := []struct {
		name               string
		request            locationDto.GetCitiesRequest
		wantGetAllResponse *dto.PaginationResponse
		want               *dto.PaginationResponse
		wantErr            error
	}{
		{
			name: "should return cities with pagination when get all success",
			request: locationDto.GetCitiesRequest{
				Limit: 10,
				Page:  1,
			},
			wantGetAllResponse: &dto.PaginationResponse{
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
			},
			want: &dto.PaginationResponse{
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
			},
			wantErr: nil,
		},
		{
			name: "should return error when get all failed",
			request: locationDto.GetCitiesRequest{
				Limit: 10,
				Page:  1,
			},
			wantGetAllResponse: &dto.PaginationResponse{
				Data:       []*model.City{},
				TotalRows:  0,
				TotalPages: 0,
			},
			want:    nil,
			wantErr: errorResponse.ErrInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := mocks.NewCityRepository(t)
			mockRepo.On("GetAll", test.request).Return(test.wantGetAllResponse.Data, test.wantGetAllResponse.TotalRows, test.wantGetAllResponse.TotalPages, test.wantErr)
			cityService := service.NewCityService(&service.CitySConfig{
				CityRepo: mockRepo,
			})

			got, err := cityService.GetCities(test.request)

			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, test.wantErr, err)
		})
	}
}

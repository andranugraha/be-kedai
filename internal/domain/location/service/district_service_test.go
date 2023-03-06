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
		beforeTests func(mockDistrictRepo *mocks.DistrictRepository)
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
			description: "should return districts and error",
			input: input{
				data: dto.GetDistrictsRequest{
					CityID: 1,
				},
				err: nil,
				beforeTests: func(mockDistrictRepo *mocks.DistrictRepository) {
					mockDistrictRepo.On("GetAll", dto.GetDistrictsRequest{
						CityID: 1,
					}).Return([]*model.District{
						{
							ID: 1,
						},
					}, nil)
				},
			},
			expected: expected{
				data: []*model.District{
					{
						ID: 1,
					},
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

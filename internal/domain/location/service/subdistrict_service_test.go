package service_test

import (
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/location/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSubdistrictByID(t *testing.T) {
	type input struct {
		data        int
		err         error
		beforeTests func(mockSubdistrictRepo *mocks.SubdistrictRepository)
	}
	type expected struct {
		data *model.Subdistrict
		err  error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return subdistrict and error",
			input: input{
				data: 1,
				err:  nil,
				beforeTests: func(mockSubdistrictRepo *mocks.SubdistrictRepository) {
					mockSubdistrictRepo.On("GetByID", 1).Return(&model.Subdistrict{
						ID: 1,
					}, nil)
				},
			},
			expected: expected{
				data: &model.Subdistrict{
					ID: 1,
				},
				err: nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockSubdistrictRepo := mocks.NewSubdistrictRepository(t)
			c.beforeTests(mockSubdistrictRepo)

			subdistrictService := service.NewSubdistrictService(&service.SubdistrictSConfig{
				SubdistrictRepo: mockSubdistrictRepo,
			})

			got, err := subdistrictService.GetSubdistrictByID(c.input.data)

			assert.Equal(t, c.expected.data, got)
			assert.ErrorIs(t, c.expected.err, err)
		})
	}

}

func TestGetSubdistricts(t *testing.T) {
	type input struct {
		data        dto.GetSubdistrictsRequest
		err         error
		beforeTests func(mockSubdistrictRepo *mocks.SubdistrictRepository)
	}
	type expected struct {
		data []*model.Subdistrict
		err  error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return subdistricts and error",
			input: input{
				data: dto.GetSubdistrictsRequest{
					DistrictID: 1,
				},
				err: nil,
				beforeTests: func(mockSubdistrictRepo *mocks.SubdistrictRepository) {
					mockSubdistrictRepo.On("GetAll", dto.GetSubdistrictsRequest{
						DistrictID: 1,
					}).Return([]*model.Subdistrict{}, nil)
				},
			},
			expected: expected{
				data: []*model.Subdistrict{},
				err:  nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockSubdistrictRepo := mocks.NewSubdistrictRepository(t)
			c.beforeTests(mockSubdistrictRepo)

			subdistrictService := service.NewSubdistrictService(&service.SubdistrictSConfig{
				SubdistrictRepo: mockSubdistrictRepo,
			})

			got, err := subdistrictService.GetSubdistricts(c.input.data)

			assert.Equal(t, c.expected.data, got)
			assert.ErrorIs(t, c.expected.err, err)
		})
	}

}

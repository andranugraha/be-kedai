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
		beforeTests func(mockSubdistrictRepo *mocks.SubdistrictRepository, mockCache *mocks.LocationCache)
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
			description: "should return subdistricts when cache is empty",
			input: input{
				data: dto.GetSubdistrictsRequest{
					DistrictID: 1,
				},
				err: nil,
				beforeTests: func(mockSubdistrictRepo *mocks.SubdistrictRepository, mockCache *mocks.LocationCache) {
					mockCache.On("GetSubdistricts", dto.GetSubdistrictsRequest{
						DistrictID: 1,
					}).Return(nil)

					mockSubdistrictRepo.On("GetAll", dto.GetSubdistrictsRequest{
						DistrictID: 1,
					}).Return([]*model.Subdistrict{
						{
							ID: 1,
						},
					}, nil)

					mockCache.On("StoreSubdistricts", dto.GetSubdistrictsRequest{
						DistrictID: 1,
					}, []*model.Subdistrict{
						{
							ID: 1,
						},
					})
				},
			},
			expected: expected{
				data: []*model.Subdistrict{
					{
						ID: 1,
					},
				},
				err: nil,
			},
		},
		{
			description: "should return subdistricts when cache hit",
			input: input{
				data: dto.GetSubdistrictsRequest{
					DistrictID: 1,
				},
				err: nil,
				beforeTests: func(mockSubdistrictRepo *mocks.SubdistrictRepository, mockCache *mocks.LocationCache) {
					mockCache.On("GetSubdistricts", dto.GetSubdistrictsRequest{
						DistrictID: 1,
					}).Return([]*model.Subdistrict{
						{
							ID: 1,
						},
					})
				},
			},
			expected: expected{
				data: []*model.Subdistrict{
					{
						ID: 1,
					},
				},
				err: nil,
			},
		},
		{
			description: "should return error when cache is empty and subdistrict repo return error",
			input: input{
				data: dto.GetSubdistrictsRequest{
					DistrictID: 1,
				},
				err: nil,
				beforeTests: func(mockSubdistrictRepo *mocks.SubdistrictRepository, mockCache *mocks.LocationCache) {
					mockCache.On("GetSubdistricts", dto.GetSubdistrictsRequest{
						DistrictID: 1,
					}).Return(nil)

					mockSubdistrictRepo.On("GetAll", dto.GetSubdistrictsRequest{
						DistrictID: 1,
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
			mockSubdistrictRepo := mocks.NewSubdistrictRepository(t)
			mockCache := mocks.NewLocationCache(t)
			c.beforeTests(mockSubdistrictRepo, mockCache)
			subdistrictService := service.NewSubdistrictService(&service.SubdistrictSConfig{
				SubdistrictRepo: mockSubdistrictRepo,
				Cache:           mockCache,
			})

			got, err := subdistrictService.GetSubdistricts(c.input.data)

			assert.Equal(t, c.expected.data, got)
			assert.ErrorIs(t, c.expected.err, err)
		})
	}

}

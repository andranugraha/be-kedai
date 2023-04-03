package service_test

import (
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/location/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProvinces(t *testing.T) {
	tests := []struct {
		name       string
		want       []*model.Province
		wantErr    error
		beforeTest func(mockProvinceRepo *mocks.ProvinceRepository, mockLocationCache *mocks.LocationCache)
	}{
		{
			name: "should return provinces when get all success",
			want: []*model.Province{
				{
					ID:   1,
					Name: "DKI Jakarta",
				},
			},
			wantErr: nil,
			beforeTest: func(mockProvinceRepo *mocks.ProvinceRepository, mockLocationCache *mocks.LocationCache) {
				mockLocationCache.On("GetProvinces").Return(nil)
				mockProvinceRepo.On("GetAll").Return([]*model.Province{
					{
						ID:   1,
						Name: "DKI Jakarta",
					},
				}, nil)
				mockLocationCache.On("StoreProvinces", []*model.Province{
					{
						ID:   1,
						Name: "DKI Jakarta",
					},
				})
			},
		},
		{
			name: "should return provinces when cache hit",
			want: []*model.Province{
				{
					ID:   1,
					Name: "DKI Jakarta",
				},
			},
			wantErr: nil,
			beforeTest: func(mockProvinceRepo *mocks.ProvinceRepository, mockLocationCache *mocks.LocationCache) {
				mockLocationCache.On("GetProvinces").Return([]*model.Province{
					{
						ID:   1,
						Name: "DKI Jakarta",
					},
				})
			},
		},
		{
			name:    "should return error when get all failed",
			want:    nil,
			wantErr: assert.AnError,
			beforeTest: func(mockProvinceRepo *mocks.ProvinceRepository, mockLocationCache *mocks.LocationCache) {
				mockLocationCache.On("GetProvinces").Return(nil)
				mockProvinceRepo.On("GetAll").Return(nil, assert.AnError)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockProvinceRepo := mocks.NewProvinceRepository(t)
			mockLocationCache := mocks.NewLocationCache(t)
			test.beforeTest(mockProvinceRepo, mockLocationCache)
			provinceService := service.NewProvinceService(&service.ProvinceSConfig{
				ProvinceRepo: mockProvinceRepo,
				Cache:        mockLocationCache,
			})

			got, err := provinceService.GetProvinces()

			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, test.wantErr, err)
		})
	}
}

func TestGetProvinceByID(t *testing.T) {
	type input struct {
		data        int
		err         error
		beforeTests func(mockProvinceRepo *mocks.ProvinceRepository)
	}
	type expected struct {
		data *model.Province
		err  error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return province and error",
			input: input{
				data: 1,
				err:  nil,
				beforeTests: func(mockProvinceRepo *mocks.ProvinceRepository) {
					mockProvinceRepo.On("GetByID", 1).Return(&model.Province{
						ID: 1,
					}, nil)
				},
			},
			expected: expected{
				data: &model.Province{
					ID: 1,
				},
				err: nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockProvinceRepo := mocks.NewProvinceRepository(t)
			c.beforeTests(mockProvinceRepo)

			provinceService := service.NewProvinceService(&service.ProvinceSConfig{
				ProvinceRepo: mockProvinceRepo,
			})

			got, err := provinceService.GetProvinceByID(c.input.data)

			assert.Equal(t, c.expected.data, got)
			assert.ErrorIs(t, c.expected.err, err)
		})
	}
}

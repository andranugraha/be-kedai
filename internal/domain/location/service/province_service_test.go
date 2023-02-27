package service_test

import (
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/location/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

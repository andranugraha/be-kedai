package service_test

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/product/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSKUByID(t *testing.T) {
	type input struct {
		skuId      int
		mockReturn *model.Sku
		mockErr    error
	}
	type expected struct {
		data *model.Sku
		err  error
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when failed to get sku",
			input: input{
				skuId:      1,
				mockReturn: nil,
				mockErr:    errors.New("failed to get sku"),
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get sku"),
			},
		},
		{
			description: "should return sku when fecthing sku succeed",
			input: input{
				skuId:      1,
				mockReturn: &model.Sku{},
				mockErr:    nil,
			},
			expected: expected{
				data: &model.Sku{},
				err:  nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			skuRepo := mocks.NewSkuRepository(t)
			skuRepo.On("GetByID", tc.input.skuId).Return(tc.input.mockReturn, tc.input.mockErr)
			skuService := service.NewSkuService(&service.SkuSConfig{
				SkuRepository: skuRepo,
			})

			actualData, actualErr := skuService.GetByID(tc.input.skuId)

			assert.Equal(t, tc.expected.data, actualData)
			assert.Equal(t, tc.expected.err, actualErr)
		})
	}
}

func TestGetSKUByVariantIDs(t *testing.T) {
	type input struct {
		request    *dto.GetSKURequest
		beforeTest func(*mocks.SkuRepository)
	}
	type expected struct {
		data *model.Sku
		err  error
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when given invalid request",
			input: input{
				request: &dto.GetSKURequest{
					VariantID: "a,b",
				},
				beforeTest: func(sr *mocks.SkuRepository) {},
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInvalidVariantID,
			},
		},
		{
			description: "should return error when failed to fetch sku",
			input: input{
				request: &dto.GetSKURequest{
					VariantID: "1,3",
				},
				beforeTest: func(sr *mocks.SkuRepository) {
					sr.On("GetByVariantIDs", []int{1, 3}).Return(nil, errors.New("failed to fetch sku"))
				},
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to fetch sku"),
			},
		},
		{
			description: "should return sku when fetching sku succeed",
			input: input{
				request: &dto.GetSKURequest{
					VariantID: "1,3",
				},
				beforeTest: func(sr *mocks.SkuRepository) {
					sr.On("GetByVariantIDs", []int{1, 3}).Return(&model.Sku{}, nil)
				},
			},
			expected: expected{
				data: &model.Sku{},
				err:  nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			skuRepo := mocks.NewSkuRepository(t)
			tc.beforeTest(skuRepo)
			skuService := service.NewSkuService(&service.SkuSConfig{
				SkuRepository: skuRepo,
			})

			actualData, actualErr := skuService.GetSKUByVariantIDs(tc.input.request)

			assert.Equal(t, tc.expected.data, actualData)
			assert.Equal(t, tc.expected.err, actualErr)
		})
	}
}

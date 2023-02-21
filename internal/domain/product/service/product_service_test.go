package service_test

import (
	"errors"
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/domain/product/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetByCodeFull(t *testing.T) {
	type input struct {
		code string
		data *model.Product
		err  error
	}
	type expected struct {
		user *model.Product
		err  error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "it should return user data if product exists",
			input: input{
				code: "PRODUCT_CODE_A",
				data: &model.Product{
					ID:   1,
					Code: "PRODUCT_CODE_A",
				},
				err: nil,
			},
			expected: expected{
				user: &model.Product{
					ID:   1,
					Code: "PRODUCT_CODE_A",
				},
				err: nil,
			},
		},
		{
			description: "it should return error if failed to get product",
			input: input{
				code: "INVALID_CODE",
				data: nil,
				err:  errors.New("failed to get product"),
			},
			expected: expected{
				user: nil,
				err:  errors.New("failed to get product"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := mocks.NewProductRepository(t)
			mockRepo.On("GetByCodeFull", tc.input.code).Return(tc.input.data, tc.input.err)
			uc := service.NewProductService(&service.ProductSConfig{
				Repository: mockRepo,
			})

			actualUser, actualErr := uc.GetByCodeFull(tc.input.code)

			assert.Equal(t, tc.expected.user, actualUser)
			assert.Equal(t, actualErr, tc.expected.err)
		})
	}
}

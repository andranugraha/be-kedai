package service_test

import (
	"errors"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCouriersByProductID(t *testing.T) {
	type input struct {
		productID  int
		mockReturn []*model.Courier
		mockErr    error
	}
	type expected struct {
		data []*model.Courier
		err  error
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when failed to get couriers",
			input: input{
				productID:  1,
				mockReturn: nil,
				mockErr:    errors.New("failed to get couriers"),
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get couriers"),
			},
		},
		{
			description: "should return error when fetching couriers succeed",
			input: input{
				productID:  1,
				mockReturn: []*model.Courier{},
				mockErr:    nil,
			},
			expected: expected{
				data: []*model.Courier{},
				err:  nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			courierRepo := mocks.NewCourierRepository(t)
			courierRepo.On("GetByProductID", tc.input.productID).Return(tc.input.mockReturn, tc.input.mockErr)
			courierService := service.NewCourierService(&service.CourierSConfig{
				CourierRepository: courierRepo,
			})

			actualData, actualErr := courierService.GetCouriersByProductID(tc.input.productID)

			assert.Equal(t, tc.expected.data, actualData)
			assert.Equal(t, tc.expected.err, actualErr)
		})
	}
}

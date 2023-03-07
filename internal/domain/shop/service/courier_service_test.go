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

func TestGetAllCouriers(t *testing.T) {
	type input struct {
		err error
	}
	type expected struct {
		result []*model.Courier
		err error
	}
	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return result and error when called",
			input: input{
				err: nil,
			},
			expected: expected{
				result: []*model.Courier{},
				err: nil,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			courierRepo := mocks.NewCourierRepository(t)
			courierRepo.On("GetAll").Return(tc.result, nil)
			courierService := service.NewCourierService(&service.CourierSConfig{
				CourierRepository: courierRepo,
			})

			result, err := courierService.GetAllCouriers()

			assert.Equal(t, tc.expected.result, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

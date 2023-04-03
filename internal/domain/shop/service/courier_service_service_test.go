package service_test

import (
	"errors"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCourierServicesByCourierIDs(t *testing.T) {
	type input struct {
		courierIDs []int
		mockData   []*model.CourierService
		mockErr    error
	}
	type expected struct {
		data []*model.CourierService
		err  error
	}

	var (
		courierIDs = []int{1, 2}
	)

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when failed to get courier services",
			input: input{
				courierIDs: courierIDs,
				mockData:   nil,
				mockErr:    errors.New("failed to get courier services"),
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get courier services"),
			},
		},
		{
			description: "should return courier services data when succeed to get courier services",
			input: input{
				courierIDs: courierIDs,
				mockData:   []*model.CourierService{},
				mockErr:    nil,
			},
			expected: expected{
				data: []*model.CourierService{},
				err:  nil,
			},
		},
	}

	for _, tc := range tests {
		courierServiceRepo := mocks.NewCourierServiceRepository(t)
		courierServiceRepo.On("GetByCourierIDs", tc.input.courierIDs).Return(tc.input.mockData, tc.input.mockErr)
		courierServiceService := service.NewCourierServiceService(&service.CourierServiceSConfig{
			CourierServiceRepository: courierServiceRepo,
		})

		data, err := courierServiceService.GetCourierServicesByCourierIDs(tc.input.courierIDs)

		assert.Equal(t, tc.expected.data, data)
		assert.Equal(t, tc.expected.err, err)
	}
}

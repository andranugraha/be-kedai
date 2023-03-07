package service_test

import (
	"errors"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInvoicesByUserID(t *testing.T) {
	type input struct {
		userID         int
		request        *dto.InvoicePerShopFilterRequest
		mockData       []*dto.InvoicePerShopDetail
		mockTotalRows  int64
		mockTotalPages int
		mockErr        error
	}
	type expected struct {
		data *commonDto.PaginationResponse
		err  error
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when failed to get invoices",
			input: input{
				userID:         1,
				request:        &dto.InvoicePerShopFilterRequest{Limit: 10, Page: 1},
				mockData:       nil,
				mockTotalRows:  0,
				mockTotalPages: 0,
				mockErr:        errors.New("failed to return invoices"),
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to return invoices"),
			},
		},
		{
			description: "should return invoices data when succeed getting invoices",
			input: input{
				userID:         1,
				request:        &dto.InvoicePerShopFilterRequest{Limit: 10, Page: 1},
				mockData:       []*dto.InvoicePerShopDetail{},
				mockTotalRows:  0,
				mockTotalPages: 0,
				mockErr:        nil,
			},
			expected: expected{
				data: &commonDto.PaginationResponse{
					Limit:      10,
					Page:       1,
					TotalRows:  0,
					TotalPages: 0,
					Data:       []*dto.InvoicePerShopDetail{},
				},
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			invoicePerShopRepo := mocks.NewInvoicePerShopRepository(t)
			invoicePerShopRepo.On("GetByUserID", tc.input.userID, tc.input.request).Return(tc.input.mockData, tc.input.mockTotalRows, tc.input.mockTotalPages, tc.input.mockErr)
			invoicePerShopService := service.NewInvoicePerShopService(&service.InvoicePerShopSConfig{
				InvoicePerShopRepo: invoicePerShopRepo,
			})

			actualData, actualErr := invoicePerShopService.GetInvoicesByUserID(tc.input.userID, tc.input.request)

			assert.Equal(t, tc.expected.data, actualData)
			assert.Equal(t, tc.expected.err, actualErr)
		})
	}
}

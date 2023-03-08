package service_test

import (
	"errors"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/order/dto"

	"kedai/backend/be-kedai/internal/domain/order/model"
	"kedai/backend/be-kedai/internal/domain/order/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	commonErr "kedai/backend/be-kedai/internal/common/error"

	"github.com/stretchr/testify/assert"
)

func Test_InvoicePerShopGetByID(t *testing.T) {
	type input struct {
		req        int
		beforeTest func(mockInvoicePerShopRepo *mocks.InvoicePerShopRepository)
	}

	type expected struct {
		data *model.InvoicePerShop
		err  error
	}

	cases := []struct {
		description string
		input       input
		expected    expected
	}{
		{
			description: "should return error when GetByID return error",
			input: input{
				req: 1,
				beforeTest: func(mockInvoicePerShopRepo *mocks.InvoicePerShopRepository) {
					mockInvoicePerShopRepo.On("GetByID", 1).Return(nil, commonErr.ErrInvoiceNotFound)
				},
			},
			expected: expected{
				data: nil,
				err:  commonErr.ErrInvoiceNotFound,
			},
		},
		{
			description: "should return invoice per shop when GetByID return invoice per shop",
			input: input{
				req: 1,
				beforeTest: func(mockInvoicePerShopRepo *mocks.InvoicePerShopRepository) {
					mockInvoicePerShopRepo.On("GetByID", 1).Return(&model.InvoicePerShop{}, nil)
				},
			},
			expected: expected{
				data: &model.InvoicePerShop{},
				err:  nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockInvoicePerShopRepo := new(mocks.InvoicePerShopRepository)
			c.input.beforeTest(mockInvoicePerShopRepo)

			service := service.NewInvoicePerShopService(&service.InvoicePerShopSConfig{
				InvoicePerShopRepo: mockInvoicePerShopRepo,
			})

			data, err := service.GetByID(c.input.req)

			assert.Equal(t, c.expected.data, data)
			assert.Equal(t, c.expected.err, err)
		})
	}

}

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

func TestGetInvoicesByUserIDAndCode(t *testing.T) {
	type input struct {
		userID   int
		code     string
		mockData *dto.InvoicePerShopDetail
		mockErr  error
	}
	type expected struct {
		data *dto.InvoicePerShopDetail
		err  error
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when failed to get invoice",
			input: input{
				userID:   1,
				code:     "INV/XX/X",
				mockData: nil,
				mockErr:  errors.New("failed to get invoice"),
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get invoice"),
			},
		},
		{
			description: "should return invoice data when fetching succeed",
			input: input{
				userID:   1,
				code:     "INV/XX/X",
				mockData: &dto.InvoicePerShopDetail{},
				mockErr:  nil,
			},
			expected: expected{
				data: &dto.InvoicePerShopDetail{},
				err:  nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			invoicePerShopRepo := mocks.NewInvoicePerShopRepository(t)
			invoicePerShopRepo.On("GetByUserIDAndCode", tc.input.userID, tc.input.code).Return(tc.input.mockData, tc.input.mockErr)
			invoicePerShopService := service.NewInvoicePerShopService(&service.InvoicePerShopSConfig{
				InvoicePerShopRepo: invoicePerShopRepo,
			})

			actualData, actualErr := invoicePerShopService.GetInvoicesByUserIDAndCode(tc.input.userID, tc.input.code)

			assert.Equal(t, tc.expected.data, actualData)
			assert.Equal(t, tc.expected.err, actualErr)
		})
	}
}

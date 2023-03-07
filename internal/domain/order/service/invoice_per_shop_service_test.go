package service_test

import (
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

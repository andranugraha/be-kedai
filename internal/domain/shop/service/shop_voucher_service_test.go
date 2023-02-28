package service_test

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	var (
		shopId  = 1
		voucher = []*model.ShopVoucher{}
	)
	type input struct {
		shopId int
		err    error
	}
	type expected struct {
		result []*model.ShopVoucher
		err    error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return list of shop ticket when success",
			input: input{
				shopId: shopId,
				err:    nil,
			},
			expected: expected{
				result: voucher,
				err:    nil,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				shopId: shopId,
				err:    errs.ErrInternalServerError,
			},
			expected: expected{
				result: nil,
				err:    errs.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.ShopVoucherRepository)
			mockRepo.On("GetShopVoucher", tc.input.shopId).Return(tc.expected.result, tc.input.err)
			service := service.NewShopVoucherService(&service.ShopVoucherSConfig{
				ShopVoucherRepository: mockRepo,
			})

			result, err := service.GetShopVoucher(tc.input.shopId)

			assert.Equal(t, tc.expected.result, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

package service_test

import (
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"kedai/backend/be-kedai/internal/domain/marketplace/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMarketplaceVoucher(t *testing.T) {
	type input struct {
		req        *dto.GetMarketplaceVoucherRequest
		err        error
		beforeTest func(*mocks.MarketplaceVoucherRepository)
	}
	type expected struct {
		result []*model.MarketplaceVoucher
		err    error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error and marketplace vouchers",
			input: input{
				err: nil,
				req: &dto.GetMarketplaceVoucherRequest{},
				beforeTest: func(m *mocks.MarketplaceVoucherRepository) {
					m.On("GetMarketplaceVoucher", &dto.GetMarketplaceVoucherRequest{}).Return(nil, nil)
				},
			},
			expected: expected{
				result: nil,
				err:    nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			m := mocks.NewMarketplaceVoucherRepository(t)
			c.beforeTest(m)

			s := service.NewMarketplaceVoucherService(&service.MarketplaceVoucherSConfig{
				MarketplaceVoucherRepository: m,
			})

			result, err := s.GetMarketplaceVoucher(c.input.req)

			assert.Equal(t, c.expected.err, err)
			assert.Equal(t, c.expected.result, result)

		})
	}

}

func TestGetMarketplaceVoucherAdmin(t *testing.T) {
	type input struct {
		req        *dto.GetMarketplaceVoucherRequest
		err        error
		beforeTest func(*mocks.MarketplaceVoucherRepository)
	}
	type expected struct {
		result []*model.MarketplaceVoucher
		err    error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error and marketplace vouchers",
			input: input{
				err: nil,
				req: &dto.GetMarketplaceVoucherRequest{},
				beforeTest: func(m *mocks.MarketplaceVoucherRepository) {
					m.On("GetMarketplaceVoucherAdmin", &dto.GetMarketplaceVoucherRequest{}).Return(nil, nil)
				},
			},
			expected: expected{
				result: nil,
				err:    nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			m := mocks.NewMarketplaceVoucherRepository(t)
			c.beforeTest(m)

			s := service.NewMarketplaceVoucherService(&service.MarketplaceVoucherSConfig{
				MarketplaceVoucherRepository: m,
			})

			result, err := s.GetMarketplaceVoucherAdmin(c.input.req)

			assert.Equal(t, c.expected.err, err)
			assert.Equal(t, c.expected.result, result)

		})
	}

}

func TestGetValidByUserID(t *testing.T) {
	type input struct {
		req        *dto.GetMarketplaceVoucherRequest
		err        error
		beforeTest func(*mocks.MarketplaceVoucherRepository)
	}
	type expected struct {
		result []*model.MarketplaceVoucher
		err    error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error and marketplace vouchers",
			input: input{
				err: nil,
				req: &dto.GetMarketplaceVoucherRequest{},
				beforeTest: func(m *mocks.MarketplaceVoucherRepository) {
					m.On("GetValidByUserID", &dto.GetMarketplaceVoucherRequest{}).Return(nil, nil)
				},
			},
			expected: expected{
				result: nil,
				err:    nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			m := mocks.NewMarketplaceVoucherRepository(t)
			c.beforeTest(m)

			s := service.NewMarketplaceVoucherService(&service.MarketplaceVoucherSConfig{
				MarketplaceVoucherRepository: m,
			})

			result, err := s.GetValidByUserID(c.input.req)

			assert.Equal(t, c.expected.err, err)
			assert.Equal(t, c.expected.result, result)

		})
	}
}

func TestGetValidForCheckout(t *testing.T) {
	type input struct {
		id, userID, PaymentMethodID int
		err                         error
		beforeTest                  func(*mocks.MarketplaceVoucherRepository)
	}
	type expected struct {
		result *model.MarketplaceVoucher
		err    error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return valid marketplace voucher when success",
			input: input{
				err:             nil,
				id:              1,
				userID:          1,
				PaymentMethodID: 1,
				beforeTest: func(m *mocks.MarketplaceVoucherRepository) {
					m.On("GetValid", 1, 1, 1).Return(nil, nil)
				},
			},
			expected: expected{
				result: nil,
				err:    nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			m := mocks.NewMarketplaceVoucherRepository(t)
			c.beforeTest(m)

			s := service.NewMarketplaceVoucherService(&service.MarketplaceVoucherSConfig{
				MarketplaceVoucherRepository: m,
			})

			result, err := s.GetValidForCheckout(c.input.id, c.input.userID, c.input.PaymentMethodID)

			assert.Equal(t, c.expected.err, err)
			assert.Equal(t, c.expected.result, result)

		})
	}
}

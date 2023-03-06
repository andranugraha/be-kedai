package service_test

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetShopVoucher(t *testing.T) {
	var (
		slug    = "shop"
		voucher = []*model.ShopVoucher{}
		shop    = &model.Shop{}
	)
	type input struct {
		slug       string
		err        error
		beforeTest func(*mocks.ShopService, *mocks.ShopVoucherRepository)
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
			description: "should return list of shop voucher when success",
			input: input{
				slug: slug,
				err:  nil,
				beforeTest: func(ss *mocks.ShopService, svr *mocks.ShopVoucherRepository) {
					ss.On("FindShopBySlug", slug).Return(shop, nil)
					svr.On("GetShopVoucher", shop.ID).Return(voucher, nil)
				},
			},
			expected: expected{
				result: voucher,
				err:    nil,
			},
		},
		{
			description: "should return error when shop not found",
			input: input{
				slug: slug,
				err:  nil,
				beforeTest: func(ss *mocks.ShopService, svr *mocks.ShopVoucherRepository) {
					ss.On("FindShopBySlug", slug).Return(nil, errs.ErrShopNotFound)
				},
			},
			expected: expected{
				result: nil,
				err:    errs.ErrShopNotFound,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				slug: slug,
				err:  errs.ErrInternalServerError,
				beforeTest: func(ss *mocks.ShopService, svr *mocks.ShopVoucherRepository) {
					ss.On("FindShopBySlug", slug).Return(shop, nil)
					svr.On("GetShopVoucher", shop.ID).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				result: nil,
				err:    errs.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.ShopVoucherRepository)
			mockService := new(mocks.ShopService)
			tc.beforeTest(mockService, mockRepo)
			service := service.NewShopVoucherService(&service.ShopVoucherSConfig{
				ShopVoucherRepository: mockRepo,
				ShopService:           mockService,
			})

			result, err := service.GetShopVoucher(tc.input.slug)

			assert.Equal(t, tc.expected.result, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestGetValidShopVoucherByUserIDAndSlug(t *testing.T) {
	var (
		slug    = "shop"
		userID  = 1
		voucher = []*model.ShopVoucher{}
		shop    = &model.Shop{}
	)
	type input struct {
		slug       string
		userID     int
		err        error
		beforeTest func(*mocks.ShopService, *mocks.ShopVoucherRepository)
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
			description: "should return list of shop voucher when success",
			input: input{
				slug:   slug,
				userID: userID,
				err:    nil,
				beforeTest: func(ss *mocks.ShopService, svr *mocks.ShopVoucherRepository) {
					ss.On("FindShopBySlug", slug).Return(shop, nil)
					svr.On("GetValidByUserIDAndShopID", userID, shop.ID).Return(voucher, nil)
				},
			},
			expected: expected{
				result: voucher,
				err:    nil,
			},
		},
		{
			description: "should return error when shop not found",
			input: input{
				slug:   slug,
				userID: userID,
				err:    nil,
				beforeTest: func(ss *mocks.ShopService, svr *mocks.ShopVoucherRepository) {
					ss.On("FindShopBySlug", slug).Return(nil, errs.ErrShopNotFound)
				},
			},
			expected: expected{
				result: nil,
				err:    errs.ErrShopNotFound,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				slug:   slug,
				userID: userID,
				err:    errs.ErrInternalServerError,
				beforeTest: func(ss *mocks.ShopService, svr *mocks.ShopVoucherRepository) {
					ss.On("FindShopBySlug", slug).Return(shop, nil)
					svr.On("GetValidByUserIDAndShopID", userID, shop.ID).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				result: nil,
				err:    errs.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.ShopVoucherRepository)
			mockService := new(mocks.ShopService)
			tc.beforeTest(mockService, mockRepo)
			service := service.NewShopVoucherService(&service.ShopVoucherSConfig{
				ShopVoucherRepository: mockRepo,
				ShopService:           mockService,
			})

			result, err := service.GetValidShopVoucherByUserIDAndSlug(int(tc.input.userID), tc.input.slug)

			assert.Equal(t, tc.expected.result, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

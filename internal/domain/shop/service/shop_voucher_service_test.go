package service_test

import (
	"errors"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
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

func TestGetSellerVoucher(t *testing.T) {
	type input struct {
		userID  int
		request *dto.SellerVoucherFilterRequest
	}
	type expected struct {
		data *commonDto.PaginationResponse
		err  error
	}

	var (
		userID     = 1
		shopID     = 1
		limit      = 20
		page       = 1
		request    = &dto.SellerVoucherFilterRequest{Limit: limit, Page: page}
		vouchers   = []*dto.SellerVoucher{}
		totalRows  = int64(0)
		totalPages = 0
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ShopService, *mocks.ShopVoucherRepository)
		expected
	}{
		{
			description: "should return error when failed to get shop",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ss *mocks.ShopService, pr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(nil, errors.New("failed to get shop"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when failed to get vouchers",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ss *mocks.ShopService, pr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{UserID: userID, ID: shopID}, nil)
				pr.On("GetSellerVoucher", shopID, request).Return(nil, int64(0), 0, errors.New("failed to get vouchers"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get vouchers"),
			},
		},
		{
			description: "should return voucher data when succeed to get vouchers",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ss *mocks.ShopService, pr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{UserID: userID, ID: shopID}, nil)
				pr.On("GetSellerVoucher", shopID, request).Return(vouchers, totalRows, totalPages, nil)
			},
			expected: expected{
				data: &commonDto.PaginationResponse{
					TotalRows:  totalRows,
					TotalPages: totalPages,
					Page:       page,
					Limit:      limit,
					Data:       vouchers,
				},
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			shopService := mocks.NewShopService(t)
			shopVoucherRepository := mocks.NewShopVoucherRepository(t)
			tc.beforeTest(shopService, shopVoucherRepository)
			shopVoucherService := service.NewShopVoucherService(&service.ShopVoucherSConfig{
				ShopService:           shopService,
				ShopVoucherRepository: shopVoucherRepository,
			})

			data, err := shopVoucherService.GetSellerVoucher(tc.input.userID, tc.input.request)

			assert.Equal(t, tc.expected.data, data)
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
		req     = dto.GetValidShopVoucherRequest{
			Slug:   slug,
			UserID: userID,
		}
	)
	type input struct {
		req        dto.GetValidShopVoucherRequest
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
				req: dto.GetValidShopVoucherRequest{
					Slug:   slug,
					UserID: userID,
				},
				err: nil,
				beforeTest: func(ss *mocks.ShopService, svr *mocks.ShopVoucherRepository) {
					ss.On("FindShopBySlug", slug).Return(shop, nil)
					svr.On("GetValidByUserIDAndShopID", req, shop.ID).Return(voucher, nil)
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
				req: dto.GetValidShopVoucherRequest{
					Slug:   slug,
					UserID: userID,
				},
				err: nil,
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
				req: dto.GetValidShopVoucherRequest{
					Slug:   slug,
					UserID: userID,
				},
				err: errs.ErrInternalServerError,
				beforeTest: func(ss *mocks.ShopService, svr *mocks.ShopVoucherRepository) {
					ss.On("FindShopBySlug", slug).Return(shop, nil)
					svr.On("GetValidByUserIDAndShopID", req, shop.ID).Return(nil, errs.ErrInternalServerError)
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

			result, err := service.GetValidShopVoucherByUserIDAndSlug(tc.input.req)

			assert.Equal(t, tc.expected.result, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestGetValidShopVoucherByIdAndUserId(t *testing.T) {
	var (
		id      = 1
		userID  = 1
		voucher = &model.ShopVoucher{
			ID: 1,
		}
	)
	type input struct {
		id         int
		userID     int
		err        error
		beforeTest func(*mocks.ShopVoucherRepository)
	}
	type expected struct {
		result *model.ShopVoucher
		err    error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return shop voucher when success",
			input: input{
				id:     id,
				userID: userID,
				err:    nil,
				beforeTest: func(svr *mocks.ShopVoucherRepository) {
					svr.On("GetValidByIdAndUserId", id, userID).Return(voucher, nil)
				},
			},
			expected: expected{
				result: voucher,
				err:    nil,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				id:     id,
				userID: userID,
				err:    errs.ErrInternalServerError,
				beforeTest: func(svr *mocks.ShopVoucherRepository) {
					svr.On("GetValidByIdAndUserId", id, userID).Return(nil, errs.ErrInternalServerError)
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
			tc.beforeTest(mockRepo)
			service := service.NewShopVoucherService(&service.ShopVoucherSConfig{
				ShopVoucherRepository: mockRepo,
			})

			result, err := service.GetValidShopVoucherByIdAndUserId(tc.input.id, tc.input.userID)

			assert.Equal(t, tc.expected.result, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

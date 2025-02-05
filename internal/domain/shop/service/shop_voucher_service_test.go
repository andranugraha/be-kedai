package service_test

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/constant"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/service"
	"kedai/backend/be-kedai/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
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
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{UserID: userID, ID: shopID}, nil)
				vr.On("GetSellerVoucher", shopID, request).Return(nil, int64(0), 0, errors.New("failed to get vouchers"))
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
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{UserID: userID, ID: shopID}, nil)
				vr.On("GetSellerVoucher", shopID, request).Return(vouchers, totalRows, totalPages, nil)
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

func TestGetVoucherByCodeAndShopId(t *testing.T) {
	type input struct {
		voucherCode string
		userID      int
	}
	type expected struct {
		data *dto.SellerVoucher
		err  error
	}

	var (
		userID      = 1
		shopID      = 1
		voucherCode = "voucher-code"
		voucher     = dto.SellerVoucher{}
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
				userID:      userID,
				voucherCode: voucherCode,
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(nil, errors.New("failed to get shop"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when failed to get voucher",
			input: input{
				userID:      userID,
				voucherCode: voucherCode,
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID, UserID: userID}, nil)
				vr.On("GetVoucherByCodeAndShopId", voucherCode, shopID).Return(nil, errors.New("failed to get voucher"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get voucher"),
			},
		},
		{
			description: "should return dto when succeed to get voucher",
			input: input{
				userID:      userID,
				voucherCode: voucherCode,
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID, UserID: userID}, nil)
				vr.On("GetVoucherByCodeAndShopId", voucherCode, shopID).Return(&voucher, nil)
			},
			expected: expected{
				data: &dto.SellerVoucher{},
				err:  nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			shopVoucherRepo := mocks.NewShopVoucherRepository(t)
			shopService := mocks.NewShopService(t)
			tc.beforeTest(shopService, shopVoucherRepo)
			shopVoucherService := service.NewShopVoucherService(&service.ShopVoucherSConfig{
				ShopVoucherRepository: shopVoucherRepo,
				ShopService:           shopService,
			})

			data, err := shopVoucherService.GetVoucherByCodeAndShopId(tc.input.voucherCode, tc.input.userID)

			assert.Equal(t, tc.expected.data, data)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestCreateVoucher(t *testing.T) {
	type input struct {
		userID  int
		request *dto.CreateVoucherRequest
	}
	type expected struct {
		data *model.ShopVoucher
		err  error
	}

	var (
		userID       = 1
		shopID       = 1
		voucherName  = "voucher name"
		voucherCode  = "voucher code"
		startFrom, _ = time.Parse(time.RFC3339Nano, "2023-03-25T10:08:28.905Z")
		expiredAt, _ = time.Parse(time.RFC3339Nano, "2023-03-25T10:08:28.905Z")
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ShopService, *mocks.ShopVoucherRepository)
		expected
	}{
		{
			description: "should return error when voucher name is invalid",
			input: input{
				userID: userID,
				request: &dto.CreateVoucherRequest{
					Name: "127.0.0.1",
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {},
			expected: expected{
				data: nil,
				err:  errs.ErrInvalidVoucherNamePattern,
			},
		},
		{
			description: "should return error when failed to get shop",
			input: input{
				userID: userID,
				request: &dto.CreateVoucherRequest{
					Name: voucherName,
					Code: voucherCode,
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(nil, errors.New("failed to get shop"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when duplicate voucher code is found",
			input: input{
				userID: userID,
				request: &dto.CreateVoucherRequest{
					Name: voucherName,
					Code: voucherCode,
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				vr.On("GetVoucherByCodeAndShopId", voucherCode, shopID).Return(nil, errors.New("failed to get voucher"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get voucher"),
			},
		},
		{
			description: "should return error when duplicate voucher code is found",
			input: input{
				userID: userID,
				request: &dto.CreateVoucherRequest{
					Name: voucherName,
					Code: voucherCode,
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				vr.On("GetVoucherByCodeAndShopId", voucherCode, shopID).Return(&dto.SellerVoucher{Status: "ongoing"}, nil)
			},
			expected: expected{
				data: nil,
				err:  errors.New("duplicate voucher code"),
			},
		},
		{
			description: "should return error when voucher date range is invalid",
			input: input{
				userID: userID,
				request: &dto.CreateVoucherRequest{
					Name:      voucherName,
					Code:      voucherCode,
					StartFrom: startFrom,
					ExpiredAt: expiredAt,
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				vr.On("GetVoucherByCodeAndShopId", voucherCode, shopID).Return(&dto.SellerVoucher{
					ShopVoucher: model.ShopVoucher{Name: voucherName, Code: voucherCode},
				}, nil)
				vr.On("ValidateVoucherDateRange", startFrom, expiredAt).Return(errors.New("invalid voucher date range"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("invalid voucher date range"),
			},
		},
		{
			description: "should return error when failed to create voucher",
			input: input{
				userID: userID,
				request: &dto.CreateVoucherRequest{
					Name:      voucherName,
					Code:      voucherCode,
					StartFrom: startFrom,
					ExpiredAt: expiredAt,
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				vr.On("GetVoucherByCodeAndShopId", voucherCode, shopID).Return(&dto.SellerVoucher{
					ShopVoucher: model.ShopVoucher{Name: voucherName, Code: voucherCode},
				}, nil)
				vr.On("ValidateVoucherDateRange", startFrom, expiredAt).Return(nil)
				vr.On("Create", shopID, &dto.CreateVoucherRequest{Name: voucherName, Code: voucherCode, StartFrom: startFrom, ExpiredAt: expiredAt}).Return(nil, errors.New("failed to create voucher"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to create voucher"),
			},
		},
		{
			description: "should return created voucher when succeed to create voucher",
			input: input{
				userID: userID,
				request: &dto.CreateVoucherRequest{
					Name:      voucherName,
					Code:      voucherCode,
					StartFrom: startFrom,
					ExpiredAt: expiredAt,
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				vr.On("GetVoucherByCodeAndShopId", voucherCode, shopID).Return(&dto.SellerVoucher{
					ShopVoucher: model.ShopVoucher{Name: voucherName, Code: voucherCode},
				}, nil)
				vr.On("ValidateVoucherDateRange", startFrom, expiredAt).Return(nil)
				vr.On("Create", shopID, &dto.CreateVoucherRequest{
					Name: voucherName, Code: voucherCode, StartFrom: startFrom, ExpiredAt: expiredAt,
				}).Return(&model.ShopVoucher{
					Name: voucherName, Code: voucherCode, StartFrom: startFrom, ExpiredAt: expiredAt,
				}, nil)
			},
			expected: expected{
				data: &model.ShopVoucher{
					Name:      voucherName,
					Code:      voucherCode,
					StartFrom: startFrom,
					ExpiredAt: expiredAt,
				},
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			shopService := mocks.NewShopService(t)
			shopVoucherRepo := mocks.NewShopVoucherRepository(t)
			tc.beforeTest(shopService, shopVoucherRepo)
			shopVoucherService := service.NewShopVoucherService(&service.ShopVoucherSConfig{
				ShopVoucherRepository: shopVoucherRepo,
				ShopService:           shopService,
			})

			data, err := shopVoucherService.CreateVoucher(tc.input.userID, tc.input.request)

			assert.Equal(t, tc.expected.data, data)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestUpdateVoucher(t *testing.T) {
	type input struct {
		userID      int
		shopID      int
		voucherCode string
		request     *dto.UpdateVoucherRequest
	}
	type expected struct {
		data *model.ShopVoucher
		err  error
	}

	var (
		userID       = 1
		shopID       = 1
		voucherName  = "voucher name"
		voucherCode  = "voucher-code"
		isHidden     = false
		startFrom, _ = time.Parse(time.RFC3339Nano, "2023-03-25T10:08:28.905Z")
		expiredAt, _ = time.Parse(time.RFC3339Nano, "2023-03-25T10:08:28.905Z")
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ShopService, *mocks.ShopVoucherRepository)
		expected
	}{
		{
			description: "should return error when voucher name is invalid",
			input: input{
				userID: userID,
				request: &dto.UpdateVoucherRequest{
					Name: "127.0.0.1",
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {},
			expected: expected{
				data: nil,
				err:  errs.ErrInvalidVoucherNamePattern,
			},
		},
		{
			description: "should return error when failed to get shop",
			input: input{
				userID: userID,
				request: &dto.UpdateVoucherRequest{
					Name: voucherName,
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(nil, errors.New("failed to get shop"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when failed to get voucher",
			input: input{
				userID:      userID,
				shopID:      shopID,
				voucherCode: voucherCode,
				request: &dto.UpdateVoucherRequest{
					Name: voucherName,
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				vr.On("GetVoucherByCodeAndShopId", voucherCode, shopID).Return(nil, errors.New("failed to get voucher"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get voucher"),
			},
		},
		{
			description: "should return updated voucher when succeed to update voucher when voucher status is upcoming",
			input: input{
				userID:      userID,
				shopID:      shopID,
				voucherCode: voucherCode,
				request: &dto.UpdateVoucherRequest{
					Name:      voucherName,
					IsHidden:  &isHidden,
					StartFrom: startFrom,
					ExpiredAt: expiredAt,
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				vr.On("GetVoucherByCodeAndShopId", voucherCode, shopID).Return(&dto.SellerVoucher{
					ShopVoucher: model.ShopVoucher{Name: voucherName},
					Status:      "upcoming",
				}, nil)
				vr.On("ValidateVoucherDateRange", startFrom, expiredAt).Return(nil)
				vr.On("Update", &model.ShopVoucher{Name: voucherName, IsHidden: isHidden, ShopId: shopID, StartFrom: startFrom, ExpiredAt: expiredAt}).Return(&model.ShopVoucher{Name: voucherName, StartFrom: startFrom, ExpiredAt: expiredAt}, nil)
			},
			expected: expected{
				data: &model.ShopVoucher{
					Name:      voucherName,
					StartFrom: startFrom,
					ExpiredAt: expiredAt,
				},
				err: nil,
			},
		},
		{
			description: "should return error when voucher date range is invalid  when voucher status is upcoming",
			input: input{
				userID:      userID,
				shopID:      shopID,
				voucherCode: voucherCode,
				request: &dto.UpdateVoucherRequest{
					Name:      voucherName,
					IsHidden:  &isHidden,
					StartFrom: startFrom,
					ExpiredAt: expiredAt,
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				vr.On("GetVoucherByCodeAndShopId", voucherCode, shopID).Return(&dto.SellerVoucher{
					ShopVoucher: model.ShopVoucher{Name: voucherName},
					Status:      "upcoming",
				}, nil)
				vr.On("ValidateVoucherDateRange", startFrom, expiredAt).Return(errors.New("invalid voucher date range"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("invalid voucher date range"),
			},
		},
		{
			description: "should return error when failed to update voucher when voucher status is upcoming",
			input: input{
				userID:      userID,
				shopID:      shopID,
				voucherCode: voucherCode,
				request: &dto.UpdateVoucherRequest{
					Name:      voucherName,
					IsHidden:  &isHidden,
					StartFrom: startFrom,
					ExpiredAt: expiredAt,
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				vr.On("GetVoucherByCodeAndShopId", voucherCode, shopID).Return(&dto.SellerVoucher{
					ShopVoucher: model.ShopVoucher{Name: voucherName},
					Status:      "upcoming",
				}, nil)
				vr.On("ValidateVoucherDateRange", startFrom, expiredAt).Return(nil)
				vr.On("Update", &model.ShopVoucher{
					Name: voucherName, IsHidden: isHidden, ShopId: shopID, StartFrom: startFrom, ExpiredAt: expiredAt,
				}).Return(nil, errors.New("failed to update voucher"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to update voucher"),
			},
		},
		{
			description: "should return updated voucher when succeed to update voucher when voucher status is ongoing",
			input: input{
				userID:      userID,
				shopID:      shopID,
				voucherCode: voucherCode,
				request: &dto.UpdateVoucherRequest{
					Name:       voucherName,
					IsHidden:   &isHidden,
					TotalQuota: 100,
					ExpiredAt:  expiredAt,
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				vr.On("GetVoucherByCodeAndShopId", voucherCode, shopID).Return(&dto.SellerVoucher{
					ShopVoucher: model.ShopVoucher{Name: voucherName},
					Status:      "ongoing",
				}, nil)
				vr.On("ValidateVoucherDateRange", mock.Anything, expiredAt).Return(nil)
				vr.On("Update", &model.ShopVoucher{
					Name:       voucherName,
					IsHidden:   isHidden,
					TotalQuota: 100,
					ExpiredAt:  expiredAt,
					ShopId:     shopID,
				}).Return(&model.ShopVoucher{
					Name:       voucherName,
					IsHidden:   isHidden,
					TotalQuota: 100,
					ExpiredAt:  expiredAt,
				}, nil)
			},
			expected: expected{
				data: &model.ShopVoucher{
					Name:       voucherName,
					IsHidden:   isHidden,
					TotalQuota: 100,
					ExpiredAt:  expiredAt,
				},
				err: nil,
			},
		},
		{
			description: "should return error when voucher date range is invalid  when voucher status is upcoming",
			input: input{
				userID:      userID,
				shopID:      shopID,
				voucherCode: voucherCode,
				request: &dto.UpdateVoucherRequest{
					Name:       voucherName,
					IsHidden:   &isHidden,
					TotalQuota: 100,
					ExpiredAt:  expiredAt,
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				vr.On("GetVoucherByCodeAndShopId", voucherCode, shopID).Return(&dto.SellerVoucher{
					ShopVoucher: model.ShopVoucher{Name: voucherName},
					Status:      "ongoing",
				}, nil)
				vr.On("ValidateVoucherDateRange", mock.Anything, expiredAt).Return(errors.New("invalid voucher date range"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("invalid voucher date range"),
			},
		},
		{
			description: "should return error when failed to update voucher when voucher status is ongoing",
			input: input{
				userID:      userID,
				shopID:      shopID,
				voucherCode: voucherCode,
				request: &dto.UpdateVoucherRequest{
					Name:       voucherName,
					IsHidden:   &isHidden,
					TotalQuota: 100,
					ExpiredAt:  expiredAt,
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				vr.On("GetVoucherByCodeAndShopId", voucherCode, shopID).Return(&dto.SellerVoucher{
					ShopVoucher: model.ShopVoucher{Name: voucherName},
					Status:      "ongoing",
				}, nil)
				vr.On("ValidateVoucherDateRange", mock.Anything, expiredAt).Return(nil)
				vr.On("Update", &model.ShopVoucher{
					Name:       voucherName,
					IsHidden:   isHidden,
					TotalQuota: 100,
					ExpiredAt:  expiredAt,
					ShopId:     shopID,
				}).Return(nil, errors.New("failed to update voucher"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to update voucher"),
			},
		},
		{
			description: "should return ErrVoucherStatusConflict when voucher status is expired",
			input: input{
				userID:      userID,
				shopID:      shopID,
				voucherCode: voucherCode,
				request: &dto.UpdateVoucherRequest{
					Name: voucherName,
				},
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(&model.Shop{ID: shopID}, nil)
				voucher := &dto.SellerVoucher{
					Status: constant.VoucherPromotionStatusExpired,
				}
				vr.On("GetVoucherByCodeAndShopId", voucherCode, shopID).Return(voucher, nil)
			},
			expected: expected{
				data: nil,
				err:  errs.ErrVoucherStatusConflict,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			shopService := mocks.NewShopService(t)
			shopVoucherRepo := mocks.NewShopVoucherRepository(t)
			tc.beforeTest(shopService, shopVoucherRepo)
			shopVoucherService := service.NewShopVoucherService(&service.ShopVoucherSConfig{
				ShopVoucherRepository: shopVoucherRepo,
				ShopService:           shopService,
			})

			data, err := shopVoucherService.UpdateVoucher(tc.input.userID, tc.input.voucherCode, tc.input.request)

			assert.Equal(t, tc.expected.data, data)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestDeleteVoucher(t *testing.T) {
	type input struct {
		voucherCode string
		userID      int
	}
	type expected struct {
		err error
	}

	var (
		userID      = 1
		shopID      = 1
		voucherCode = "voucher-code"
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
				userID:      userID,
				voucherCode: voucherCode,
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				ss.On("FindShopByUserId", userID).Return(nil, errors.New("failed to get shop"))
			},
			expected: expected{
				err: errors.New("failed to get shop"),
			},
		}, {
			description: "should return error when failed to delete voucher",
			input: input{
				userID:      userID,
				voucherCode: voucherCode,
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				shop := &model.Shop{ID: shopID}
				ss.On("FindShopByUserId", userID).Return(shop, nil)
				vr.On("Delete", shopID, voucherCode).Return(errors.New("failed to delete voucher"))
			},
			expected: expected{
				err: errors.New("failed to delete voucher"),
			},
		},
		{
			description: "should return nil when voucher is successfully deleted",
			input: input{
				userID:      userID,
				voucherCode: voucherCode,
			},
			beforeTest: func(ss *mocks.ShopService, vr *mocks.ShopVoucherRepository) {
				shop := &model.Shop{ID: shopID}
				ss.On("FindShopByUserId", userID).Return(shop, nil)
				vr.On("Delete", shopID, voucherCode).Return(nil)
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			shopVoucherRepo := mocks.NewShopVoucherRepository(t)
			shopService := mocks.NewShopService(t)
			tc.beforeTest(shopService, shopVoucherRepo)
			shopVoucherService := service.NewShopVoucherService(&service.ShopVoucherSConfig{
				ShopVoucherRepository: shopVoucherRepo,
				ShopService:           shopService,
			})

			err := shopVoucherService.DeleteVoucher(tc.input.userID, tc.input.voucherCode)

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

package service_test

import (
	"errors"
	"kedai/backend/be-kedai/internal/common/constant"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/order/dto"

	"kedai/backend/be-kedai/internal/domain/order/model"
	"kedai/backend/be-kedai/internal/domain/order/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	walletModel "kedai/backend/be-kedai/internal/domain/user/model"

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

func TestGetInvoicesByShopId(t *testing.T) {
	var (
		shop = &shopModel.Shop{
			ID: 1,
		}
		userId     = 1
		invoice    = []*dto.InvoicePerShopDetail{}
		req        = &dto.InvoicePerShopFilterRequest{}
		pagination = &commonDto.PaginationResponse{
			Data: invoice,
		}
	)
	type input struct {
		userId     int
		req        *dto.InvoicePerShopFilterRequest
		err        error
		beforeTest func(*mocks.ShopService, *mocks.InvoicePerShopRepository)
	}
	type expected struct {
		result *commonDto.PaginationResponse
		err    error
	}
	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return list of shop invoices when success",
			input: input{
				userId: userId,
				req:    req,
				err:    nil,
				beforeTest: func(ss *mocks.ShopService, ipsr *mocks.InvoicePerShopRepository) {
					ss.On("FindShopByUserId", userId).Return(shop, nil)
					ipsr.On("GetByShopId", shop.ID, req).Return(invoice, int64(0), 0, nil)
				},
			},
			expected: expected{
				result: pagination,
				err:    nil,
			},
		},
		{
			description: "should return error when user shop not found",
			input: input{
				userId: userId,
				req:    req,
				err:    commonErr.ErrShopNotFound,
				beforeTest: func(ss *mocks.ShopService, ipsr *mocks.InvoicePerShopRepository) {
					ss.On("FindShopByUserId", userId).Return(nil, commonErr.ErrShopNotFound)
				},
			},
			expected: expected{
				result: nil,
				err:    commonErr.ErrShopNotFound,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				userId: userId,
				req:    req,
				err:    commonErr.ErrInternalServerError,
				beforeTest: func(ss *mocks.ShopService, ipsr *mocks.InvoicePerShopRepository) {
					ss.On("FindShopByUserId", userId).Return(shop, nil)
					ipsr.On("GetByShopId", shop.ID, req).Return(nil, int64(0), 0, commonErr.ErrInternalServerError)
				},
			},
			expected: expected{
				result: nil,
				err:    commonErr.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			invoicePerShopRepo := mocks.NewInvoicePerShopRepository(t)
			shopService := mocks.NewShopService(t)
			tc.beforeTest(shopService, invoicePerShopRepo)
			invoicePerShopService := service.NewInvoicePerShopService(&service.InvoicePerShopSConfig{
				InvoicePerShopRepo: invoicePerShopRepo,
				ShopService:        shopService,
			})

			result, err := invoicePerShopService.GetInvoicesByShopId(tc.input.userId, tc.input.req)

			assert.Equal(t, tc.expected.result, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestGetInvoiceByUserIdAndId(t *testing.T) {
	type input struct {
		userID     int
		id         int
		beforeTest func(*mocks.InvoicePerShopRepository, *mocks.ShopService)
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
			description: "should return error when FindShopByUserId failed",
			input: input{
				userID: 1,
				id:     1,
				beforeTest: func(ipsr *mocks.InvoicePerShopRepository, ss *mocks.ShopService) {
					ss.On("FindShopByUserId", 1).Return(nil, errors.New("failed to get shop"))
				},
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when GetByShopIdAndId failed",
			input: input{
				userID: 1,
				id:     1,
				beforeTest: func(ipsr *mocks.InvoicePerShopRepository, ss *mocks.ShopService) {
					ss.On("FindShopByUserId", 1).Return(&shopModel.Shop{ID: 1}, nil)
					ipsr.On("GetByShopIdAndId", 1, 1).Return(nil, errors.New("failed to get invoice"))
				},
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get invoice"),
			},
		},
		{
			description: "should return invoice detail and no error when success",
			input: input{
				userID: 1,
				id:     1,
				beforeTest: func(ipsr *mocks.InvoicePerShopRepository, ss *mocks.ShopService) {
					ss.On("FindShopByUserId", 1).Return(&shopModel.Shop{ID: 1}, nil)
					ipsr.On("GetByShopIdAndId", 1, 1).Return(&dto.InvoicePerShopDetail{}, nil)
				},
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
			shopService := mocks.NewShopService(t)
			tc.beforeTest(invoicePerShopRepo, shopService)
			invoicePerShopService := service.NewInvoicePerShopService(&service.InvoicePerShopSConfig{
				InvoicePerShopRepo: invoicePerShopRepo,
				ShopService:        shopService,
			})

			data, err := invoicePerShopService.GetInvoiceByUserIdAndId(tc.input.userID, tc.input.id)

			assert.Equal(t, tc.expected.data, data)
			assert.Equal(t, tc.expected.err, err)
		})
	}

}

func TestWithdrawFromInvoice(t *testing.T) {
	type input struct {
		userID     int
		id         []int
		beforeTest func(*mocks.InvoicePerShopRepository, *mocks.ShopService, *mocks.WalletService)
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
			description: "should return error when FindShopByUserId failed",
			input: input{
				userID: 1,
				id:     []int{1},
				beforeTest: func(ipsr *mocks.InvoicePerShopRepository, ss *mocks.ShopService, ws *mocks.WalletService) {
					ss.On("FindShopByUserId", 1).Return(nil, errors.New("failed to get shop"))
				},
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when GetWalletByUserID failed",
			input: input{
				userID: 1,
				id:     []int{1},
				beforeTest: func(ipsr *mocks.InvoicePerShopRepository, ss *mocks.ShopService, ws *mocks.WalletService) {
					ss.On("FindShopByUserId", 1).Return(&shopModel.Shop{ID: 1}, nil)
					ws.On("GetWalletByUserID", 1).Return(nil, errors.New("failed to get wallet"))
				},
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get wallet"),
			},
		},
		{
			description: "should return error when WithdrawFromInvoice failed",
			input: input{
				userID: 1,
				id:     []int{1},
				beforeTest: func(ipsr *mocks.InvoicePerShopRepository, ss *mocks.ShopService, ws *mocks.WalletService) {
					ss.On("FindShopByUserId", 1).Return(&shopModel.Shop{ID: 1}, nil)
					ws.On("GetWalletByUserID", 1).Return(&walletModel.Wallet{ID: 1}, nil)
					ipsr.On("WithdrawFromInvoice", []int{1}, 1, 1).Return(errors.New("failed to get invoice"))
				},
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get invoice"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			invoicePerShopRepo := mocks.NewInvoicePerShopRepository(t)
			shopService := mocks.NewShopService(t)
			walletService := mocks.NewWalletService(t)
			tc.beforeTest(invoicePerShopRepo, shopService, walletService)
			invoicePerShopService := service.NewInvoicePerShopService(&service.InvoicePerShopSConfig{
				InvoicePerShopRepo: invoicePerShopRepo,
				ShopService:        shopService,
				WalletService:      walletService,
			})

			err := invoicePerShopService.WithdrawFromInvoice(tc.input.id, tc.input.userID)

			assert.Equal(t, tc.expected.err, err)
		})
	}

}

func TestGetShopOrder(t *testing.T) {
	var (
		shop = &shopModel.Shop{
			ID: 1,
		}
		userId     = 1
		invoice    = []*dto.InvoicePerShopDetail{}
		req        = &dto.InvoicePerShopFilterRequest{}
		pagination = &commonDto.PaginationResponse{
			Data: invoice,
		}
	)
	type input struct {
		userId     int
		req        *dto.InvoicePerShopFilterRequest
		err        error
		beforeTest func(*mocks.ShopService, *mocks.InvoicePerShopRepository)
	}
	type expected struct {
		result *commonDto.PaginationResponse
		err    error
	}
	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return list of shop invoices order when success",
			input: input{
				userId: userId,
				req:    req,
				err:    nil,
				beforeTest: func(ss *mocks.ShopService, ipsr *mocks.InvoicePerShopRepository) {
					ss.On("FindShopByUserId", userId).Return(shop, nil)
					ipsr.On("GetShopOrder", shop.ID, req).Return(invoice, int64(0), 0, nil)
				},
			},
			expected: expected{
				result: pagination,
				err:    nil,
			},
		},
		{
			description: "should return error when user shop not found",
			input: input{
				userId: userId,
				req:    req,
				err:    commonErr.ErrShopNotFound,
				beforeTest: func(ss *mocks.ShopService, ipsr *mocks.InvoicePerShopRepository) {
					ss.On("FindShopByUserId", userId).Return(nil, commonErr.ErrShopNotFound)
				},
			},
			expected: expected{
				result: nil,
				err:    commonErr.ErrShopNotFound,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				userId: userId,
				req:    req,
				err:    commonErr.ErrInternalServerError,
				beforeTest: func(ss *mocks.ShopService, ipsr *mocks.InvoicePerShopRepository) {
					ss.On("FindShopByUserId", userId).Return(shop, nil)
					ipsr.On("GetShopOrder", shop.ID, req).Return(nil, int64(0), 0, commonErr.ErrInternalServerError)
				},
			},
			expected: expected{
				result: nil,
				err:    commonErr.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			invoicePerShopRepo := mocks.NewInvoicePerShopRepository(t)
			shopService := mocks.NewShopService(t)
			tc.beforeTest(shopService, invoicePerShopRepo)
			invoicePerShopService := service.NewInvoicePerShopService(&service.InvoicePerShopSConfig{
				InvoicePerShopRepo: invoicePerShopRepo,
				ShopService:        shopService,
			})

			result, err := invoicePerShopService.GetShopOrder(tc.input.userId, tc.input.req)

			assert.Equal(t, tc.expected.result, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestUpdateStatusToDelivery(t *testing.T) {
	var (
		shop = &shopModel.Shop{
			ID: 1,
		}
		invoiceStatuses = []*model.InvoiceStatus{
			{
				InvoicePerShopID: 1,
				Status:           "ON_DELIVERY",
			},
		}
		userId  = 1
		orderId = 1
	)
	type input struct {
		userId     int
		orderId    int
		beforeTest func(*mocks.ShopService, *mocks.InvoicePerShopRepository)
	}
	type expected struct {
		err error
	}
	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return nil error when success",
			input: input{
				userId:  userId,
				orderId: orderId,
				beforeTest: func(ss *mocks.ShopService, ipsr *mocks.InvoicePerShopRepository) {
					ss.On("FindShopByUserId", userId).Return(shop, nil)
					ipsr.On("UpdateStatusToDelivery", shop.ID, orderId, invoiceStatuses).Return(nil)
				},
			},
			expected: expected{
				err: nil,
			},
		},
		{
			description: "should return error when shop not found",
			input: input{
				userId:  userId,
				orderId: orderId,
				beforeTest: func(ss *mocks.ShopService, ipsr *mocks.InvoicePerShopRepository) {
					ss.On("FindShopByUserId", userId).Return(nil, commonErr.ErrShopNotFound)
				},
			},
			expected: expected{
				err: commonErr.ErrShopNotFound,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				userId:  userId,
				orderId: orderId,
				beforeTest: func(ss *mocks.ShopService, ipsr *mocks.InvoicePerShopRepository) {
					ss.On("FindShopByUserId", userId).Return(shop, nil)
					ipsr.On("UpdateStatusToDelivery", shop.ID, orderId, invoiceStatuses).Return(commonErr.ErrInternalServerError)
				},
			},
			expected: expected{
				err: commonErr.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			invoicePerShopRepo := mocks.NewInvoicePerShopRepository(t)
			shopService := mocks.NewShopService(t)
			tc.beforeTest(shopService, invoicePerShopRepo)
			invoicePerShopService := service.NewInvoicePerShopService(&service.InvoicePerShopSConfig{
				InvoicePerShopRepo: invoicePerShopRepo,
				ShopService:        shopService,
			})

			err := invoicePerShopService.UpdateStatusToDelivery(tc.userId, tc.orderId)

			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestUpdateStatusToCanceled(t *testing.T) {
	var (
		shop = &shopModel.Shop{
			ID: 1,
		}
		invoiceStatuses = []*model.InvoiceStatus{
			{
				InvoicePerShopID: 1,
				Status:           "CANCELED",
			},
		}
		userId  = 1
		orderId = 1
	)
	type input struct {
		userId     int
		orderId    int
		beforeTest func(*mocks.ShopService, *mocks.InvoicePerShopRepository)
	}
	type expected struct {
		err error
	}
	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return nil error when success",
			input: input{
				userId:  userId,
				orderId: orderId,
				beforeTest: func(ss *mocks.ShopService, ipsr *mocks.InvoicePerShopRepository) {
					ss.On("FindShopByUserId", userId).Return(shop, nil)
					ipsr.On("UpdateStatusToCanceled", shop.ID, orderId, invoiceStatuses).Return(nil)
				},
			},
			expected: expected{
				err: nil,
			},
		},
		{
			description: "should return error when shop not found",
			input: input{
				userId:  userId,
				orderId: orderId,
				beforeTest: func(ss *mocks.ShopService, ipsr *mocks.InvoicePerShopRepository) {
					ss.On("FindShopByUserId", userId).Return(nil, commonErr.ErrShopNotFound)
				},
			},
			expected: expected{
				err: commonErr.ErrShopNotFound,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				userId:  userId,
				orderId: orderId,
				beforeTest: func(ss *mocks.ShopService, ipsr *mocks.InvoicePerShopRepository) {
					ss.On("FindShopByUserId", userId).Return(shop, nil)
					ipsr.On("UpdateStatusToCanceled", shop.ID, orderId, invoiceStatuses).Return(commonErr.ErrInternalServerError)
				},
			},
			expected: expected{
				err: commonErr.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			invoicePerShopRepo := mocks.NewInvoicePerShopRepository(t)
			shopService := mocks.NewShopService(t)
			tc.beforeTest(shopService, invoicePerShopRepo)
			invoicePerShopService := service.NewInvoicePerShopService(&service.InvoicePerShopSConfig{
				InvoicePerShopRepo: invoicePerShopRepo,
				ShopService:        shopService,
			})

			err := invoicePerShopService.UpdateStatusToCanceled(tc.userId, tc.orderId)

			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestUpdateStatusToReceived(t *testing.T) {
	var (
		orderCode = "code"
		order     = &dto.InvoicePerShopDetail{
			InvoicePerShop: model.InvoicePerShop{ID: 1, ShopID: 1},
		}
		userId        = 1
		invoiceStatus = []*model.InvoiceStatus{
			{
				InvoicePerShopID: order.ID,
				Status:           constant.TransactionStatusReceived,
			},
		}
	)
	type input struct {
		id         int
		code       string
		beforeTest func(*mocks.InvoicePerShopRepository)
	}
	type expected struct {
		err error
	}
	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return nil error when success",
			input: input{
				id:   userId,
				code: orderCode,
				beforeTest: func(ipsr *mocks.InvoicePerShopRepository) {
					ipsr.On("GetByUserIDAndCode", userId, orderCode).Return(order, nil)
					ipsr.On("UpdateStatusToReceived", order.ShopID, order.ID, invoiceStatus).Return(nil)
				},
			},
			expected: expected{
				err: nil,
			},
		},
		{
			description: "should return error when invoice not found",
			input: input{
				id:   userId,
				code: orderCode,
				beforeTest: func(ipsr *mocks.InvoicePerShopRepository) {
					ipsr.On("GetByUserIDAndCode", userId, orderCode).Return(nil, commonErr.ErrInvoiceNotFound)
				},
			},
			expected: expected{
				err: commonErr.ErrInvoiceNotFound,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				id:   userId,
				code: orderCode,
				beforeTest: func(ipsr *mocks.InvoicePerShopRepository) {
					ipsr.On("GetByUserIDAndCode", userId, orderCode).Return(order, nil)
					ipsr.On("UpdateStatusToReceived", order.ShopID, order.ID, invoiceStatus).Return(commonErr.ErrInternalServerError)
				},
			},
			expected: expected{
				err: commonErr.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.InvoicePerShopRepository)
			tc.beforeTest(mockRepo)
			service := service.NewInvoicePerShopService(&service.InvoicePerShopSConfig{
				InvoicePerShopRepo: mockRepo,
			})

			err := service.UpdateStatusToReceived(tc.id, tc.code)

			assert.Equal(t, tc.err, err)
		})
	}
}

func TestUpdateStatusToCompleted(t *testing.T) {
	var (
		orderCode = "code"
		order     = &dto.InvoicePerShopDetail{
			InvoicePerShop: model.InvoicePerShop{ID: 1, ShopID: 1},
		}
		userId        = 1
		invoiceStatus = []*model.InvoiceStatus{
			{
				InvoicePerShopID: order.ID,
				Status:           constant.TransactionStatusCompleted,
			},
		}
	)
	type input struct {
		id         int
		code       string
		beforeTest func(*mocks.InvoicePerShopRepository)
	}
	type expected struct {
		err error
	}
	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return nil error when success",
			input: input{
				id:   userId,
				code: orderCode,
				beforeTest: func(ipsr *mocks.InvoicePerShopRepository) {
					ipsr.On("GetByUserIDAndCode", userId, orderCode).Return(order, nil)
					ipsr.On("UpdateStatusToCompleted", order.ShopID, order.ID, invoiceStatus).Return(nil)
				},
			},
			expected: expected{
				err: nil,
			},
		},
		{
			description: "should return error when invoice not found",
			input: input{
				id:   userId,
				code: orderCode,
				beforeTest: func(ipsr *mocks.InvoicePerShopRepository) {
					ipsr.On("GetByUserIDAndCode", userId, orderCode).Return(nil, commonErr.ErrInvoiceNotFound)
				},
			},
			expected: expected{
				err: commonErr.ErrInvoiceNotFound,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				id:   userId,
				code: orderCode,
				beforeTest: func(ipsr *mocks.InvoicePerShopRepository) {
					ipsr.On("GetByUserIDAndCode", userId, orderCode).Return(order, nil)
					ipsr.On("UpdateStatusToCompleted", order.ShopID, order.ID, invoiceStatus).Return(commonErr.ErrInternalServerError)
				},
			},
			expected: expected{
				err: commonErr.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.InvoicePerShopRepository)
			tc.beforeTest(mockRepo)
			service := service.NewInvoicePerShopService(&service.InvoicePerShopSConfig{
				InvoicePerShopRepo: mockRepo,
			})

			err := service.UpdateStatusToCompleted(tc.id, tc.code)

			assert.Equal(t, tc.err, err)
		})
	}
}

func TestRefundRequest(t *testing.T) {
	var (
		code    = "code"
		userId  = 1
		invoice = &dto.InvoicePerShopDetail{
			InvoicePerShop: model.InvoicePerShop{
				ID: 1,
			},
		}
		status = []*model.InvoiceStatus{
			{
				Status:           constant.TransactionStatusComplained,
				InvoicePerShopID: 1,
			},
		}
		req = &model.RefundRequest{
			Status:    constant.TransactionStatusComplained,
			InvoiceId: 1,
			Invoice:   &invoice.InvoicePerShop,
		}
	)
	type input struct {
		code       string
		userId     int
		beforeTest func(*mocks.InvoicePerShopRepository)
	}
	type expected struct {
		result *model.RefundRequest
		err    error
	}
	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return refund request when success",
			input: input{
				code: code,
				userId: userId,
				beforeTest: func(ipsr *mocks.InvoicePerShopRepository) {
					ipsr.On("GetByUserIDAndCode", userId, code).Return(invoice, nil)
					ipsr.On("RefundRequest", req, status).Return(req, nil)
				},
			},
			expected: expected{
				result: req,
				err: nil,
			},
		},
		{
			description: "should return error when invoice not found",
			input: input{
				code: code,
				userId: userId,
				beforeTest: func(ipsr *mocks.InvoicePerShopRepository) {
					ipsr.On("GetByUserIDAndCode", userId, code).Return(nil, commonErr.ErrInvoiceNotFound)
				},
			},
			expected: expected{
				result: nil,
				err: commonErr.ErrInvoiceNotFound,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				code: code,
				userId: userId,
				beforeTest: func(ipsr *mocks.InvoicePerShopRepository) {
					ipsr.On("GetByUserIDAndCode", userId, code).Return(invoice, nil)
					ipsr.On("RefundRequest", req, status).Return(nil, commonErr.ErrInternalServerError)
				},
			},
			expected: expected{
				result: nil,
				err: commonErr.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.InvoicePerShopRepository)
			tc.beforeTest(mockRepo)
			service := service.NewInvoicePerShopService(&service.InvoicePerShopSConfig{
				InvoicePerShopRepo: mockRepo,
			})

			result, err := service.RefundRequest(code, userId)

			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.result, result)
		})
	}
}

func TestUpdateStatusCRONJob(t *testing.T) {
	t.Run("should return number of calls when called", func(t *testing.T) {
		mockRepo := new(mocks.InvoicePerShopRepository)
		mockRepo.On("UpdateStatusCRONJob").Return(nil)
		service := service.NewInvoicePerShopService(&service.InvoicePerShopSConfig{
			InvoicePerShopRepo: mockRepo,
		})

		service.UpdateStatusCRONJob()

		mockRepo.AssertNumberOfCalls(t, "UpdateStatusCRONJob", 1)
	})
}

func TestAutoReceivedCRONJob(t *testing.T) {
	t.Run("should return number of calls when called", func(t *testing.T) {
		mockRepo := new(mocks.InvoicePerShopRepository)
		mockRepo.On("AutoReceivedCRONJob").Return(nil)
		service := service.NewInvoicePerShopService(&service.InvoicePerShopSConfig{
			InvoicePerShopRepo: mockRepo,
		})

		service.AutoReceivedCRONJob()

		mockRepo.AssertNumberOfCalls(t, "AutoReceivedCRONJob", 1)
	})
}

func TestAutoCompletedCRONJob(t *testing.T) {
	t.Run("should return number of calls when called", func(t *testing.T) {
		mockRepo := new(mocks.InvoicePerShopRepository)
		mockRepo.On("AutoCompletedCRONJob").Return(nil)
		service := service.NewInvoicePerShopService(&service.InvoicePerShopSConfig{
			InvoicePerShopRepo: mockRepo,
		})

		service.AutoCompletedCRONJob()

		mockRepo.AssertNumberOfCalls(t, "AutoCompletedCRONJob", 1)
	})
}

package service_test

import (
	"kedai/backend/be-kedai/internal/common/constant"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/model"
	"kedai/backend/be-kedai/internal/domain/order/service"
	productDto "kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	var emptyString = ""
	type input struct {
		req        dto.TransactionReviewRequest
		beforeTest func(mockTransactionService *mocks.TransactionService, mockInvoicePerShopService *mocks.InvoicePerShopService, mockTransactionReviewRepo *mocks.TransactionReviewRepository)
	}

	type expected struct {
		data *model.TransactionReview
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
				req: dto.TransactionReviewRequest{
					UserId:        1,
					TransactionId: 1,
				},
				beforeTest: func(mockTransactionService *mocks.TransactionService, mockInvoicePerShopService *mocks.InvoicePerShopService, mockTransactionReviewRepo *mocks.TransactionReviewRepository) {
					mockTransactionService.On("GetByID", 1).Return(nil, commonErr.ErrTransactionNotFound)
				},
			},
			expected: expected{
				data: nil,
				err:  commonErr.ErrTransactionNotFound,
			},
		},
		{
			description: "should return error when transaction user id not match with request user id",
			input: input{
				req: dto.TransactionReviewRequest{
					UserId:        1,
					TransactionId: 1,
				},
				beforeTest: func(mockTransactionService *mocks.TransactionService, mockInvoicePerShopService *mocks.InvoicePerShopService, mockTransactionReviewRepo *mocks.TransactionReviewRepository) {
					mockTransactionService.On("GetByID", 1).Return(&model.Transaction{
						UserID: 2,
					}, nil)
				},
			},
			expected: expected{
				data: nil,
				err:  commonErr.ErrTransactionNotFound,
			},
		},
		{
			description: "should return error when GetByID return error",
			input: input{
				req: dto.TransactionReviewRequest{
					UserId:        1,
					TransactionId: 1,
				},
				beforeTest: func(mockTransactionService *mocks.TransactionService, mockInvoicePerShopService *mocks.InvoicePerShopService, mockTransactionReviewRepo *mocks.TransactionReviewRepository) {
					mockTransactionService.On("GetByID", 1).Return(&model.Transaction{
						UserID:    1,
						ID:        1,
						InvoiceID: 1,
					}, nil)
					mockInvoicePerShopService.On("GetByID", 1).Return(nil, commonErr.ErrInvoiceNotFound)
				},
			},
			expected: expected{
				data: nil,
				err:  commonErr.ErrInvoiceNotFound,
			},
		},
		{
			description: "should return error when invoice status is not completed",
			input: input{
				req: dto.TransactionReviewRequest{
					UserId:        1,
					TransactionId: 1,
				},
				beforeTest: func(mockTransactionService *mocks.TransactionService, mockInvoicePerShopService *mocks.InvoicePerShopService, mockTransactionReviewRepo *mocks.TransactionReviewRepository) {
					mockTransactionService.On("GetByID", 1).Return(&model.Transaction{
						UserID:    1,
						ID:        1,
						InvoiceID: 1,
					}, nil)
					mockInvoicePerShopService.On("GetByID", 1).Return(&model.InvoicePerShop{
						Status: "not completed",
						ID:     1,
					}, nil)
				},
			},
			expected: expected{
				data: nil,
				err:  commonErr.ErrInvoiceNotCompleted,
			},
		},
		{
			description: "should return error when Create return error",
			input: input{
				req: dto.TransactionReviewRequest{
					UserId:        1,
					TransactionId: 1,
					Rating:        5,
					Description:   &emptyString,
				},
				beforeTest: func(mockTransactionService *mocks.TransactionService, mockInvoicePerShopService *mocks.InvoicePerShopService, mockTransactionReviewRepo *mocks.TransactionReviewRepository) {
					mockTransactionService.On("GetByID", 1).Return(&model.Transaction{
						UserID:    1,
						ID:        1,
						InvoiceID: 1,
					}, nil)
					mockInvoicePerShopService.On("GetByID", 1).Return(&model.InvoicePerShop{
						Status: constant.TransactionStatusCompleted,
						ID:     1,
					}, nil)
					mockTransactionReviewRepo.On("Create", &model.TransactionReview{
						TransactionId: 1,
						Rating:        5,
					}).Return(nil, commonErr.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  commonErr.ErrInternalServerError,
			},
		},
		{
			description: "should return review and nil error when Create success",
			input: input{
				req: dto.TransactionReviewRequest{
					UserId:        1,
					TransactionId: 1,
					Rating:        5,
					Description:   &emptyString,
				},
				beforeTest: func(mockTransactionService *mocks.TransactionService, mockInvoicePerShopService *mocks.InvoicePerShopService, mockTransactionReviewRepo *mocks.TransactionReviewRepository) {
					mockTransactionService.On("GetByID", 1).Return(&model.Transaction{
						UserID:    1,
						ID:        1,
						InvoiceID: 1,
					}, nil)
					mockInvoicePerShopService.On("GetByID", 1).Return(&model.InvoicePerShop{
						Status: constant.TransactionStatusCompleted,
						ID:     1,
					}, nil)
					mockTransactionReviewRepo.On("Create", &model.TransactionReview{
						TransactionId: 1,
						Rating:        5,
					}).Return(&model.TransactionReview{TransactionId: 1, Rating: 5}, nil)
				},
			},
			expected: expected{
				data: &model.TransactionReview{TransactionId: 1, Rating: 5},
				err:  nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockTransactionService := new(mocks.TransactionService)
			mockInvoicePerShopService := new(mocks.InvoicePerShopService)
			mockTransactionReviewRepo := new(mocks.TransactionReviewRepository)

			c.input.beforeTest(mockTransactionService, mockInvoicePerShopService, mockTransactionReviewRepo)

			transactionReviewService := service.NewTransactionReviewService(&service.TransactionReviewSConfig{
				TransactionReviewRepo: mockTransactionReviewRepo,
				TransactionService:    mockTransactionService,
				InvoicePerShopService: mockInvoicePerShopService,
			})

			data, err := transactionReviewService.Create(c.input.req)

			assert.Equal(t, c.expected.err, err)
			assert.Equal(t, c.expected.data, data)
		})
	}

}

func TestGetReviews(t *testing.T) {
	type input struct {
		req        productDto.GetReviewRequest
		beforeTest func(mockTransactionReviewRepo *mocks.TransactionReviewRepository, mockProductService *mocks.ProductService)
	}

	type expected struct {
		data *commonDto.PaginationResponse
		err  error
	}

	cases := []struct {
		description string
		input       input
		expected    expected
	}{
		{
			description: "should return error when GetByCode return error",
			input: input{
				req: productDto.GetReviewRequest{
					ProductCode: "",
				},
				beforeTest: func(mockTransactionReviewRepo *mocks.TransactionReviewRepository, mockProductService *mocks.ProductService) {
					mockProductService.On("GetByCode", "").Return(nil, commonErr.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  commonErr.ErrInternalServerError,
			},
		},

		{
			description: "should return error when GetReviews return error",
			input: input{
				req: productDto.GetReviewRequest{
					ProductCode: "",
				},
				beforeTest: func(mockTransactionReviewRepo *mocks.TransactionReviewRepository, mockProductService *mocks.ProductService) {
					mockProductService.On("GetByCode", "").Return(&productDto.ProductDetail{}, nil)
					mockTransactionReviewRepo.On("GetReviews", productDto.GetReviewRequest{}).Return(nil, int64(0), 0, commonErr.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  commonErr.ErrInternalServerError,
			},
		},
		{
			description: "should return pagination response and nil error when GetReviews success",
			input: input{
				req: productDto.GetReviewRequest{
					ProductCode: "",
				},
				beforeTest: func(mockTransactionReviewRepo *mocks.TransactionReviewRepository, mockProductService *mocks.ProductService) {
					mockProductService.On("GetByCode", "").Return(&productDto.ProductDetail{}, nil)
					mockTransactionReviewRepo.On("GetReviews", productDto.GetReviewRequest{}).Return([]*model.TransactionReview{}, int64(0), 0, nil)
				},
			},
			expected: expected{
				data: &commonDto.PaginationResponse{
					Data: []*dto.ReviewResponse{},
				},
				err: nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockTransactionReviewRepo := new(mocks.TransactionReviewRepository)
			mockProductService := new(mocks.ProductService)

			c.input.beforeTest(mockTransactionReviewRepo, mockProductService)

			transactionReviewService := service.NewTransactionReviewService(&service.TransactionReviewSConfig{
				TransactionReviewRepo: mockTransactionReviewRepo,
				ProductService:        mockProductService,
			})

			data, err := transactionReviewService.GetReviews(c.input.req)

			assert.Equal(t, c.expected.err, err)
			assert.Equal(t, c.expected.data, data)
		})
	}

}

func TestGetReviewStats(t *testing.T) {

	type input struct {
		req        string
		beforeTest func(mockTransactionReviewRepo *mocks.TransactionReviewRepository, mockProductService *mocks.ProductService)
	}

	type expected struct {
		data *productDto.GetReviewStatsResponse
		err  error
	}

	cases := []struct {
		description string
		input       input
		expected    expected
	}{
		{
			description: "should return error when GetReviewStats return error",
			input: input{
				req: "",
				beforeTest: func(mockTransactionReviewRepo *mocks.TransactionReviewRepository, mockProductService *mocks.ProductService) {
					mockProductService.On("GetByCode", "").Return(&productDto.ProductDetail{}, nil)
					mockTransactionReviewRepo.On("GetReviewStats", "").Return(nil, commonErr.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  commonErr.ErrInternalServerError,
			},
		},
		{
			description: "should return error when GetByCode return error",
			input: input{
				req: "",
				beforeTest: func(mockTransactionReviewRepo *mocks.TransactionReviewRepository, mockProductService *mocks.ProductService) {
					mockProductService.On("GetByCode", "").Return(nil, commonErr.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  commonErr.ErrInternalServerError,
			},
		},
		{
			description: "should return review stats and nil error when GetReviewStats success",
			input: input{
				req: "",
				beforeTest: func(mockTransactionReviewRepo *mocks.TransactionReviewRepository, mockProductService *mocks.ProductService) {
					mockProductService.On("GetByCode", "").Return(&productDto.ProductDetail{}, nil)
					mockTransactionReviewRepo.On("GetReviewStats", "").Return(&productDto.GetReviewStatsResponse{}, nil)
				},
			},
			expected: expected{
				data: &productDto.GetReviewStatsResponse{},
				err:  nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockTransactionReviewRepo := new(mocks.TransactionReviewRepository)
			mockProductService := new(mocks.ProductService)

			c.input.beforeTest(mockTransactionReviewRepo, mockProductService)

			transactionReviewService := service.NewTransactionReviewService(&service.TransactionReviewSConfig{
				TransactionReviewRepo: mockTransactionReviewRepo,
				ProductService:        mockProductService,
			})

			data, err := transactionReviewService.GetReviewStats(c.input.req)

			assert.Equal(t, c.expected.err, err)
			assert.Equal(t, c.expected.data, data)
		})
	}

}

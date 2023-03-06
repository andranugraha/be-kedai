package service_test

import (
	"kedai/backend/be-kedai/internal/common/constant"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/model"
	"kedai/backend/be-kedai/internal/domain/order/service"
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

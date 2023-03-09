package service_test

import (
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/model"
	"kedai/backend/be-kedai/internal/domain/order/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TransactionReviewGetByID(t *testing.T) {
	type input struct {
		req        int
		beforeTest func(mockTransactionRepo *mocks.TransactionRepository)
	}

	type expected struct {
		data *model.Transaction
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
				beforeTest: func(mockTransactionRepo *mocks.TransactionRepository) {
					mockTransactionRepo.On("GetByID", 1).Return(nil, commonErr.ErrTransactionNotFound)
				},
			},
			expected: expected{
				data: nil,
				err:  commonErr.ErrTransactionNotFound,
			},
		},
		{
			description: "should return transaction when GetByID return transaction",
			input: input{
				req: 1,
				beforeTest: func(mockTransactionRepo *mocks.TransactionRepository) {
					mockTransactionRepo.On("GetByID", 1).Return(&model.Transaction{}, nil)
				},
			},
			expected: expected{
				data: &model.Transaction{},
				err:  nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			mockTransactionRepo := new(mocks.TransactionRepository)
			c.input.beforeTest(mockTransactionRepo)

			service := service.NewTransactionService(&service.TransactionSConfig{
				TransactionRepo: mockTransactionRepo,
			})

			data, err := service.GetByID(c.input.req)

			assert.Equal(t, c.expected.data, data)
			assert.Equal(t, c.expected.err, err)
		})
	}

}

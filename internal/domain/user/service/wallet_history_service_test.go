package service_test

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWalletHistoryById(t *testing.T) {
	var (
		userId   = 1
		walletId = 1
		wallet   = &model.Wallet{
			ID: 1,
		}
		history  = []*model.WalletHistory{}
	)
	type input struct {
		userId     int
		err        error
		beforeTest func(*mocks.WalletHistoryRepository, *mocks.WalletService)
	}
	type expected struct {
		result []*model.WalletHistory
		err    error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return list of wallet transaction histories when success",
			input: input{
				userId: userId,
				err:    nil,
				beforeTest: func(whr *mocks.WalletHistoryRepository, ws *mocks.WalletService) {
					ws.On("GetWalletByUserID", userId).Return(wallet, nil)
					whr.On("GetWalletHistoryById", walletId).Return(history, nil)
				},
			},
			expected: expected{
				result: history,
				err:    nil,
			},
		},
		{
			description: "should return error when wallet no found",
			input: input{
				userId: userId,
				err:    errs.ErrWalletDoesNotExist,
				beforeTest: func(whr *mocks.WalletHistoryRepository, ws *mocks.WalletService) {
					ws.On("GetWalletByUserID", userId).Return(nil, errs.ErrWalletDoesNotExist)
				},
			},
			expected: expected{
				result: nil,
				err:    errs.ErrWalletDoesNotExist,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockWalletService := new(mocks.WalletService)
			mockWalletHistoryRepo := new(mocks.WalletHistoryRepository)
			tc.beforeTest(mockWalletHistoryRepo, mockWalletService)
			service := service.NewWalletHistoryService(&service.WalletHistorySConfig{
				WalletHistoryRepository: mockWalletHistoryRepo,
				WalletService: mockWalletService,
			})

			result, err := service.GetWalletHistoryById(tc.input.userId)

			assert.Equal(t, tc.expected.err, err)
			assert.Equal(t, tc.expected.result, result)
		})
	}
}

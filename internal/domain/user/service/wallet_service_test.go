package service_test

import (
	"errors"
	errRes "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterWallet(t *testing.T) {
	var (
		wallet = &model.Wallet{
			UserID:  1,
			Pin:     "123456",
			Balance: 0,
		}
	)

	tests := []struct {
		name    string
		want    *model.Wallet
		wantErr error
	}{
		{
			name:    "should return wallet when wallet registered successfully",
			want:    wallet,
			wantErr: nil,
		},
		{
			name:    "should return error when wallet already exist",
			want:    nil,
			wantErr: errRes.ErrWalletAlreadyExist,
		},
		{
			name:    "should return error when create wallet failed",
			want:    nil,
			wantErr: errRes.ErrInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := mocks.NewWalletRepository(t)
			mockRepo.On("Create", mock.Anything).Return(test.want, test.wantErr)
			walletService := service.NewWalletService(&service.WalletSConfig{
				WalletRepo: mockRepo,
			})

			got, err := walletService.RegisterWallet(wallet.UserID, wallet.Pin)

			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, test.wantErr, err)
		})
	}
}

func TestGetWalletByUserID(t *testing.T) {
	type input struct {
		userId int
		data   *model.Wallet
		err    error
	}
	type expected struct {
		wallet *model.Wallet
		err    error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when failed to get wallet",
			input: input{
				userId: 1,
				data:   nil,
				err:    errors.New("failed to get wallet"),
			},
			expected: expected{
				wallet: nil,
				err:    errors.New("failed to get wallet"),
			},
		},
		{
			description: "should return wallet data when successed fetching user wallet",
			input: input{
				userId: 1,
				data:   &model.Wallet{},
				err:    nil,
			},
			expected: expected{
				wallet: &model.Wallet{},
				err:    nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			walletRepo := mocks.NewWalletRepository(t)
			walletRepo.On("GetByUserID", tc.input.userId).Return(tc.input.data, tc.input.err)
			walletService := service.NewWalletService(&service.WalletSConfig{
				WalletRepo: walletRepo,
			})

			actualWallet, actualErr := walletService.GetWalletByUserID(tc.userId)

			assert.Equal(t, tc.expected.wallet, actualWallet)
			assert.Equal(t, tc.expected.err, actualErr)
		})
	}
}

func TestTopUp(t *testing.T) {
	var (
		userId  = 1
		history = &model.WalletHistory{
			ID:        0,
			Type:      "Top-up",
			Reference: "15602",
			Amount:    50000,
			WalletId:  1,
		}
		wallet = &model.Wallet{
			ID: 1,
		}
		req = dto.TopUpRequest{
			Amount:     50000,
			TxnId:      "15602",
			Signature:  "027e361f8776a9fb6f25961687e3cf0879af6fdfe00e03335fbb5764a3763d40",
			CardNumber: "2793765051084376",
		}
		invalidReq = dto.TopUpRequest{
			Amount:     50000,
			TxnId:      "15602",
			Signature:  "d714a01f755b9f5f2c6fdb1d41107cace0e220154b3edc8603c0800e32b479e0",
			CardNumber: "2793765051084376",
		}
	)
	type input struct {
		userId      int
		request     dto.TopUpRequest
		history     *model.WalletHistory
		wallet      *model.Wallet
		err         error
		beforeTests func(mockWalletRepo *mocks.WalletRepository)
	}

	type expected struct {
		data *model.WalletHistory
		err  error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return wallet top-up history when success",
			input: input{
				userId:  userId,
				request: req,
				history: history,
				wallet:  wallet,
				err:     nil,
				beforeTests: func(mockWalletRepo *mocks.WalletRepository) {
					mockWalletRepo.On("GetByUserID", userId).Return(wallet, nil)
					mockWalletRepo.On("TopUp", history, wallet).Return(history, nil)
				},
			},
			expected: expected{
				data: history,
				err:  nil,
			},
		},
		{
			description: "should return error when user wallet does not exist",
			input: input{
				userId:  userId,
				request: req,
				history: history,
				wallet:  nil,
				err:     errRes.ErrWalletDoesNotExist,
				beforeTests: func(mockWalletRepo *mocks.WalletRepository) {
					mockWalletRepo.On("GetByUserID", userId).Return(nil, errRes.ErrWalletDoesNotExist)
				},
			},
			expected: expected{
				data: nil,
				err:  errRes.ErrWalletDoesNotExist,
			},
		},
		{
			description: "should return error when signature is invalid",
			input: input{
				userId:  userId,
				request: invalidReq,
				history: nil,
				wallet:  nil,
				err:     errRes.ErrInvalidSignature,
				beforeTests: func(mockWalletRepo *mocks.WalletRepository) {
				},
			},
			expected: expected{
				data: nil,
				err:  errRes.ErrInvalidSignature,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				userId:  1,
				request: req,
				history: history,
				wallet:  wallet,
				err:     errRes.ErrInternalServerError,
				beforeTests: func(mockWalletRepo *mocks.WalletRepository) {
					mockWalletRepo.On("GetByUserID", userId).Return(wallet, nil)
					mockWalletRepo.On("TopUp", history, wallet).Return(nil, errRes.ErrInternalServerError)
				},
			},
			expected: expected{
				data: nil,
				err:  errRes.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockWalletRepo := mocks.NewWalletRepository(t)
			tc.beforeTests(mockWalletRepo)
			service := service.NewWalletService(&service.WalletSConfig{
				WalletRepo: mockWalletRepo,
			})

			result, err := service.TopUp(tc.input.userId, tc.input.request)

			assert.Equal(t, tc.expected.err, err)
			assert.Equal(t, tc.expected.data, result)
		})
	}
}

func TestCheckIsWalletBloced(t *testing.T) {
	tests := []struct {
		name       string
		wantErr    error
		beforeTest func(mockWalletRepo *mocks.WalletRepository, walletCache *mocks.WalletCache)
	}{
		{
			name:    "should return error when wallet is blocked",
			wantErr: errRes.ErrWalletTemporarilyBlocked,
			beforeTest: func(mockWalletRepo *mocks.WalletRepository, walletCache *mocks.WalletCache) {
				mockWalletRepo.On("GetByUserID", 1).Return(&model.Wallet{
					ID: 1,
				}, nil)
				walletCache.On("CheckIsWalletBlocked", 1).Return(errRes.ErrWalletTemporarilyBlocked)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockWalletRepo := mocks.NewWalletRepository(t)
			mockWalletCache := mocks.NewWalletCache(t)
			test.beforeTest(mockWalletRepo, mockWalletCache)
			walletService := service.NewWalletService(&service.WalletSConfig{
				WalletRepo:  mockWalletRepo,
				WalletCache: mockWalletCache,
			})

			err := walletService.CheckIsWalletBlocked(1)

			assert.Equal(t, test.wantErr, err)
		})
	}
}

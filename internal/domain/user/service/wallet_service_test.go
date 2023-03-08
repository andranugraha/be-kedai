package service_test

import (
	"errors"
	errRes "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/service"
	"kedai/backend/be-kedai/internal/utils/hash"
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
			Amount: 50000,
			TxnId:  "15602",
		}
	)
	type input struct {
		userId      int
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
			description: "should return error when internal server error",
			input: input{
				userId:  1,
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

			result, err := service.TopUp(tc.input.userId, req)

			assert.Equal(t, tc.expected.err, err)
			assert.Equal(t, tc.expected.data, result)
		})
	}
}

func TestRequestPinChange(t *testing.T) {
	type input struct {
		userID  int
		request *dto.ChangePinRequest
	}
	type expected struct {
		err error
	}

	var (
		userID          = 1
		oldPin          = "123456"
		hashedOldPin, _ = hash.HashAndSalt(oldPin)
		// codeLength       = 6
		// verificationCode = "a1b2c3"
		// email            = "test@email.com"
	)

	tests := []struct {
		description string
		beforeTest  func(*mocks.WalletRepository, *mocks.UserService, *mocks.WalletCache, *mocks.RandomUtils, *mocks.MailUtils)
		input
		expected
	}{
		{
			description: "should return error when failed to get wallet",
			beforeTest: func(wr *mocks.WalletRepository, us *mocks.UserService, wc *mocks.WalletCache, ru *mocks.RandomUtils, mu *mocks.MailUtils) {
				wr.On("GetByUserID", userID).Return(nil, errors.New("failed to get wallet"))
			},
			input: input{
				userID:  userID,
				request: &dto.ChangePinRequest{},
			},
			expected: expected{
				err: errors.New("failed to get wallet"),
			},
		},
		{
			description: "should return error when current pin is invalid",
			beforeTest: func(wr *mocks.WalletRepository, us *mocks.UserService, wc *mocks.WalletCache, ru *mocks.RandomUtils, mu *mocks.MailUtils) {
				wr.On("GetByUserID", userID).Return(&model.Wallet{Pin: hashedOldPin}, nil)
			},
			input: input{
				userID: userID,
				request: &dto.ChangePinRequest{
					CurrentPin: "789012",
				},
			},
			expected: expected{
				err: errRes.ErrPinMismatch,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			walletRepo := mocks.NewWalletRepository(t)
			randomUtils := mocks.NewRandomUtils(t)
			mailUtils := mocks.NewMailUtils(t)
			userService := mocks.NewUserService(t)
			walletCache := mocks.NewWalletCache(t)
			tc.beforeTest(walletRepo, userService, walletCache, randomUtils, mailUtils)
			walletService := service.NewWalletService(&service.WalletSConfig{
				WalletRepo:  walletRepo,
				UserService: userService,
				WalletCache: walletCache,
				RandomUtils: randomUtils,
				MailUtils:   mailUtils,
			})

			err := walletService.RequestPinChange(tc.input.userID, tc.input.request)

			assert.Equal(t, tc.expected.err, err)
		})
	}
}

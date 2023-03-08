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
		userID     int
		request    *dto.ChangePinRequest
		beforeTest func(*mocks.WalletRepository, *mocks.UserService, *mocks.WalletCache, *mocks.RandomUtils, *mocks.MailUtils)
	}
	type expected struct {
		err error
	}

	var (
		userID           = 1
		oldPin           = "123456"
		hashedOldPin, _  = hash.HashAndSalt(oldPin)
		codeLength       = 6
		verificationCode = "a1b2c3"
		email            = "test@email.com"
	)

	tests := []struct {
		description string

		input
		expected
	}{
		{
			description: "should return error when failed to get wallet",
			input: input{
				userID:  userID,
				request: &dto.ChangePinRequest{},
				beforeTest: func(wr *mocks.WalletRepository, us *mocks.UserService, wc *mocks.WalletCache, ru *mocks.RandomUtils, mu *mocks.MailUtils) {
					wr.On("GetByUserID", userID).Return(nil, errors.New("failed to get wallet"))
				},
			},
			expected: expected{
				err: errors.New("failed to get wallet"),
			},
		},
		{
			description: "should return error when current pin is invalid",
			input: input{
				userID: userID,
				request: &dto.ChangePinRequest{
					CurrentPin: "789012",
				},
				beforeTest: func(wr *mocks.WalletRepository, us *mocks.UserService, wc *mocks.WalletCache, ru *mocks.RandomUtils, mu *mocks.MailUtils) {
					wr.On("GetByUserID", userID).Return(&model.Wallet{Pin: hashedOldPin}, nil)
				},
			},
			expected: expected{
				err: errRes.ErrPinMismatch,
			},
		},
		{
			description: "should return error when failed to store token",
			input: input{
				userID: userID,
				request: &dto.ChangePinRequest{
					CurrentPin: oldPin,
					NewPin:     "098765",
				},
				beforeTest: func(wr *mocks.WalletRepository, us *mocks.UserService, wc *mocks.WalletCache, ru *mocks.RandomUtils, mu *mocks.MailUtils) {
					wr.On("GetByUserID", userID).Return(&model.Wallet{Pin: hashedOldPin}, nil)
					ru.On("GenerateAlphanumericString", codeLength).Return(verificationCode)
					wc.On("StorePinAndVerificationCode", userID, mock.Anything, verificationCode).Return(errors.New("failed to store token"))
				},
			},
			expected: expected{
				err: errors.New("failed to store token"),
			},
		},
		{
			description: "should return error when failed to get user",
			input: input{
				userID: userID,
				request: &dto.ChangePinRequest{
					CurrentPin: oldPin,
					NewPin:     "098765",
				},
				beforeTest: func(wr *mocks.WalletRepository, us *mocks.UserService, wc *mocks.WalletCache, ru *mocks.RandomUtils, mu *mocks.MailUtils) {
					wr.On("GetByUserID", userID).Return(&model.Wallet{Pin: hashedOldPin}, nil)
					ru.On("GenerateAlphanumericString", codeLength).Return(verificationCode)
					wc.On("StorePinAndVerificationCode", userID, mock.Anything, verificationCode).Return(nil)
					us.On("GetByID", userID).Return(nil, errors.New("failed to get user"))
				},
			},
			expected: expected{
				err: errors.New("failed to get user"),
			},
		},
		{
			description: "should return error when failed to send email",
			input: input{
				userID: userID,
				request: &dto.ChangePinRequest{
					CurrentPin: oldPin,
					NewPin:     "098765",
				},
				beforeTest: func(wr *mocks.WalletRepository, us *mocks.UserService, wc *mocks.WalletCache, ru *mocks.RandomUtils, mu *mocks.MailUtils) {
					wr.On("GetByUserID", userID).Return(&model.Wallet{Pin: hashedOldPin}, nil)
					ru.On("GenerateAlphanumericString", codeLength).Return(verificationCode)
					wc.On("StorePinAndVerificationCode", userID, mock.Anything, verificationCode).Return(nil)
					us.On("GetByID", userID).Return(&model.User{Email: email}, nil)
					mu.On("SendUpdatePinEmail", email, verificationCode).Return(errors.New("failed to send email"))
				},
			},
			expected: expected{
				err: errors.New("failed to send email"),
			},
		},
		{
			description: "should return nil when suceed to send email",
			input: input{
				userID: userID,
				request: &dto.ChangePinRequest{
					CurrentPin: oldPin,
					NewPin:     "098765",
				},
				beforeTest: func(wr *mocks.WalletRepository, us *mocks.UserService, wc *mocks.WalletCache, ru *mocks.RandomUtils, mu *mocks.MailUtils) {
					wr.On("GetByUserID", userID).Return(&model.Wallet{Pin: hashedOldPin}, nil)
					ru.On("GenerateAlphanumericString", codeLength).Return(verificationCode)
					wc.On("StorePinAndVerificationCode", userID, mock.Anything, verificationCode).Return(nil)
					us.On("GetByID", userID).Return(&model.User{Email: email}, nil)
					mu.On("SendUpdatePinEmail", email, verificationCode).Return(nil)
				},
			},
			expected: expected{
				err: nil,
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

func TestCompletePinChange(t *testing.T) {
	type input struct {
		userID  int
		request *dto.CompleteChangePinRequest
	}
	type expected struct {
		err error
	}

	var (
		userID           = 1
		verificationCode = "a1b2c3"
		newPin           = "098765"
		hashedPin, _     = hash.HashAndSalt(newPin)
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.WalletRepository, *mocks.WalletCache)
		expected
	}{
		{
			description: "should return error when failed to find token",
			input: input{
				userID:  userID,
				request: &dto.CompleteChangePinRequest{},
			},
			beforeTest: func(wr *mocks.WalletRepository, wc *mocks.WalletCache) {
				wc.On("FindPinAndVerificationCode", userID).Return("", "", errors.New("failed to find token"))
			},
			expected: expected{
				err: errors.New("failed to find token"),
			},
		},
		{
			description: "should return error when given invalid verification code",
			input: input{
				userID: userID,
				request: &dto.CompleteChangePinRequest{
					VerificationCode: "d4e5f6",
				},
			},
			beforeTest: func(wr *mocks.WalletRepository, wc *mocks.WalletCache) {
				wc.On("FindPinAndVerificationCode", userID).Return(hashedPin, verificationCode, nil)
			},
			expected: expected{
				err: errRes.ErrIncorrectVerificationCode,
			},
		},
		{
			description: "should return error when failed to change pin",
			input: input{
				userID: userID,
				request: &dto.CompleteChangePinRequest{
					VerificationCode: verificationCode,
				},
			},
			beforeTest: func(wr *mocks.WalletRepository, wc *mocks.WalletCache) {
				wc.On("FindPinAndVerificationCode", userID).Return(hashedPin, verificationCode, nil)
				wr.On("ChangePin", userID, hashedPin).Return(errors.New("failed to change pin"))
			},
			expected: expected{
				err: errors.New("failed to change pin"),
			},
		},
		{
			description: "should return nil when succeed to change pin",
			input: input{
				userID: userID,
				request: &dto.CompleteChangePinRequest{
					VerificationCode: verificationCode,
				},
			},
			beforeTest: func(wr *mocks.WalletRepository, wc *mocks.WalletCache) {
				wc.On("FindPinAndVerificationCode", userID).Return(hashedPin, verificationCode, nil)
				wr.On("ChangePin", userID, hashedPin).Return(nil)
				wc.On("DeletePinAndVerificationCode", userID).Return(nil)
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			walletRepo := mocks.NewWalletRepository(t)
			walletCache := mocks.NewWalletCache(t)
			tc.beforeTest(walletRepo, walletCache)
			walletService := service.NewWalletService(&service.WalletSConfig{
				WalletRepo:  walletRepo,
				WalletCache: walletCache,
			})

			err := walletService.CompletePinChange(tc.input.userID, tc.input.request)

			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestRequestPinReset(t *testing.T) {
	type input struct {
		userID int
	}
	type expected struct {
		err error
	}

	var (
		userID           = 1
		email            = "test@email.com"
		codeLength       = 6
		verificationCode = "a1b2c3"
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.UserService, *mocks.RandomUtils, *mocks.MailUtils, *mocks.WalletCache)
		expected
	}{
		{
			description: "should return error when failed to get user",
			input: input{
				userID: userID,
			},
			beforeTest: func(us *mocks.UserService, ru *mocks.RandomUtils, mu *mocks.MailUtils, wc *mocks.WalletCache) {
				us.On("GetByID", userID).Return(nil, errors.New("failed to get user"))
			},
			expected: expected{
				err: errors.New("failed to get user"),
			},
		},
		{
			description: "should return error when failed to store token",
			input: input{
				userID: userID,
			},
			beforeTest: func(us *mocks.UserService, ru *mocks.RandomUtils, mu *mocks.MailUtils, wc *mocks.WalletCache) {
				us.On("GetByID", userID).Return(&model.User{ID: userID, Email: email}, nil)
				ru.On("GenerateAlphanumericString", codeLength).Return(verificationCode)
				wc.On("StoreResetPinToken", userID, verificationCode).Return(errors.New("failed to store token"))
			},
			expected: expected{
				err: errors.New("failed to store token"),
			},
		},
		{
			description: "should return error when failed to send email",
			input: input{
				userID: userID,
			},
			beforeTest: func(us *mocks.UserService, ru *mocks.RandomUtils, mu *mocks.MailUtils, wc *mocks.WalletCache) {
				us.On("GetByID", userID).Return(&model.User{ID: userID, Email: email}, nil)
				ru.On("GenerateAlphanumericString", codeLength).Return(verificationCode)
				wc.On("StoreResetPinToken", userID, verificationCode).Return(nil)
				mu.On("SendResetPinEmail", email, verificationCode).Return(errors.New("failed to send email"))
			},
			expected: expected{
				err: errors.New("failed to send email"),
			},
		},
		{
			description: "should return nil when succeed to send email",
			input: input{
				userID: userID,
			},
			beforeTest: func(us *mocks.UserService, ru *mocks.RandomUtils, mu *mocks.MailUtils, wc *mocks.WalletCache) {
				us.On("GetByID", userID).Return(&model.User{ID: userID, Email: email}, nil)
				ru.On("GenerateAlphanumericString", codeLength).Return(verificationCode)
				wc.On("StoreResetPinToken", userID, verificationCode).Return(nil)
				mu.On("SendResetPinEmail", email, verificationCode).Return(nil)
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			userService := mocks.NewUserService(t)
			randomUtils := mocks.NewRandomUtils(t)
			mailUtils := mocks.NewMailUtils(t)
			walletCache := mocks.NewWalletCache(t)
			tc.beforeTest(userService, randomUtils, mailUtils, walletCache)
			walletService := service.NewWalletService(&service.WalletSConfig{
				UserService: userService,
				RandomUtils: randomUtils,
				MailUtils:   mailUtils,
				WalletCache: walletCache,
			})

			err := walletService.RequestPinReset(tc.input.userID)

			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestCompletePinReset(t *testing.T) {
	type input struct {
		userID  int
		request *dto.CompleteResetPinRequest
	}
	type expected struct {
		err error
	}

	var (
		userID = 1
		token  = "a1b2c3"
		newPin = "098765"
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.WalletRepository, *mocks.WalletCache)
		expected
	}{
		{
			description: "should return error when failed to find token",
			input: input{
				userID: userID,
				request: &dto.CompleteResetPinRequest{
					Token: token,
				},
			},
			beforeTest: func(wr *mocks.WalletRepository, wc *mocks.WalletCache) {
				wc.On("FindResetPinToken", token).Return(errors.New("failed to find token"))
			},
			expected: expected{
				err: errors.New("failed to find token"),
			},
		},
		{
			description: "should return error when failed to change pin",
			input: input{
				userID: userID,
				request: &dto.CompleteResetPinRequest{
					Token:  token,
					NewPin: newPin,
				},
			},
			beforeTest: func(wr *mocks.WalletRepository, wc *mocks.WalletCache) {
				wc.On("FindResetPinToken", token).Return(nil)
				wr.On("ChangePin", userID, mock.Anything).Return(errors.New("failed to change pin"))
			},
			expected: expected{
				err: errors.New("failed to change pin"),
			},
		},
		{
			description: "should return nil when succeed to change pin",
			input: input{
				userID: userID,
				request: &dto.CompleteResetPinRequest{
					Token:  token,
					NewPin: newPin,
				},
			},
			beforeTest: func(wr *mocks.WalletRepository, wc *mocks.WalletCache) {
				wc.On("FindResetPinToken", token).Return(nil)
				wr.On("ChangePin", userID, mock.Anything).Return(nil)
				wc.On("DeleteResetPinToken", token).Return(nil)
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			walletRepo := mocks.NewWalletRepository(t)
			walletCache := mocks.NewWalletCache(t)
			tc.beforeTest(walletRepo, walletCache)
			walletService := service.NewWalletService(&service.WalletSConfig{
				WalletRepo:  walletRepo,
				WalletCache: walletCache,
			})

			err := walletService.CompletePinReset(tc.input.userID, tc.input.request)

			assert.Equal(t, tc.expected.err, err)
		})
	}
}

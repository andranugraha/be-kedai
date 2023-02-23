package service_test

import (
	"errors"
	errRes "kedai/backend/be-kedai/internal/common/error"
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

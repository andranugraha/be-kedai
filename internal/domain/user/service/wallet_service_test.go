package service_test

import (
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

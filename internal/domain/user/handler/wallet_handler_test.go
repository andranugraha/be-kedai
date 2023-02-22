package handler_test

import (
	"encoding/json"
	"kedai/backend/be-kedai/internal/common/code"
	errRes "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/handler"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/utils/response"
	testutil "kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterWallet(t *testing.T) {
	var (
		req = &dto.RegisterWalletRequest{
			Pin: "123456",
		}
		userId = 1
		wallet = &model.Wallet{
			UserID:  userId,
			Balance: 0,
		}
	)

	tests := []struct {
		name                  string
		req                   *dto.RegisterWalletRequest
		wantRegisterWallet    *model.Wallet
		wantRegisterWalletErr error
		code                  int
		want                  response.Response
		wantErr               error
		beforeTest            func(*mocks.WalletService)
	}{
		{
			name:                  "should return wallet when wallet registered successfully",
			req:                   req,
			wantRegisterWallet:    wallet,
			wantRegisterWalletErr: nil,
			code:                  http.StatusCreated,
			want: response.Response{
				Code:    code.CREATED,
				Message: "wallet registered successfully",
				Data:    wallet,
			},
			wantErr: nil,
			beforeTest: func(mockWalletService *mocks.WalletService) {
				mockWalletService.On("RegisterWallet", userId, req.Pin).Return(wallet, nil)
			},
		},
		{
			name:                  "should return error when pin less than 6 characters or not numeric",
			req:                   &dto.RegisterWalletRequest{Pin: "12345"},
			wantRegisterWallet:    nil,
			wantRegisterWalletErr: nil,
			code:                  http.StatusBadRequest,
			want: response.Response{
				Code:    code.BAD_REQUEST,
				Message: "Pin must be 6 characters",
				Data:    nil,
			},
			wantErr:    errRes.ErrInvalidPin,
			beforeTest: func(mockWalletService *mocks.WalletService) {},
		},
		{
			name:                  "should return error when wallet already exist",
			req:                   req,
			wantRegisterWallet:    nil,
			wantRegisterWalletErr: errRes.ErrWalletAlreadyExist,
			code:                  http.StatusConflict,
			want: response.Response{
				Code:    code.WALLET_ALREADY_EXIST,
				Message: errRes.ErrWalletAlreadyExist.Error(),
				Data:    nil,
			},
			wantErr: errRes.ErrWalletAlreadyExist,
			beforeTest: func(mockWalletService *mocks.WalletService) {
				mockWalletService.On("RegisterWallet", userId, req.Pin).Return(nil, errRes.ErrWalletAlreadyExist)
			},
		},
		{
			name:                  "should return error when create wallet failed",
			req:                   req,
			wantRegisterWallet:    nil,
			wantRegisterWalletErr: errRes.ErrInternalServerError,
			code:                  http.StatusInternalServerError,
			want: response.Response{
				Code:    code.INTERNAL_SERVER_ERROR,
				Message: errRes.ErrInternalServerError.Error(),
				Data:    nil,
			},
			wantErr: errRes.ErrInternalServerError,
			beforeTest: func(mockWalletService *mocks.WalletService) {
				mockWalletService.On("RegisterWallet", userId, req.Pin).Return(nil, errRes.ErrInternalServerError)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			payload := testutil.MakeRequestBody(test.req)
			jsonRes, _ := json.Marshal(test.want)
			mockWalletService := mocks.NewWalletService(t)
			test.beforeTest(mockWalletService)
			h := handler.New(&handler.HandlerConfig{
				WalletService: mockWalletService,
			})
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", userId)

			c.Request, _ = http.NewRequest("POST", "/v1/users/wallets", payload)
			h.RegisterWallet(c)

			assert.Equal(t, test.code, rec.Code)
			assert.Equal(t, string(jsonRes), rec.Body.String())
		})
	}
}

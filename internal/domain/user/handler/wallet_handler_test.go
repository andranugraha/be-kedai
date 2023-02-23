package handler_test

import (
	"encoding/json"
	"errors"
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

func TestGetWalletByUserID(t *testing.T) {
	type input struct {
		userId int
		data   *model.Wallet
		err    error
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return wallet with status code 200 when successed fetching user wallet",
			input: input{
				userId: 1,
				data: &model.Wallet{
					UserID:  1,
					Balance: 0,
				},
				err: nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data: &model.Wallet{
						UserID:  1,
						Balance: 0,
					},
				},
			},
		},
		{
			description: "should return error with status code 404 when user does not have any wallet yet",
			input: input{
				userId: 1,
				data:   nil,
				err:    errRes.ErrWalletDoesNotExist,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.WALLET_DOES_NOT_EXIST,
					Message: errRes.ErrWalletDoesNotExist.Error(),
					Data:    nil,
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to get user wallet",
			input: input{
				userId: 1,
				data:   nil,
				err:    errRes.ErrInternalServerError,
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errRes.ErrInternalServerError.Error(),
					Data:    nil,
				},
			},
		},
	}

	for _, tc := range cases {
		expectedRes, _ := json.Marshal(tc.expected.response)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Set("userId", tc.input.userId)
		walletService := mocks.NewWalletService(t)
		walletService.On("GetWalletByUserID", tc.input.userId).Return(tc.input.data, tc.input.err)
		cfg := handler.HandlerConfig{
			WalletService: walletService,
		}
		h := handler.New(&cfg)
		c.Request, _ = http.NewRequest("GET", "/users/wallets", nil)

		h.GetWalletByUserID(c)

		assert.Equal(t, tc.expected.statusCode, rec.Code)
		assert.Equal(t, string(expectedRes), rec.Body.String())
	}
}

func TestTopUp(t *testing.T) {
	var (
		userId = 1
		validRequest = dto.TopUpRequest{
			Amount: 50000,
		}
		invalidRequest = dto.TopUpRequest{
			Amount: 5000,
		}
		res = &model.WalletHistory{
			Amount: 50000,
		}
	)
	type input struct {
		userId int
		data    dto.TopUpRequest
		response *model.WalletHistory
		err    error
		beforeTest func(mockWalletService *mocks.WalletService)
	}

	type expected struct {
		statusCode int
		response   response.Response
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return code 200 with top-up wallet history when success",
			input: input{
				userId: 1,
				data: validRequest,
				response: &model.WalletHistory{
					Amount: 50000,
				},
				err: nil,
				beforeTest: func(mockWalletService *mocks.WalletService) {
					mockWalletService.On("TopUp", userId, validRequest.Amount).Return(res, nil)
				},
			},
			expected: expected{
				statusCode: 200,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data: res,
				},
			},
		},
		{
			description: "should return code 400 when input condition doesn't met",
			input: input{
				userId: 1,
				data: invalidRequest,
				response: nil,
				err: errors.New("error"),
				beforeTest: func(mockWalletService *mocks.WalletService) {},
			},
			expected: expected{
				statusCode: 400,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Amount must be greater than 10000",
				},
			},
		},
		{
			description: "should return code 500 when internal server error",
			input: input{
				userId: 1,
				data: dto.TopUpRequest{
					Amount: 50000,
				},
				response: nil,
				err: errRes.ErrInternalServerError,
				beforeTest: func(mockWalletService *mocks.WalletService) {
					mockWalletService.On("TopUp", userId, validRequest.Amount).Return(nil, errRes.ErrInternalServerError)
				},
			},
			expected: expected{
				statusCode: 500,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errRes.ErrInternalServerError.Error(),
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userId)
			walletService := mocks.NewWalletService(t)
			tc.beforeTest(walletService)
			handler := handler.New(&handler.HandlerConfig{
				WalletService: walletService,
			})
			c.Request, _ = http.NewRequest("POST", "/users/wallets/top-up", testutil.MakeRequestBody(tc.data))

			handler.TopUp(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

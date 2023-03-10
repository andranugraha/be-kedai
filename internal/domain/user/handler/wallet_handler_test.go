package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
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
		data   *dto.GetWalletResponse
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
				data: &dto.GetWalletResponse{
					Balance: 0,
				},
				err: nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data: &dto.GetWalletResponse{
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
		walletService.On("GetWalletDetailByUserID", tc.input.userId).Return(tc.input.data, tc.input.err)
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
		userId       = 1
		validRequest = dto.TopUpRequest{
			Amount:     50000,
			TxnId:      "50400",
			Signature:  "1243asdkjaisdw",
			CardNumber: "12388394834",
		}
		invalidRequest = dto.TopUpRequest{
			Amount: 5000,
		}
		res = &model.WalletHistory{
			Amount: 50000,
		}
	)
	type input struct {
		userId     int
		data       dto.TopUpRequest
		response   *model.WalletHistory
		query      string
		err        error
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
				data:   validRequest,
				response: &model.WalletHistory{
					Amount: 50000,
				},
				query: "txnId=50400&amount=50000&cardNumber=12388394834&signature=1243asdkjaisdw",
				err:   nil,
				beforeTest: func(mockWalletService *mocks.WalletService) {
					mockWalletService.On("TopUp", userId, validRequest).Return(res, nil)
				},
			},
			expected: expected{
				statusCode: 200,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    res,
				},
			},
		},
		{
			description: "should return code 422 with error when signature is invalid",
			input: input{
				userId: 1,
				data:   validRequest,
				response: &model.WalletHistory{
					Amount: 50000,
				},
				query: "txnId=50400&amount=50000&cardNumber=12388394834&signature=1243asdkjaisdw",
				err:   errRes.ErrInvalidSignature,
				beforeTest: func(mockWalletService *mocks.WalletService) {
					mockWalletService.On("TopUp", userId, validRequest).Return(nil, errRes.ErrInvalidSignature)
				},
			},
			expected: expected{
				statusCode: 422,
				response: response.Response{
					Code:    code.INVALID_SIGNATURE,
					Message: errRes.ErrInvalidSignature.Error(),
				},
			},
		},
		{
			description: "should return code 404 with error when wallet does not exist",
			input: input{
				userId: 1,
				data:   validRequest,
				response: &model.WalletHistory{
					Amount: 50000,
				},
				query: "txnId=50400&amount=50000&cardNumber=12388394834&signature=1243asdkjaisdw",
				err:   errRes.ErrWalletDoesNotExist,
				beforeTest: func(mockWalletService *mocks.WalletService) {
					mockWalletService.On("TopUp", userId, validRequest).Return(nil, errRes.ErrWalletDoesNotExist)
				},
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
			description: "should return code 400 when input condition doesn't met",
			input: input{
				userId:     1,
				data:       invalidRequest,
				response:   nil,
				err:        errors.New("error"),
				beforeTest: func(mockWalletService *mocks.WalletService) {},
			},
			expected: expected{
				statusCode: 400,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "TxnId is required",
				},
			},
		},
		{
			description: "should return code 500 when internal server error",
			input: input{
				userId:   1,
				data:     validRequest,
				response: nil,
				query:    "txnId=50400&amount=50000&cardNumber=12388394834&signature=1243asdkjaisdw",
				err:      errRes.ErrInternalServerError,
				beforeTest: func(mockWalletService *mocks.WalletService) {
					mockWalletService.On("TopUp", userId, validRequest).Return(nil, errRes.ErrInternalServerError)
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
			c.Request, _ = http.NewRequest("POST", fmt.Sprintf("/users/wallets/top-up?%s", tc.query), nil)

			handler.TopUp(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestRequestWalletPinChange(t *testing.T) {
	type input struct {
		userID  int
		request *dto.ChangePinRequest
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID     = 1
		currentPin = "123456"
		NewPin     = "098765"
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.WalletService)
		expected
	}{
		{
			description: "should return error with status code 400 when given invalid request body",
			input: input{
				userID: userID,
				request: &dto.ChangePinRequest{
					CurrentPin: currentPin,
					NewPin:     "12",
				},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "NewPin must be 6 characters",
				},
			},
			beforeTest: func(ws *mocks.WalletService) {},
		},
		{
			description: "should return error with status code 404 when wallet does not exist",
			input: input{
				userID: userID,
				request: &dto.ChangePinRequest{
					CurrentPin: currentPin,
					NewPin:     NewPin,
				},
			},
			beforeTest: func(ws *mocks.WalletService) {
				ws.On("RequestPinChange", userID, &dto.ChangePinRequest{
					CurrentPin: currentPin,
					NewPin:     NewPin,
				}).Return(errRes.ErrWalletDoesNotExist)
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.NOT_FOUND,
					Message: errRes.ErrWalletDoesNotExist.Error(),
				},
			},
		},
		{
			description: "should return error with status code 400 when given wrong pin",
			input: input{
				userID: userID,
				request: &dto.ChangePinRequest{
					CurrentPin: "102938",
					NewPin:     NewPin,
				},
			},
			beforeTest: func(ws *mocks.WalletService) {
				ws.On("RequestPinChange", userID, &dto.ChangePinRequest{
					CurrentPin: "102938",
					NewPin:     NewPin,
				}).Return(errRes.ErrPinMismatch)
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.WRONG_PIN,
					Message: errRes.ErrPinMismatch.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to generate token",
			input: input{
				userID: userID,
				request: &dto.ChangePinRequest{
					CurrentPin: currentPin,
					NewPin:     NewPin,
				},
			},
			beforeTest: func(ws *mocks.WalletService) {
				ws.On("RequestPinChange", userID, &dto.ChangePinRequest{
					CurrentPin: currentPin,
					NewPin:     NewPin,
				}).Return(errors.New("failed to generate token"))
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errRes.ErrInternalServerError.Error(),
				},
			},
		},
		{
			description: "should return nil with status code 200 when succeed to generate token",
			input: input{
				userID: userID,
				request: &dto.ChangePinRequest{
					CurrentPin: currentPin,
					NewPin:     NewPin,
				},
			},
			beforeTest: func(ws *mocks.WalletService) {
				ws.On("RequestPinChange", userID, &dto.ChangePinRequest{
					CurrentPin: currentPin,
					NewPin:     NewPin,
				}).Return(nil)
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    nil,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			walletService := mocks.NewWalletService(t)
			tc.beforeTest(walletService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			handler := handler.New(&handler.HandlerConfig{
				WalletService: walletService,
			})
			payload := testutil.MakeRequestBody(tc.input.request)
			c.Request, _ = http.NewRequest(http.MethodPost, "/v1/users/wallets/pins/change-requests", payload)

			handler.RequestWalletPinChange(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestCompleteChangeWalletPin(t *testing.T) {
	type input struct {
		userID  int
		request *dto.CompleteChangePinRequest
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID           = 1
		verificationCode = "a1b2c3"
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.WalletService)
		expected
	}{
		{
			description: "should return error with status code 400 when given invalid request body",
			input: input{
				userID: userID,
				request: &dto.CompleteChangePinRequest{
					VerificationCode: "aa",
				},
			},
			beforeTest: func(ws *mocks.WalletService) {},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "VerificationCode must be 6 characters",
				},
			},
		},
		{
			description: "should return error with status code 404 when verification code not found",
			input: input{
				userID: userID,
				request: &dto.CompleteChangePinRequest{
					VerificationCode: verificationCode,
				},
			},
			beforeTest: func(ws *mocks.WalletService) {
				ws.On("CompletePinChange", userID, &dto.CompleteChangePinRequest{
					VerificationCode: verificationCode,
				}).Return(errRes.ErrVerificationCodeNotFound)
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.NOT_FOUND,
					Message: errRes.ErrVerificationCodeNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with status code 400 when verification code is incorrect",
			input: input{
				userID: userID,
				request: &dto.CompleteChangePinRequest{
					VerificationCode: verificationCode,
				},
			},
			beforeTest: func(ws *mocks.WalletService) {
				ws.On("CompletePinChange", userID, &dto.CompleteChangePinRequest{
					VerificationCode: verificationCode,
				}).Return(errRes.ErrIncorrectVerificationCode)
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.INCORRECT_VERIFICATION_CODE,
					Message: errRes.ErrIncorrectVerificationCode.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to change pin",
			input: input{
				userID: userID,
				request: &dto.CompleteChangePinRequest{
					VerificationCode: verificationCode,
				},
			},
			beforeTest: func(ws *mocks.WalletService) {
				ws.On("CompletePinChange", userID, &dto.CompleteChangePinRequest{
					VerificationCode: verificationCode,
				}).Return(errors.New("failed to complete change pin"))
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errRes.ErrInternalServerError.Error(),
				},
			},
		},
		{
			description: "should return nil with status code 200 when succeed to change pin",
			input: input{
				userID: userID,
				request: &dto.CompleteChangePinRequest{
					VerificationCode: verificationCode,
				},
			},
			beforeTest: func(ws *mocks.WalletService) {
				ws.On("CompletePinChange", userID, &dto.CompleteChangePinRequest{
					VerificationCode: verificationCode,
				}).Return(nil)
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			walletService := mocks.NewWalletService(t)
			tc.beforeTest(walletService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			handler := handler.New(&handler.HandlerConfig{
				WalletService: walletService,
			})
			payload := testutil.MakeRequestBody(tc.input.request)
			c.Request, _ = http.NewRequest(http.MethodPost, "/v1/users/wallets/pins/change-confirmations", payload)

			handler.CompleteChangeWalletPin(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestRequestWalletPinReset(t *testing.T) {
	type input struct {
		userID int
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID = 1
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.WalletService)
		expected
	}{
		{
			description: "should return error with status code 500 when failed to generate token",
			input: input{
				userID: userID,
			},
			beforeTest: func(ws *mocks.WalletService) {
				ws.On("RequestPinReset", userID).Return(errors.New("failed to generate token"))
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errRes.ErrInternalServerError.Error(),
				},
			},
		},
		{
			description: "should return nil with status code 200 when succeed to generate token",
			input: input{
				userID: userID,
			},
			beforeTest: func(ws *mocks.WalletService) {
				ws.On("RequestPinReset", userID).Return(nil)
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    nil,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			walletService := mocks.NewWalletService(t)
			tc.beforeTest(walletService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			handler := handler.New(&handler.HandlerConfig{
				WalletService: walletService,
			})
			c.Request, _ = http.NewRequest(http.MethodPost, "/v1/users/wallets/pins/reset-requests", nil)

			handler.RequestWalletPinReset(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestCompleteResetWalletPin(t *testing.T) {
	type input struct {
		userID  int
		request *dto.CompleteResetPinRequest
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID = 1
		token  = "a1b2c3"
		newPin = "123456"
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.WalletService)
		expected
	}{
		{
			description: "should return error with status code 400 when given invalid request body",
			input: input{
				userID: userID,
				request: &dto.CompleteResetPinRequest{
					Token:  "aa",
					NewPin: newPin,
				},
			},
			beforeTest: func(ws *mocks.WalletService) {},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Token must be 6 characters",
				},
			},
		},
		{
			description: "should return error with status code 404 when token not found",
			input: input{
				userID: userID,
				request: &dto.CompleteResetPinRequest{
					Token:  token,
					NewPin: newPin,
				},
			},
			beforeTest: func(ws *mocks.WalletService) {
				ws.On("CompletePinReset", userID, &dto.CompleteResetPinRequest{
					Token:  token,
					NewPin: newPin,
				}).Return(errRes.ErrResetPinTokenNotFound)
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.NOT_FOUND,
					Message: errRes.ErrResetPinTokenNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to change pin",
			input: input{
				userID: userID,
				request: &dto.CompleteResetPinRequest{
					Token:  token,
					NewPin: newPin,
				},
			},
			beforeTest: func(ws *mocks.WalletService) {
				ws.On("CompletePinReset", userID, &dto.CompleteResetPinRequest{
					Token:  token,
					NewPin: newPin,
				}).Return(errors.New("failed to complete change pin"))
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errRes.ErrInternalServerError.Error(),
				},
			},
		},
		{
			description: "should return nil with status code 200 when succeed to change pin",
			input: input{
				userID: userID,
				request: &dto.CompleteResetPinRequest{
					Token:  token,
					NewPin: newPin,
				},
			},
			beforeTest: func(ws *mocks.WalletService) {
				ws.On("CompletePinReset", userID, &dto.CompleteResetPinRequest{
					Token:  token,
					NewPin: newPin,
				}).Return(nil)
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			walletService := mocks.NewWalletService(t)
			tc.beforeTest(walletService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			handler := handler.New(&handler.HandlerConfig{
				WalletService: walletService,
			})
			payload := testutil.MakeRequestBody(tc.input.request)
			c.Request, _ = http.NewRequest(http.MethodPost, "/v1/users/wallets/pins/reset-confirmations", payload)

			handler.CompleteResetWalletPin(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

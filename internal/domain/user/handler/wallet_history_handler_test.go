package handler_test

import (
	"encoding/json"
	"kedai/backend/be-kedai/internal/common/code"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/handler"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetWalletHistory(t *testing.T) {
	var (
		userId        = 1
		request       = dto.WalletHistoryRequest{
			Page: 1,
			Limit: 10,
		}
		walletHistory = []*model.WalletHistory{}
		pagination    = &commonDto.PaginationResponse{
			Data: walletHistory,
		}
	)
	type input struct {
		result *commonDto.PaginationResponse
		req    dto.WalletHistoryRequest
		userId int
		err    error
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
			description: "should return list of wallet transaction histories with code 200 when success",
			input: input{
				result: pagination,
				userId: userId,
				req:    request,
				err:    nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    pagination,
				},
			},
		},
		{
			description: "should return error with code 404 when wallet does not exist",
			input: input{
				result: nil,
				req:    request,
				userId: userId,
				err:    errs.ErrWalletDoesNotExist,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.NOT_FOUND,
					Message: errs.ErrWalletDoesNotExist.Error(),
				},
			},
		},
		{
			description: "should return error with code 500 when internal server error",
			input: input{
				result: nil,
				req:    request,
				userId: userId,
				err:    errs.ErrInternalServerError,
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errs.ErrInternalServerError.Error(),
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userId)
			mockWalletHistoryService := new(mocks.WalletHistoryService)
			mockWalletHistoryService.On("GetWalletHistoryById", tc.input.req, tc.input.userId).Return(tc.input.result, tc.input.err)
			handler := handler.New(&handler.HandlerConfig{
				WalletHistoryService: mockWalletHistoryService,
			})
			c.Request, _ = http.NewRequest("POST", "/users/wallets/histories", nil)

			handler.GetWalletHistory(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

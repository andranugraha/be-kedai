package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"kedai/backend/be-kedai/internal/common/code"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/handler"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetInvoicesByUserID(t *testing.T) {
	type input struct {
		userID     int
		request    *dto.InvoicePerShopFilterRequest
		beforeTest func(*mocks.InvoicePerShopService)
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error with status code 400 when given invalid queries",
			input: input{
				userID: 1,
				request: &dto.InvoicePerShopFilterRequest{
					StartDate: "test",
				},
				beforeTest: func(ipss *mocks.InvoicePerShopService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "EndDate is required",
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to get invoices",
			input: input{
				userID:  1,
				request: &dto.InvoicePerShopFilterRequest{},
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("GetInvoicesByUserID", 1, &dto.InvoicePerShopFilterRequest{Limit: 10, Page: 1}).Return(nil, errors.New("failed to get invoices"))
				},
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errs.ErrInternalServerError.Error(),
				},
			},
		},
		{
			description: "should return error with status code 200 when suceed fetching invoices",
			input: input{
				userID:  1,
				request: &dto.InvoicePerShopFilterRequest{},
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("GetInvoicesByUserID", 1, &dto.InvoicePerShopFilterRequest{Limit: 10, Page: 1}).Return(&commonDto.PaginationResponse{Page: 1, Limit: 10, Data: []*dto.InvoicePerShopDetail{}}, nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    &commonDto.PaginationResponse{Page: 1, Limit: 10, Data: []*dto.InvoicePerShopDetail{}},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedJson, _ := json.Marshal(tc.expected.response)
			invoicePerShopService := mocks.NewInvoicePerShopService(t)
			tc.beforeTest(invoicePerShopService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/orders/invoices?startDate=%s&endDate=%s&s=%s&status=%s", tc.input.request.StartDate, tc.input.request.EndDate, tc.input.request.S, tc.input.request.Status), nil)
			handler := handler.New(&handler.Config{
				InvoicePerShopService: invoicePerShopService,
			})

			handler.GetInvoicePerShopsByUserID(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		})
	}
}

func TestGetInvoiceByCode(t *testing.T) {
	type input struct {
		userID   int
		code     string
		mockData *dto.InvoicePerShopDetail
		mockErr  error
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error with status code 404 when invoice does not exist",
			input: input{
				userID:   1,
				code:     "INV/XX/X",
				mockData: nil,
				mockErr:  errs.ErrInvoiceNotFound,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.INVOICE_NOT_FOUND,
					Message: errs.ErrInvoiceNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to get invoice",
			input: input{
				userID:   1,
				code:     "INV-XX-X",
				mockData: nil,
				mockErr:  errors.New("failed to get invoice"),
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errs.ErrInternalServerError.Error(),
				},
			},
		},
		{
			description: "should return invoice with status code 200 when fetching invoice succeed",
			input: input{
				userID:   1,
				code:     "INV/XX/X",
				mockData: &dto.InvoicePerShopDetail{},
				mockErr:  nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    &dto.InvoicePerShopDetail{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedJson, _ := json.Marshal(tc.expected.response)
			invoicePerShopService := mocks.NewInvoicePerShopService(t)
			invoicePerShopService.On("GetInvoicesByUserIDAndCode", tc.input.userID, tc.input.code).Return(tc.input.mockData, tc.input.mockErr)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			c.AddParam("code", tc.input.code)
			c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/orders/invoices/%s", tc.input.code), nil)
			handler := handler.New(&handler.Config{
				InvoicePerShopService: invoicePerShopService,
			})

			handler.GetInvoiceByCode(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		})
	}
}

func TestGetInvoicePerShopsByShopId(t *testing.T) {
	var (
		req = &dto.InvoicePerShopFilterRequest{
			Page:  1,
			Limit: 10,
		}
		invoices = &commonDto.PaginationResponse{}
		userId   = 1
	)
	type input struct {
		req        *dto.InvoicePerShopFilterRequest
		result     *commonDto.PaginationResponse
		err        error
		beforeTest func(*mocks.InvoicePerShopService)
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
			description: "should return list of invoices with code 200 when success",
			input: input{
				req:    req,
				result: invoices,
				err:    nil,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("GetInvoicesByShopId", userId, req).Return(invoices, nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    invoices,
				},
			},
		},
		{
			description: "should return error with code 400 when queries invalid",
			input: input{
				req: &dto.InvoicePerShopFilterRequest{
					StartDate: "test",
				},
				result:     nil,
				err:        nil,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "EndDate is required",
				},
			},
		},
		{
			description: "should return error with code 404 when user shop not found",
			input: input{
				req:    req,
				result: nil,
				err:    errs.ErrShopNotFound,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("GetInvoicesByShopId", userId, req).Return(nil, errs.ErrShopNotFound)
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.SHOP_NOT_REGISTERED,
					Message: errs.ErrShopNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with code 500 when internal server error",
			input: input{
				req:    req,
				result: nil,
				err:    errs.ErrInternalServerError,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("GetInvoicesByShopId", userId, req).Return(nil, errs.ErrInternalServerError)
				},
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
			expectedJson, _ := json.Marshal(tc.expected.response)
			invoicePerShopService := mocks.NewInvoicePerShopService(t)
			tc.beforeTest(invoicePerShopService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", userId)
			c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/sellers/finances/incomes?startDate=%s&endDate=%s&s=%s&status=%s", tc.input.req.StartDate, tc.input.req.EndDate, tc.input.req.S, tc.input.req.Status), nil)
			handler := handler.New(&handler.Config{
				InvoicePerShopService: invoicePerShopService,
			})

			handler.GetInvoicePerShopsByShopId(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		})
	}
}

func TestGetInvoiceByShopIdAndOrderId(t *testing.T) {
	var (
		userId  = 1
		shopId  = 1
		orderId = 1
	)
	type input struct {
		shopId     int
		orderId    int
		result     *dto.InvoicePerShopDetail
		err        error
		beforeTest func(*mocks.InvoicePerShopService)
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
			description: "should return error with code 404 when invoice not found",
			input: input{
				shopId:  shopId,
				orderId: orderId,
				result:  nil,
				err:     errs.ErrInvoiceNotFound,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("GetInvoiceByUserIdAndId", shopId, orderId).Return(nil, errs.ErrInvoiceNotFound)
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.INVOICE_NOT_FOUND,
					Message: errs.ErrInvoiceNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with code 404 when shop not found",
			input: input{
				shopId:  shopId,
				orderId: orderId,
				result:  nil,
				err:     errs.ErrShopNotFound,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("GetInvoiceByUserIdAndId", shopId, orderId).Return(nil, errs.ErrShopNotFound)
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.SHOP_NOT_REGISTERED,
					Message: errs.ErrShopNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with code 500 when internal server error",
			input: input{
				shopId:  shopId,
				orderId: orderId,
				result:  nil,
				err:     errs.ErrInternalServerError,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("GetInvoiceByUserIdAndId", shopId, orderId).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errs.ErrInternalServerError.Error(),
				},
			},
		},

		{
			description: "should return success",
			input: input{
				shopId:  shopId,
				orderId: orderId,
				result:  &dto.InvoicePerShopDetail{},
				err:     nil,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("GetInvoiceByUserIdAndId", shopId, orderId).Return(&dto.InvoicePerShopDetail{}, nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    &dto.InvoicePerShopDetail{},
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			expectedJson, _ := json.Marshal(tc.expected.response)
			invoicePerShopService := mocks.NewInvoicePerShopService(t)
			tc.beforeTest(invoicePerShopService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", userId)
			c.Params = []gin.Param{
				{
					Key:   "orderId",
					Value: strconv.Itoa(orderId),
				},
			}

			c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/sellers/finances/incomes/%d/%d", tc.input.shopId, tc.input.orderId), nil)
			handler := handler.New(&handler.Config{
				InvoicePerShopService: invoicePerShopService,
			})

			handler.GetInvoiceByShopIdAndOrderId(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		})
	}
}

func TestWithdrawFromInvoice(t *testing.T) {
	var (
		userId  = 1
		orderId = 1
	)
	type input struct {
		req        dto.WithdrawInvoiceRequest
		err        error
		beforeTest func(*mocks.InvoicePerShopService)
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
			description: "should return error with code 400 when request body invalid",
			input: input{
				req: dto.WithdrawInvoiceRequest{},
				err: errs.ErrInvoiceNotFound,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
				},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "OrderID is required",
				},
			},
		},
		{
			description: "should return error with code 404 when invoice not found",
			input: input{
				req: dto.WithdrawInvoiceRequest{
					OrderID: orderId,
				},
				err: errs.ErrInvoiceNotFound,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("WithdrawFromInvoice", userId, orderId).Return(errs.ErrInvoiceNotFound)
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.INVOICE_NOT_FOUND,
					Message: errs.ErrInvoiceNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with code 404 when shop not found",
			input: input{
				req: dto.WithdrawInvoiceRequest{
					OrderID: orderId,
				},
				err: errs.ErrShopNotFound,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("WithdrawFromInvoice", userId, orderId).Return(errs.ErrShopNotFound)
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.SHOP_NOT_REGISTERED,
					Message: errs.ErrShopNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with code 404 when wallet not found",
			input: input{
				req: dto.WithdrawInvoiceRequest{
					OrderID: orderId,
				},
				err: errs.ErrWalletDoesNotExist,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("WithdrawFromInvoice", userId, orderId).Return(errs.ErrWalletDoesNotExist)
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.WALLET_DOES_NOT_EXIST,
					Message: errs.ErrWalletDoesNotExist.Error(),
				},
			},
		},
		{
			description: "should return error with code 500 when withdraw from invoice failed",
			input: input{
				req: dto.WithdrawInvoiceRequest{
					OrderID: orderId,
				},
				err: errs.ErrInternalServerError,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("WithdrawFromInvoice", userId, orderId).Return(errs.ErrInternalServerError)
				},
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errs.ErrInternalServerError.Error(),
				},
			},
		},
		{
			description: "should return success with code 200 when withdraw from invoice success",
			input: input{
				req: dto.WithdrawInvoiceRequest{
					OrderID: orderId,
				},
				err: nil,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("WithdrawFromInvoice", userId, orderId).Return(nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
				},
			},
		},
	} {

		t.Run(tc.description, func(t *testing.T) {
			expectedJson, _ := json.Marshal(tc.expected.response)
			invoicePerShopService := mocks.NewInvoicePerShopService(t)
			tc.beforeTest(invoicePerShopService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", userId)

			payload := test.MakeRequestBody(tc.input.req)
			c.Request, _ = http.NewRequest(http.MethodPost, "/sellers/finances/incomes/withdrawals", payload)
			handler := handler.New(&handler.Config{
				InvoicePerShopService: invoicePerShopService,
			})

			handler.WithdrawFromInvoice(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		},
		)
	}

}

func TestGetShopOrder(t *testing.T) {
	var (
		req = &dto.InvoicePerShopFilterRequest{
			Page:  1,
			Limit: 10,
		}
		invoices = &commonDto.PaginationResponse{}
		userId   = 1
	)
	type input struct {
		req        *dto.InvoicePerShopFilterRequest
		result     *commonDto.PaginationResponse
		err        error
		beforeTest func(*mocks.InvoicePerShopService)
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
			description: "should return list of invoices with code 200 when success",
			input: input{
				req:    req,
				result: invoices,
				err:    nil,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("GetShopOrder", userId, req).Return(invoices, nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    invoices,
				},
			},
		},
		{
			description: "should return error with code 400 when queries invalid",
			input: input{
				req: &dto.InvoicePerShopFilterRequest{
					StartDate: "test",
				},
				result:     nil,
				err:        nil,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "EndDate is required",
				},
			},
		},
		{
			description: "should return error with code 404 when user shop not found",
			input: input{
				req:    req,
				result: nil,
				err:    errs.ErrShopNotFound,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("GetShopOrder", userId, req).Return(nil, errs.ErrShopNotFound)
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.SHOP_NOT_REGISTERED,
					Message: errs.ErrShopNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with code 500 when internal server error",
			input: input{
				req:    req,
				result: nil,
				err:    errs.ErrInternalServerError,
				beforeTest: func(ipss *mocks.InvoicePerShopService) {
					ipss.On("GetShopOrder", userId, req).Return(nil, errs.ErrInternalServerError)
				},
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
			expectedJson, _ := json.Marshal(tc.expected.response)
			invoicePerShopService := mocks.NewInvoicePerShopService(t)
			tc.beforeTest(invoicePerShopService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", userId)
			c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/sellers/orders?startDate=%s&endDate=%s&user=%s", tc.input.req.StartDate, tc.input.req.EndDate, tc.input.req.Username), nil)
			handler := handler.New(&handler.Config{
				InvoicePerShopService: invoicePerShopService,
			})

			handler.GetShopOrder(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		})
	}
}

func TestUpdateToDeliver(t *testing.T) {
	var (
		userId  = 1
		orderId = 1
	)
	type input struct {
		userId  int
		orderId int
		err     error
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
			description: "should return nil error with code 200 when success",
			input: input{
				userId:  userId,
				orderId: orderId,
				err:     nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
				},
			},
		},
		{
			description: "should return error with code 404 when shop not found",
			input: input{
				userId:  userId,
				orderId: orderId,
				err:     errs.ErrShopNotFound,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.SHOP_NOT_REGISTERED,
					Message: errs.ErrShopNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with code 404 when invoice not found",
			input: input{
				userId:  1,
				orderId: 1,
				err:     errs.ErrInvoiceNotFound,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.INVOICE_NOT_FOUND,
					Message: errs.ErrInvoiceNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with code 500 when internal server error",
			input: input{
				userId:  1,
				orderId: 1,
				err:     errs.ErrInternalServerError,
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
			expectedJson, _ := json.Marshal(tc.expected.response)
			invoicePerShopService := mocks.NewInvoicePerShopService(t)
			invoicePerShopService.On("UpdateStatusToDelivery", tc.input.userId, tc.input.orderId).Return(tc.input.err)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", userId)
			c.AddParam("orderId", "1")
			c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/sellers/orders/{%d}/delivery", tc.orderId), nil)
			handler := handler.New(&handler.Config{
				InvoicePerShopService: invoicePerShopService,
			})

			handler.UpdateToDelivery(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		})
	}
}

func TestUpdateToCancelled(t *testing.T) {
	var (
		userId  = 1
		orderId = 1
	)
	type input struct {
		userId  int
		orderId int
		err     error
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
			description: "should return nil error with code 200 when success",
			input: input{
				userId:  userId,
				orderId: orderId,
				err:     nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
				},
			},
		},
		{
			description: "should return error with code 404 when shop not found",
			input: input{
				userId:  userId,
				orderId: orderId,
				err:     errs.ErrShopNotFound,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.SHOP_NOT_REGISTERED,
					Message: errs.ErrShopNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with code 404 when invoice not found",
			input: input{
				userId:  1,
				orderId: 1,
				err:     errs.ErrInvoiceNotFound,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.INVOICE_NOT_FOUND,
					Message: errs.ErrInvoiceNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with code 500 when internal server error",
			input: input{
				userId:  1,
				orderId: 1,
				err:     errs.ErrInternalServerError,
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
			expectedJson, _ := json.Marshal(tc.expected.response)
			invoicePerShopService := mocks.NewInvoicePerShopService(t)
			invoicePerShopService.On("UpdateStatusToCancelled", tc.input.userId, tc.input.orderId).Return(tc.input.err)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", userId)
			c.AddParam("orderId", "1")
			c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/sellers/orders/{%d}/cancel", tc.orderId), nil)
			handler := handler.New(&handler.Config{
				InvoicePerShopService: invoicePerShopService,
			})

			handler.UpdateToCancelled(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		})
	}
}

func TestUpdateToReceived(t *testing.T) {
	var (
		userId    = 1
		orderCode = "code"
	)
	type input struct {
		userId    int
		orderCode string
		err       error
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
			description: "should return nil error with code 200 when success",
			input: input{
				userId:    userId,
				orderCode: orderCode,
				err:       nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
				},
			},
		},
		{
			description: "should return error with code 404 when invoice not found",
			input: input{
				userId:    1,
				orderCode: orderCode,
				err:       errs.ErrInvoiceNotFound,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.INVOICE_NOT_FOUND,
					Message: errs.ErrInvoiceNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with code 500 when internal server error",
			input: input{
				userId:    1,
				orderCode: orderCode,
				err:       errs.ErrInternalServerError,
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
			expectedJson, _ := json.Marshal(tc.expected.response)
			invoicePerShopService := mocks.NewInvoicePerShopService(t)
			invoicePerShopService.On("UpdateStatusToReceived", tc.input.userId, tc.input.orderCode).Return(tc.input.err)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", userId)
			c.AddParam("code", orderCode)
			c.Request, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/orders/invoices/{%s}/receive", tc.orderCode), nil)
			handler := handler.New(&handler.Config{
				InvoicePerShopService: invoicePerShopService,
			})

			handler.UpdateToReceived(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		})
	}
}

func TestUpdateCronJob(t *testing.T) {
	t.Run("should return nothing when called whether its error or success", func(t *testing.T) {
		mockService := new(mocks.InvoicePerShopService)
		mockService.On("UpdateStatusCRONJob").Return(nil)
		mockService.On("AutoReceivedCRONJob").Return(nil)
		handler := handler.New(&handler.Config{
			InvoicePerShopService: mockService,
		})
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		handler.UpdateCronJob(c)

		mockService.AssertNumberOfCalls(t, "UpdateStatusCRONJob", 1)
		mockService.AssertNumberOfCalls(t, "AutoReceivedCRONJob", 1)
	})
}

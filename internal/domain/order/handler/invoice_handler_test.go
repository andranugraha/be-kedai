package handler_test

import (
	"encoding/json"
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/handler"
	userDto "kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/utils/response"
	testutil "kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCheckout(t *testing.T) {
	var (
		userId   = 1
		one      = 1
		products = []dto.CheckoutProduct{
			{
				CartItemID: 1,
				Quantity:   1,
			},
		}
		items = []dto.CheckoutItem{
			{
				ShopID:           1,
				VoucherID:        &one,
				CourierServiceID: 1,
				ShippingCost:     1000,
				Products:         products,
			},
		}
		req = dto.CheckoutRequest{
			AddressID:       1,
			TotalPrice:      10000,
			VoucherID:       &one,
			UserID:          1,
			PaymentMethodID: 2,
			SealabsPayID:    &one,
			Items:           items,
		}
	)

	tests := []struct {
		name            string
		req             dto.CheckoutRequest
		wantCheckoutRes *dto.CheckoutResponse
		wantCheckoutErr error
		want            response.Response
		code            int
		beforeTest      func(*mocks.InvoiceService)
	}{
		{
			name: "should return 200 when checkout success",
			req:  req,
			wantCheckoutRes: &dto.CheckoutResponse{
				ID: 1,
			},
			wantCheckoutErr: nil,
			want: response.Response{
				Code:    code.CREATED,
				Message: "checkout success",
				Data: dto.CheckoutResponse{
					ID: 1,
				},
			},
			code: http.StatusCreated,
			beforeTest: func(is *mocks.InvoiceService) {
				is.On("Checkout", req).Return(&dto.CheckoutResponse{
					ID: 1,
				}, nil)
			},
		},
		{
			name:            "should return 400 when request is invalid",
			req:             dto.CheckoutRequest{},
			wantCheckoutRes: nil,
			wantCheckoutErr: errors.New("PaymentMethodID is required"),
			want: response.Response{
				Code:    code.BAD_REQUEST,
				Message: "PaymentMethodID is required",
			},
			code:       http.StatusBadRequest,
			beforeTest: func(is *mocks.InvoiceService) {},
		},
		{
			name: "should return 400 when request failed validation check",
			req: dto.CheckoutRequest{
				AddressID:       req.AddressID,
				TotalPrice:      req.TotalPrice,
				PaymentMethodID: req.PaymentMethodID,
				Items:           req.Items,
			},
			wantCheckoutRes: nil,
			wantCheckoutErr: errs.ErrSealabsPayIdIsRequired,
			want: response.Response{
				Code:    code.BAD_REQUEST,
				Message: errs.ErrSealabsPayIdIsRequired.Error(),
			},
			code:       http.StatusBadRequest,
			beforeTest: func(is *mocks.InvoiceService) {},
		},
		{
			name:            "should return 400 when address not found",
			req:             req,
			wantCheckoutRes: nil,
			wantCheckoutErr: errs.ErrAddressNotFound,
			want: response.Response{
				Code:    code.BAD_REQUEST,
				Message: errs.ErrAddressNotFound.Error(),
			},
			code: http.StatusBadRequest,
			beforeTest: func(is *mocks.InvoiceService) {
				is.On("Checkout", req).Return(nil, errs.ErrAddressNotFound)
			},
		},
		{
			name:            "should return 400 when stock is insufficient",
			req:             req,
			wantCheckoutRes: nil,
			wantCheckoutErr: errs.ErrProductQuantityNotEnough,
			want: response.Response{
				Code:    code.QUANTITY_NOT_ENOUGH,
				Message: errs.ErrProductQuantityNotEnough.Error(),
			},
			code: http.StatusBadRequest,
			beforeTest: func(is *mocks.InvoiceService) {
				is.On("Checkout", req).Return(nil, errs.ErrProductQuantityNotEnough)
			},
		},
		{
			name:            "should return 400 when quantity is not match",
			req:             req,
			wantCheckoutRes: nil,
			wantCheckoutErr: errs.ErrQuantityNotMatch,
			want: response.Response{
				Code:    code.CART_ITEM_MISMATCH,
				Message: errs.ErrQuantityNotMatch.Error(),
			},
			code: http.StatusBadRequest,
			beforeTest: func(is *mocks.InvoiceService) {
				is.On("Checkout", req).Return(nil, errs.ErrQuantityNotMatch)
			},
		},
		{
			name:            "should return 400 when total spend is not enough to use voucher",
			req:             req,
			wantCheckoutRes: nil,
			wantCheckoutErr: errs.ErrTotalSpentBelowMinimumSpendingRequirement,
			want: response.Response{
				Code:    code.MINIMUM_SPEND_REQUIREMENT_NOT_MET,
				Message: errs.ErrTotalSpentBelowMinimumSpendingRequirement.Error(),
			},
			code: http.StatusBadRequest,
			beforeTest: func(is *mocks.InvoiceService) {
				is.On("Checkout", req).Return(nil, errs.ErrTotalSpentBelowMinimumSpendingRequirement)
			},
		},
		{
			name:            "should return 500 when internal server error",
			req:             req,
			wantCheckoutRes: nil,
			wantCheckoutErr: errs.ErrInternalServerError,
			want: response.Response{
				Code:    code.INTERNAL_SERVER_ERROR,
				Message: errs.ErrInternalServerError.Error(),
			},
			code: http.StatusInternalServerError,
			beforeTest: func(is *mocks.InvoiceService) {
				is.On("Checkout", req).Return(nil, errs.ErrInternalServerError)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedJson, _ := json.Marshal(test.want)
			payload := testutil.MakeRequestBody(test.req)
			service := mocks.NewInvoiceService(t)
			test.beforeTest(service)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", userId)

			c.Request, _ = http.NewRequest(http.MethodPost, "/v1/orders", payload)
			handler := handler.New(&handler.Config{
				InvoiceService: service,
			})
			handler.Checkout(c)

			assert.Equal(t, test.code, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		})
	}
}

func TestPayInvoice(t *testing.T) {
	var (
		token  = "token"
		userId = 1
		req    = dto.PayInvoiceRequest{
			InvoiceID:       1,
			UserID:          1,
			PaymentMethodID: 1,
		}
		res = userDto.Token{
			AccessToken:  token,
			RefreshToken: token,
		}
	)

	tests := []struct {
		name       string
		req        dto.PayInvoiceRequest
		want       response.Response
		code       int
		beforeTest func(*mocks.InvoiceService)
	}{
		{
			name: "should return 200 when pay success",
			req:  req,
			want: response.Response{
				Code:    code.OK,
				Message: "pay invoice success",
				Data:    res,
			},
			code: http.StatusOK,
			beforeTest: func(is *mocks.InvoiceService) {
				is.On("PayInvoice", req, token).Return(&res, nil)
			},
		},
		{
			name: "should return 400 when request is invalid",
			req:  dto.PayInvoiceRequest{},
			want: response.Response{
				Code:    code.BAD_REQUEST,
				Message: "TxnID is required",
			},
			code:       http.StatusBadRequest,
			beforeTest: func(is *mocks.InvoiceService) {},
		},
		{
			name: "should return 400 when invoice not found",
			req:  req,
			want: response.Response{
				Code:    code.BAD_REQUEST,
				Message: errs.ErrInvoiceNotFound.Error(),
			},
			code: http.StatusBadRequest,
			beforeTest: func(is *mocks.InvoiceService) {
				is.On("PayInvoice", req, token).Return(nil, errs.ErrInvoiceNotFound)
			},
		},
		{
			name: "should return 400 when balance is insufficient",
			req:  req,
			want: response.Response{
				Code:    code.INSUFFICIENT_BALANCE,
				Message: errs.ErrInsufficientBalance.Error(),
			},
			code: http.StatusBadRequest,
			beforeTest: func(is *mocks.InvoiceService) {
				is.On("PayInvoice", req, token).Return(nil, errs.ErrInsufficientBalance)
			},
		},
		{
			name: "should return 500 when internal server error",
			req:  req,
			want: response.Response{
				Code:    code.INTERNAL_SERVER_ERROR,
				Message: errs.ErrInternalServerError.Error(),
			},
			code: http.StatusInternalServerError,
			beforeTest: func(is *mocks.InvoiceService) {
				is.On("PayInvoice", req, token).Return(nil, errs.ErrInternalServerError)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedJson, _ := json.Marshal(test.want)
			payload := testutil.MakeRequestBody(test.req)
			service := mocks.NewInvoiceService(t)
			test.beforeTest(service)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", userId)
			c.Set("level", 1)

			c.Request, _ = http.NewRequest(http.MethodPost, "/v1/orders/invoices", payload)
			c.Request.Header.Add("authorization", "Bearer "+token)
			handler := handler.New(&handler.Config{
				InvoiceService: service,
			})
			handler.PayInvoice(c)

			assert.Equal(t, test.code, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		})
	}
}

func TestCancelCheckout(t *testing.T) {
	var (
		userId = 1
		req    = dto.CancelCheckoutRequest{
			InvoiceID: 1,
			UserID:    1,
		}
	)

	tests := []struct {
		name       string
		req        dto.CancelCheckoutRequest
		want       response.Response
		code       int
		beforeTest func(*mocks.InvoiceService)
	}{
		{
			name: "should return 200 when cancel success",
			req:  req,
			want: response.Response{
				Code:    code.OK,
				Message: "cancel checkout success",
			},
			code: http.StatusOK,
			beforeTest: func(is *mocks.InvoiceService) {
				is.On("CancelCheckout", req).Return(nil)
			},
		},
		{
			name: "should return 400 when request is invalid",
			req:  dto.CancelCheckoutRequest{},
			want: response.Response{
				Code:    code.BAD_REQUEST,
				Message: "InvoiceID is required",
			},
			code: http.StatusBadRequest,
			beforeTest: func(is *mocks.InvoiceService) {
			},
		},
		{
			name: "should return 400 when invoice not found",
			req:  req,
			want: response.Response{
				Code:    code.BAD_REQUEST,
				Message: errs.ErrInvoiceNotFound.Error(),
			},
			code: http.StatusBadRequest,
			beforeTest: func(is *mocks.InvoiceService) {
				is.On("CancelCheckout", req).Return(errs.ErrInvoiceNotFound)
			},
		},
		{
			name: "should return 500 when internal server error",
			req:  req,
			want: response.Response{
				Code:    code.INTERNAL_SERVER_ERROR,
				Message: errs.ErrInternalServerError.Error(),
			},
			code: http.StatusInternalServerError,
			beforeTest: func(is *mocks.InvoiceService) {
				is.On("CancelCheckout", req).Return(errs.ErrInternalServerError)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedJson, _ := json.Marshal(test.want)
			payload := testutil.MakeRequestBody(test.req)
			service := mocks.NewInvoiceService(t)
			test.beforeTest(service)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", userId)

			c.Request, _ = http.NewRequest(http.MethodPost, "/v1/orders/invoices/cancel", payload)
			handler := handler.New(&handler.Config{
				InvoiceService: service,
			})
			handler.CancelCheckout(c)

			assert.Equal(t, test.code, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())
		})
	}
}

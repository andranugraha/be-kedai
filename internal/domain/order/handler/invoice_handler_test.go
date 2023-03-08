package handler_test

import (
	"encoding/json"
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/domain/order/dto"
	"kedai/backend/be-kedai/internal/domain/order/handler"
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
		wantCheckoutRes dto.CheckoutResponse
		wantCheckoutErr error
		want            response.Response
		code            int
		beforeTest      func(*mocks.InvoiceService)
	}{
		{
			name: "should return 200 when checkout success",
			req:  req,
			wantCheckoutRes: dto.CheckoutResponse{
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

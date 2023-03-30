package handler_test

import (
	"encoding/json"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
	"kedai/backend/be-kedai/internal/domain/marketplace/handler"
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"kedai/backend/be-kedai/internal/server"
	"kedai/backend/be-kedai/internal/utils/response"
	testutil "kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetMarketplaceVoucher(t *testing.T) {
	type input struct {
		err        error
		req        dto.GetMarketplaceVoucherRequest
		beforeTest func(*mocks.MarketplaceVoucherService)
	}
	type expected struct {
		response   response.Response
		statusCode int
	}
	cases := []struct {
		description string
		input       input
		expected    expected
	}{
		{
			description: "should return response error internal server error and status code 500",
			input: input{
				err: errs.ErrInternalServerError,
				req: dto.GetMarketplaceVoucherRequest{},
				beforeTest: func(m *mocks.MarketplaceVoucherService) {
					m.On("GetMarketplaceVoucher", &dto.GetMarketplaceVoucherRequest{}).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errs.ErrInternalServerError.Error(),
				},
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "should return response marketplace vouchers, code ok and status code 200",
			input: input{
				err: nil,
				req: dto.GetMarketplaceVoucherRequest{},
				beforeTest: func(m *mocks.MarketplaceVoucherService) {
					m.On("GetMarketplaceVoucher", &dto.GetMarketplaceVoucherRequest{}).Return([]*model.MarketplaceVoucher{}, nil)
				},
			},
			expected: expected{
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    []*model.MarketplaceVoucher{},
				},
				statusCode: http.StatusOK,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			jsonRes, _ := json.Marshal(c.expected.response)
			m := mocks.NewMarketplaceVoucherService(t)
			c.input.beforeTest(m)

			cfg := &server.RouterConfig{
				MarketplaceHandler: handler.New(&handler.HandlerConfig{
					MarketplaceVoucherService: m,
				}),
			}

			req, _ := http.NewRequest("GET", "/v1/marketplaces/vouchers", nil)
			_, rec := testutil.ServeReq(cfg, req)

			assert.Equal(t, c.expected.statusCode, rec.Code)
			assert.Equal(t, string(jsonRes), rec.Body.String())

		})
	}

}

func TestGetMarketplaceVoucherAdmin(t *testing.T) {
	type input struct {
		err        error
		req        dto.GetMarketplaceVoucherRequest
		beforeTest func(*mocks.MarketplaceVoucherService)
	}
	type expected struct {
		response   response.Response
		statusCode int
	}
	cases := []struct {
		description string
		input       input
		expected    expected
	}{
		{
			description: "should return response error internal server error and status code 500",
			input: input{
				err: errs.ErrInternalServerError,
				req: dto.GetMarketplaceVoucherRequest{},
				beforeTest: func(m *mocks.MarketplaceVoucherService) {
					m.On("GetMarketplaceVoucherAdmin", &dto.GetMarketplaceVoucherRequest{}).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errs.ErrInternalServerError.Error(),
				},
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "should return response marketplace vouchers, code ok and status code 200",
			input: input{
				err: nil,
				req: dto.GetMarketplaceVoucherRequest{},
				beforeTest: func(m *mocks.MarketplaceVoucherService) {
					m.On("GetMarketplaceVoucherAdmin", &dto.GetMarketplaceVoucherRequest{}).Return([]*model.MarketplaceVoucher{}, nil)
				},
			},
			expected: expected{
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    []*model.MarketplaceVoucher{},
				},
				statusCode: http.StatusOK,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			jsonRes, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			mockMarketplaceVoucherService := mocks.NewMarketplaceVoucherService(t)
			tc.input.beforeTest(mockMarketplaceVoucherService)

			handler := handler.New(&handler.HandlerConfig{
				MarketplaceVoucherService: mockMarketplaceVoucherService,
			})

			c.Request, _ = http.NewRequest("GET", "/v1/admins/marketplaces/vouchers", nil)

			handler.GetMarketplaceVoucherAdmin(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(jsonRes), rec.Body.String())

		})
	}
}

func TestGetValidMarketplaceVoucher(t *testing.T) {
	type input struct {
		err        error
		req        dto.GetMarketplaceVoucherRequest
		beforeTest func(*mocks.MarketplaceVoucherService)
	}
	type expected struct {
		response   response.Response
		statusCode int
	}
	cases := []struct {
		description string
		input       input
		expected    expected
	}{
		{
			description: "should return response error internal server error and status code 500",
			input: input{
				err: errs.ErrInternalServerError,
				req: dto.GetMarketplaceVoucherRequest{},
				beforeTest: func(m *mocks.MarketplaceVoucherService) {
					m.On("GetValidByUserID", &dto.GetMarketplaceVoucherRequest{UserId: 1}).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errs.ErrInternalServerError.Error(),
				},
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "should return ok response and status code 200",
			input: input{
				err: nil,
				req: dto.GetMarketplaceVoucherRequest{
					UserId: 1,
				},
				beforeTest: func(m *mocks.MarketplaceVoucherService) {
					m.On("GetValidByUserID", &dto.GetMarketplaceVoucherRequest{
						UserId: 1,
					}).Return([]*model.MarketplaceVoucher{}, nil)
				},
			},
			expected: expected{
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    []*model.MarketplaceVoucher{},
				},
				statusCode: http.StatusOK,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			jsonRes, _ := json.Marshal(tc.expected.response)
			m := mocks.NewMarketplaceVoucherService(t)
			tc.input.beforeTest(m)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", 1)

			c.Request, _ = http.NewRequest(http.MethodGet, "v1//marketplace/vouchers/valid", nil)

			h := handler.New(&handler.HandlerConfig{
				MarketplaceVoucherService: m,
			})

			h.GetValidMarketplaceVoucher(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(jsonRes), rec.Body.String())

		})
	}
}

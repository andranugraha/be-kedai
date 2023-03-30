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
	"time"

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

func TestCreateMarketplaceVoucher(t *testing.T) {
	var (
		boolValue = true
		value     = 1
		time, _   = time.Parse("2006-03-02", "2022-11-22")
		req       = dto.CreateMarketplaceVoucherRequest{
			Code:            "A",
			Name:            "A",
			Amount:          1,
			Type:            "A",
			IsHidden:        &boolValue,
			Description:     "A",
			MinimumSpend:    1,
			ExpiredAt:       time,
			CategoryID:      &value,
			PaymentMethodID: &value,
		}
		res = &model.MarketplaceVoucher{
			Code:            "A",
			Name:            "A",
			Amount:          1,
			Type:            "A",
			IsHidden:        boolValue,
			Description:     "A",
			MinimumSpend:    1,
			ExpiredAt:       time,
			CategoryID:      &value,
			PaymentMethodID: &value,
		}
	)
	type input struct {
		req        dto.CreateMarketplaceVoucherRequest
		result     *model.MarketplaceVoucher
		beforeTest func(*mocks.MarketplaceVoucherService)
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
			description: "should return created voucher with code 201 when success",
			input: input{
				req:    req,
				result: res,
				beforeTest: func(mvs *mocks.MarketplaceVoucherService) {
					mvs.On("CreateMarketplaceVoucher", &req).Return(res, nil)
				},
			},
			expected: expected{
				statusCode: http.StatusCreated,
				response: response.Response{
					Code:    code.CREATED,
					Message: "created",
					Data:    res,
				},
			},
		},
		{
			description: "should return error with code 400 when invalid input requested",
			input: input{
				req:        dto.CreateMarketplaceVoucherRequest{},
				result:     nil,
				beforeTest: func(mvs *mocks.MarketplaceVoucherService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "ExpiredAt is required",
				},
			},
		},
		{
			description: "should return error with code 409 when voucher code duplicate",
			input: input{
				req:    req,
				result: nil,
				beforeTest: func(mvs *mocks.MarketplaceVoucherService) {
					mvs.On("CreateMarketplaceVoucher", &req).Return(nil, errs.ErrDuplicateVoucherCode)
				},
			},
			expected: expected{
				statusCode: http.StatusConflict,
				response: response.Response{
					Code:    code.DUPLICATE_VOUCHER_CODE,
					Message: errs.ErrDuplicateVoucherCode.Error(),
				},
			},
		},
		{
			description: "should return error with code 500 when internal server error",
			input: input{
				req:    req,
				result: nil,
				beforeTest: func(mvs *mocks.MarketplaceVoucherService) {
					mvs.On("CreateMarketplaceVoucher", &req).Return(nil, errs.ErrInternalServerError)
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
			jsonRes, _ := json.Marshal(tc.expected.response)
			m := mocks.NewMarketplaceVoucherService(t)
			tc.input.beforeTest(m)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request, _ = http.NewRequest(http.MethodPost, "v1/admins/marketplaces/vouchers", testutil.MakeRequestBody(tc.input.req))
			h := handler.New(&handler.HandlerConfig{
				MarketplaceVoucherService: m,
			})

			h.CreateMarketplaceVoucher(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(jsonRes), rec.Body.String())
		})
	}
}

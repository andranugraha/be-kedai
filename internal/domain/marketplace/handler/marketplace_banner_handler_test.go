package handler_test

import (
	"encoding/json"
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
	"kedai/backend/be-kedai/internal/domain/marketplace/handler"
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"kedai/backend/be-kedai/internal/utils/response"
	testutil "kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetMarketplaceBanner(t *testing.T) {
	banners := []*model.MarketplaceBanner{}
	type input struct {
		beforeTest func(*mocks.MarketplaceBannerService)
	}
	type expected struct {
		response   response.Response
		statusCode int
	}
	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error with status code 500 when fails to fetch banners",
			input: input{
				beforeTest: func(m *mocks.MarketplaceBannerService) {
					m.On("GetMarketplaceBanner").Return(nil, errors.New("internal server error"))
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
			description: "should return banners with status code 200 on success",
			input: input{
				beforeTest: func(m *mocks.MarketplaceBannerService) {
					m.On("GetMarketplaceBanner").Return(banners, nil)
				},
			},
			expected: expected{
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    banners,
				},
				statusCode: http.StatusOK,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			jsonRes, _ := json.Marshal(tc.expected.response)
			m := mocks.NewMarketplaceBannerService(t)
			tc.input.beforeTest(m)

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)

			c.Request, _ = http.NewRequest(http.MethodGet, "/marketplaces/banners", nil)

			h := handler.New(&handler.HandlerConfig{
				MarketplaceBannerService: m,
			})

			h.GetMarketplaceBanner(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(jsonRes), rec.Body.String())
		})
	}
}

func TestAddMarketplaceBanner(t *testing.T) {
	var (
		startDate = time.Now()
		body      = &dto.MarketplaceBannerRequest{
			MediaUrl:  "https://a.com",
			StartDate: startDate.Format(time.RFC3339Nano),
			EndDate:   startDate.AddDate(0, 0, 14).Format(time.RFC3339Nano),
		}
		banner = &model.MarketplaceBanner{
			ID:        1,
			MediaUrl:  "https://a.com",
			StartDate: startDate,
			EndDate:   startDate.AddDate(0, 0, 14),
		}
	)

	type input struct {
		body   *dto.MarketplaceBannerRequest
		result *model.MarketplaceBanner
		err    error
	}

	type expected struct {
		statusCode int
		response   response.Response
	}

	type cases struct {
		description string
		input
		beforeTests func(cs *mocks.MarketplaceBannerService)
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return banner response with code 201 when success",
			input: input{
				body:   body,
				result: banner,
				err:    nil,
			},
			beforeTests: func(cs *mocks.MarketplaceBannerService) {
				cs.On("AddMarketplaceBanner", body).Return(banner, nil)
			},
			expected: expected{
				statusCode: http.StatusCreated,
				response: response.Response{
					Code:    code.CREATED,
					Message: "success",
					Data:    banner,
				},
			},
		},
		{
			description: "should return error with code 400 when bad request",
			input: input{
				body:   body,
				result: nil,
				err:    errs.ErrBackDate,
			},
			beforeTests: func(cs *mocks.MarketplaceBannerService) {
				cs.On("AddMarketplaceBanner", body).Return(nil, errs.ErrBackDate)
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: errs.ErrBackDate.Error(),
				},
			},
		},
		{
			description: "should return error with code 400 when bad params",
			input: input{
				body:   &dto.MarketplaceBannerRequest{},
				result: nil,
				err:    nil,
			},
			beforeTests: func(cs *mocks.MarketplaceBannerService) {},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: errors.New("EndDate is required").Error(),
				},
			},
		},
		{
			description: "should return error with code 500 when error",
			input: input{
				body:   body,
				result: nil,
				err:    errors.New("error"),
			},
			beforeTests: func(cs *mocks.MarketplaceBannerService) {
				cs.On("AddMarketplaceBanner", body).Return(nil, errors.New("error"))
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
			expectedBody, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			mockMarketplaceBannerService := mocks.NewMarketplaceBannerService(t)
			tc.beforeTests(mockMarketplaceBannerService)
			handler := handler.New(&handler.HandlerConfig{
				MarketplaceBannerService: mockMarketplaceBannerService,
			})

			c.Request, _ = http.NewRequest("POST", "/admins/marketplaces/banners", testutil.MakeRequestBody(tc.input.body))

			handler.AddMarketplaceBanner(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

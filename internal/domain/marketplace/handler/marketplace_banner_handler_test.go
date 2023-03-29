package handler_test

import (
	"encoding/json"
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/marketplace/handler"
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

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

package handler_test

import (
	"encoding/json"
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/handler"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetShipmentList(t *testing.T) {
	var shopId = 1
	type input struct {
		result []*dto.ShipmentCourierResponse
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
			description: "should return list of couriers with code 200 when success",
			input: input{
				result: []*dto.ShipmentCourierResponse{},
				err:    nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    []*model.Courier{},
				},
			},
		},
		{
			description: "should return error with code 404 when shop not found",
			input: input{
				result: nil,
				err:    errs.ErrShopNotFound,
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
				result: nil,
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
			expectedBody, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", 1)
			mockService := new(mocks.CourierService)
			mockService.On("GetShipmentList", shopId).Return(tc.input.result, tc.input.err)
			handler := handler.New(&handler.HandlerConfig{
				CourierService: mockService,
			})
			c.Request, _ = http.NewRequest("GET", "/sellers/couriers?shopId=1", nil)

			handler.GetShipmentList(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestGetAllCouriers(t *testing.T) {
	type input struct {
		mockData []*model.Courier
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
			description: "should return error with status code 500 when failed to get couriers",
			input: input{
				mockData: nil,
				mockErr:  errors.New("failed to get couriers"),
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
			description: "should return couriers with status code 200 when succeed to get couriers",
			input: input{
				mockData: []*model.Courier{},
				mockErr:  nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    []*model.Courier{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			courierService := mocks.NewCourierService(t)
			courierService.On("GetAllCouriers").Return(tc.input.mockData, tc.input.mockErr)
			expectedRes, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			handler := handler.New(&handler.HandlerConfig{
				CourierService: courierService,
			})
			c.Request, _ = http.NewRequest(http.MethodGet, "/couriers", nil)

			handler.GetAllCouriers(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

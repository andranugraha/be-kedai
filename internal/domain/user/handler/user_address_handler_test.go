package handler_test

import (
	"encoding/json"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"
	"kedai/backend/be-kedai/internal/domain/user/handler"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAddUserAddress(t *testing.T) {
	var (
		trueValue = true
	)
	type input struct {
		data        *dto.AddAddressRequest
		beforeTests func(mockUserAddressService *mocks.UserAddressService)
	}
	type expected struct {
		data       *response.Response
		statusCode int
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "response status bad request when request is invalid",
			input: input{
				data: &dto.AddAddressRequest{
					PhoneNumber:   "asd",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("AddUserAddress", &dto.AddAddressRequest{
						PhoneNumber: "asd",
					}).Return(nil, nil)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.BAD_REQUEST,
					Message: "PhoneNumber must be numeric",
				},
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description: "response status Not Found when error is ErrProvinceNotFound or ErrCityNotFound or ErrSubdistrictNotFound or ErrDistrictNotFound",
			input: input{
				data: &dto.AddAddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					UserID:        1,
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("AddUserAddress", &dto.AddAddressRequest{
						PhoneNumber:   "123456789123",
						SubdistrictID: 1,
						Street:        "asd",
						Name:          "asd",
						UserID:        1,
						IsDefault:     &trueValue,
					}).Return(nil, errs.ErrProvinceNotFound)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.NOT_FOUND,
					Message: errs.ErrProvinceNotFound.Error(),
				},
				statusCode: http.StatusNotFound,
			},
		},
		{
			description: "response status conflic when error is ErrMaxAddress",
			input: input{
				data: &dto.AddAddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					UserID:        1,
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("AddUserAddress", &dto.AddAddressRequest{
						PhoneNumber:   "123456789123",
						SubdistrictID: 1,
						Street:        "asd",
						Name:          "asd",
						UserID:        1,
						IsDefault:     &trueValue,
					}).Return(nil, errs.ErrMaxAddress)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.MAX_ADDRESS_REACHED,
					Message: errs.ErrMaxAddress.Error(),
				},
				statusCode: http.StatusConflict,
			},
		},
		{
			description: "response status Internal Server Error when error is not ErrProvinceNotFound or ErrCityNotFound or ErrSubdistrictNotFound or ErrDistrictNotFound",
			input: input{
				data: &dto.AddAddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					UserID:        1,
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("AddUserAddress", &dto.AddAddressRequest{
						PhoneNumber:   "123456789123",
						SubdistrictID: 1,
						Street:        "asd",
						Name:          "asd",
						UserID:        1,
						IsDefault:     &trueValue,
					}).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errs.ErrInternalServerError.Error(),
				},
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			description: "response status CREATED when request is valid",
			input: input{
				data: &dto.AddAddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					UserID:        1,
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("AddUserAddress", &dto.AddAddressRequest{
						PhoneNumber:   "123456789123",
						SubdistrictID: 1,
						Street:        "asd",
						Name:          "asd",
						UserID:        1,
						IsDefault:     &trueValue,
					}).Return(&model.UserAddress{
						PhoneNumber:   "123456789123",
						SubdistrictID: 1,
						Street:        "asd",
						Name:          "asd",
						UserID:        1,
					}, nil)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.CREATED,
					Message: "created",
					Data: &model.UserAddress{
						PhoneNumber:   "123456789123",
						SubdistrictID: 1,
						Street:        "asd",
						Name:          "asd",
						UserID:        1,
					},
				},
				statusCode: http.StatusCreated,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", 1)

			payload := test.MakeRequestBody(tc.input.data)
			c.Request, _ = http.NewRequest(http.MethodGet, "/users/addresses", payload)

			mockUserAddressService := new(mocks.UserAddressService)
			tc.input.beforeTests(mockUserAddressService)

			handler := handler.New(&handler.HandlerConfig{
				UserAddressService: mockUserAddressService,
			})
			handler.AddUserAddress(c)

			expectedJson, _ := json.Marshal(tc.expected.data)
			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())

		})
	}

}

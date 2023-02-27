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
		data        *dto.AddressRequest
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
				data: &dto.AddressRequest{
					PhoneNumber:   "asd",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("AddUserAddress", &dto.AddressRequest{
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
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					UserID:        1,
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("AddUserAddress", &dto.AddressRequest{
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
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					UserID:        1,
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("AddUserAddress", &dto.AddressRequest{
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
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					UserID:        1,
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("AddUserAddress", &dto.AddressRequest{
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
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					UserID:        1,
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("AddUserAddress", &dto.AddressRequest{
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

func TestGetAllUserAddress(t *testing.T) {
	type input struct {
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
			description: "response status Internal Server Error when server error",
			input: input{
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("GetAllUserAddress", 1).Return(nil, errs.ErrInternalServerError)
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
			description: "response status OK when request is valid",
			input: input{
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("GetAllUserAddress", 1).Return([]*model.UserAddress{
						{
							ID:          1,
							PhoneNumber: "123456789123",
							Street:      "asd",
							Name:        "asd",
							UserID:      1,
							IsDefault:   false,
						},
					}, nil)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.OK,
					Message: "success",
					Data: []*model.UserAddress{
						{
							ID:          1,
							PhoneNumber: "123456789123",
							Street:      "asd",
							Name:        "asd",
							UserID:      1,
							IsDefault:   false,
						},
					},
				},
				statusCode: http.StatusOK,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", 1)

			c.Request, _ = http.NewRequest(http.MethodGet, "/users/addresses", nil)

			mockUserAddressService := new(mocks.UserAddressService)
			tc.input.beforeTests(mockUserAddressService)

			handler := handler.New(&handler.HandlerConfig{
				UserAddressService: mockUserAddressService,
			})
			handler.GetAllUserAddress(c)

			expectedJson, _ := json.Marshal(tc.expected.data)
			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())

		})
	}
}

func TestUpdateUserAddress(t *testing.T) {
	var (
		trueValue = true
	)
	type input struct {
		data        *dto.AddressRequest
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
			description: "response status Bad Request when request is invalid",
			input: input{
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "", // invalid
					UserID:        1,
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Name is required",
				},
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description: "response status Internal Server Error when server error",
			input: input{
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("UpdateUserAddress", &dto.AddressRequest{
						Name:          "asd",
						PhoneNumber:   "123456789123",
						Street:        "asd",
						SubdistrictID: 1,
						UserID:        1,
						IsDefault:     &trueValue,
						ID:            1,
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
			description: "response status not found when address not found",
			input: input{
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("UpdateUserAddress", &dto.AddressRequest{
						Name:          "asd",
						PhoneNumber:   "123456789123",
						Street:        "asd",
						SubdistrictID: 1,
						UserID:        1,
						IsDefault:     &trueValue,
						ID:            1,
					}).Return(nil, errs.ErrAddressNotFound)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.NOT_FOUND,
					Message: errs.ErrAddressNotFound.Error(),
				},
				statusCode: http.StatusNotFound,
			},
		},
		{
			description: "response status conflict when error ErrMustHaveAtLeastOneDefaultAddress when update address",
			input: input{
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("UpdateUserAddress", &dto.AddressRequest{
						Name:          "asd",
						PhoneNumber:   "123456789123",
						Street:        "asd",
						SubdistrictID: 1,
						UserID:        1,
						IsDefault:     &trueValue,
						ID:            1,
					}).Return(nil, errs.ErrMustHaveAtLeastOneDefaultAddress)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.MUST_HAVE_AT_LEAST_ONE_DEFAULT_ADDRESS,
					Message: errs.ErrMustHaveAtLeastOneDefaultAddress.Error(),
				},
				statusCode: http.StatusConflict,
			},
		},
		{
			description: "response status Internal Server Error when UpdateUserAddress return other error",
			input: input{
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("UpdateUserAddress", &dto.AddressRequest{
						Name:          "asd",
						PhoneNumber:   "123456789123",
						Street:        "asd",
						SubdistrictID: 1,
						UserID:        1,
						IsDefault:     &trueValue,
						ID:            1,
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
		// response status OK when update address success
		{
			description: "response status OK when update address success",
			input: input{
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					IsDefault:     &trueValue,
				},
				beforeTests: func(mockUserAddressService *mocks.UserAddressService) {
					mockUserAddressService.On("UpdateUserAddress", &dto.AddressRequest{
						Name:          "asd",
						PhoneNumber:   "123456789123",
						Street:        "asd",
						SubdistrictID: 1,
						UserID:        1,
						IsDefault:     &trueValue,
						ID:            1,
					}).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						Name:          "asd",
						PhoneNumber:   "123456789123",
						Street:        "asd",
						SubdistrictID: 1,
						IsDefault:     trueValue,
					}, nil)

				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.OK,
					Message: "success",
					Data: &model.UserAddress{
						ID:            1,
						UserID:        1,
						Name:          "asd",
						PhoneNumber:   "123456789123",
						Street:        "asd",
						SubdistrictID: 1,
						IsDefault:     trueValue,
					},
				},
				statusCode: http.StatusOK,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", 1)
			c.Params = gin.Params{
				{
					Key:   "addressId",
					Value: "1",
				},
			}

			payload := test.MakeRequestBody(tc.input.data)
			c.Request, _ = http.NewRequest(http.MethodPut, "/users/addresses", payload)

			mockUserAddressService := new(mocks.UserAddressService)
			tc.input.beforeTests(mockUserAddressService)

			handler := handler.New(&handler.HandlerConfig{
				UserAddressService: mockUserAddressService,
			})
			handler.UpdateUserAddress(c)

			expectedJson, _ := json.Marshal(tc.expected.data)
			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())

		})
	}
}

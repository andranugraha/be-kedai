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
		beforeTests func(mockAddressService *mocks.AddressService)
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
					IsPickup:      &trueValue,
				},
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("AddUserAddress", &dto.AddressRequest{
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
					IsPickup:      &trueValue,
				},
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("AddUserAddress", &dto.AddressRequest{
						PhoneNumber:   "123456789123",
						SubdistrictID: 1,
						Street:        "asd",
						Name:          "asd",
						UserID:        1,
						IsDefault:     &trueValue,
						IsPickup:      &trueValue,
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
			description: "response status conflict when error is ErrMaxAddress",
			input: input{
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					UserID:        1,
					IsDefault:     &trueValue,
					IsPickup:      &trueValue,
				},
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("AddUserAddress", &dto.AddressRequest{
						PhoneNumber:   "123456789123",
						SubdistrictID: 1,
						Street:        "asd",
						Name:          "asd",
						UserID:        1,
						IsDefault:     &trueValue,
						IsPickup:      &trueValue,
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
			description: "response status not found when error is ErrShopNotFound",
			input: input{
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					UserID:        1,
					IsDefault:     &trueValue,
					IsPickup:      &trueValue,
				},
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("AddUserAddress", &dto.AddressRequest{
						PhoneNumber:   "123456789123",
						SubdistrictID: 1,
						Street:        "asd",
						Name:          "asd",
						UserID:        1,
						IsDefault:     &trueValue,
						IsPickup:      &trueValue,
					}).Return(nil, errs.ErrShopNotFound)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.SHOP_NOT_REGISTERED,
					Message: errs.ErrShopNotFound.Error(),
				},
				statusCode: http.StatusNotFound,
			},
		},
		{
			description: "response status Internal Server Error",
			input: input{
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					UserID:        1,
					IsDefault:     &trueValue,
					IsPickup:      &trueValue,
				},
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("AddUserAddress", &dto.AddressRequest{
						PhoneNumber:   "123456789123",
						SubdistrictID: 1,
						Street:        "asd",
						Name:          "asd",
						UserID:        1,
						IsDefault:     &trueValue,
						IsPickup:      &trueValue,
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
					IsPickup:      &trueValue,
				},
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("AddUserAddress", &dto.AddressRequest{
						PhoneNumber:   "123456789123",
						SubdistrictID: 1,
						Street:        "asd",
						Name:          "asd",
						UserID:        1,
						IsDefault:     &trueValue,
						IsPickup:      &trueValue,
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

			mockAddressService := new(mocks.AddressService)
			tc.input.beforeTests(mockAddressService)

			handler := handler.New(&handler.HandlerConfig{
				AddressService: mockAddressService,
			})
			handler.AddUserAddress(c)

			expectedJson, _ := json.Marshal(tc.expected.data)
			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())

		})
	}

}

func TestGetAllUserAddress(t *testing.T) {
	var falseValue = false
	type input struct {
		beforeTests func(mockAddressService *mocks.AddressService)
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
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("GetAllUserAddress", 1).Return(nil, errs.ErrInternalServerError)
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
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("GetAllUserAddress", 1).Return([]*model.UserAddress{
						{
							ID:          1,
							PhoneNumber: "123456789123",
							Street:      "asd",
							Name:        "asd",
							UserID:      1,
							IsDefault:   &falseValue,
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
							IsDefault:   &falseValue,
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

			mockAddressService := new(mocks.AddressService)
			tc.input.beforeTests(mockAddressService)

			handler := handler.New(&handler.HandlerConfig{
				AddressService: mockAddressService,
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
		beforeTests func(mockAddressService *mocks.AddressService)
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
					IsPickup:      &trueValue,
				},
				beforeTests: func(mockAddressService *mocks.AddressService) {
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
					IsPickup:      &trueValue,
				},
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("UpdateUserAddress", &dto.AddressRequest{
						Name:          "asd",
						PhoneNumber:   "123456789123",
						Street:        "asd",
						SubdistrictID: 1,
						UserID:        1,
						IsDefault:     &trueValue,
						ID:            1,
						IsPickup:      &trueValue,
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
					IsPickup:      &trueValue,
				},
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("UpdateUserAddress", &dto.AddressRequest{
						Name:          "asd",
						PhoneNumber:   "123456789123",
						Street:        "asd",
						SubdistrictID: 1,
						UserID:        1,
						IsDefault:     &trueValue,
						IsPickup:      &trueValue,
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
					IsPickup:      &trueValue,
				},
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("UpdateUserAddress", &dto.AddressRequest{
						Name:          "asd",
						PhoneNumber:   "123456789123",
						Street:        "asd",
						SubdistrictID: 1,
						UserID:        1,
						IsDefault:     &trueValue,
						IsPickup:      &trueValue,
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
					IsPickup:      &trueValue,
				},
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("UpdateUserAddress", &dto.AddressRequest{
						Name:          "asd",
						PhoneNumber:   "123456789123",
						Street:        "asd",
						SubdistrictID: 1,
						UserID:        1,
						IsDefault:     &trueValue,
						IsPickup:      &trueValue,
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
			description: "response status Not Found when UpdateUserAddress return ErrShopNotFound",
			input: input{
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					IsDefault:     &trueValue,
					IsPickup:      &trueValue,
				},
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("UpdateUserAddress", &dto.AddressRequest{
						Name:          "asd",
						PhoneNumber:   "123456789123",
						Street:        "asd",
						SubdistrictID: 1,
						UserID:        1,
						IsDefault:     &trueValue,
						IsPickup:      &trueValue,
						ID:            1,
					}).Return(nil, errs.ErrShopNotFound)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.SHOP_NOT_REGISTERED,
					Message: errs.ErrShopNotFound.Error(),
				},
				statusCode: http.StatusNotFound,
			},
		},
		{
			description: "response status Conflict when UpdateUserAddress return ErrMustHaveAtLeastOnePickupAddress",
			input: input{
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					IsDefault:     &trueValue,
					IsPickup:      &trueValue,
				},
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("UpdateUserAddress", &dto.AddressRequest{
						Name:          "asd",
						PhoneNumber:   "123456789123",
						Street:        "asd",
						SubdistrictID: 1,
						UserID:        1,
						IsDefault:     &trueValue,
						IsPickup:      &trueValue,
						ID:            1,
					}).Return(nil, errs.ErrMustHaveAtLeastOnePickupAddress)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.MUST_HAVE_AT_LEAST_ONE_PICKUP_ADDRESS,
					Message: errs.ErrMustHaveAtLeastOnePickupAddress.Error(),
				},
				statusCode: http.StatusConflict,
			},
		},
		{
			description: "response status OK when update address success",
			input: input{
				data: &dto.AddressRequest{
					PhoneNumber:   "123456789123",
					SubdistrictID: 1,
					Street:        "asd",
					Name:          "asd",
					IsDefault:     &trueValue,
					IsPickup:      &trueValue,
				},
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("UpdateUserAddress", &dto.AddressRequest{
						Name:          "asd",
						PhoneNumber:   "123456789123",
						Street:        "asd",
						SubdistrictID: 1,
						UserID:        1,
						IsDefault:     &trueValue,
						IsPickup:      &trueValue,
						ID:            1,
					}).Return(&model.UserAddress{
						ID:            1,
						UserID:        1,
						Name:          "asd",
						PhoneNumber:   "123456789123",
						Street:        "asd",
						SubdistrictID: 1,
						IsDefault:     &trueValue,
						IsPickup:      &trueValue,
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
						IsDefault:     &trueValue,
						IsPickup:      &trueValue,
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

			mockAddressService := new(mocks.AddressService)
			tc.input.beforeTests(mockAddressService)

			handler := handler.New(&handler.HandlerConfig{
				AddressService: mockAddressService,
			})
			handler.UpdateUserAddress(c)

			expectedJson, _ := json.Marshal(tc.expected.data)
			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())

		})
	}
}

func TestDeleteUserAddress(t *testing.T) {
	type input struct {
		beforeTests func(mockAddressService *mocks.AddressService)
	}

	type expected struct {
		data       *response.Response
		statusCode int
	}

	type testCase struct {
		description string
		input       input
		expected    expected
	}

	cases := []testCase{
		{
			description: "response status not found when error ErrAddressNotFound when delete address",
			input: input{
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("DeleteUserAddress", 1, 1).Return(errs.ErrAddressNotFound)
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
			description: "response status conflict when error ErrMustHaveAtLeastOneDefaultAddress when delete address",
			input: input{
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("DeleteUserAddress", 1, 1).Return(errs.ErrMustHaveAtLeastOneDefaultAddress)
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
			description: "response status conflict when error ErrMustHaveAtLeastOnePickupAddress when delete address",
			input: input{
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("DeleteUserAddress", 1, 1).Return(errs.ErrMustHaveAtLeastOnePickupAddress)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.MUST_HAVE_AT_LEAST_ONE_PICKUP_ADDRESS,
					Message: errs.ErrMustHaveAtLeastOnePickupAddress.Error(),
				},
				statusCode: http.StatusConflict,
			},
		},
		{
			description: "response status Internal Server Error when DeleteUserAddress return other error",
			input: input{
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("DeleteUserAddress", 1, 1).Return(errs.ErrInternalServerError)
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
			description: "response status OK when delete address success",
			input: input{
				beforeTests: func(mockAddressService *mocks.AddressService) {
					mockAddressService.On("DeleteUserAddress", 1, 1).Return(nil)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.OK,
					Message: "success",
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

			c.Request, _ = http.NewRequest(http.MethodDelete, "/users/addresses", nil)

			mockAddressService := new(mocks.AddressService)
			tc.input.beforeTests(mockAddressService)

			handler := handler.New(&handler.HandlerConfig{
				AddressService: mockAddressService,
			})
			handler.DeleteUserAddress(c)

			expectedJson, _ := json.Marshal(tc.expected.data)
			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())

		})
	}
}

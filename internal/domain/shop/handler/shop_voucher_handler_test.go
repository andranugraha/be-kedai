package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"kedai/backend/be-kedai/internal/common/code"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/handler"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetShopVoucher(t *testing.T) {
	var (
		slug    = "shop"
		voucher = []*model.ShopVoucher{}
	)
	type input struct {
		slug    string
		voucher []*model.ShopVoucher
		err     error
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
			description: "should return list of voucher with code 200 when successful",
			input: input{
				slug:    slug,
				voucher: voucher,
				err:     nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    voucher,
				},
			},
		},
		{
			description: "should return error with code 404 when shop not found",
			input: input{
				slug:    slug,
				voucher: nil,
				err:     errs.ErrShopNotFound,
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
				slug:    slug,
				voucher: nil,
				err:     errs.ErrInternalServerError,
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
			c.Params = gin.Params{
				{
					Key:   "slug",
					Value: slug,
				},
			}
			mockService := new(mocks.ShopVoucherService)
			mockService.On("GetShopVoucher", slug).Return(tc.input.voucher, tc.input.err)
			handler := handler.New(&handler.HandlerConfig{
				ShopVoucherService: mockService,
			})
			c.Request, _ = http.NewRequest("GET", "/shops/:slug", nil)

			handler.GetShopVoucher(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestGetSellerVoucher(t *testing.T) {
	type input struct {
		userID   int
		request  *dto.SellerVoucherFilterRequest
		mockData *commonDto.PaginationResponse
		mockErr  error
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID     = 1
		page       = 2
		limit      = 10
		totalRows  = int64(0)
		vouchers   = []*dto.SellerVoucher{}
		totalPages = 0
		request    = &dto.SellerVoucherFilterRequest{
			Page:  page,
			Limit: limit,
		}
	)

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error with status code 404 when shop does not exist",
			input: input{
				userID:   userID,
				request:  request,
				mockData: nil,
				mockErr:  errs.ErrShopNotFound,
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
			description: "should return error with status code 500 when something went wrong",
			input: input{
				userID:   userID,
				request:  request,
				mockData: nil,
				mockErr:  errors.New("something went wrong"),
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
			description: "should return data with status code 200 when succeed fetching vouchers",
			input: input{
				userID:  userID,
				request: request,
				mockData: &commonDto.PaginationResponse{
					TotalRows:  totalRows,
					TotalPages: totalPages,
					Page:       page,
					Limit:      limit,
					Data:       vouchers,
				},
				mockErr: nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data: &commonDto.PaginationResponse{
						TotalRows:  totalRows,
						TotalPages: totalPages,
						Page:       page,
						Limit:      limit,
						Data:       vouchers,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		expectedRes, _ := json.Marshal(tc.expected.response)
		shopVoucherService := mocks.NewShopVoucherService(t)
		shopVoucherService.On("GetSellerVoucher", tc.input.userID, tc.input.request).Return(tc.input.mockData, tc.input.mockErr)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Set("userId", tc.input.userID)
		h := handler.New(&handler.HandlerConfig{
			ShopVoucherService: shopVoucherService,
		})
		c.Request = httptest.NewRequest("GET", fmt.Sprintf("/v1/sellers/vouchers?page=%d&limit=%d", tc.input.request.Page, tc.input.request.Limit), nil)

		h.GetSellerVoucher(c)

		assert.Equal(t, tc.expected.statusCode, rec.Code)
		assert.Equal(t, string(expectedRes), rec.Body.String())
	}
}

func TestGetVoucherByCodeAndShopId(t *testing.T) {
	type input struct {
		userID      int
		voucherCode string
		mockData    *dto.SellerVoucher
		mockErr     error
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID      = 1
		voucherCode = "voucher-code"
	)

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error with status code 404 when shop not found",
			input: input{
				userID:      userID,
				voucherCode: voucherCode,
				mockData:    nil,
				mockErr:     errs.ErrShopNotFound,
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
			description: "should return error with status code 404 when voucher not found",
			input: input{
				userID:      userID,
				voucherCode: voucherCode,
				mockData:    nil,
				mockErr:     errs.ErrVoucherNotFound,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.VOUCHER_NOT_FOUND,
					Message: errs.ErrVoucherNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when something went wrong",
			input: input{
				userID:      userID,
				voucherCode: voucherCode,
				mockData:    nil,
				mockErr:     errs.ErrInternalServerError,
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
			description: "should return voucher detail with status code 200 when succeed to get voucher",
			input: input{
				userID:      userID,
				voucherCode: voucherCode,
				mockData:    &dto.SellerVoucher{},
				mockErr:     nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    &dto.SellerVoucher{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			shopVoucherService := mocks.NewShopVoucherService(t)
			shopVoucherService.On("GetVoucherByCodeAndShopId", tc.input.voucherCode, tc.input.userID).Return(tc.input.mockData, tc.input.mockErr)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			c.AddParam("code", tc.input.voucherCode)
			h := handler.New(&handler.HandlerConfig{
				ShopVoucherService: shopVoucherService,
			})
			c.Request = httptest.NewRequest("GET", fmt.Sprintf("/sellers/vouchers?%s", tc.input.voucherCode), nil)

			h.GetVoucherByCodeAndShopId(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestCreateVoucher(t *testing.T) {
	type input struct {
		userID  int
		request *dto.CreateVoucherRequest
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID               = 1
		voucherName          = "voucher name"
		voucherCode          = "VOUC123AB"
		amount       float64 = 15
		voucherType          = "percent"
		isHidden             = false
		description          = "description"
		minimumSpend float64 = 1000
		totalQuota           = 10
		startFrom, _         = time.Parse("2006-01-02", "2006-01-02")
		expiredAt, _         = time.Parse("2006-01-02", "2006-01-14")

		request = &dto.CreateVoucherRequest{
			Name:         voucherName,
			Code:         voucherCode,
			Amount:       amount,
			Type:         voucherType,
			IsHidden:     &isHidden,
			Description:  description,
			MinimumSpend: minimumSpend,
			TotalQuota:   totalQuota,
			StartFrom:    startFrom,
			ExpiredAt:    expiredAt,
		}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ShopVoucherService)
		expected
	}{
		{
			description: "should return error with status code 400 when given invalid request body",
			input: input{
				userID: userID,
				request: &dto.CreateVoucherRequest{
					Name:         voucherName,
					Code:         voucherCode,
					Amount:       amount,
					Type:         voucherType,
					IsHidden:     &isHidden,
					Description:  "desc",
					MinimumSpend: minimumSpend,
					TotalQuota:   totalQuota,
					StartFrom:    startFrom,
					ExpiredAt:    expiredAt,
				},
			},
			beforeTest: func(vs *mocks.ShopVoucherService) {},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Description must be greater than 5",
				},
			},
		},
		{
			description: "should return error with status code 404 when failed to get shop",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(vs *mocks.ShopVoucherService) {
				vs.On("CreateVoucher", userID, request).Return(nil, errs.ErrShopNotFound)
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
			description: "should return error with status code 422 when voucher name is invalid",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(vs *mocks.ShopVoucherService) {
				vs.On("CreateVoucher", userID, request).Return(nil, errs.ErrInvalidVoucherNamePattern)
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_VOUCHER_NAME,
					Message: errs.ErrInvalidVoucherNamePattern.Error(),
				},
			},
		},
		{
			description: "should return error with status code 422 when voucher date range is invalid",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(vs *mocks.ShopVoucherService) {
				vs.On("CreateVoucher", userID, request).Return(nil, errs.ErrInvalidVoucherDateRange)
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_DATE_RANGE,
					Message: errs.ErrInvalidVoucherDateRange.Error(),
				},
			},
		},
		{
			description: "should return error with status code 409 when duplicate voucher code",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(vs *mocks.ShopVoucherService) {
				vs.On("CreateVoucher", userID, request).Return(nil, errs.ErrDuplicateVoucherCode)
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
			description: "should return error with status code 500 when failed to create voucher",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(vs *mocks.ShopVoucherService) {
				vs.On("CreateVoucher", userID, request).Return(nil, errors.New("failed to create voucher"))
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
			description: "should return voucher data with status code 201 when succeed to create voucher",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(vs *mocks.ShopVoucherService) {
				vs.On("CreateVoucher", userID, request).Return(&model.ShopVoucher{}, nil)
			},
			expected: expected{
				statusCode: http.StatusCreated,
				response: response.Response{
					Code:    code.CREATED,
					Message: "voucher created",
					Data:    &model.ShopVoucher{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			shopVoucherService := mocks.NewShopVoucherService(t)
			tc.beforeTest(shopVoucherService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			h := handler.New(&handler.HandlerConfig{
				ShopVoucherService: shopVoucherService,
			})
			payload := test.MakeRequestBody(tc.input.request)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/sellers/vouchers", payload)

			h.CreateVoucher(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestDeleteVoucher(t *testing.T) {
	type input struct {
		beforeTests func(mockShopVoucherService *mocks.ShopVoucherService)
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
			description: "should return error with status code 404 when failed to get shop",
			input: input{
				beforeTests: func(vs *mocks.ShopVoucherService) {
					vs.On("DeleteVoucher", 1, "BAKM12a").Return(errs.ErrShopNotFound)
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
			description: "response status not found when error ErrVoucherNotFound when delete voucher",
			input: input{
				beforeTests: func(vs *mocks.ShopVoucherService) {
					vs.On("DeleteVoucher", 1, "BAKM12a").Return(errs.ErrVoucherNotFound)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.VOUCHER_NOT_FOUND,
					Message: errs.ErrVoucherNotFound.Error(),
				},
				statusCode: http.StatusNotFound,
			},
		},
		{
			description: "response status conflict when error ErrVoucherStatusConflict when delete voucher",
			input: input{
				beforeTests: func(vs *mocks.ShopVoucherService) {
					vs.On("DeleteVoucher", 1, "BAKM12a").Return(errs.ErrVoucherStatusConflict)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.VOUCHER_STATUS_CONFLICT,
					Message: errs.ErrVoucherStatusConflict.Error(),
				},
				statusCode: http.StatusConflict,
			},
		},
		{
			description: "response status Internal Server Error when DeleteVoucher return other error",
			input: input{
				beforeTests: func(vs *mocks.ShopVoucherService) {
					vs.On("DeleteVoucher", 1, "BAKM12a").Return(errs.ErrInternalServerError)
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
			description: "response status OK when delete voucher success",
			input: input{
				beforeTests: func(vs *mocks.ShopVoucherService) {
					vs.On("DeleteVoucher", 1, "BAKM12a").Return(nil)
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
					Key:   "code",
					Value: "BAKM12a",
				},
			}

			c.Request, _ = http.NewRequest(http.MethodDelete, "/selleres/vouchers", nil)

			mockShopVoucherService := new(mocks.ShopVoucherService)
			tc.input.beforeTests(mockShopVoucherService)

			handler := handler.New(&handler.HandlerConfig{
				ShopVoucherService: mockShopVoucherService,
			})
			handler.DeleteVoucher(c)

			expectedJson, _ := json.Marshal(tc.expected.data)
			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())

		})
	}
}

func TestUpdateVoucher(t *testing.T) {
	type input struct {
		userID      int
		voucherCode string
		request     *dto.UpdateVoucherRequest
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID      = 1
		voucherCode = "voucher-code"

		voucherName          = "voucher name"
		amount       float64 = 15
		voucherType          = "percent"
		isHidden             = false
		description          = "description"
		minimumSpend float64 = 1000
		totalQuota           = 10
		startFrom, _         = time.Parse("2006-01-02", "2006-01-02")
		expiredAt, _         = time.Parse("2006-01-02", "2006-01-14")

		request = &dto.UpdateVoucherRequest{
			Name:         voucherName,
			Amount:       amount,
			Type:         voucherType,
			IsHidden:     &isHidden,
			Description:  description,
			MinimumSpend: minimumSpend,
			TotalQuota:   totalQuota,
			StartFrom:    startFrom,
			ExpiredAt:    expiredAt,
		}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ShopVoucherService)
		expected
	}{
		{
			description: "should return error with status code 400 when given invalid request body",
			input: input{
				userID: userID,
				request: &dto.UpdateVoucherRequest{
					Name:         voucherName,
					Amount:       amount,
					Type:         voucherType,
					IsHidden:     &isHidden,
					Description:  "desc",
					MinimumSpend: minimumSpend,
					TotalQuota:   totalQuota,
					StartFrom:    startFrom,
					ExpiredAt:    expiredAt,
				},
			},
			beforeTest: func(vs *mocks.ShopVoucherService) {},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Description must be greater than 5",
				},
			},
		},
		{
			description: "should return error with status code 404 when failed to get shop",
			input: input{
				userID:      userID,
				voucherCode: voucherCode,
				request:     request,
			},
			beforeTest: func(vs *mocks.ShopVoucherService) {
				vs.On("UpdateVoucher", userID, voucherCode, request).Return(nil, errs.ErrShopNotFound)
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
			description: "should return error with status code 404 and VOUCHER_NOT_FOUND code when voucher is not found",
			input: input{
				userID:      userID,
				voucherCode: voucherCode,
				request:     request,
			},
			beforeTest: func(vs *mocks.ShopVoucherService) {
				vs.On("UpdateVoucher", userID, voucherCode, request).Return(nil, errs.ErrVoucherNotFound)
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.VOUCHER_NOT_FOUND,
					Message: errs.ErrVoucherNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with status code 422 when voucher name is invalid",
			input: input{
				userID:      userID,
				voucherCode: voucherCode,
				request:     request,
			},
			beforeTest: func(vs *mocks.ShopVoucherService) {
				vs.On("UpdateVoucher", userID, voucherCode, request).Return(nil, errs.ErrInvalidVoucherNamePattern)
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_VOUCHER_NAME,
					Message: errs.ErrInvalidVoucherNamePattern.Error(),
				},
			},
		},
		{
			description: "should return error with status code 422 when voucher date range is invalid",
			input: input{
				userID:      userID,
				voucherCode: voucherCode,
				request:     request,
			},
			beforeTest: func(vs *mocks.ShopVoucherService) {
				vs.On("UpdateVoucher", userID, voucherCode, request).Return(nil, errs.ErrInvalidVoucherDateRange)
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_DATE_RANGE,
					Message: errs.ErrInvalidVoucherDateRange.Error(),
				},
			},
		},
		{
			description: "should return error with status code 409 and VOUCHER_STATUS_CONFLICT code when voucher status conflict",
			input: input{
				userID:      userID,
				voucherCode: voucherCode,
				request:     request,
			},
			beforeTest: func(vs *mocks.ShopVoucherService) {
				vs.On("UpdateVoucher", userID, voucherCode, request).Return(nil, errs.ErrVoucherStatusConflict)
			},
			expected: expected{
				statusCode: http.StatusConflict,
				response: response.Response{
					Code:    code.VOUCHER_STATUS_CONFLICT,
					Message: errs.ErrVoucherStatusConflict.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when something went wrong",
			input: input{
				userID:      userID,
				voucherCode: voucherCode,
				request:     request,
			},
			beforeTest: func(vs *mocks.ShopVoucherService) {
				vs.On("UpdateVoucher", userID, voucherCode, request).Return(nil, errors.New("something went wrong"))
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
			description: "should return success with status code 200 and UPDATED code when voucher is updated successfully",
			input: input{
				userID:      userID,
				voucherCode: voucherCode,
				request:     request,
			},
			beforeTest: func(vs *mocks.ShopVoucherService) {
				vs.On("UpdateVoucher", userID, voucherCode, request).Return(&model.ShopVoucher{}, nil)
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.UPDATED,
					Message: "update voucher succesful",
					Data:    &model.ShopVoucher{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			shopVoucherService := mocks.NewShopVoucherService(t)
			tc.beforeTest(shopVoucherService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			c.AddParam("code", tc.input.voucherCode)
			h := handler.New(&handler.HandlerConfig{
				ShopVoucherService: shopVoucherService,
			})
			payload := test.MakeRequestBody(tc.input.request)
			c.Request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/v1/sellers/vouchers?%s", tc.input.voucherCode), payload)

			h.UpdateVoucher(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestGetValidShopVoucher(t *testing.T) {
	var (
		slug    = "shop"
		voucher = []*model.ShopVoucher{}
	)
	type input struct {
		slug    string
		voucher []*model.ShopVoucher
		err     error
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
			description: "should return list of voucher with code 200 when successful",
			input: input{
				slug:    slug,
				voucher: voucher,
				err:     nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    voucher,
				},
			},
		},
		{
			description: "should return error with code 404 when shop not found",
			input: input{
				slug:    slug,
				voucher: nil,
				err:     errs.ErrShopNotFound,
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
				slug:    slug,
				voucher: nil,
				err:     errs.ErrInternalServerError,
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
			c.Params = gin.Params{
				{
					Key:   "slug",
					Value: slug,
				},
			}
			mockService := new(mocks.ShopVoucherService)
			mockService.On("GetValidShopVoucherByUserIDAndSlug", dto.GetValidShopVoucherRequest{
				Slug: slug,
			}).Return(tc.input.voucher, tc.input.err)
			handler := handler.New(&handler.HandlerConfig{
				ShopVoucherService: mockService,
			})
			c.Request, _ = http.NewRequest("GET", "/shops/:slug/vouchers/valid", nil)

			handler.GetValidShopVoucher(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

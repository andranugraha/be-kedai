package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"kedai/backend/be-kedai/internal/common/code"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	productDto "kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/handler"
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

func TestGetSellerPromotion(t *testing.T) {
	type input struct {
		userID   int
		request  *dto.SellerPromotionFilterRequest
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
		promotions = []*dto.SellerPromotion{}
		totalPages = 0
		request    = &dto.SellerPromotionFilterRequest{
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
			description: "should return error with status code 200 when suceed fetching promotions",
			input: input{
				userID:  userID,
				request: request,
				mockData: &commonDto.PaginationResponse{
					TotalRows:  totalRows,
					TotalPages: totalPages,
					Page:       page,
					Limit:      limit,
					Data:       promotions,
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
						Data:       promotions,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		expectedRes, _ := json.Marshal(tc.expected.response)
		shopPromotionService := mocks.NewShopPromotionService(t)
		shopPromotionService.On("GetSellerPromotions", tc.input.userID, tc.input.request).Return(tc.input.mockData, tc.input.mockErr)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Set("userId", tc.input.userID)
		h := handler.New(&handler.HandlerConfig{
			ShopPromotionService: shopPromotionService,
		})
		c.Request = httptest.NewRequest("GET", fmt.Sprintf("/v1/sellers/promotions?page=%d&limit=%d", tc.input.request.Page, tc.input.request.Limit), nil)

		h.GetSellerPromotions(c)

		assert.Equal(t, tc.expected.statusCode, rec.Code)
		assert.Equal(t, string(expectedRes), rec.Body.String())
	}
}

func TestGetSellerPromotionById(t *testing.T) {
	type input struct {
		userID      int
		promotionId int
		mockData    *dto.SellerPromotion
		mockErr     error
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID      = 1
		promotionId = 1
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
				promotionId: promotionId,
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
			description: "should return error with status code 404 when promotion not found",
			input: input{
				userID:      userID,
				promotionId: promotionId,
				mockData:    nil,
				mockErr:     errs.ErrPromotionNotFound,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.PROMOTION_NOT_FOUND,
					Message: errs.ErrPromotionNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when something went wrong",
			input: input{
				userID:      userID,
				promotionId: promotionId,
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
			description: "should return promotion detail with status code 200 when succeed to get promotion",
			input: input{
				userID:      userID,
				promotionId: promotionId,
				mockData:    &dto.SellerPromotion{},
				mockErr:     nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    &dto.SellerPromotion{},
				},
			},
		},
	}

	for _, tc := range tests {
		expectedRes, _ := json.Marshal(tc.expected.response)
		shopPromotionService := mocks.NewShopPromotionService(t)
		shopPromotionService.On("GetSellerPromotionById", tc.input.userID, tc.input.promotionId).Return(tc.input.mockData, tc.input.mockErr)
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Set("userId", tc.input.userID)
		c.AddParam("promotionId", "1")
		h := handler.New(&handler.HandlerConfig{
			ShopPromotionService: shopPromotionService,
		})
		c.Request = httptest.NewRequest("GET", fmt.Sprintf("/v1/sellers/promotions?%d", tc.input.promotionId), nil)

		h.GetSellerPromotionById(c)

		assert.Equal(t, tc.expected.statusCode, rec.Code)
		assert.Equal(t, string(expectedRes), rec.Body.String())
	}
}

func TestUpdatePromotion(t *testing.T) {
	type input struct {
		userID      int
		promotionID int
		request     dto.UpdateShopPromotionRequest
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID         = 1
		promotionID    = 1
		promotionName  = "promotion name"
		startPeriod, _ = time.Parse("2006-01-02", "2006-01-02")
		endPeriod, _   = time.Parse("2006-01-02", "2006-01-14")

		request = dto.UpdateShopPromotionRequest{
			Name:              promotionName,
			StartPeriod:       startPeriod,
			EndPeriod:         endPeriod,
			ProductPromotions: []*productDto.UpdateProductPromotionRequest{},
		}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ShopPromotionService)
		expected
	}{
		{
			description: "should return error with status code 400 when given invalid request body",
			input: input{
				userID:      userID,
				promotionID: promotionID,
				request: dto.UpdateShopPromotionRequest{
					Name:              "",
					StartPeriod:       startPeriod,
					EndPeriod:         endPeriod,
					ProductPromotions: request.ProductPromotions,
				},
			},
			beforeTest: func(vs *mocks.ShopPromotionService) {},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Name must be greater than 1",
				},
			},
		},
		{
			description: "should return error with status code 404 when failed to get shop",
			input: input{
				userID:      userID,
				promotionID: promotionID,
				request:     request,
			},
			beforeTest: func(ps *mocks.ShopPromotionService) {
				ps.On("UpdatePromotion", userID, promotionID, request).Return(errs.ErrShopNotFound)
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
			description: "should return error with status code 404 and PROMOTION_NOT_FOUND code when promotion is not found",
			input: input{
				userID:      userID,
				promotionID: promotionID,
				request:     request,
			},
			beforeTest: func(ps *mocks.ShopPromotionService) {
				ps.On("UpdatePromotion", userID, promotionID, request).Return(errs.ErrPromotionNotFound)
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.PROMOTION_NOT_FOUND,
					Message: errs.ErrPromotionNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with status code 422 when promotion name is invalid",
			input: input{
				userID:      userID,
				promotionID: promotionID,
				request:     request,
			},
			beforeTest: func(ps *mocks.ShopPromotionService) {
				ps.On("UpdatePromotion", userID, promotionID, request).Return(errs.ErrInvalidPromotionNamePattern)
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_PROMOTION_NAME,
					Message: errs.ErrInvalidPromotionNamePattern.Error(),
				},
			},
		},
		{
			description: "should return error with status code 422 when promotion fields cant be edited",
			input: input{
				userID:      userID,
				promotionID: promotionID,
				request:     request,
			},
			beforeTest: func(vs *mocks.ShopPromotionService) {
				vs.On("UpdatePromotion", userID, promotionID, request).Return(errs.ErrPromotionFieldsCantBeEdited)
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.PROMOTION_FIELDS_CANT_BE_EDITED,
					Message: errs.ErrPromotionFieldsCantBeEdited.Error(),
				},
			},
		},
		{
			description: "should return error with status code 422 when voucher date range is invalid",
			input: input{
				userID:      userID,
				promotionID: promotionID,
				request:     request,
			},
			beforeTest: func(ps *mocks.ShopPromotionService) {
				ps.On("UpdatePromotion", userID, promotionID, request).Return(errs.ErrInvalidPromotionDateRange)
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_DATE_RANGE,
					Message: errs.ErrInvalidPromotionDateRange.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when something went wrong",
			input: input{
				userID:      userID,
				promotionID: promotionID,
				request:     request,
			},
			beforeTest: func(ps *mocks.ShopPromotionService) {
				ps.On("UpdatePromotion", userID, promotionID, request).Return(errors.New("something went wrong"))
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
			description: "should return success with status code 200 and UPDATED code when promotion is updated successfully",
			input: input{
				userID:      userID,
				promotionID: promotionID,
				request:     request,
			},
			beforeTest: func(ps *mocks.ShopPromotionService) {
				ps.On("UpdatePromotion", userID, promotionID, request).Return(nil)
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.UPDATED,
					Message: "update promotion succesful",
					Data:    nil,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			shopPromotionService := mocks.NewShopPromotionService(t)
			tc.beforeTest(shopPromotionService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			c.AddParam("promotionId", "1")
			h := handler.New(&handler.HandlerConfig{
				ShopPromotionService: shopPromotionService,
			})
			payload := test.MakeRequestBody(tc.input.request)
			c.Request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/v1/sellers/promotions?%d", tc.input.promotionID), payload)

			h.UpdatePromotion(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestCreatePromotion(t *testing.T) {
	type input struct {
		userID  int
		request *dto.CreateShopPromotionRequest
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID         = 1
		promotionName  = "promotion name"
		startPeriod, _ = time.Parse("2006-01-02", "2006-01-02")
		endPeriod, _   = time.Parse("2006-01-02", "2006-01-14")

		request = &dto.CreateShopPromotionRequest{
			Name:              promotionName,
			StartPeriod:       startPeriod,
			EndPeriod:         endPeriod,
			ProductPromotions: []*productDto.CreateProductPromotionRequest{},
		}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ShopPromotionService)
		expected
	}{
		{
			description: "should return error with status code 400 when given invalid request body",
			input: input{
				userID: userID,
				request: &dto.CreateShopPromotionRequest{
					Name:              "",
					StartPeriod:       startPeriod,
					EndPeriod:         endPeriod,
					ProductPromotions: request.ProductPromotions,
				},
			},
			beforeTest: func(vs *mocks.ShopPromotionService) {},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Name is required",
				},
			},
		},
		{
			description: "should return error with status code 404 when failed to get shop",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ps *mocks.ShopPromotionService) {
				ps.On("CreateShopPromotion", userID, request).Return(nil, errs.ErrShopNotFound)
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
			description: "should return error with status code 422 when promotion name is invalid",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ps *mocks.ShopPromotionService) {
				ps.On("CreateShopPromotion", userID, request).Return(nil, errs.ErrInvalidPromotionNamePattern)
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_PROMOTION_NAME,
					Message: errs.ErrInvalidPromotionNamePattern.Error(),
				},
			},
		},
		{
			description: "should return error with status code 422 when promotion date range is invalid",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ps *mocks.ShopPromotionService) {
				ps.On("CreateShopPromotion", userID, request).Return(nil, errs.ErrInvalidPromotionDateRange)
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_DATE_RANGE,
					Message: errs.ErrInvalidPromotionDateRange.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to create promotion",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(ps *mocks.ShopPromotionService) {
				ps.On("CreateShopPromotion", userID, request).Return(nil, errors.New("failed to create promotion"))
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
			description: "should return promotion data with status code 201 when succeed to create promotion",
			input: input{
				userID:  userID,
				request: request,
			},
			beforeTest: func(vs *mocks.ShopPromotionService) {
				vs.On("CreateShopPromotion", userID, request).Return(&dto.CreateShopPromotionResponse{}, nil)
			},
			expected: expected{
				statusCode: http.StatusCreated,
				response: response.Response{
					Code:    code.CREATED,
					Message: "promotion created",
					Data:    &dto.CreateShopPromotionResponse{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			shopPromotionService := mocks.NewShopPromotionService(t)
			tc.beforeTest(shopPromotionService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			h := handler.New(&handler.HandlerConfig{
				ShopPromotionService: shopPromotionService,
			})
			payload := test.MakeRequestBody(tc.input.request)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/sellers/promotions", payload)

			h.CreateShopPromotion(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestDeletePromotion(t *testing.T) {
	type input struct {
		beforeTests func(mockShopPromotionService *mocks.ShopPromotionService)
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
				beforeTests: func(ps *mocks.ShopPromotionService) {
					ps.On("DeletePromotion", 1, 1).Return(errs.ErrShopNotFound)
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
			description: "response status not found when error ErrPromotionNotFound when delete promotion",
			input: input{
				beforeTests: func(vs *mocks.ShopPromotionService) {
					vs.On("DeletePromotion", 1, 1).Return(errs.ErrPromotionNotFound)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.PROMOTION_NOT_FOUND,
					Message: errs.ErrPromotionNotFound.Error(),
				},
				statusCode: http.StatusNotFound,
			},
		},
		{
			description: "response status conflict when error ErrPromotionStatusConflict when delete promotion",
			input: input{
				beforeTests: func(vs *mocks.ShopPromotionService) {
					vs.On("DeletePromotion", 1, 1).Return(errs.ErrPromotionStatusConflict)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.PROMOTION_STATUS_CONFLICT,
					Message: errs.ErrPromotionStatusConflict.Error(),
				},
				statusCode: http.StatusConflict,
			},
		},
		{
			description: "response status Internal Server Error when DeletePromotion return other error",
			input: input{
				beforeTests: func(vs *mocks.ShopPromotionService) {
					vs.On("DeletePromotion", 1, 1).Return(errs.ErrInternalServerError)
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
			description: "response status OK when delete promotion success",
			input: input{
				beforeTests: func(vs *mocks.ShopPromotionService) {
					vs.On("DeletePromotion", 1, 1).Return(nil)
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
					Key:   "promotionId",
					Value: "1",
				},
			}

			c.Request, _ = http.NewRequest(http.MethodDelete, "/selleres/promotions", nil)

			mockShopPromotionService := new(mocks.ShopPromotionService)
			tc.input.beforeTests(mockShopPromotionService)

			handler := handler.New(&handler.HandlerConfig{
				ShopPromotionService: mockShopPromotionService,
			})
			handler.DeletePromotion(c)

			expectedJson, _ := json.Marshal(tc.expected.data)
			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedJson), rec.Body.String())

		})
	}
}

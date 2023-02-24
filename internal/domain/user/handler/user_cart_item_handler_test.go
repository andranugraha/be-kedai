package handler_test

import (
	"encoding/json"
	"kedai/backend/be-kedai/internal/common/code"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/handler"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCartItem(t *testing.T) {
	var (
		invalidReq = &dto.UserCartItemRequest{
			Quantity: 0,
			Notes:    "test",
			UserId:   1,
			SkuId:    1,
		}
		validReq = &dto.UserCartItemRequest{
			Quantity: 1,
			Notes:    "test",
			UserId:   1,
			SkuId:    1,
		}
	)

	type input struct {
		data        *dto.UserCartItemRequest
		beforeTests func(mockUserCartItemService *mocks.UserCartItemService)
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
				data: invalidReq,
				beforeTests: func(mockUserCartItemService *mocks.UserCartItemService) {
					mockUserCartItemService.On("CreateCartItem", invalidReq).Return(nil, nil)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Quantity is required",
				},
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description: "response status not found when product not found",
			input: input{
				data: validReq,
				beforeTests: func(mockUserCartItemService *mocks.UserCartItemService) {
					mockUserCartItemService.On("CreateCartItem", validReq).Return(nil, errs.ErrProductDoesNotExist)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.PRODUCT_NOT_EXISTS,
					Message: errs.ErrProductDoesNotExist.Error(),
				},
				statusCode: http.StatusNotFound,
			},
		},
		{
			description: "response status not found when sku not found",
			input: input{
				data: validReq,
				beforeTests: func(mockUserCartItemService *mocks.UserCartItemService) {
					mockUserCartItemService.On("CreateCartItem", validReq).Return(nil, errs.ErrProductDoesNotExist)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.PRODUCT_NOT_EXISTS,
					Message: errs.ErrProductDoesNotExist.Error(),
				},
				statusCode: http.StatusNotFound,
			},
		},
		{
			description: "response status conflict when quantity not enough",
			input: input{
				data: validReq,
				beforeTests: func(mockUserCartItemService *mocks.UserCartItemService) {
					mockUserCartItemService.On("CreateCartItem", validReq).Return(nil, errs.ErrProductQuantityNotEnough)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.QUANTITY_NOT_ENOUGH,
					Message: errs.ErrProductQuantityNotEnough.Error(),
				},
				statusCode: http.StatusConflict,
			},
		},
		{
			description: "response status forbidden when shop owner is user",
			input: input{
				data: validReq,
				beforeTests: func(mockUserCartItemService *mocks.UserCartItemService) {
					mockUserCartItemService.On("CreateCartItem", validReq).Return(nil, errs.ErrUserIsShopOwner)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.FORBIDDEN,
					Message: errs.ErrUserIsShopOwner.Error(),
				},
				statusCode: http.StatusForbidden,
			},
		},
		{
			description: "response status internal server error when create cart item failed",
			input: input{
				data: validReq,
				beforeTests: func(mockUserCartItemService *mocks.UserCartItemService) {
					mockUserCartItemService.On("CreateCartItem", validReq).Return(nil, errs.ErrInternalServerError)
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
			description: "response status created when create cart item success",
			input: input{
				data: validReq,
				beforeTests: func(mockUserCartItemService *mocks.UserCartItemService) {
					mockUserCartItemService.On("CreateCartItem", validReq).Return(&model.CartItem{
						Quantity: 1,
						Notes:    "test",
						UserId:   1,
						SkuId:    1,
					}, nil)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.CREATED,
					Message: "create cart item succesful",
					Data: &model.CartItem{
						Quantity: 1,
						Notes:    "test",
						UserId:   1,
						SkuId:    1,
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
			c.Set("userId", tc.input.data.UserId)

			payload := test.MakeRequestBody(tc.input.data)
			c.Request, _ = http.NewRequest(http.MethodGet, "/users/carts", payload)

			mocCartItemService := new(mocks.UserCartItemService)
			tc.beforeTests(mocCartItemService)

			handler := handler.New(&handler.HandlerConfig{
				UserCartItemService: mocCartItemService,
			})

			handler.CreateCartItem(c)

			expectedJson, _ := json.Marshal(tc.expected.data)
			assert.Equal(t, expectedJson, rec.Body.Bytes())
			assert.Equal(t, tc.expected.statusCode, rec.Code)
		})
	}
}

func TestGetAllCartItem(t *testing.T) {
	type input struct {
		data        *dto.GetCartItemsRequest
		beforeTests func(mockUserCartItemService *mocks.UserCartItemService)
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
			description: "response status internal server error when get cart item failed",
			input: input{
				data: &dto.GetCartItemsRequest{
					UserId: 1,
				},
				beforeTests: func(mockUserCartItemService *mocks.UserCartItemService) {
					mockUserCartItemService.On("GetAllCartItem", mock.AnythingOfType("*dto.GetCartItemsRequest")).Return(nil, errs.ErrInternalServerError)
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
			description: "response status ok when get cart item success",
			input: input{
				data: &dto.GetCartItemsRequest{
					UserId: 1,
				},
				beforeTests: func(mockUserCartItemService *mocks.UserCartItemService) {
					mockUserCartItemService.On("GetAllCartItem", mock.AnythingOfType("*dto.GetCartItemsRequest")).Return(
						&commonDto.PaginationResponse{}, nil)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.OK,
					Message: "get all cart item successful",
					Data:    &commonDto.PaginationResponse{},
				},
				statusCode: http.StatusOK,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.data.UserId)

			payload := test.MakeRequestBody(tc.input.data)
			c.Request, _ = http.NewRequest(http.MethodGet, "/users/carts", payload)

			mocCartItemService := new(mocks.UserCartItemService)
			tc.beforeTests(mocCartItemService)

			handler := handler.New(&handler.HandlerConfig{
				UserCartItemService: mocCartItemService,
			})

			handler.GetAllCartItem(c)

			expectedJson, _ := json.Marshal(tc.expected.data)
			assert.Equal(t, expectedJson, rec.Body.Bytes())
			assert.Equal(t, tc.expected.statusCode, rec.Code)
		})
	}

}

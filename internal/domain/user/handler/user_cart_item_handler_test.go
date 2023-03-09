package handler_test

import (
	"encoding/json"
	"errors"
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
	"strconv"
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
			description: "response status conflict when cart item quantity exceed limit",
			input: input{
				data: validReq,
				beforeTests: func(mockUserCartItemService *mocks.UserCartItemService) {
					mockUserCartItemService.On("CreateCartItem", validReq).Return(nil, errs.ErrCartItemLimitExceeded)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.CART_ITEM_EXCEED_LIMIT,
					Message: errs.ErrCartItemLimitExceeded.Error(),
				},
				statusCode: http.StatusConflict,
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

func TestUpdateCartItem(t *testing.T) {
	type input struct {
		userID     int
		skuID      string
		request    *dto.UpdateCartItemRequest
		beforeTest func(*mocks.UserCartItemService)
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
			description: "should return error with status code 400 when sku ID is invalid",
			input: input{
				userID: 1,
				skuID:  "-1",
				request: &dto.UpdateCartItemRequest{
					SkuID: 2,
				},
				beforeTest: func(ucis *mocks.UserCartItemService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "sku ID must be a number and greater than or equal 1",
				},
			},
		},
		{
			description: "should return error with status code 400 when request body is empty",
			input: input{
				userID:     1,
				skuID:      "2",
				request:    &dto.UpdateCartItemRequest{},
				beforeTest: func(ucis *mocks.UserCartItemService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Notes is required",
				},
			},
		},
		{
			description: "should return error with status code 400 when given invalid request body",
			input: input{
				userID: 1,
				skuID:  "2",
				request: &dto.UpdateCartItemRequest{
					Notes: "a veryveryveryveryveryveryveryveryveryveryveryveryveryveryveryveryveryveryvery long notes",
				},
				beforeTest: func(ucis *mocks.UserCartItemService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Notes must be shorter than 50",
				},
			},
		},
		{
			description: "should return error with status code 404 when product does not exist or inactive",
			input: input{
				userID: 1,
				skuID:  "2",
				request: &dto.UpdateCartItemRequest{
					Quantity: 3,
					Notes:    "test",
				},
				beforeTest: func(ucis *mocks.UserCartItemService) {
					ucis.On("UpdateCartItem", 1, &dto.UpdateCartItemRequest{
						SkuID:    2,
						Quantity: 3,
						Notes:    "test",
					}).Return(nil, errs.ErrProductDoesNotExist)
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.PRODUCT_NOT_EXISTS,
					Message: errs.ErrProductDoesNotExist.Error(),
				},
			},
		},
		{
			description: "should return error with status code 404 when cart item does not exist",
			input: input{
				userID: 1,
				skuID:  "2",
				request: &dto.UpdateCartItemRequest{
					Quantity: 3,
					Notes:    "test",
				},
				beforeTest: func(ucis *mocks.UserCartItemService) {
					ucis.On("UpdateCartItem", 1, &dto.UpdateCartItemRequest{
						SkuID:    2,
						Quantity: 3,
						Notes:    "test",
					}).Return(nil, errs.ErrCartItemNotFound)
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.CART_ITEM_NOT_FOUND,
					Message: errs.ErrCartItemNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with status code 409 when product quantity is not enough",
			input: input{
				userID: 1,
				skuID:  "2",
				request: &dto.UpdateCartItemRequest{
					Quantity: 3,
					Notes:    "test",
				},
				beforeTest: func(ucis *mocks.UserCartItemService) {
					ucis.On("UpdateCartItem", 1, &dto.UpdateCartItemRequest{
						SkuID:    2,
						Quantity: 3,
						Notes:    "test",
					}).Return(nil, errs.ErrProductQuantityNotEnough)
				},
			},
			expected: expected{
				statusCode: http.StatusConflict,
				response: response.Response{
					Code:    code.QUANTITY_NOT_ENOUGH,
					Message: errs.ErrProductQuantityNotEnough.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to update cart item",
			input: input{
				userID: 1,
				skuID:  "2",
				request: &dto.UpdateCartItemRequest{
					Quantity: 3,
					Notes:    "test",
				},
				beforeTest: func(ucis *mocks.UserCartItemService) {
					ucis.On("UpdateCartItem", 1, &dto.UpdateCartItemRequest{
						SkuID:    2,
						Quantity: 3,
						Notes:    "test",
					}).Return(nil, errors.New("failed to update cart item"))
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
		{
			description: "should return error with status code 200 when update cart item succeed",
			input: input{
				userID: 1,
				skuID:  "2",
				request: &dto.UpdateCartItemRequest{
					SkuID:    2,
					Quantity: 3,
					Notes:    "test",
				},
				beforeTest: func(ucis *mocks.UserCartItemService) {
					ucis.On("UpdateCartItem", 1, &dto.UpdateCartItemRequest{
						SkuID:    2,
						Quantity: 3,
						Notes:    "test",
					}).Return(&dto.UpdateCartItemResponse{SkuID: 2, Quantity: 3, Notes: "test"}, nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.UPDATED,
					Message: "update cart item succesful",
					Data:    &dto.UpdateCartItemResponse{SkuID: 2, Quantity: 3, Notes: "test"},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			cartItemService := mocks.NewUserCartItemService(t)
			tc.beforeTest(cartItemService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userID)
			c.AddParam("skuId", tc.input.skuID)
			payload := test.MakeRequestBody(tc.input.request)
			c.Request, _ = http.NewRequest(http.MethodPut, "v1/users/carts", payload)
			handler := handler.New(&handler.HandlerConfig{
				UserCartItemService: cartItemService,
			})

			handler.UpdateCartItem(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestDeleteCartItem(t *testing.T) {
	type input struct {
		req        *dto.DeleteCartItemRequest
		beforeTest func(ucis *mocks.UserCartItemService)
	}
	type expected struct {
		statusCode int
		response   response.Response
	}
	tests := []struct {
		description string
		input       input
		expected    expected
	}{
		{
			description: "should return error with status code 400 when sku id is invalid",
			input: input{
				req: &dto.DeleteCartItemRequest{
					UserId:     1,
					CartItemId: 1,
				},
				beforeTest: func(ucis *mocks.UserCartItemService) {
					ucis.On("DeleteCartItem", &dto.DeleteCartItemRequest{
						UserId:     1,
						CartItemId: 1,
					}).Return(errs.ErrCartItemNotFound)
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.CART_ITEM_NOT_FOUND,
					Message: errs.ErrCartItemNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to delete cart item",
			input: input{
				req: &dto.DeleteCartItemRequest{
					UserId:     1,
					CartItemId: 1,
				},
				beforeTest: func(ucis *mocks.UserCartItemService) {
					ucis.On("DeleteCartItem", &dto.DeleteCartItemRequest{
						UserId:     1,
						CartItemId: 1,
					}).Return(errs.ErrInternalServerError)
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
		{
			description: "should return error with status code 200 when delete cart item succeed",
			input: input{
				req: &dto.DeleteCartItemRequest{
					UserId:     1,
					CartItemId: 1,
				},
				beforeTest: func(ucis *mocks.UserCartItemService) {
					ucis.On("DeleteCartItem", &dto.DeleteCartItemRequest{
						UserId:     1,
						CartItemId: 1,
					}).Return(nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.DELETED,
					Message: "delete cart item succesful",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			cartItemService := mocks.NewUserCartItemService(t)
			tc.input.beforeTest(cartItemService)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.req.UserId)
			c.AddParam("cartItemId", strconv.Itoa(tc.input.req.CartItemId))
			c.Request, _ = http.NewRequest(http.MethodDelete, "v1/users/carts", nil)
			handler := handler.New(&handler.HandlerConfig{
				UserCartItemService: cartItemService,
			})

			handler.DeleteCartItem(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

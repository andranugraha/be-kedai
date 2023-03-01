package handler_test

import (
	"fmt"
	"kedai/backend/be-kedai/internal/common/code"
	"kedai/backend/be-kedai/internal/common/constant"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
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
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

func TestUserWishlist_GetUserWishlist(t *testing.T) {
	var (
		invalidProductId = 0
		validProductId   = 1
		validProductCode = "product-code"

		product = &productModel.Product{
			ID:   validProductId,
			Code: validProductCode,
		}
		user = &model.User{
			ID: 1,
		}
		wishlist = &model.UserWishlist{
			UserID:    user.ID,
			ProductID: product.ID,
		}
	)
	type input struct {
		data        *dto.UserWishlistRequest
		beforeTests func(mockUserWishlistService *mocks.UserWishlistService)
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
			description: "it should return error product id required and status code bad request if product id is empty",
			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    user.ID,
					ProductId: invalidProductId,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("GetUserWishlist", &dto.UserWishlistRequest{
						UserId:    user.ID,
						ProductId: invalidProductId,
					}).Return(nil, errs.ErrProductIdRequired)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.BAD_REQUEST,
					Message: errs.ErrProductIdRequired.Error(),
				},
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description: "it should return error user not exist and status code not found if user is not registered",
			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    user.ID,
					ProductId: product.ID,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("GetUserWishlist", &dto.UserWishlistRequest{
						UserId:    user.ID,
						ProductId: product.ID,
					}).Return(nil, errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.USER_NOT_REGISTERED,
					Message: errs.ErrUserDoesNotExist.Error(),
				},
				statusCode: http.StatusNotFound,
			},
		},
		{
			description: "it should return error product not exist and status code not found if product is not found",
			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    user.ID,
					ProductId: product.ID,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("GetUserWishlist", &dto.UserWishlistRequest{
						UserId:    user.ID,
						ProductId: product.ID,
					}).Return(nil, errs.ErrProductDoesNotExist)
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
			description: "it should return error product not exist in wishlist and status code not found if product is not found in wishlist",

			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    user.ID,
					ProductId: product.ID,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("GetUserWishlist", &dto.UserWishlistRequest{
						UserId:    user.ID,
						ProductId: product.ID,
					}).Return(nil, errs.ErrProductNotInWishlist)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.PRODUCT_NOT_IN_WISHLIST,
					Message: errs.ErrProductNotInWishlist.Error(),
				},
				statusCode: http.StatusNotFound,
			},
		},
		{
			description: "it should return error internal server and status code internal server error if error is not expected",

			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    user.ID,
					ProductId: product.ID,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("GetUserWishlist", &dto.UserWishlistRequest{
						UserId:    user.ID,
						ProductId: product.ID,
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
			description: "it should return success message and status code ok if product is found in wishlist",

			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    user.ID,
					ProductId: product.ID,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("GetUserWishlist", &dto.UserWishlistRequest{
						UserId:    user.ID,
						ProductId: product.ID,
					}).Return(wishlist, nil)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.OK,
					Message: "wishlist retrieved successfully",
					Data:    wishlist,
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
			productId := strconv.Itoa(tc.input.data.ProductId)
			c.Params = gin.Params{
				{
					Key:   "productId",
					Value: productId,
				},
			}

			c.Request, _ = http.NewRequest(http.MethodGet, "/users/wishlists", nil)

			mockWishlistService := new(mocks.UserWishlistService)
			tc.beforeTests(mockWishlistService)

			handler := handler.New(&handler.HandlerConfig{
				UserWishlistService: mockWishlistService,
			})

			handler.GetUserWishlist(c)

			expectedJson, _ := json.Marshal(tc.expected.data)
			assert.Equal(t, expectedJson, rec.Body.Bytes())
			assert.Equal(t, tc.expected.statusCode, rec.Code)
		})
	}
}

func TestUserWishlist_AddUserWishlist(t *testing.T) {
	var (
		userId         = 1
		productId      = 1
		invalidRequest = &dto.UserWishlistRequest{}
		validRequest   = &dto.UserWishlistRequest{
			ProductId: productId,
			UserId:    userId,
		}
		wishlist = &model.UserWishlist{
			ProductID: 1,
			UserID:    validRequest.UserId,
		}
	)
	type input struct {
		data        *dto.UserWishlistRequest
		beforeTests func(mockUserWishlistService *mocks.UserWishlistService)
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
			description: "it should return error product code required and bad request if product code is empty",
			input: input{
				data: invalidRequest,
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("AddUserWishlist", invalidRequest).Return(nil, fmt.Errorf("%s is required", "ProductId"))
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.BAD_REQUEST,
					Message: "ProductId is required",
				},
				statusCode: http.StatusBadRequest,
			},
		},

		{
			description: "it should return error user not exist and status code not found if user is not registered",
			input: input{
				data: validRequest,
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("AddUserWishlist", validRequest).Return(nil, errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.USER_NOT_REGISTERED,
					Message: errs.ErrUserDoesNotExist.Error(),
				},
				statusCode: http.StatusNotFound,
			},
		},

		{
			description: "it should return error product not exist and status code not found if product is not found",
			input: input{
				data: validRequest,
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("AddUserWishlist", validRequest).Return(nil, errs.ErrProductDoesNotExist)
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
			description: "it should return error product in wishlist and status code conflict if product already in wishlist",
			input: input{
				data: validRequest,
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("AddUserWishlist", validRequest).Return(nil, errs.ErrProductInWishlist)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.PRODUCT_ALREADY_IN_WISHLIST,
					Message: errs.ErrProductInWishlist.Error(),
				},
				statusCode: http.StatusConflict,
			},
		},

		{
			description: "it should return error internal server",
			input: input{
				data: validRequest,
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("AddUserWishlist", validRequest).Return(nil, errs.ErrInternalServerError)
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
			description: "it should return user wishlist and status code created if success",
			input: input{
				data: validRequest,
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("AddUserWishlist", validRequest).Return(wishlist, nil)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.CREATED,
					Message: "wishlist success created successfully",
					Data:    wishlist,
				},
				statusCode: http.StatusCreated,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", userId)

			mockWishlistService := new(mocks.UserWishlistService)
			tc.beforeTests(mockWishlistService)

			handler := handler.New(&handler.HandlerConfig{
				UserWishlistService: mockWishlistService,
			})

			payload := test.MakeRequestBody(tc.input.data)
			c.Request, _ = http.NewRequest(http.MethodPost, "/users/wishlists", payload)

			handler.AddUserWishlist(c)

			expectedJson, _ := json.Marshal(tc.expected.data)
			assert.Equal(t, expectedJson, rec.Body.Bytes())
			assert.Equal(t, tc.expected.statusCode, rec.Code)
		})
	}
}

func TestUserWishlist_RemoveUserWishlist(t *testing.T) {
	var (
		userId           = 1
		invalidProductId = 0
		validProductId   = 1
	)
	type input struct {
		data        *dto.UserWishlistRequest
		beforeTests func(mockUserWishlistService *mocks.UserWishlistService)
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
			description: "it should return error product id required and status code bad request if product id is empty",
			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    userId,
					ProductId: invalidProductId,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("RemoveUserWishlist", &dto.UserWishlistRequest{
						UserId:    userId,
						ProductId: invalidProductId,
					}).Return(errs.ErrProductIdRequired)
				},
			},

			expected: expected{
				data: &response.Response{
					Code:    code.BAD_REQUEST,
					Message: "product id is required",
				},
				statusCode: http.StatusBadRequest,
			},
		},
		{
			description: "it should return error user not exist and status code not found if user is not registered",
			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    userId,
					ProductId: validProductId,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("RemoveUserWishlist", &dto.UserWishlistRequest{
						UserId:    userId,
						ProductId: validProductId,
					}).Return(errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.USER_NOT_REGISTERED,
					Message: errs.ErrUserDoesNotExist.Error(),
				},
				statusCode: http.StatusNotFound,
			},
		},
		{
			description: "it should return error product not exist and status code not found if product is not registered",
			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    userId,
					ProductId: validProductId,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("RemoveUserWishlist", &dto.UserWishlistRequest{
						UserId:    userId,
						ProductId: validProductId,
					}).Return(errs.ErrProductDoesNotExist)
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
			description: "it should return error product not exist in wishlist and status code not found if product is not found in wishlist",
			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    userId,
					ProductId: validProductId,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("RemoveUserWishlist", &dto.UserWishlistRequest{
						UserId:    userId,
						ProductId: validProductId,
					}).Return(errs.ErrProductNotInWishlist)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.PRODUCT_NOT_IN_WISHLIST,
					Message: errs.ErrProductNotInWishlist.Error(),
				},
				statusCode: http.StatusNotFound,
			},
		},
		{
			description: "it should return error internal server and status code internal server error if error is not expected",
			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    userId,
					ProductId: validProductId,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("RemoveUserWishlist", &dto.UserWishlistRequest{
						UserId:    userId,
						ProductId: validProductId,
					}).Return(errs.ErrInternalServerError)
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
			description: "it should return success and status code ok if product is removed from wishlist",
			input: input{
				data: &dto.UserWishlistRequest{
					UserId:    userId,
					ProductId: validProductId,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("RemoveUserWishlist", &dto.UserWishlistRequest{
						UserId:    userId,
						ProductId: validProductId,
					}).Return(nil)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.OK,
					Message: "wishlist removed successfully",
				},
				statusCode: http.StatusOK,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", userId)
			productId := strconv.Itoa(tc.input.data.ProductId)
			c.Params = gin.Params{
				{
					Key:   "productId",
					Value: productId,
				},
			}

			mockWishlistService := new(mocks.UserWishlistService)
			tc.beforeTests(mockWishlistService)

			handler := handler.New(&handler.HandlerConfig{
				UserWishlistService: mockWishlistService,
			})

			c.Request, _ = http.NewRequest(http.MethodDelete, "/users/wishlists", nil)

			handler.RemoveUserWishlist(c)

			expectedJson, _ := json.Marshal(tc.expected.data)
			assert.Equal(t, expectedJson, rec.Body.Bytes())
			assert.Equal(t, tc.expected.statusCode, rec.Code)
		})
	}
}

func TestGetUserWishlists(t *testing.T) {
	var (
		userId = 1
		req    = dto.GetUserWishlistsRequest{
			UserId: userId,
			Sort:   constant.SortByRecommended,
			Limit:  10,
			Page:   1,
		}
		res = &commonDto.PaginationResponse{
			Data: []*model.UserWishlist{
				{
					UserID:    userId,
					ProductID: 1,
				},
			},
			TotalRows:  1,
			TotalPages: 1,
			Limit:      req.Limit,
			Page:       req.Page,
		}
	)

	tests := []struct {
		name                string
		getUserWishlistsRes *commonDto.PaginationResponse
		getUserWishlistsErr error
		want                response.Response
		wantStatusCode      int
	}{
		{
			name:                "should return 200 when get user wishlists success",
			getUserWishlistsRes: res,
			getUserWishlistsErr: nil,
			want: response.Response{
				Code:    code.OK,
				Message: "wishlist retrieved successfully",
				Data:    res,
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:                "should return 500 when get user wishlists failed",
			getUserWishlistsRes: nil,
			getUserWishlistsErr: errs.ErrInternalServerError,
			want: response.Response{
				Code:    code.INTERNAL_SERVER_ERROR,
				Message: errs.ErrInternalServerError.Error(),
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockWishlistService := new(mocks.UserWishlistService)
			mockWishlistService.On("GetUserWishlists", req).Return(test.getUserWishlistsRes, test.getUserWishlistsErr)

			handler := handler.New(&handler.HandlerConfig{
				UserWishlistService: mockWishlistService,
			})

			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", userId)
			c.Request, _ = http.NewRequest(http.MethodGet, "/users/wishlists", nil)

			handler.GetUserWishlists(c)

			expectedJson, _ := json.Marshal(test.want)
			assert.Equal(t, expectedJson, rec.Body.Bytes())
			assert.Equal(t, test.wantStatusCode, rec.Code)
		})
	}
}

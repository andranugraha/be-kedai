package handler_test

import (
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/handler"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/utils/response"
	testutils "kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

func TestUserWishlist_AddUserWishlist(t *testing.T) {
	var (
		userId         = 1
		invalidRequest = &dto.UserWishlistRequest{}
		validRequest   = &dto.UserWishlistRequest{
			ProductCode: "123",
			UserID:      userId,
		}
		wishlist = &model.UserWishlist{
			ProductID: 1,
			UserID:    validRequest.UserID,
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
					mockWishlistService.On("AddUserWishlist", invalidRequest).Return(nil, errs.ErrProductCodeRequired)
				},
			},
			expected: expected{
				data: &response.Response{
					Code:    code.PRODUCT_CODE_IS_REQUIRED,
					Message: errs.ErrProductCodeRequired.Error(),
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

			handler := handler.NewHandler(&handler.HandlerConfig{
				UserWishlistService: mockWishlistService,
			})

			payload := testutils.MakeRequestBody(tc.input.data)
			c.Request, _ = http.NewRequest(http.MethodPost, "/users/wishlists", payload)

			handler.AddUserWishlist(c)

			expectedJson, _ := json.Marshal(tc.expected.data)
			assert.Equal(t, expectedJson, rec.Body.Bytes())
			assert.Equal(t, tc.expected.statusCode, rec.Code)

		})
	}
}

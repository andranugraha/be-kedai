package handler_test

import (
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/handler"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

func TestUserWishlist_RemoveUserWishlist(t *testing.T) {
	var (
		userId             = 1
		invalidProductCode = ""
		validProductCode   = "ITEM-001"
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
			description: "it should return error product code required and status code bad request if product code is empty",
			input: input{
				data: &dto.UserWishlistRequest{
					UserID:      userId,
					ProductCode: invalidProductCode,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("RemoveUserWishlist", &dto.UserWishlistRequest{
						UserID:      userId,
						ProductCode: invalidProductCode,
					}).Return(errs.ErrProductCodeRequired)
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
				data: &dto.UserWishlistRequest{
					UserID:      userId,
					ProductCode: validProductCode,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("RemoveUserWishlist", &dto.UserWishlistRequest{
						UserID:      userId,
						ProductCode: validProductCode,
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
					UserID:      userId,
					ProductCode: validProductCode,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("RemoveUserWishlist", &dto.UserWishlistRequest{
						UserID:      userId,
						ProductCode: validProductCode,
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
					UserID:      userId,
					ProductCode: validProductCode,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("RemoveUserWishlist", &dto.UserWishlistRequest{
						UserID:      userId,
						ProductCode: validProductCode,
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
					UserID:      userId,
					ProductCode: validProductCode,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("RemoveUserWishlist", &dto.UserWishlistRequest{
						UserID:      userId,
						ProductCode: validProductCode,
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
					UserID:      userId,
					ProductCode: validProductCode,
				},
				beforeTests: func(mockWishlistService *mocks.UserWishlistService) {
					mockWishlistService.On("RemoveUserWishlist", &dto.UserWishlistRequest{
						UserID:      userId,
						ProductCode: validProductCode,
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
			c.Params = gin.Params{
				{
					Key:   "productCode",
					Value: tc.input.data.ProductCode,
				},
			}

			mockWishlistService := new(mocks.UserWishlistService)
			tc.beforeTests(mockWishlistService)

			handler := handler.NewHandler(&handler.HandlerConfig{
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

package handler_test

import (
	"encoding/json"
	"errors"
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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetShopRating(t *testing.T) {
	var (
		userId           = 1
		filter           = &dto.GetShopRatingFilterRequest{

		}
		shopRating = dto.GetShopRatingResponse{
			ShopRating: 1,
			Data: &commonDto.PaginationResponse{
				Data:       []dto.ProductItem{},
				TotalRows:  0,
				TotalPages: 0,
				Limit:      0,
				Page:       0,
			},
		}
	)

	type input struct {
		userId int
		filter *dto.GetShopRatingFilterRequest
		result *dto.GetShopRatingResponse
		err    error
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
			description: "should return shop rating information with code 200 when success",
			input: input{
				userId: userId,
				filter: filter,
				result: &shopRating,
				err:    nil},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    shopRating,
				},
			},
		},
		{
			description: "should return error with code 500 when error",
			input: input{
				userId: userId,
				filter: filter,
				result: nil,
				err:    errors.New("error"),
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
			mockService := new(mocks.ShopService)
			mockService.On("GetShopRating", tc.input.userId,*tc.input.filter).Return(tc.input.result, tc.input.err)
			handler := handler.New(&handler.HandlerConfig{
				ShopService: mockService,
			})
			c.Set("userId", 1)

			c.Request, _ = http.NewRequest("GET", "/sellers/ratings", nil)

			handler.GetShopRating(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}

}

func TestFindShopBySlug(t *testing.T) {
	var (
		slug       = "shop"
		shopResult = &model.Shop{
			ID: 1,
		}
	)

	type input struct {
		shop *model.Shop
		err  error
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
			description: "should return shop information with code 200 when success",
			input: input{
				shop: shopResult,
				err:  nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    shopResult,
				},
			},
		},
		{
			description: "should return error with code 404 when shop not found",
			input: input{
				shop: nil,
				err:  errs.ErrShopNotFound,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.NOT_FOUND,
					Message: "shop not found",
				},
			},
		},
		{
			description: "should return error with code 500 when internal server error",
			input: input{
				shop: nil,
				err:  errs.ErrInternalServerError,
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: "something went wrong in the server",
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
			mockService := new(mocks.ShopService)
			mockService.On("FindShopBySlug", slug).Return(tc.input.shop, tc.input.err)
			handler := handler.New(&handler.HandlerConfig{
				ShopService: mockService,
			})
			c.Request, _ = http.NewRequest("GET", "/shops/:slug", nil)

			handler.FindShopBySlug(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestFindShopByKeyword(t *testing.T) {
	var (
		shopList   = []*model.Shop{}
		pagination = &commonDto.PaginationResponse{
			Data:  shopList,
			Limit: 10,
			Page:  1,
		}
		req = dto.FindShopRequest{
			Keyword: "test",
			Page:    1,
			Limit:   10,
		}
	)
	type input struct {
		dto    dto.FindShopRequest
		result *commonDto.PaginationResponse
		err    error
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
			description: "should return shop list with code 200 when success",
			input: input{
				dto:    req,
				result: pagination,
				err:    nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    pagination,
				},
			},
		},
		{
			description: "should return error with code 500 when internal server error",
			input: input{
				dto:    req,
				result: nil,
				err:    errs.ErrInternalServerError,
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
			mockService := new(mocks.ShopService)
			mockService.On("FindShopByKeyword", tc.input.dto).Return(tc.input.result, tc.input.err)
			handler := handler.New(&handler.HandlerConfig{
				ShopService: mockService,
			})
			c.Request, _ = http.NewRequest("GET", "/shops?keyword=test", nil)

			handler.FindShopByKeyword(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestGetShopFinanceOverview(t *testing.T) {
	var (
		userId   = 1
		overview = &dto.ShopFinanceOverviewResponse{}
	)

	type input struct {
		shopID int
		result *dto.ShopFinanceOverviewResponse
		err    error
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
			description: "should return shop finance overview with code 200 when success",
			input: input{
				shopID: userId,
				result: overview,
				err:    nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    overview,
				},
			},
		},
		{
			description: "should return error with code 404 when shop not found",
			input: input{
				shopID: userId,
				result: nil,
				err:    errs.ErrShopNotFound,
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
				shopID: userId,
				result: nil,
				err:    errs.ErrInternalServerError,
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
			mockService := new(mocks.ShopService)
			mockService.On("GetShopFinanceOverview", tc.input.shopID).Return(tc.input.result, tc.input.err)
			handler := handler.New(&handler.HandlerConfig{
				ShopService: mockService,
			})
			c.Set("userId", 1)

			c.Request, _ = http.NewRequest("GET", "/shops/1/finance/overview", nil)

			handler.GetShopFinanceOverview(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestGetShopStats(t *testing.T) {
	var (
		userId = 1
		stats  = &dto.GetShopStatsResponse{}
	)

	type input struct {
		userId int
		result *dto.GetShopStatsResponse
		err    error
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
			description: "should return shop stats with code 200 when success",
			input: input{
				userId: userId,
				result: stats,
				err:    nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    stats,
				},
			},
		},
		{
			description: "should return error with code 404 when shop not found",
			input: input{
				userId: userId,
				result: nil,
				err:    errs.ErrShopNotFound,
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
				userId: userId,
				result: nil,
				err:    errs.ErrInternalServerError,
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
			mockService := new(mocks.ShopService)
			mockService.On("GetShopStats", tc.input.userId).Return(tc.input.result, tc.input.err)
			handler := handler.New(&handler.HandlerConfig{
				ShopService: mockService,
			})
			c.Set("userId", 1)

			c.Request, _ = http.NewRequest("GET", "/v1/sellers/stats", nil)
			handler.GetShopStats(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestAddShopGuest(t *testing.T) {

	type input struct {
		req        *dto.AddShopGuestRequest
		err        error
		beforeTest func(mockShopGuestService *mocks.ShopGuestService)
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
			description: "should return bad request with code 400 when request is invalid",
			input: input{
				beforeTest: func(mockShopGuestService *mocks.ShopGuestService) {
				},
				err: errs.ErrBadRequest,
				req: &dto.AddShopGuestRequest{},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "ShopId is required",
				},
			},
		},
		{
			description: "should return not found with code 404 when shop not found",
			input: input{
				beforeTest: func(mockShopGuestService *mocks.ShopGuestService) {
					mockShopGuestService.On("CreateShopGuest", 1).Return(nil, errs.ErrShopNotFound)
				},
				err: errs.ErrShopNotFound,
				req: &dto.AddShopGuestRequest{
					ShopId: 1,
				},
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
			description: "should return internal server error with code 500 when internal server error",
			input: input{
				beforeTest: func(mockShopGuestService *mocks.ShopGuestService) {
					mockShopGuestService.On("CreateShopGuest", 1).Return(nil, errs.ErrInternalServerError)
				},
				err: errs.ErrInternalServerError,
				req: &dto.AddShopGuestRequest{
					ShopId: 1,
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
			description: "should return success with code 200 when success",
			input: input{
				beforeTest: func(mockShopGuestService *mocks.ShopGuestService) {
					mockShopGuestService.On("CreateShopGuest", 1).Return(&model.ShopGuest{}, nil)
				},
				err: nil,
				req: &dto.AddShopGuestRequest{
					ShopId: 1,
				},
			},
			expected: expected{
				statusCode: http.StatusCreated,
				response: response.Response{
					Code:    code.CREATED,
					Message: "created",
					Data:    &model.ShopGuest{},
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			expectedBody, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			mockService := new(mocks.ShopGuestService)
			tc.input.beforeTest(mockService)

			handler := handler.New(&handler.HandlerConfig{
				ShopGuestService: mockService,
			})

			payload := test.MakeRequestBody(tc.input.req)
			c.Request, _ = http.NewRequest("POST", "/v1/shops/visitors", payload)

			handler.AddShopGuest(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestGetShopInsight(t *testing.T) {
	var (
		userId = 1
		req    = dto.GetShopInsightRequest{
			Timeframe: dto.ShopInsightTimeframeDay,
			UserId:    userId,
		}
		insight = &dto.GetShopInsightResponse{}
	)

	type input struct {
		req    dto.GetShopInsightRequest
		result *dto.GetShopInsightResponse
		err    error
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
			description: "should return shop insight with code 200 when success",
			input: input{
				req:    req,
				result: insight,
				err:    nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    insight,
				},
			},
		},
		{
			description: "should return error with code 404 when shop not found",
			input: input{
				req:    req,
				result: nil,
				err:    errs.ErrShopNotFound,
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
				req:    req,
				result: nil,
				err:    errs.ErrInternalServerError,
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
			mockService := new(mocks.ShopService)
			mockService.On("GetShopInsight", tc.input.req).Return(tc.input.result, tc.input.err)
			handler := handler.New(&handler.HandlerConfig{
				ShopService: mockService,
			})
			c.Set("userId", 1)

			c.Request, _ = http.NewRequest("GET", "/v1/sellers/insights", nil)
			handler.GetShopInsights(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestGetShopProfile(t *testing.T) {
	var (
		userId = 1
		shop   = &dto.ShopProfile{
			Name: "shop name",
		}
	)

	type input struct {
		userId int
		result *dto.ShopProfile
		err    error
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
			description: "should return shop profile with code 200 when success",
			input: input{
				userId: userId,
				result: shop,
				err:    nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    shop,
				},
			},
		},
		{
			description: "should return error with code 404 when shop not found",
			input: input{
				userId: userId,
				result: nil,
				err:    errs.ErrShopNotFound,
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
				userId: userId,
				result: nil,
				err:    errs.ErrInternalServerError,
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
			mockService := new(mocks.ShopService)
			mockService.On("GetShopProfile", tc.input.userId).Return(tc.input.result, tc.input.err)
			handler := handler.New(&handler.HandlerConfig{
				ShopService: mockService,
			})
			c.Set("userId", 1)

			c.Request, _ = http.NewRequest("GET", "/v1/shops/profile", nil)
			handler.GetShopProfile(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestUpdateShopProfile(t *testing.T) {
	var (
		userId      = 1
		logoUrl     = "https://logo.com"
		bannerUrl   = "https://banner.com"
		description = "description"
		req         = dto.ShopProfile{
			Name:        "shop name",
			LogoUrl:     &logoUrl,
			BannerUrl:   &bannerUrl,
			Description: &description,
		}
	)

	type input struct {
		userId     int
		req        dto.ShopProfile
		beforeTest func(*mocks.ShopService)
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
			description: "should return shop profile with code 200 when success",
			input: input{
				userId: userId,
				req:    req,
				beforeTest: func(mockService *mocks.ShopService) {
					mockService.On("UpdateShopProfile", userId, req).Return(nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
				},
			},
		},
		{
			description: "should return error with code 400 when request invalid",
			input: input{
				userId:     userId,
				req:        dto.ShopProfile{},
				beforeTest: func(mockService *mocks.ShopService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Description is required",
				},
			},
		},
		{
			description: "should return error with code 404 when shop not found",
			input: input{
				userId: userId,
				req:    req,
				beforeTest: func(mockService *mocks.ShopService) {
					mockService.On("UpdateShopProfile", userId, req).Return(errs.ErrShopNotFound)
				},
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
				userId: userId,
				req:    req,
				beforeTest: func(mockService *mocks.ShopService) {
					mockService.On("UpdateShopProfile", userId, req).Return(errs.ErrInternalServerError)
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
	} {
		t.Run(tc.description, func(t *testing.T) {
			expectedBody, _ := json.Marshal(tc.expected.response)
			payload := test.MakeRequestBody(tc.input.req)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			mockService := new(mocks.ShopService)
			tc.beforeTest(mockService)
			handler := handler.New(&handler.HandlerConfig{
				ShopService: mockService,
			})
			c.Set("userId", 1)

			c.Request, _ = http.NewRequest("PUT", "/v1/shops/profile", payload)
			handler.UpdateShopProfile(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestCreateShop(t *testing.T) {
	type input struct {
		userID  int
		request *dto.CreateShopRequest
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		userID     = 1
		shopName   = "New shop"
		addressID  = 1
		courierIDs = []int{1, 2}
		request    = &dto.CreateShopRequest{
			Name:       shopName,
			AddressID:  addressID,
			CourierIDs: courierIDs,
		}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ShopService)
		expected
	}{
		{
			description: "should return error with status code 400 when given invalid request body",
			input: input{
				userID: userID,
				request: &dto.CreateShopRequest{
					Name:       "a",
					AddressID:  addressID,
					CourierIDs: courierIDs,
				},
			},
			beforeTest: func(ss *mocks.ShopService) {},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Name must be greater than 5",
				},
			},
		},
		{
			description: "should return error with status code 409 when user already has shop",
			input: input{
				userID: userID,
				request: &dto.CreateShopRequest{
					Name:       shopName,
					AddressID:  addressID,
					CourierIDs: courierIDs,
				},
			},
			beforeTest: func(ss *mocks.ShopService) {
				ss.On("CreateShop", userID, request).Return(nil, errs.ErrUserHasShop)
			},
			expected: expected{
				statusCode: http.StatusConflict,
				response: response.Response{
					Code:    code.HAVE_SHOP,
					Message: errs.ErrUserHasShop.Error(),
				},
			},
		},
		{
			description: "should return error with status code 409 when shop name already taken",
			input: input{
				userID: userID,
				request: &dto.CreateShopRequest{
					Name:       shopName,
					AddressID:  addressID,
					CourierIDs: courierIDs,
				},
			},
			beforeTest: func(ss *mocks.ShopService) {
				ss.On("CreateShop", userID, request).Return(nil, errs.ErrShopRegistered)
			},
			expected: expected{
				statusCode: http.StatusConflict,
				response: response.Response{
					Code:    code.SHOP_REGISTERED,
					Message: errs.ErrShopRegistered.Error(),
				},
			},
		},
		{
			description: "should return error with status code 422 when shop name is invalid",
			input: input{
				userID: userID,
				request: &dto.CreateShopRequest{
					Name:       "invalid_shop_name",
					AddressID:  addressID,
					CourierIDs: courierIDs,
				},
			},
			beforeTest: func(ss *mocks.ShopService) {
				ss.On("CreateShop", userID, &dto.CreateShopRequest{
					Name:       "invalid_shop_name",
					AddressID:  addressID,
					CourierIDs: courierIDs,
				},
				).Return(nil, errs.ErrInvalidShopName)
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_SHOP_NAME,
					Message: errs.ErrInvalidShopName.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to create shop",
			input: input{
				userID: userID,
				request: &dto.CreateShopRequest{
					Name:       shopName,
					AddressID:  addressID,
					CourierIDs: courierIDs,
				},
			},
			beforeTest: func(ss *mocks.ShopService) {
				ss.On("CreateShop", userID, request).Return(nil, errors.New("failed to create shop"))
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
			description: "should return error with status code 201 when succeed to create shop",
			input: input{
				userID: userID,
				request: &dto.CreateShopRequest{
					Name:       shopName,
					AddressID:  addressID,
					CourierIDs: courierIDs,
				},
			},
			beforeTest: func(ss *mocks.ShopService) {
				ss.On("CreateShop", userID, request).Return(&model.Shop{
					Name:      shopName,
					AddressID: addressID,
					UserID:    userID,
				}, nil)
			},
			expected: expected{
				statusCode: http.StatusCreated,
				response: response.Response{
					Code:    code.CREATED,
					Message: "shop created",
					Data: &model.Shop{
						Name:      shopName,
						AddressID: addressID,
						UserID:    userID,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			payload := test.MakeRequestBody(tc.input.request)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			mockService := new(mocks.ShopService)
			tc.beforeTest(mockService)
			handler := handler.New(&handler.HandlerConfig{
				ShopService: mockService,
			})
			c.Set("userId", tc.input.userID)

			c.Request, _ = http.NewRequest("POST", "/v1/sellers/register", payload)
			handler.CreateShop(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

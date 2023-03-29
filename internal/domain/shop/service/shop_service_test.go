package service_test

import (
	"errors"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/service"
	stringUtils "kedai/backend/be-kedai/internal/utils/strings"
	"kedai/backend/be-kedai/mocks"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetShopRating(t *testing.T) {

	var (
		userId     = 1
		shopRating = &dto.ShopRatingResponse{
			ShopRating:  1,
			ProductItem: []*dto.ProductItem{},
		}

		filter = &dto.GetShopRatingFilterRequest{
			Page:  1,
			Limit: 10,
		}
	)

	type input struct {
		userId int
		filter *dto.GetShopRatingFilterRequest
		err    error
	}

	type expected struct {
		shopRatingResponse *dto.ShopRatingResponse
		paginationResponse *dto.GetShopRatingResponse
		totalRows          int64
		totalPage          int
		err                error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return shop rating when success",
			input: input{
				userId: userId,
				filter: filter,
				err:    nil,
			},
			expected: expected{
				shopRatingResponse: shopRating,
				paginationResponse: &dto.GetShopRatingResponse{
					ShopRating: 1,
					Data: &commonDto.PaginationResponse{
						Data:      []*dto.ProductItem{},
						Limit:     10,
						Page:      1,
						TotalRows: 0,
					},
				},
				err: nil,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.ShopRepository)
			mockRepo.On("GetShopRating", tc.input.userId, *tc.input.filter).Return(tc.expected.shopRatingResponse, tc.expected.totalRows, tc.expected.totalPage, tc.expected.err)
			service := service.NewShopService(&service.ShopSConfig{
				ShopRepository: mockRepo,
			})

			result, err := service.GetShopRating(tc.input.userId, *tc.input.filter)

			assert.Equal(t, tc.expected.paginationResponse, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestFindShopById(t *testing.T) {
	var (
		shopId     = 1
		shopResult = &model.Shop{
			ID: 1,
		}
	)
	type input struct {
		id  int
		err error
	}

	type expected struct {
		shop *model.Shop
		err  error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return shop when success",
			input: input{
				id:  shopId,
				err: nil,
			},
			expected: expected{
				shop: shopResult,
				err:  nil,
			},
		},
		{
			description: "should return error when shop not found",
			input: input{
				id:  shopId,
				err: errs.ErrShopNotFound,
			},
			expected: expected{
				shop: nil,
				err:  errs.ErrShopNotFound,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.ShopRepository)
			mockRepo.On("FindShopById", tc.input.id).Return(tc.expected.shop, tc.expected.err)
			service := service.NewShopService(&service.ShopSConfig{
				ShopRepository: mockRepo,
			})

			result, err := service.FindShopById(tc.input.id)

			assert.Equal(t, tc.expected.shop, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestFindShopByUserId(t *testing.T) {
	var (
		shopId     = 1
		shopResult = &model.Shop{
			ID: 1,
		}
	)
	type input struct {
		id  int
		err error
	}

	type expected struct {
		shop *model.Shop
		err  error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return shop when success",
			input: input{
				id:  shopId,
				err: nil,
			},
			expected: expected{
				shop: shopResult,
				err:  nil,
			},
		},
		{
			description: "should return error when shop not found",
			input: input{
				id:  shopId,
				err: errs.ErrShopNotFound,
			},
			expected: expected{
				shop: nil,
				err:  errs.ErrShopNotFound,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.ShopRepository)
			mockRepo.On("FindShopByUserId", tc.input.id).Return(tc.expected.shop, tc.expected.err)
			service := service.NewShopService(&service.ShopSConfig{
				ShopRepository: mockRepo,
			})

			result, err := service.FindShopByUserId(tc.input.id)

			assert.Equal(t, tc.expected.shop, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestFindShopBySlug(t *testing.T) {
	var (
		shopSlug   = "shop"
		shopResult = &model.Shop{
			ID: 1,
		}
	)
	type input struct {
		slug       string
		err        error
		beforeTest func(*mocks.ShopRepository)
	}

	type expected struct {
		shop *model.Shop
		err  error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return shop when success",
			input: input{
				slug: shopSlug,
				err:  nil,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopBySlug", shopSlug).Return(shopResult, nil)
				},
			},
			expected: expected{
				shop: shopResult,
				err:  nil,
			},
		},
		{
			description: "should return error when shop not found",
			input: input{
				slug: shopSlug,
				err:  errs.ErrShopNotFound,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopBySlug", shopSlug).Return(nil, errs.ErrShopNotFound)
				},
			},
			expected: expected{
				shop: nil,
				err:  errs.ErrShopNotFound,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.ShopRepository)
			tc.beforeTest(mockRepo)
			service := service.NewShopService(&service.ShopSConfig{
				ShopRepository: mockRepo,
			})

			result, err := service.FindShopBySlug(tc.input.slug)

			assert.Equal(t, tc.expected.shop, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestFindShopByKeyword(t *testing.T) {
	var (
		shopList   = []*dto.FindShopResponse{}
		rows       = int64(1)
		pages      = 1
		limit      = 10
		pagination = &commonDto.PaginationResponse{
			Data:       shopList,
			TotalRows:  rows,
			TotalPages: pages,
			Page:       pages,
			Limit:      limit,
		}
		emptyPagination = &commonDto.PaginationResponse{
			Data:  []*dto.FindShopResponse{},
			Page:  pages,
			Limit: limit,
		}
		req = dto.FindShopRequest{
			Limit:   limit,
			Page:    pages,
			Keyword: "test",
		}
		invalidReq = dto.FindShopRequest{
			Limit: limit,
			Page:  pages,
		}
	)
	type input struct {
		dto        dto.FindShopRequest
		err        error
		beforeTest func(*mocks.ShopRepository)
	}
	type expected struct {
		result *commonDto.PaginationResponse
		err    error
	}
	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return shop list and pagination when success",
			input: input{
				dto: req,
				err: nil,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByKeyword", req).Return(shopList, rows, pages, nil)
				},
			},
			expected: expected{
				result: pagination,
				err:    nil,
			},
		},
		{
			description: "should return empty shop list when keyword is empty",
			input: input{
				dto:        invalidReq,
				err:        nil,
				beforeTest: func(sr *mocks.ShopRepository) {},
			},
			expected: expected{
				result: emptyPagination,
				err:    nil,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				dto: req,
				err: errs.ErrInternalServerError,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByKeyword", req).Return(nil, int64(0), 0, errors.New("error"))
				},
			},
			expected: expected{
				result: nil,
				err:    errors.New("error"),
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.ShopRepository)
			tc.beforeTest(mockRepo)
			service := service.NewShopService(&service.ShopSConfig{
				ShopRepository: mockRepo,
			})

			result, err := service.FindShopByKeyword(tc.input.dto)

			assert.Equal(t, tc.expected.result, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestGetShopFinanceOverview(t *testing.T) {

	type input struct {
		userId     int
		beforeTest func(*mocks.ShopRepository)
	}

	type expected struct {
		result *dto.ShopFinanceOverviewResponse
		err    error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return error when FindShopByUserId return error",
			input: input{
				userId: 1,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByUserId", 1).Return(nil, errors.New("error"))
				},
			},
			expected: expected{
				result: nil,
				err:    errors.New("error"),
			},
		},
		{
			description: "should return error when GetShopFinanceOverview return error",
			input: input{
				userId: 1,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByUserId", 1).Return(&model.Shop{
						ID: 1,
					}, nil)
					sr.On("GetShopFinanceOverview", 1).Return(nil, errors.New("error"))
				},
			},
			expected: expected{
				result: nil,
				err:    errors.New("error"),
			},
		},
		{
			description: "should return shop finance overview when success",
			input: input{
				userId: 1,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByUserId", 1).Return(&model.Shop{
						ID: 1,
					}, nil)
					sr.On("GetShopFinanceOverview", 1).Return(&dto.ShopFinanceOverviewResponse{}, nil)
				},
			},
			expected: expected{
				result: &dto.ShopFinanceOverviewResponse{},
				err:    nil,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.ShopRepository)
			tc.beforeTest(mockRepo)
			service := service.NewShopService(&service.ShopSConfig{
				ShopRepository: mockRepo,
			})

			result, err := service.GetShopFinanceOverview(tc.input.userId)

			assert.Equal(t, tc.expected.result, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestGetShopStats(t *testing.T) {
	var (
		userId = 1
		shopId = 1
	)
	type input struct {
		userId     int
		beforeTest func(*mocks.ShopRepository)
	}

	type expected struct {
		result *dto.GetShopStatsResponse
		err    error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return error when shop not found",
			input: input{
				userId: userId,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByUserId", userId).Return(nil, errs.ErrShopNotFound)
				},
			},
			expected: expected{
				result: nil,
				err:    errs.ErrShopNotFound,
			},
		},
		{
			description: "should return error when GetShopStats return error",
			input: input{
				userId: userId,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByUserId", userId).Return(&model.Shop{
						ID: shopId,
					}, nil)
					sr.On("GetShopStats", shopId).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				result: nil,
				err:    errs.ErrInternalServerError,
			},
		},
		{
			description: "should return shop stats when success",
			input: input{
				userId: userId,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByUserId", userId).Return(&model.Shop{
						ID: shopId,
					}, nil)
					sr.On("GetShopStats", shopId).Return(&dto.GetShopStatsResponse{}, nil)
				},
			},
			expected: expected{
				result: &dto.GetShopStatsResponse{},
				err:    nil,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.ShopRepository)
			tc.beforeTest(mockRepo)
			service := service.NewShopService(&service.ShopSConfig{
				ShopRepository: mockRepo,
			})

			result, err := service.GetShopStats(tc.input.userId)

			assert.Equal(t, tc.expected.result, result)
			assert.ErrorIs(t, tc.expected.err, err)
		})
	}
}

func TestGetShopInsight(t *testing.T) {
	var (
		userId = 1
		shopId = 1
		req    = dto.GetShopInsightRequest{
			UserId: userId,
		}
	)
	type input struct {
		req        dto.GetShopInsightRequest
		beforeTest func(*mocks.ShopRepository)
	}

	type expected struct {
		result *dto.GetShopInsightResponse
		err    error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return error when shop not found",
			input: input{
				req: req,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByUserId", userId).Return(nil, errs.ErrShopNotFound)
				},
			},
			expected: expected{
				result: nil,
				err:    errs.ErrShopNotFound,
			},
		},
		{
			description: "should return error when GetShopInsight return error",
			input: input{
				req: req,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByUserId", userId).Return(&model.Shop{
						ID: shopId,
					}, nil)
					sr.On("GetShopInsight", shopId, req).Return(nil, errs.ErrInternalServerError)
				},
			},
			expected: expected{
				result: nil,
				err:    errs.ErrInternalServerError,
			},
		},
		{
			description: "should return shop insight when success",
			input: input{
				req: req,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByUserId", userId).Return(&model.Shop{
						ID: shopId,
					}, nil)
					sr.On("GetShopInsight", shopId, req).Return(&dto.GetShopInsightResponse{}, nil)
				},
			},
			expected: expected{
				result: &dto.GetShopInsightResponse{},
				err:    nil,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.ShopRepository)
			tc.beforeTest(mockRepo)
			service := service.NewShopService(&service.ShopSConfig{
				ShopRepository: mockRepo,
			})

			result, err := service.GetShopInsight(tc.input.req)

			assert.Equal(t, tc.expected.result, result)
			assert.ErrorIs(t, tc.expected.err, err)
		})
	}
}

func TestGetShopProfile(t *testing.T) {
	var (
		userId = 1
		shopId = 1
	)
	type input struct {
		userId     int
		beforeTest func(*mocks.ShopRepository)
	}

	type expected struct {
		result *dto.ShopProfile
		err    error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return error when shop not found",
			input: input{
				userId: userId,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByUserId", userId).Return(nil, errs.ErrShopNotFound)
				},
			},
			expected: expected{
				result: nil,
				err:    errs.ErrShopNotFound,
			},
		},
		{
			description: "should return shop profile when success",
			input: input{
				userId: userId,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByUserId", userId).Return(&model.Shop{
						ID: shopId,
					}, nil)
				},
			},
			expected: expected{
				result: &dto.ShopProfile{},
				err:    nil,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.ShopRepository)
			tc.beforeTest(mockRepo)
			service := service.NewShopService(&service.ShopSConfig{
				ShopRepository: mockRepo,
			})

			result, err := service.GetShopProfile(tc.input.userId)

			assert.Equal(t, tc.expected.result, result)
			assert.ErrorIs(t, tc.expected.err, err)
		})
	}
}

func TestUpdateShopProfile(t *testing.T) {
	var (
		userId = 1
		shopId = 1
		req    = dto.ShopProfile{
			Name: "shop name",
		}
	)
	type input struct {
		req        dto.ShopProfile
		beforeTest func(*mocks.ShopRepository)
	}

	type expected struct {
		result *dto.ShopProfile
		err    error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return error when shop not found",
			input: input{
				req: req,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByUserId", userId).Return(nil, errs.ErrShopNotFound)
				},
			},
			expected: expected{
				result: nil,
				err:    errs.ErrShopNotFound,
			},
		},
		{
			description: "should return error when UpdateShopProfile return error",
			input: input{
				req: req,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByUserId", userId).Return(&model.Shop{
						ID: shopId,
					}, nil)
					sr.On("UpdateShop", mock.Anything).Return(errs.ErrInternalServerError)
				},
			},
			expected: expected{
				result: nil,
				err:    errs.ErrInternalServerError,
			},
		},
		{
			description: "should return shop profile when success",
			input: input{
				req: req,
				beforeTest: func(sr *mocks.ShopRepository) {
					sr.On("FindShopByUserId", userId).Return(&model.Shop{
						ID: shopId,
					}, nil)
					sr.On("UpdateShop", mock.Anything).Return(nil)
				},
			},
			expected: expected{
				result: &dto.ShopProfile{},
				err:    nil,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.ShopRepository)
			tc.beforeTest(mockRepo)
			service := service.NewShopService(&service.ShopSConfig{
				ShopRepository: mockRepo,
			})

			err := service.UpdateShopProfile(userId, tc.input.req)

			assert.ErrorIs(t, tc.expected.err, err)
		})
	}
}

func TestCreateShop(t *testing.T) {
	type input struct {
		userID  int
		request *dto.CreateShopRequest
	}
	type expected struct {
		data *model.Shop
		err  error
	}

	var (
		userID          = 1
		shopName        = "shop name 123"
		AddressID       = 1
		courierIDs      = []int{1, 2}
		photoUrl        = "http://image.png"
		courierServices = []*model.CourierService{}
	)

	tests := []struct {
		description string
		input
		beforeTest func(*mocks.ShopRepository, *mocks.CourierServiceService)
		expected
	}{
		{
			description: "should return error when failed to get previous shop",
			input: input{
				userID: userID,
				request: &dto.CreateShopRequest{
					Name:       shopName,
					AddressID:  AddressID,
					CourierIDs: courierIDs,
					PhotoUrl:   &photoUrl,
				},
			},
			beforeTest: func(sr *mocks.ShopRepository, css *mocks.CourierServiceService) {
				sr.On("FindShopByUserId", userID).Return(nil, errors.New("failed to get shop"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get shop"),
			},
		},
		{
			description: "should return error when user already has shop",
			input: input{
				userID: userID,
				request: &dto.CreateShopRequest{
					Name:       shopName,
					AddressID:  AddressID,
					CourierIDs: courierIDs,
					PhotoUrl:   &photoUrl,
				},
			},
			beforeTest: func(sr *mocks.ShopRepository, css *mocks.CourierServiceService) {
				sr.On("FindShopByUserId", userID).Return(&model.Shop{}, nil)
			},
			expected: expected{
				data: nil,
				err:  errs.ErrUserHasShop,
			},
		},
		{
			description: "should return error when shop name is invalid",
			input: input{
				userID: userID,
				request: &dto.CreateShopRequest{
					Name:       "invalid_shop_name",
					AddressID:  AddressID,
					CourierIDs: courierIDs,
					PhotoUrl:   &photoUrl,
				},
			},
			beforeTest: func(sr *mocks.ShopRepository, css *mocks.CourierServiceService) {
				sr.On("FindShopByUserId", userID).Return(nil, errs.ErrShopNotFound)
			},
			expected: expected{
				data: nil,
				err:  errs.ErrInvalidShopName,
			},
		},
		{
			description: "should return error when failed to get courier services",
			input: input{
				userID: userID,
				request: &dto.CreateShopRequest{
					Name:       shopName,
					AddressID:  AddressID,
					CourierIDs: courierIDs,
					PhotoUrl:   &photoUrl,
				},
			},
			beforeTest: func(sr *mocks.ShopRepository, css *mocks.CourierServiceService) {
				sr.On("FindShopByUserId", userID).Return(nil, errs.ErrShopNotFound)
				css.On("GetCourierServicesByCourierIDs", courierIDs).Return(nil, errors.New("failed to get courier services"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to get courier services"),
			},
		},
		{
			description: "should return error when failed to create shop",
			input: input{
				userID: userID,
				request: &dto.CreateShopRequest{
					Name:       shopName,
					AddressID:  AddressID,
					CourierIDs: courierIDs,
					PhotoUrl:   &photoUrl,
				},
			},
			beforeTest: func(sr *mocks.ShopRepository, css *mocks.CourierServiceService) {
				sr.On("FindShopByUserId", userID).Return(nil, errs.ErrShopNotFound)
				css.On("GetCourierServicesByCourierIDs", courierIDs).Return(courierServices, nil)
				sr.On("Create", mock.Anything).Return(errors.New("failed to create shop"))
			},
			expected: expected{
				data: nil,
				err:  errors.New("failed to create shop"),
			},
		},
		{
			description: "should return created shop when succeed to create shop",
			input: input{
				userID: userID,
				request: &dto.CreateShopRequest{
					Name:       shopName,
					AddressID:  AddressID,
					CourierIDs: courierIDs,
					PhotoUrl:   &photoUrl,
				},
			},
			beforeTest: func(sr *mocks.ShopRepository, css *mocks.CourierServiceService) {
				sr.On("FindShopByUserId", userID).Return(nil, errs.ErrShopNotFound)
				css.On("GetCourierServicesByCourierIDs", courierIDs).Return(courierServices, nil)
				sr.On("Create", mock.Anything).Return(nil)
			},
			expected: expected{
				data: &model.Shop{
					UserID:         userID,
					Name:           shopName,
					AddressID:      AddressID,
					CourierService: courierServices,
					PhotoUrl:       &photoUrl,
					Slug:           stringUtils.GenerateSlug(strings.ToLower(shopName)),
				},
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			shopRepo := mocks.NewShopRepository(t)
			courierServiceService := mocks.NewCourierServiceService(t)
			tc.beforeTest(shopRepo, courierServiceService)
			shopService := service.NewShopService(&service.ShopSConfig{
				ShopRepository:        shopRepo,
				CourierServiceService: courierServiceService,
			})

			data, err := shopService.CreateShop(tc.input.userID, tc.input.request)

			assert.Equal(t, tc.expected.data, data)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

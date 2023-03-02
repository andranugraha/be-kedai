package service_test

import (
	"errors"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	var(
		shopList = []*model.Shop{}
		rows = int64(1)
		pages = 1
		limit = 10
		pagination = &commonDto.PaginationResponse{
			Data: shopList,
			TotalRows: rows,
			TotalPages: pages,
			Page: pages,
			Limit: limit,
		}
		emptyPagination = &commonDto.PaginationResponse{
			Page: pages,
			Limit: limit,
		}
		req = &dto.FindShopRequest{
			Limit: limit,
			Page: pages,
			Keyword: "test",
		}
		invalidReq = &dto.FindShopRequest{
			Limit: limit,
			Page: pages,
		}
	)
	type input struct {
		dto *dto.FindShopRequest
		err error
		beforeTest func(*mocks.ShopRepository)
	}
	type expected struct {
		result *commonDto.PaginationResponse
		err error
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
				err: nil,
			},
		},
		{
			description: "should return empty shop list when keyword is empty",
			input: input{
				dto: invalidReq,
				err: nil,
				beforeTest: func(sr *mocks.ShopRepository) {},
			},
			expected: expected{
				result: emptyPagination,
				err: nil,
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
				err: errors.New("error"),
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

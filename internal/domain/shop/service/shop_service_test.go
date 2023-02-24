package service_test

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindShopById(t *testing.T) {
	var(
		shopId = 1
		shopResult = &model.Shop{
			ID: 1,
		}
	)
	type input struct {
		id int
		err error
	}

	type expected struct {
		shop *model.Shop
		err error
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
				id: shopId,
				err: nil,
			},
			expected: expected{
				shop: shopResult,
				err: nil,
			},
		},
		{
			description: "should return error when shop not found",
			input: input{
				id: shopId,
				err: errs.ErrShopNotFound,
			},
			expected: expected{
				shop: nil,
				err: errs.ErrShopNotFound,
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
	var(
		shopId = 1
		shopResult = &model.Shop{
			ID: 1,
		}
	)
	type input struct {
		id int
		err error
	}

	type expected struct {
		shop *model.Shop
		err error
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
				id: shopId,
				err: nil,
			},
			expected: expected{
				shop: shopResult,
				err: nil,
			},
		},
		{
			description: "should return error when shop not found",
			input: input{
				id: shopId,
				err: errs.ErrShopNotFound,
			},
			expected: expected{
				shop: nil,
				err: errs.ErrShopNotFound,
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
	var(
		shopSlug = "shop"
		shopResult = &model.Shop{
			ID: 1,
		}
	)
	type input struct {
		slug string
		err error
	}

	type expected struct {
		shop *model.Shop
		err error
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
				err: nil,
			},
			expected: expected{
				shop: shopResult,
				err: nil,
			},
		},
		{
			description: "should return error when shop not found",
			input: input{
				slug: shopSlug,
				err: errs.ErrShopNotFound,
			},
			expected: expected{
				shop: nil,
				err: errs.ErrShopNotFound,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.ShopRepository)
			mockRepo.On("FindShopBySlug", tc.input.slug).Return(tc.expected.shop, tc.expected.err)
			service := service.NewShopService(&service.ShopSConfig{
				ShopRepository: mockRepo,
			})
			
			result, err := service.FindShopBySlug(tc.input.slug)

			assert.Equal(t, tc.expected.shop, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}
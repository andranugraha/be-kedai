package service_test

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errorResponse "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/service"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDisscussionByProductID(t *testing.T) {
	var (
		productID  = 1
		data       = []*dto.Discussion{}
		limit      = 0
		totalRows  = 0
		TotalPages = 0
		page       = 0
		discussion = commonDto.PaginationResponse{
			Data:       data,
			Limit:      limit,
			Page:       page,
			TotalRows:  int64(totalRows),
			TotalPages: TotalPages,
		}
	)

	type input struct {
		productID  int
		req        dto.GetDiscussionReq
		data       []*dto.Discussion
		limit      int
		totalRows  int
		totalPages int
		page       int
		err        error
	}

	type expected struct {
		discussion *commonDto.PaginationResponse
		err        error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return discussion when success",
			input: input{
				productID:  productID,
				data:       data,
				limit:      limit,
				totalRows:  totalRows,
				totalPages: TotalPages,
				page:       page,
				req: dto.GetDiscussionReq{
					Limit: limit,
					Page:  page,
				},
				err: nil,
			},
			expected: expected{
				discussion: &discussion,
				err:        nil,
			},
		},
		{
			description: "should return error when failed",
			input: input{
				productID:  productID,
				data:       nil,
				limit:      limit,
				totalRows:  totalRows,
				totalPages: TotalPages,
				page:       page,
				req: dto.GetDiscussionReq{
					Limit: limit,
					Page:  page,
				},
				err: errorResponse.ErrInternalServerError,
			},
			expected: expected{
				discussion: nil,
				err:        errorResponse.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockDiscussionRepository := new(mocks.DiscussionRepository)
			mockDiscussionRepository.On("GetDiscussionByProductID", tc.input.productID, tc.input.req).Return(tc.input.data, tc.input.limit, tc.input.page, tc.input.totalRows, tc.input.totalPages, tc.input.err)

			discussionService := service.NewDiscussionService(&service.DiscussionSConfig{
				DiscussionRepository: mockDiscussionRepository,
			})

			discussion, err := discussionService.GetDiscussionByProductID(tc.input.productID, tc.input.req)
			assert.Equal(t, tc.expected.err, err)
			assert.Equal(t, tc.expected.discussion, discussion)
		})
	}

}

func TestGetChildDiscussionByParentID(t *testing.T) {
	var (
		parentID   = 1
		discussion = []*dto.DiscussionReply{}
	)

	type input struct {
		parentID int
		err      error
	}

	type expected struct {
		discussion []*dto.DiscussionReply
		err        error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return discussion when success",
			input: input{
				parentID: parentID,
				err:      nil,
			},
			expected: expected{
				discussion: discussion,
				err:        nil,
			},
		},
		{
			description: "should return error when failed",
			input: input{
				parentID: parentID,
				err:      errorResponse.ErrInternalServerError,
			},
			expected: expected{
				discussion: nil,
				err:        errorResponse.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockDiscussionRepository := new(mocks.DiscussionRepository)
			mockDiscussionRepository.On("GetChildDiscussionByParentID", tc.input.parentID).Return(tc.expected.discussion, tc.input.err)

			discussionService := service.NewDiscussionService(&service.DiscussionSConfig{
				DiscussionRepository: mockDiscussionRepository,
			})

			discussion, err := discussionService.GetChildDiscussionByParentID(tc.input.parentID)
			assert.Equal(t, tc.expected.err, err)
			assert.Equal(t, tc.expected.discussion, discussion)
		})
	}

}

func TestPostDiscussion(t *testing.T) {

	var (
		IsSellerTrue = true
		shop         = model.Shop{}
	)

	type input struct {
		discussion *dto.DiscussionReq
		err        error
		beforeTest func(*mocks.ShopService)
	}

	type expected struct {
		err error
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return success when success",
			input: input{
				beforeTest: func(ss *mocks.ShopService) {
					ss.On("FindShopByUserId", 1).Return(&shop, nil)
				},
				discussion: &dto.DiscussionReq{
					ProductID: 1,
					UserID:    1,
					Message:   "test",
					IsSeller:  IsSellerTrue,
				},
				err: nil,
			},
			expected: expected{
				err: nil,
			},
		},
		{
			description: "should return error when failed shop not found",
			input: input{
				beforeTest: func(ss *mocks.ShopService) {
					ss.On("FindShopByUserId", 1).Return(&shop, errorResponse.ErrInternalServerError)
				},
				discussion: &dto.DiscussionReq{
					ProductID: 1,
					UserID:    1,
					Message:   "test",
					IsSeller:  IsSellerTrue,
				},
				err: errorResponse.ErrInternalServerError,
			},
			expected: expected{
				err: errorResponse.ErrInternalServerError,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			mockDiscussionRepository := new(mocks.DiscussionRepository)
			mockDiscussionRepository.On("PostDiscussion", tc.input.discussion).Return(tc.input.err)

			mockShopService := new(mocks.ShopService)
			tc.input.beforeTest(mockShopService)

			discussionService := service.NewDiscussionService(&service.DiscussionSConfig{
				DiscussionRepository: mockDiscussionRepository,
				ShopService:          mockShopService,
			})

			err := discussionService.PostDiscussion(tc.input.discussion)
			assert.Equal(t, tc.expected.err, err)
		})
	}

}

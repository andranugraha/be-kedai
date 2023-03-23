package service_test

import (
	errorResponse "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDisscussionByProductID(t *testing.T) {
	var (
		productID  = 1
		discussion = []*dto.Discussion{}
	)

	type input struct {
		productID int
		err       error
	}

	type expected struct {
		discussion []*dto.Discussion
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
				productID: productID,
				err:       nil,
			},
			expected: expected{
				discussion: discussion,
				err:        nil,
			},
		},
		{
			description: "should return error when failed",
			input: input{
				productID: productID,
				err:       errorResponse.ErrInternalServerError,
			},
			expected: expected{
				discussion: nil,
				err:        errorResponse.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockDiscussionRepository := new(mocks.DiscussionRepository)
			mockDiscussionRepository.On("GetDiscussionByProductID", tc.input.productID).Return(tc.expected.discussion, tc.input.err)

			discussionService := service.NewDiscussionService(&service.DiscussionSConfig{
				DiscussionRepository: mockDiscussionRepository,
			})

			discussion, err := discussionService.GetDiscussionByProductID(tc.input.productID)
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

// func TestPostDiscussion(t *testing.T) {
// 	var(
// 		IsSeller = false
// 	)

// 	type input struct {
// 		discussion *dto.DiscussionReq
// 		err        error
// 	}

// 	type expected struct {
// 		err error
// 	}

// 	tests := []struct {
// 		description string
// 		input
// 		expected
// 	}{
// 		{
// 			description: "should return success when success",
// 			input: input{
// 				discussion: &dto.DiscussionReq{
// 					ProductID: 1,
// 					UserID:    1,
// 					Message:   "test",
// 					IsSeller:  &IsSeller,
// 				},
// 				err: nil,
// 			},
// 			expected: expected{
// 				err: nil,
// 			},
// 		},
// 	}

// }

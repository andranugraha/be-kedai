package handler_test

import (
	"encoding/json"
	"fmt"
	"kedai/backend/be-kedai/internal/common/code"
	errorResponse "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/handler"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetDiscussionByProductID(t *testing.T) {
	var (
		productID  = 1
		discussion = []*dto.Discussion{}
	)

	type input struct {
		productID int
		err       error
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
			description: "should return discussion when success",
			input: input{
				productID: productID,
				err:       nil,
			},
			expected: expected{
				statusCode: 200,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    discussion,
				},
			},
		},
		{
			description: "should return error when failed",
			input: input{
				productID: productID,
				err:       errorResponse.ErrInternalServerError,
			},
			expected: expected{
				statusCode: 500,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errorResponse.ErrInternalServerError.Error(),
				},
			},
		},
		{
			description: "should return error when bad request",
			input: input{

				err: errorResponse.ErrBadRequest,
			},
			expected: expected{
				statusCode: 500,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: errorResponse.ErrBadRequest.Error(),
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			expectedBody, _ := json.Marshal(tc.expected.response)
			mockService := new(mocks.DiscussionService)
			mockService.On("GetDiscussionByProductID", tc.input.productID).Return(tc.expected.response.Data, tc.input.err)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			h := handler.New(&handler.Config{
				DiscussionService: mockService,
			})
			c.AddParam("productId", fmt.Sprintf("%d", tc.input.productID))
			c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/products/discussions/%d", tc.input.productID), nil)
			h.GetDiscussionByProductID(c)
			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, expectedBody, rec.Body.Bytes())
		})

	}

}

func TestGetDiscussionByParentID(t *testing.T) {
	var (
		parentID   = 1
		discussion = []*dto.DiscussionReply{}
	)

	type input struct {
		parentID int
		err      error
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
			description: "should return discussion when success",
			input: input{
				parentID: parentID,
				err:      nil,
			},
			expected: expected{
				statusCode: 200,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    discussion,
				},
			},
		},
		{
			description: "should return error when failed",
			input: input{
				parentID: parentID,
				err:      errorResponse.ErrInternalServerError,
			},
			expected: expected{
				statusCode: 500,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errorResponse.ErrInternalServerError.Error(),
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			expectedBody, _ := json.Marshal(tc.expected.response)
			mockService := new(mocks.DiscussionService)
			mockService.On("GetChildDiscussionByParentID", tc.input.parentID).Return(tc.expected.response.Data, tc.input.err)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			h := handler.New(&handler.Config{
				DiscussionService: mockService,
			})
			c.AddParam("parentId", fmt.Sprintf("%d", tc.input.parentID))
			c.Request, _ = http.NewRequest("GET", fmt.Sprintf("/products/discussions/replies/%d", tc.input.parentID), nil)
			h.GetDiscussionByParentID(c)
			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, expectedBody, rec.Body.Bytes())
		})
	}

}

func TestPostDiscussion(t *testing.T) {

	var (
		discussion = &dto.DiscussionReq{
			ProductID: 1,
			Message:   "test",
		}
		userId = 0
	)

	type input struct {
		discussion *dto.DiscussionReq
		userId     int
		err        error
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
			description: "should return discussion when success",
			input: input{
				discussion: discussion,
				userId:     userId,
				err:        nil,
			},
			expected: expected{
				statusCode: 200,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
				},
			},
		},
		{
			description: "should return error when failed",
			input: input{
				discussion: discussion,
				userId:     userId,
				err:        errorResponse.ErrInternalServerError,
			},
			expected: expected{
				statusCode: 500,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errorResponse.ErrInternalServerError.Error(),
				},
			},
		},
		{
			description: "should return bad request when body is empty",
			input: input{
				discussion: &dto.DiscussionReq{},
				userId:     userId,
				err:        errorResponse.ErrBadRequest,
			},
			expected: expected{
				statusCode: 400,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Message is required",
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			expectedBody, _ := json.Marshal(tc.expected.response)
			mockService := new(mocks.DiscussionService)
			mockService.On("PostDiscussion", tc.input.discussion).Return(tc.input.err)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			h := handler.New(&handler.Config{
				DiscussionService: mockService,
			})
			payload := test.MakeRequestBody(tc.input.discussion)
			c.Request, _ = http.NewRequest("POST", "/products/discussions", payload)
			c.Request.Header.Set("Content-Type", "application/json")
			c.Set("userId", tc.input.userId)
			h.PostDiscussion(c)
			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, expectedBody, rec.Body.Bytes())
		})
	}

}

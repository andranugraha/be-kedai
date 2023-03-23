package handler_test

import (
	"encoding/json"
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/chat/dto"
	"kedai/backend/be-kedai/internal/domain/chat/handler"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	testutil "kedai/backend/be-kedai/internal/utils/test"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserAddChat(t *testing.T) {
	var (
		userId   = 1
		shopSlug = "shop-A"
		body     = &dto.SendChatBodyRequest{
			Message: "Hai sayang",
			Type:    "text",
		}
		chat = &dto.ChatResponse{
			ID:         1,
			Message:    "Hai sayang",
			Time:       time.Now(),
			Type:       "text",
			IsIncoming: false,
		}
	)

	type input struct {
		userId   int
		shopSlug string
		body     *dto.SendChatBodyRequest
		result   *dto.ChatResponse
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
			description: "should return chat response with code 201 when success",
			input: input{
				userId:   userId,
				shopSlug: shopSlug,
				body:     body,
				result:   chat,
				err:      nil,
			},
			expected: expected{
				statusCode: http.StatusCreated,
				response: response.Response{
					Code:    code.CREATED,
					Message: "success",
					Data:    chat,
				},
			},
		},
		{
			description: "should return error with code 400 when bad request",
			input: input{
				userId:   userId,
				shopSlug: shopSlug,
				body:     body,
				result:   nil,
				err:      errs.ErrShopNotFound,
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: errs.ErrShopNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with code 500 when error",
			input: input{
				userId:   userId,
				shopSlug: shopSlug,
				body:     body,
				result:   nil,
				err:      errors.New("error"),
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
			mockChatService := mocks.NewChatService(t)
			mockChatService.On("UserAddChat", tc.input.body, tc.input.userId, tc.input.shopSlug).Return(tc.input.result, tc.input.err)
			handler := handler.New(&handler.Config{
				ChatService: mockChatService,
			})
			c.Set("userId", 1)
			c.Params = gin.Params{
				{
					Key:   "shopSlug",
					Value: shopSlug,
				},
			}

			c.Request, _ = http.NewRequest("POST", "/users/chats", testutil.MakeRequestBody(tc.input.body))

			handler.UserAddChat(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}

}

func TestSellerAddChat(t *testing.T) {
	var (
		userId   = 1
		username = "usernameA"
		body     = &dto.SendChatBodyRequest{
			Message: "Hai sayang",
			Type:    "text",
		}
		chat = &dto.ChatResponse{
			ID:         1,
			Message:    "Hai sayang",
			Time:       time.Now(),
			Type:       "text",
			IsIncoming: false,
		}
	)

	type input struct {
		userId   int
		username string
		body     *dto.SendChatBodyRequest
		result   *dto.ChatResponse
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
			description: "should return chat response with code 201 when success",
			input: input{
				userId:   userId,
				username: username,
				body:     body,
				result:   chat,
				err:      nil,
			},
			expected: expected{
				statusCode: http.StatusCreated,
				response: response.Response{
					Code:    code.CREATED,
					Message: "success",
					Data:    chat,
				},
			},
		},
		{
			description: "should return error with code 400 when bad request",
			input: input{
				userId:   userId,
				username: username,
				body:     body,
				result:   nil,
				err:      errs.ErrShopNotFound,
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: errs.ErrShopNotFound.Error(),
				},
			},
		},
		{
			description: "should return error with code 500 when error",
			input: input{
				userId:   userId,
				username: username,
				body:     body,
				result:   nil,
				err:      errors.New("error"),
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
			mockChatService := mocks.NewChatService(t)
			mockChatService.On("SellerAddChat", tc.input.body, tc.input.userId, tc.input.username).Return(tc.input.result, tc.input.err)
			handler := handler.New(&handler.Config{
				ChatService: mockChatService,
			})
			c.Set("userId", 1)
			c.Params = gin.Params{
				{
					Key:   "username",
					Value: username,
				},
			}

			c.Request, _ = http.NewRequest("POST", "/sellers/chats", testutil.MakeRequestBody(tc.input.body))

			handler.SellerAddChat(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}

}

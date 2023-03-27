package handler_test

import (
	"encoding/json"
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/chat/dto"
	"kedai/backend/be-kedai/internal/domain/chat/handler"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	testutil "kedai/backend/be-kedai/internal/utils/test"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserGetListOfChats(t *testing.T) {
	var (
		userId        = 1
		param         = &dto.ListOfChatsParamRequest{}
		chatResponses = []*dto.UserListOfChatResponse{}
	)

	type input struct {
		userId int
		param  *dto.ListOfChatsParamRequest
		result []*dto.UserListOfChatResponse
		err    error
	}

	type expected struct {
		statusCode int
		response   response.Response
	}

	type cases struct {
		description string
		input
		beforeTests func(cs *mocks.ChatService)
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return chat response with code 200 when success",
			input: input{
				userId: userId,
				param:  param,
				result: chatResponses,
				err:    nil,
			},
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("UserGetListOfChats", param, userId).Return(chatResponses, nil)
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    chatResponses,
				},
			},
		},
		{
			description: "should return error with code 500 when error",
			input: input{
				userId: userId,
				param:  param,
				result: nil,
				err:    errors.New("error"),
			},
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("UserGetListOfChats", param, userId).Return(nil, errors.New("error"))
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
			tc.beforeTests(mockChatService)
			handler := handler.New(&handler.Config{
				ChatService: mockChatService,
			})
			c.Set("userId", 1)

			c.Request, _ = http.NewRequest("GET", "/users/chats", nil)

			handler.UserGetListOfChats(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestSellerGetListOfChats(t *testing.T) {
	var (
		userId        = 1
		param         = &dto.ListOfChatsParamRequest{}
		chatResponses = []*dto.SellerListOfChatResponse{}
	)

	type input struct {
		userId int
		param  *dto.ListOfChatsParamRequest
		result []*dto.SellerListOfChatResponse
		err    error
	}

	type expected struct {
		statusCode int
		response   response.Response
	}

	type cases struct {
		description string
		input
		beforeTests func(cs *mocks.ChatService)
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return chat response with code 200 when success",
			input: input{
				userId: userId,
				param:  param,
				result: chatResponses,
				err:    nil,
			},
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("SellerGetListOfChats", param, userId).Return(chatResponses, nil)
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    chatResponses,
				},
			},
		},
		{
			description: "should return error with code 400 when bad request",
			input: input{
				userId: userId,
				param:  param,
				result: nil,
				err:    errs.ErrShopNotFound,
			},
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("SellerGetListOfChats", param, userId).Return(nil, errs.ErrShopNotFound)
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
				userId: userId,
				param:  param,
				result: nil,
				err:    errors.New("error"),
			},
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("SellerGetListOfChats", param, userId).Return(nil, errors.New("error"))
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
			tc.beforeTests(mockChatService)
			handler := handler.New(&handler.Config{
				ChatService: mockChatService,
			})
			c.Set("userId", 1)

			c.Request, _ = http.NewRequest("GET", "/users/chats", nil)

			handler.SellerGetListOfChats(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestUserGetChat(t *testing.T) {
	var (
		userId   = 1
		shopSlug = "shop-A"
		param    = &dto.ChatParamRequest{
			Page:       1,
			LimitByDay: 366,
		}
		chat = &dto.ChatResponse{
			ID:         1,
			Message:    "Hai sayang",
			Time:       time.Now(),
			Type:       "text",
			IsIncoming: false,
		}
		paginatedChats = &commonDto.PaginationResponse{Data: chat}
	)

	type input struct {
		userId   int
		shopSlug string
		param    *dto.ChatParamRequest
		result   *commonDto.PaginationResponse
		err      error
	}

	type expected struct {
		statusCode int
		response   response.Response
	}

	type cases struct {
		description string
		input
		beforeTests func(cs *mocks.ChatService)
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return chat response with code 200 when success",
			input: input{
				userId:   userId,
				shopSlug: shopSlug,
				param:    param,
				result:   paginatedChats,
				err:      nil,
			},
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("UserGetChat", param, userId, shopSlug).Return(paginatedChats, nil)
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    paginatedChats,
				},
			},
		},
		{
			description: "should return error with code 400 when bad request",
			input: input{
				userId:   userId,
				shopSlug: shopSlug,
				param:    param,
				result:   nil,
				err:      errs.ErrShopNotFound,
			},
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("UserGetChat", param, userId, shopSlug).Return(nil, errs.ErrShopNotFound)
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
				param:    param,
				result:   nil,
				err:      errors.New("error"),
			},
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("UserGetChat", param, userId, shopSlug).Return(nil, errors.New("error"))
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
			tc.beforeTests(mockChatService)
			handler := handler.New(&handler.Config{
				ChatService: mockChatService,
			})
			c.Set("userId", 1)
			c.Params = gin.Params{
				{
					Key:   "shopSlug",
					Value: shopSlug,
				},
				{
					Key:   "page",
					Value: strconv.Itoa(param.Page),
				},
				{
					Key:   "limitByDay",
					Value: strconv.Itoa(param.LimitByDay),
				},
			}

			c.Request, _ = http.NewRequest("GET", "/users/chats", nil)

			handler.UserGetChat(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestSellerGetChat(t *testing.T) {
	var (
		userId   = 1
		username = "usernameA"
		param    = &dto.ChatParamRequest{
			Page:       1,
			LimitByDay: 366,
		}
		chat = &dto.ChatResponse{
			ID:         1,
			Message:    "Hai sayang",
			Time:       time.Now(),
			Type:       "text",
			IsIncoming: false,
		}
		paginatedChats = &commonDto.PaginationResponse{Data: chat}
	)

	type input struct {
		userId   int
		username string
		param    *dto.ChatParamRequest
		result   *commonDto.PaginationResponse
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
		beforeTests func(cs *mocks.ChatService)
	}

	for _, tc := range []cases{
		{
			description: "should return chat response with code 200 when success",
			input: input{
				userId:   userId,
				username: username,
				param:    param,
				result:   paginatedChats,
				err:      nil,
			},
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("SellerGetChat", param, userId, username).Return(paginatedChats, nil)
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    paginatedChats,
				},
			},
		},
		{
			description: "should return error with code 400 when bad request",
			input: input{
				userId:   userId,
				username: username,
				param:    param,
				result:   nil,
				err:      errs.ErrShopNotFound,
			},
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("SellerGetChat", param, userId, username).Return(nil, errs.ErrShopNotFound)
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
				param:    param,
				result:   nil,
				err:      errors.New("error"),
			},
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("SellerGetChat", param, userId, username).Return(nil, errors.New("error"))
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
			tc.beforeTests(mockChatService)
			handler := handler.New(&handler.Config{
				ChatService: mockChatService,
			})
			c.Set("userId", 1)
			c.Params = gin.Params{
				{
					Key:   "username",
					Value: username,
				},
				{
					Key:   "page",
					Value: strconv.Itoa(param.Page),
				},
				{
					Key:   "limitByDay",
					Value: strconv.Itoa(param.LimitByDay),
				},
			}

			c.Request, _ = http.NewRequest("GET", "/sellers/chats", nil)

			handler.SellerGetChat(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

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
		beforeTests func(cs *mocks.ChatService)
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
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("UserAddChat", body, userId, shopSlug).Return(chat, nil)
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
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("UserAddChat", body, userId, shopSlug).Return(nil, errs.ErrShopNotFound)
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
			description: "should return error with code 400 when bad params",
			input: input{
				userId:   userId,
				shopSlug: shopSlug,
				body:     &dto.SendChatBodyRequest{},
				result:   nil,
				err:      errs.ErrShopNotFound,
			},
			beforeTests: func(cs *mocks.ChatService) {},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: errors.New("Message is required").Error(),
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
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("UserAddChat", body, userId, shopSlug).Return(nil, errors.New("error"))
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
			tc.beforeTests(mockChatService)
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
		beforeTests func(cs *mocks.ChatService)
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
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("SellerAddChat", body, userId, username).Return(chat, nil)
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
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("SellerAddChat", body, userId, username).Return(nil, errs.ErrShopNotFound)
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
			description: "should return error with code 400 when bad params",
			input: input{
				userId:   userId,
				username: username,
				body:     &dto.SendChatBodyRequest{},
				result:   nil,
				err:      errs.ErrShopNotFound,
			},
			beforeTests: func(cs *mocks.ChatService) {},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: errors.New("Message is required").Error(),
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
			beforeTests: func(cs *mocks.ChatService) {
				cs.On("SellerAddChat", body, userId, username).Return(nil, errors.New("error"))
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
			tc.beforeTests(mockChatService)
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

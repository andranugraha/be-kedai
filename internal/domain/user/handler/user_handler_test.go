package handler_test

import (
	"encoding/json"
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/handler"
	"kedai/backend/be-kedai/internal/utils/response"
	"kedai/backend/be-kedai/internal/utils/test"
	"kedai/backend/be-kedai/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserRegister(t *testing.T) {
	type input struct {
		user *dto.UserRegistration
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
			description: "should return created user data when successfully registered",
			input: input{
				user: &dto.UserRegistration{
					Email:    "user@mail.com",
					Password: "password",
				},
				err: nil,
			},
			expected: expected{
				statusCode: http.StatusCreated,
				response: response.Response{
					Code:    code.CREATED,
					Message: "created",
					Data: &dto.UserRegistration{
						Email:    "user@mail.com",
						Password: "password",
					},
				},
			},
		},
		{
			description: "should return error when required input not met condition",
			input: input{
				user: &dto.UserRegistration{
					Email: "user@mail.com",
				},
				err: errors.New("bad request"),
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Password is required",
					Data:    nil,
				},
			},
		},
		{
			description: "should return error when email already registered",
			input: input{
				user: &dto.UserRegistration{
					Email:    "user@mail.com",
					Password: "password",
				},
				err: errs.ErrUserAlreadyExist,
			},
			expected: expected{
				statusCode: http.StatusConflict,
				response: response.Response{
					Code:    code.EMAIL_ALREADY_REGISTERED,
					Message: "user already exist",
					Data:    nil,
				},
			},
		},
		{
			description: "should return error when server internal error",
			input: input{
				user: &dto.UserRegistration{
					Email:    "user@mail.com",
					Password: "password",
				},
				err: errs.ErrInternalServerError,
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: "something went wrong in the server",
					Data:    nil,
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			expectedBody, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			mockService := new(mocks.UserService)
			mockService.On("SignUp", tc.input.user).Return(tc.expected.response.Data, tc.input.err)
			h := handler.New(&handler.HandlerConfig{
				UserService: mockService,
			})
			c.Request = httptest.NewRequest("POST", "/users", test.MakeRequestBody(tc.input.user))

			h.UserRegistration(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestUserHandler_GetUserByID(t *testing.T) {
	type input struct {
		userId int
		data   *model.User
		err    error
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "it should return user data with status code 200 if successed getting user data",
			input: input{
				userId: 1,
				data: &model.User{
					Email:    "user@email.com",
					Username: "user_name",
					Profile: &model.UserProfile{
						UserID: 1,
					},
				},
				err: nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data: &model.User{
						Email:    "user@email.com",
						Username: "user_name",
						Profile: &model.UserProfile{
							UserID: 1,
						},
					},
				},
			},
		},
		{
			description: "it should return status code 404 when user does not exist",
			input: input{
				userId: 1,
				data:   nil,
				err:    errs.ErrUserDoesNotExist,
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.USER_NOT_REGISTERED,
					Message: errs.ErrUserDoesNotExist.Error(),
				},
			},
		},
		{
			description: "it should return status code 500 when something went wrong went trying to get user data",
			input: input{
				userId: 1,
				data:   nil,
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
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userId)
			userServiceMock := mocks.NewUserService(t)
			userServiceMock.On("GetByID", tc.input.userId).Return(tc.input.data, tc.input.err)
			cfg := handler.HandlerConfig{
				UserService: userServiceMock,
			}
			h := handler.New(&cfg)
			c.Request, _ = http.NewRequest("GET", "/users", nil)

			h.GetUserByID(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

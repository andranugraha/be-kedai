package handler_test

import (
	"encoding/json"
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
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

func TestUserLogin(t *testing.T) {
	type input struct {
		dto *dto.UserLogin
		err error
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
			description: "should return access token when user log in accepted",
			input: input{
				dto: &dto.UserLogin{
					Email:    "user@mail.com",
					Password: "password",
				},
				err: nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
					Data:    &dto.Token{},
				},
			},
		},
		{
			description: "should return error when required input not met condition",
			input: input{
				dto: &dto.UserLogin{
					Email: "user@mail.com",
				},
				err: errors.New("bad request"),
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Password is required",
				},
			},
		},
		{
			description: "should return error when user input wrong password",
			input: input{
				dto: &dto.UserLogin{
					Email:    "user@mail.com",
					Password: "password",
				},
				err: errs.ErrInvalidCredential,
			},
			expected: expected{
				statusCode: http.StatusUnauthorized,
				response: response.Response{
					Code:    code.UNAUTHORIZED,
					Message: errs.ErrInvalidCredential.Error(),
				},
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				dto: &dto.UserLogin{
					Email:    "user@mail.com",
					Password: "password",
				},
				err: errs.ErrInternalServerError,
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
			mockService := new(mocks.UserService)
			mockService.On("SignIn", tc.input.dto, tc.input.dto.Password).Return(tc.expected.response.Data, tc.input.err)
			h := handler.New(&handler.HandlerConfig{
				UserService: mockService,
			})
			c.Request = httptest.NewRequest("POST", "/users/login", test.MakeRequestBody(tc.input.dto))

			h.UserLogin(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}

}

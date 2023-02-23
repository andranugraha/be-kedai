package handler_test

import (
	"encoding/json"
	"errors"
	"kedai/backend/be-kedai/internal/common/code"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/handler"
	"kedai/backend/be-kedai/internal/domain/user/model"
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
		user *dto.UserRegistrationRequest
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
				user: &dto.UserRegistrationRequest{
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
					Data: &dto.UserRegistrationResponse{
						Email: "user@mail.com",
					},
				},
			},
		},
		{
			description: "should return error when required input not met condition",
			input: input{
				user: &dto.UserRegistrationRequest{
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
				user: &dto.UserRegistrationRequest{
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
			description: "should return error when invalid password pattern",
			input: input{
				user: &dto.UserRegistrationRequest{
					Email:    "user@mail.com",
					Password: "password",
				},
				err: errs.ErrInvalidPasswordPattern,
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_PASSWORD_PATTERN,
					Message: "invalid password pattern",
					Data:    nil,
				},
			},
		},
		{
			description: "should return error when password contain email address",
			input: input{
				user: &dto.UserRegistrationRequest{
					Email:    "user@mail.com",
					Password: "passworD1user",
				},
				err: errs.ErrContainEmail,
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.PASSWORD_CONTAIN_EMAIL,
					Message: "password cannot contain email address",
					Data:    nil,
				},
			},
		},
		{
			description: "should return error when server internal error",
			input: input{
				user: &dto.UserRegistrationRequest{
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

func TestGetSession(t *testing.T) {
	type input struct {
		userId int
		token  string
		err    error
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
			description: "should return error when a session is unavailable",
			input: input{
				userId: 1,
				token:  "",
				err:    errors.New("error"),
			},
			expected: expected{
				statusCode: 401,
				response: response.Response{
					Code:    code.UNAUTHORIZED,
					Message: "error",
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			expectedBody, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userId)
			mockService := new(mocks.UserService)
			mockService.On("GetSession", tc.input.userId, tc.input.token).Return(tc.input.err)
			h := handler.New(&handler.HandlerConfig{
				UserService: mockService,
			})
			c.Request = httptest.NewRequest("POST", "/users/login", nil)

			h.GetSession(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedBody), rec.Body.String())
		})
	}
}

func TestGetUserByID(t *testing.T) {
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
			description: "it should return status code 500 if something went wrong when trying to get user data",
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

func TestUpdateEmail(t *testing.T) {
	type input struct {
		userId     int
		request    *dto.UpdateEmailRequest
		beforeTest func(*mocks.UserService)
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
			description: "should return error with status code 400 when given bad request body",
			input: input{
				userId: 1,
				request: &dto.UpdateEmailRequest{
					Email: "test",
				},
				beforeTest: func(us *mocks.UserService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Email must be an email format",
				},
			},
		},
		{
			description: "should return error with status code 409 when email is already used",
			input: input{
				userId: 1,
				request: &dto.UpdateEmailRequest{
					Email: "used.email@mail.com",
				},
				beforeTest: func(us *mocks.UserService) {
					us.On("UpdateEmail", 1, &dto.UpdateEmailRequest{
						Email: "used.email@mail.com",
					}).Return(nil, errs.ErrEmailUsed)
				},
			},
			expected: expected{
				statusCode: http.StatusConflict,
				response: response.Response{
					Code:    code.EMAIL_ALREADY_REGISTERED,
					Message: errs.ErrEmailUsed.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to update email",
			input: input{
				userId: 1,
				request: &dto.UpdateEmailRequest{
					Email: "new.email@mail.com",
				},
				beforeTest: func(us *mocks.UserService) {
					us.On("UpdateEmail", 1, &dto.UpdateEmailRequest{
						Email: "new.email@mail.com",
					}).Return(nil, errors.New("failed to update email"))
				},
			},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errs.ErrInternalServerError.Error(),
				},
			},
		},
		{
			description: "should return updated user data with status code 200 when update email successed",
			input: input{
				userId: 1,
				request: &dto.UpdateEmailRequest{
					Email: "new.email@mail.com",
				},
				beforeTest: func(us *mocks.UserService) {
					us.On("UpdateEmail", 1, &dto.UpdateEmailRequest{
						Email: "new.email@mail.com",
					}).Return(&dto.UpdateEmailResponse{Email: "new.email@mail.com"}, nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.UPDATED,
					Message: "updated",
					Data:    &dto.UpdateEmailResponse{Email: "new.email@mail.com"},
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
			userService := mocks.NewUserService(t)
			tc.input.beforeTest(userService)
			cfg := handler.HandlerConfig{
				UserService: userService,
			}
			h := handler.New(&cfg)
			c.Request, _ = http.NewRequest("PUT", "/v1/users/emails", test.MakeRequestBody(tc.input.request))

			h.UpdateUserEmail(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

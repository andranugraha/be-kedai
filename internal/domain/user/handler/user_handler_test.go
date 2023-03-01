package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
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

func TestRenewSession(t *testing.T) {
	type input struct {
		userId     int
		token      string
		mockReturn *dto.Token
		err        error
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error with status code 401 when refresh token is expired",
			input: input{
				userId:     1,
				token:      "",
				mockReturn: nil,
				err:        errs.ErrExpiredToken,
			},
			expected: expected{
				statusCode: http.StatusUnauthorized,
				response: response.Response{
					Code:    code.TOKEN_EXPIRED,
					Message: errs.ErrExpiredToken.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to renew session",
			input: input{
				userId:     1,
				token:      "",
				mockReturn: nil,
				err:        errors.New("failed to renew token"),
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
			description: "should return error with status code 200 when token successfully renewed",
			input: input{
				userId:     1,
				token:      "",
				mockReturn: &dto.Token{},
				err:        nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "success",
					Data:    &dto.Token{},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", tc.input.userId)
			mockService := new(mocks.UserService)
			mockService.On("RenewToken", tc.input.userId, tc.input.token).Return(tc.input.mockReturn, tc.input.err)
			h := handler.New(&handler.HandlerConfig{
				UserService: mockService,
			})
			c.Request = httptest.NewRequest("POST", "/users/tokens/refresh", nil)

			h.RenewSession(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
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
			description: "it should return user data with status code 200 if succeed getting user data",
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

func TestUserRegistrationWithGoogle(t *testing.T) {
	var (
		validCredential   = "test"
		invalidCredential = ""
	)
	type input struct {
		dto         *dto.UserRegistrationWithGoogleRequest
		beforeTests func(mockUserService *mocks.UserService)
		err         error
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
			description: "it should return credential required and bad request status code if credential is empty",
			input: input{
				dto: &dto.UserRegistrationWithGoogleRequest{
					Credential: invalidCredential,
					Username:   "testasd",
					Password:   "password123",
				},
				err: nil,
				beforeTests: func(mockUserService *mocks.UserService) {
				},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Credential is required",
				},
			},
		},
		{
			description: "it should return ErrUserAlreadyExist and conflict status code if user already exist",
			input: input{
				dto: &dto.UserRegistrationWithGoogleRequest{
					Credential: validCredential,
					Username:   "testasd",
					Password:   "password123",
				},
				err: errs.ErrUserAlreadyExist,
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("SignUpWithGoogle", &dto.UserRegistrationWithGoogleRequest{
						Credential: validCredential,
						Username:   "testasd",
						Password:   "password123",
					}).Return(nil, errs.ErrUserAlreadyExist)
				},
			},
			expected: expected{
				statusCode: http.StatusConflict,
				response: response.Response{
					Code:    code.EMAIL_ALREADY_REGISTERED,
					Message: errs.ErrUserAlreadyExist.Error(),
				},
			},
		},
		{
			description: "it should return error ErrUsernameUsed and conflict status code if username already used",
			input: input{
				dto: &dto.UserRegistrationWithGoogleRequest{
					Credential: validCredential,
					Username:   "testasd",
					Password:   "password123",
				},
				err: errs.ErrUsernameUsed,
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("SignUpWithGoogle", &dto.UserRegistrationWithGoogleRequest{
						Credential: validCredential,
						Username:   "testasd",
						Password:   "password123",
					}).Return(nil, errs.ErrUsernameUsed)
				},
			},
			expected: expected{
				statusCode: http.StatusConflict,
				response: response.Response{
					Code:    code.USERNAME_ALREADY_REGISTERED,
					Message: errs.ErrUsernameUsed.Error(),
				},
			},
		},
		{
			description: "it should return error ErrInvalidUsernamePattern and StatusUnprocessableEntity status code if username is invalid",
			input: input{
				dto: &dto.UserRegistrationWithGoogleRequest{
					Credential: validCredential,
					Username:   "testasd",
					Password:   "password123",
				},
				err: errs.ErrInvalidUsernamePattern,
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("SignUpWithGoogle", &dto.UserRegistrationWithGoogleRequest{
						Credential: validCredential,
						Username:   "testasd",
						Password:   "password123",
					}).Return(nil, errs.ErrInvalidUsernamePattern)
				},
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_USERNAME_PATTERN,
					Message: errs.ErrInvalidUsernamePattern.Error(),
				},
			},
		},
		{
			description: "it should return error ErrInvalidPasswordPattern and StatusUnprocessableEntity status code if password is invalid",
			input: input{
				dto: &dto.UserRegistrationWithGoogleRequest{
					Credential: validCredential,
					Username:   "testasd",
					Password:   "password123",
				},
				err: errs.ErrUsernameUsed,
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("SignUpWithGoogle", &dto.UserRegistrationWithGoogleRequest{
						Credential: validCredential,
						Username:   "testasd",
						Password:   "password123",
					}).Return(nil, errs.ErrInvalidPasswordPattern)
				},
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_PASSWORD_PATTERN,
					Message: errs.ErrInvalidPasswordPattern.Error(),
				},
			},
		},
		{
			description: "it should return error ErrContainEmail and StatusUnprocessableEntity status code if password is invalid",
			input: input{
				dto: &dto.UserRegistrationWithGoogleRequest{
					Credential: validCredential,
					Username:   "testasd",
					Password:   "password123",
				},
				err: errs.ErrUsernameUsed,
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("SignUpWithGoogle", &dto.UserRegistrationWithGoogleRequest{
						Credential: validCredential,
						Username:   "testasd",
						Password:   "password123",
					}).Return(nil, errs.ErrContainEmail)
				},
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.PASSWORD_CONTAIN_EMAIL,
					Message: errs.ErrContainEmail.Error(),
				},
			},
		},
		{
			description: "it should return error ErrUnauthorized and StatusUnauthorized status code if credential is invalid",
			input: input{
				dto: &dto.UserRegistrationWithGoogleRequest{
					Credential: validCredential,
					Username:   "testasd",
					Password:   "password123",
				},
				err: errs.ErrUsernameUsed,
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("SignUpWithGoogle", &dto.UserRegistrationWithGoogleRequest{
						Credential: validCredential,
						Username:   "testasd",
						Password:   "password123",
					}).Return(nil, errs.ErrUnauthorized)
				},
			},
			expected: expected{
				statusCode: http.StatusUnauthorized,
				response: response.Response{
					Code:    code.UNAUTHORIZED,
					Message: errs.ErrUnauthorized.Error(),
				},
			},
		},
		{
			description: "it should return error ErrInternalServer and StatusInternalServerError status code if error is not handled",
			input: input{
				dto: &dto.UserRegistrationWithGoogleRequest{
					Credential: validCredential,
					Username:   "testasd",
					Password:   "password123",
				},
				err: errs.ErrUsernameUsed,
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("SignUpWithGoogle", &dto.UserRegistrationWithGoogleRequest{
						Credential: validCredential,
						Username:   "testasd",
						Password:   "password123",
					}).Return(nil, errs.ErrInternalServerError)
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
			description: "it should return nil error and StatusOK status code if success",
			input: input{
				dto: &dto.UserRegistrationWithGoogleRequest{
					Credential: validCredential,
					Username:   "testasd",
					Password:   "password123",
				},
				err: nil,
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("SignUpWithGoogle", &dto.UserRegistrationWithGoogleRequest{
						Credential: validCredential,
						Username:   "testasd",
						Password:   "password123",
					}).Return(&dto.Token{}, nil)
				},
			},
			expected: expected{
				statusCode: http.StatusCreated,
				response: response.Response{
					Code:    code.CREATED,
					Message: "Sign up with google successful",
					Data:    &dto.Token{},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			userServiceMock := mocks.NewUserService(t)
			tc.beforeTests(userServiceMock)
			cfg := handler.HandlerConfig{
				UserService: userServiceMock,
			}

			h := handler.New(&cfg)
			c.Request, _ = http.NewRequest("POST", "/users", nil)
			c.Request.Header.Set("Content-Type", "application/json")
			body, _ := json.Marshal(tc.input.dto)
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

			h.UserRegistrationWithGoogle(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestUserLoginWithGoogle(t *testing.T) {
	var (
		validCredential   = "test"
		invalidCredential = ""
	)
	type input struct {
		dto         *dto.UserLoginWithGoogleRequest
		beforeTests func(mockUserService *mocks.UserService)
		err         error
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
			description: "it should return credential required and bad request status code if credential is empty",
			input: input{
				dto: &dto.UserLoginWithGoogleRequest{
					Credential: invalidCredential,
				},
				err: nil,
				beforeTests: func(mockUserService *mocks.UserService) {

				},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Credential is required",
				},
			},
		},
		{
			description: "it should return status code 404 when user does not exist",
			input: input{
				dto: &dto.UserLoginWithGoogleRequest{
					Credential: validCredential,
				},
				err: errs.ErrInvalidCredential,
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("SignInWithGoogle", &dto.UserLoginWithGoogleRequest{
						Credential: validCredential,
					}).Return(nil, errs.ErrInvalidCredential)

				}},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.USER_NOT_REGISTERED,
					Message: errs.ErrInvalidCredential.Error(),
				},
			},
		},
		{
			description: "it should return status code 401 when google jwt token is invalid",
			input: input{
				dto: &dto.UserLoginWithGoogleRequest{
					Credential: validCredential,
				},
				err: errs.ErrUnauthorized,
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("SignInWithGoogle", &dto.UserLoginWithGoogleRequest{
						Credential: validCredential,
					}).Return(nil, errs.ErrUnauthorized)
				},
			},
			expected: expected{
				statusCode: http.StatusUnauthorized,
				response: response.Response{
					Code:    code.UNAUTHORIZED,
					Message: errs.ErrUnauthorized.Error(),
				},
			},
		},
		{
			description: "it should return status code 500 if something went wrong when trying to get user data",
			input: input{
				dto: &dto.UserLoginWithGoogleRequest{
					Credential: validCredential,
				},
				err: errs.ErrInternalServerError,
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("SignInWithGoogle", &dto.UserLoginWithGoogleRequest{
						Credential: validCredential,
					}).Return(nil, errs.ErrInternalServerError)
				}},
			expected: expected{
				statusCode: http.StatusInternalServerError,
				response: response.Response{
					Code:    code.INTERNAL_SERVER_ERROR,
					Message: errs.ErrInternalServerError.Error(),
				},
			},
		},

		{
			description: "it should return status code 200 if succeed login with google",
			input: input{
				dto: &dto.UserLoginWithGoogleRequest{
					Credential: validCredential,
				},
				err: nil,
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("SignInWithGoogle", &dto.UserLoginWithGoogleRequest{
						Credential: validCredential,
					}).Return(&dto.Token{}, nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "Sign in with google successful",
					Data:    &dto.Token{},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request, _ = http.NewRequest("POST", "/v1/users/google-login", nil)
			c.Request.Header.Set("Content-Type", "application/json")
			body, _ := json.Marshal(tc.input.dto)
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

			userServiceMock := mocks.NewUserService(t)
			tc.beforeTests(userServiceMock)
			cfg := handler.HandlerConfig{
				UserService: userServiceMock,
			}
			h := handler.New(&cfg)

			h.UserLoginWithGoogle(c)

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
			description: "should return updated user data with status code 200 when update email succeed",
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

func TestUpdateUsername(t *testing.T) {
	type input struct {
		userId     int
		request    *dto.UpdateUsernameRequest
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
				request: &dto.UpdateUsernameRequest{
					Username: "a_veryveryveryveryveryveryveryveryveryveryveryveryveryveryveryvery_long_username",
				},
				beforeTest: func(us *mocks.UserService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "Username must be shorter than 30",
				},
			},
		},
		{
			description: "should return error with status code 422 when username is invalid",
			input: input{
				userId: 1,
				request: &dto.UpdateUsernameRequest{
					Username: "inval1d_u$ername",
				},
				beforeTest: func(us *mocks.UserService) {
					us.On("UpdateUsername", 1, &dto.UpdateUsernameRequest{
						Username: "inval1d_u$ername",
					}).Return(nil, errs.ErrInvalidUsernamePattern)
				},
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_USERNAME_PATTERN,
					Message: errs.ErrInvalidUsernamePattern.Error(),
				},
			},
		},
		{
			description: "should return error with status code 409 when username is already used",
			input: input{
				userId: 1,
				request: &dto.UpdateUsernameRequest{
					Username: "valid123username",
				},
				beforeTest: func(us *mocks.UserService) {
					us.On("UpdateUsername", 1, &dto.UpdateUsernameRequest{
						Username: "valid123username",
					}).Return(nil, errs.ErrUsernameUsed)
				},
			},
			expected: expected{
				statusCode: http.StatusConflict,
				response: response.Response{
					Code:    code.USERNAME_ALREADY_REGISTERED,
					Message: errs.ErrUsernameUsed.Error(),
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to update email",
			input: input{
				userId: 1,
				request: &dto.UpdateUsernameRequest{
					Username: "validUsername",
				},
				beforeTest: func(us *mocks.UserService) {
					us.On("UpdateUsername", 1, &dto.UpdateUsernameRequest{
						Username: "validUsername",
					}).Return(nil, errors.New("failed to update username"))
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
			description: "should return updated user data with status code 200 when update username succeed",
			input: input{
				userId: 1,
				request: &dto.UpdateUsernameRequest{
					Username: "newUsername123",
				},
				beforeTest: func(us *mocks.UserService) {
					us.On("UpdateUsername", 1, &dto.UpdateUsernameRequest{
						Username: "newUsername123",
					}).Return(&dto.UpdateUsernameResponse{Username: "newUsername123"}, nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.UPDATED,
					Message: "updated",
					Data:    &dto.UpdateUsernameResponse{Username: "newUsername123"},
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
			c.Request, _ = http.NewRequest("PUT", "/v1/users/usernames", test.MakeRequestBody(tc.input.request))

			h.UpdateUsername(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

func TestSignOut(t *testing.T) {
	type input struct {
		dto         *dto.UserLogoutRequest
		beforeTests func(mockUserService *mocks.UserService)
		err         error
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
				dto: &dto.UserLogoutRequest{
					RefreshToken: "",
				},
				beforeTests: func(mockUserService *mocks.UserService) {},
				err:         nil,
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "RefreshToken is required",
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to logout",
			input: input{
				dto: &dto.UserLogoutRequest{
					RefreshToken: "token",
				},
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("SignOut", &dto.UserLogoutRequest{
						RefreshToken: "token",
						AccessToken:  "token",
						UserId:       1,
					}).Return(errs.ErrInternalServerError)
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
		{
			description: "should return success with status code 200 when logout succeed",
			input: input{
				dto: &dto.UserLogoutRequest{
					RefreshToken: "token",
				},
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("SignOut", &dto.UserLogoutRequest{
						RefreshToken: "token",
						AccessToken:  "token",
						UserId:       1,
					}).Return(nil)
				},
				err: nil,
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)

			c.Set("userId", 1)

			userService := mocks.NewUserService(t)
			tc.input.beforeTests(userService)
			cfg := handler.HandlerConfig{
				UserService: userService,
			}
			h := handler.New(&cfg)
			c.Request, _ = http.NewRequest("POST", "/v1/users/logout", test.MakeRequestBody(tc.input.dto))
			c.Request.Header.Set("authorization", "Bearer "+"token")

			h.SignOut(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}

}

func TestRequestPasswordChange(t *testing.T) {
	type input struct {
		request     *dto.RequestPasswordChangeRequest
		beforeTests func(mockUserService *mocks.UserService)
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
				request: &dto.RequestPasswordChangeRequest{
					CurrentPassword: "",
					NewPassword:     "asdasdasdsa",
				},
				beforeTests: func(mockUserService *mocks.UserService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "CurrentPassword is required",
				},
			},
		},
		{
			description: "should return ErrUserDoesNotExist and status not found when RequestPasswordChange failed",
			input: input{
				request: &dto.RequestPasswordChangeRequest{
					CurrentPassword: "asdasdasd",
					NewPassword:     "asdasdasdsa",
				},
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("RequestPasswordChange", &dto.RequestPasswordChangeRequest{
						CurrentPassword: "asdasdasd",
						NewPassword:     "asdasdasdsa",
						UserId:          1,
					}).Return(errs.ErrUserDoesNotExist)
				},
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
			description: "should return ErrInvalidCredential and status Bad Request when RequestPasswordChange failed",
			input: input{
				request: &dto.RequestPasswordChangeRequest{
					CurrentPassword: "asdasdasd",
					NewPassword:     "asdasdasdsa",
				},
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("RequestPasswordChange", &dto.RequestPasswordChangeRequest{
						CurrentPassword: "asdasdasd",
						NewPassword:     "asdasdasdsa",
						UserId:          1,
					}).Return(errs.ErrInvalidCredential)
				},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.WRONG_PASSWORD,
					Message: errs.ErrInvalidCredential.Error(),
				},
			},
		},
		{
			description: "should return ErrSamePassword and status Bad Request when RequestPasswordChange failed",
			input: input{
				request: &dto.RequestPasswordChangeRequest{
					CurrentPassword: "asdasdasd",
					NewPassword:     "asdasdasdsa",
				},
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("RequestPasswordChange", &dto.RequestPasswordChangeRequest{
						CurrentPassword: "asdasdasd",
						NewPassword:     "asdasdasdsa",
						UserId:          1,
					}).Return(errs.ErrSamePassword)
				},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.SAME_PASSWORD,
					Message: errs.ErrSamePassword.Error(),
				},
			},
		},
		{
			description: "should return ErrInvalidPasswordPattern and status StatusUnprocessableEntity when RequestPasswordChange failed",
			input: input{
				request: &dto.RequestPasswordChangeRequest{
					CurrentPassword: "asdasdasd",
					NewPassword:     "asdasdasdsa",
				},
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("RequestPasswordChange", &dto.RequestPasswordChangeRequest{
						CurrentPassword: "asdasdasd",
						NewPassword:     "asdasdasdsa",
						UserId:          1,
					}).Return(errs.ErrInvalidPasswordPattern)
				},
			},
			expected: expected{
				statusCode: http.StatusUnprocessableEntity,
				response: response.Response{
					Code:    code.INVALID_PASSWORD_PATTERN,
					Message: errs.ErrInvalidPasswordPattern.Error(),
				},
			},
		},
		{
			description: "should return ErrInternalServerError and status InternalServerError when RequestPasswordChange failed",
			input: input{
				request: &dto.RequestPasswordChangeRequest{
					CurrentPassword: "asdasdasd",
					NewPassword:     "asdasdasdsa",
				},
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("RequestPasswordChange", &dto.RequestPasswordChangeRequest{
						CurrentPassword: "asdasdasd",
						NewPassword:     "asdasdasdsa",
						UserId:          1,
					}).Return(errs.ErrInternalServerError)
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
			description: "should return success and status OK when RequestPasswordChange success",
			input: input{
				request: &dto.RequestPasswordChangeRequest{
					CurrentPassword: "asdasdasd",
					NewPassword:     "asdasdasdsa",
				},
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("RequestPasswordChange", &dto.RequestPasswordChangeRequest{
						CurrentPassword: "asdasdasd",
						NewPassword:     "asdasdasdsa",
						UserId:          1,
					}).Return(nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", 1)
			userService := mocks.NewUserService(t)
			tc.input.beforeTests(userService)
			cfg := handler.HandlerConfig{
				UserService: userService,
			}
			h := handler.New(&cfg)
			c.Request, _ = http.NewRequest("POST", "/v1/users/request-password-change", test.MakeRequestBody(tc.input.request))

			h.RequestPasswordChange(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}

}

func TestCompletePasswordChange(t *testing.T) {
	type input struct {
		request     *dto.CompletePasswordChangeRequest
		beforeTests func(mockUserService *mocks.UserService)
	}
	type expected struct {
		statusCode int
		response   response.Response
	}
	cases := []struct {
		description string
		input       input
		expected    expected
	}{
		{
			description: "should return error bad request when request is nil",
			input: input{
				request: &dto.CompletePasswordChangeRequest{
					VerificationCode: "",
				},
				beforeTests: func(mockUserService *mocks.UserService) {
				},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "VerificationCode is required",
				},
			},
		},
		{
			description: "should return error ErrIncorrectVerificationCode and status Bad Request when CompletePasswordChange failed",
			input: input{
				request: &dto.CompletePasswordChangeRequest{
					VerificationCode: "asdasdasd",
				},
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("CompletePasswordChange", &dto.CompletePasswordChangeRequest{
						VerificationCode: "asdasdasd",
						UserId:           1,
					}).Return(errs.ErrIncorrectVerificationCode)
				},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.INCORRECT_VERIFICATION_CODE,
					Message: errs.ErrIncorrectVerificationCode.Error(),
				},
			},
		},
		{
			description: "should return ErrVerficationCodeNotFound and status Bad Request when CompletePasswordChange failed",
			input: input{
				request: &dto.CompletePasswordChangeRequest{
					VerificationCode: "asdasdasd",
				},
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("CompletePasswordChange", &dto.CompletePasswordChangeRequest{
						VerificationCode: "asdasdasd",
						UserId:           1,
					}).Return(errs.ErrVerficationCodeNotFound)
				},
			},
			expected: expected{
				statusCode: http.StatusNotFound,
				response: response.Response{
					Code:    code.NOT_FOUND,
					Message: errs.ErrVerficationCodeNotFound.Error(),
				},
			},
		},
		{
			description: "should return error and status Internal Server Error when CompletePasswordChange failed",
			input: input{
				request: &dto.CompletePasswordChangeRequest{
					VerificationCode: "asdasdasd",
				},
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("CompletePasswordChange", &dto.CompletePasswordChangeRequest{
						VerificationCode: "asdasdasd",
						UserId:           1,
					}).Return(errs.ErrInternalServerError)
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
			description: "should return ok and status OK when CompletePasswordChange success",
			input: input{
				request: &dto.CompletePasswordChangeRequest{
					VerificationCode: "asdasdasd",
				},
				beforeTests: func(mockUserService *mocks.UserService) {
					mockUserService.On("CompletePasswordChange", &dto.CompletePasswordChangeRequest{
						VerificationCode: "asdasdasd",
						UserId:           1,
					}).Return(nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.OK,
					Message: "ok",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			expectedRes, _ := json.Marshal(tc.expected.response)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Set("userId", 1)
			userService := mocks.NewUserService(t)
			tc.input.beforeTests(userService)
			cfg := handler.HandlerConfig{
				UserService: userService,
			}
			h := handler.New(&cfg)
			c.Request, _ = http.NewRequest("POST", "/v1/users/complete-password-change", test.MakeRequestBody(tc.input.request))

			h.CompletePasswordChange(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}

}

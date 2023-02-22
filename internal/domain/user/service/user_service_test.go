package service_test

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/service"
	"kedai/backend/be-kedai/internal/utils/hash"
	mocks "kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	type input struct {
		user *model.User
		dto  *dto.UserRegistration
		err  error
	}

	type expected struct {
		user *model.User
		dto  *dto.UserRegistrationResponse
		err  error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return created user data when called",
			input: input{
				user: &model.User{
					Email: "user@mail.com",
					Password: "Password2",
				},
				dto: &dto.UserRegistration{
					Email: "user@mail.com",
					Password: "Password2",
				},
				err: nil,
			},
			expected: expected{
				user: &model.User{
					Email: "user@mail.com",
				},
				dto: &dto.UserRegistrationResponse{
					Email: "user@mail.com",
				},
				err: nil,
			},
		},
		{
			description: "should return error when server error",
			input: input{
				user: &model.User{
					Email: "user@mail.com",
					Password: "Password1",
				},
				dto: &dto.UserRegistration{
					Email: "user@mail.com",
					Password: "Password1",
				},
				err: errors.New("server internal error"),
			},
			expected: expected{
				user: nil,
				dto:  nil,
				err:  errors.New("server internal error"),
			},
		},
		{
			description: "should return error when invalid password pattern",
			input: input{
				user: &model.User{
					Email: "user@mail.com",
					Password: "Password",
				},
				dto: &dto.UserRegistration{
					Email: "user@mail.com",
					Password: "Password",
				},
				err: errs.ErrInvalidPasswordPattern,
			},
			expected: expected{
				user: nil,
				dto:  nil,
				err:  errs.ErrInvalidPasswordPattern,
			},
		},
		{
			description: "should return error when password contain email address",
			input: input{
				user: &model.User{
					Email: "user@mail.com",
					Password: "Password1user",
				},
				dto: &dto.UserRegistration{
					Email: "user@mail.com",
					Password: "Password1user",
				},
				err: errs.ErrContainEmail,
			},
			expected: expected{
				user: nil,
				dto:  nil,
				err:  errs.ErrContainEmail,
			},
		},
		{
			description: "should return error when registering same user 2 times",
			input: input{
				user: &model.User{
					Email: "user@mail.com",
					Password: "Password1",
				},
				dto: &dto.UserRegistration{
					Email: "user@mail.com",
					Password: "Password1",
				},
				err: errs.ErrUserAlreadyExist,
			},
			expected: expected{
				user: nil,
				dto:  nil,
				err:  errs.ErrUserAlreadyExist,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)
			service := service.NewUserService(&service.UserSConfig{
				Repository: mockRepo,
			})
			mockRepo.On("SignUp", tc.input.user).Return(tc.expected.user, tc.expected.err)

			result, err := service.SignUp(tc.input.dto)

			assert.Equal(t, tc.expected.dto, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestSignIn(t *testing.T) {
	t.Run("should return error when invalid credential", func(t *testing.T) {
		hashedPw, _ := hash.HashAndSalt("password")
		user := &model.User{
			Email:    "user@mail.com",
			Password: "password1",
		}
		dto := &dto.UserLogin{
			Email:    "user@mail.com",
			Password: "password1",
		}
		expectedUser := &model.User{
			Email:    "user@mail.com",
			Password: hashedPw,
		}
		mockRepo := new(mocks.UserRepository)
		service := service.NewUserService(&service.UserSConfig{
			Repository: mockRepo,
		})
		mockRepo.On("SignIn", user).Return(expectedUser, nil)

		_, err := service.SignIn(dto, dto.Password)

		assert.Error(t, errs.ErrInvalidCredential, err)
	})

	type input struct {
		user *model.User
		dto  *dto.UserLogin
		err  error
	}

	type expected struct {
		user *model.User
		dto  *dto.Token
		err  error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return error when user input invalid credential",
			input: input{
				user: &model.User{
					Email:    "user@mail.com",
					Password: "password",
				},
				dto: &dto.UserLogin{
					Email:    "user@mail.com",
					Password: "password",
				},
				err: errs.ErrInvalidCredential,
			},
			expected: expected{
				user: nil,
				dto:  nil,
				err:  errs.ErrInvalidCredential,
			},
		},
		{
			description: "should return error when internal server error",
			input: input{
				user: &model.User{
					Email:    "user@mail.com",
					Password: "password",
				},
				dto: &dto.UserLogin{
					Email:    "user@mail.com",
					Password: "password",
				},
				err: errs.ErrInternalServerError,
			},
			expected: expected{
				user: nil,
				dto:  nil,
				err:  errs.ErrInternalServerError,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)
			service := service.NewUserService(&service.UserSConfig{
				Repository: mockRepo,
			})
			mockRepo.On("SignIn", tc.input.user).Return(tc.expected.user, tc.expected.err)

			result, err := service.SignIn(tc.input.dto, tc.input.dto.Password)

			assert.Equal(t, tc.expected.dto, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestGetByID(t *testing.T) {
	type input struct {
		id   int
		data *model.User
		err  error
	}
	type expected struct {
		user *model.User
		err  error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "it should return user data if user exists",
			input: input{
				id: 1,
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
				user: &model.User{
					Email:    "user@email.com",
					Username: "user_name",
					Profile: &model.UserProfile{
						UserID: 1,
					},
				},
				err: nil,
			},
		},
		{
			description: "it should return error if failed to get user",
			input: input{
				id:   1,
				data: nil,
				err:  errors.New("failed to get user"),
			},
			expected: expected{
				user: nil,
				err:  errors.New("failed to get user"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := mocks.NewUserRepository(t)
			mockRepo.On("GetByID", tc.input.id).Return(tc.input.data, tc.input.err)
			uc := service.NewUserService(&service.UserSConfig{
				Repository: mockRepo,
			})

			actualUser, actualErr := uc.GetByID(tc.input.id)

			assert.Equal(t, tc.expected.user, actualUser)
			assert.Equal(t, actualErr, tc.expected.err)
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
		err error
	}

	type cases struct {
		description string
		input
		expected
	}

	for _, tc := range []cases{
		{
			description: "should return nil error when session available",
			input: input{
				userId: 1,
				token:  "token",
				err:    nil,
			},
			expected: expected{
				err: nil,
			},
		},
		{
			description: "should return error when session unavailable",
			input: input{
				userId: 1,
				token:  "token",
				err:    errors.New("error"),
			},
			expected: expected{
				err: errors.New("error"),
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRedis := new(mocks.UserCache)
			service := service.NewUserService(&service.UserSConfig{
				Redis: mockRedis,
			})
			mockRedis.On("FindToken", tc.input.userId, tc.input.token).Return(tc.expected.err)

			result := service.GetSession(tc.input.userId, tc.input.token)

			assert.Equal(t, tc.expected.err, result)
		})
	}
}

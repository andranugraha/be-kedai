package service_test

import (
	"errors"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/service"
	errs "kedai/backend/be-kedai/internal/common/error"
	mocks "kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	type input struct {
		user *model.User
		dto *dto.UserRegistration
		err error
	}

	type expected struct {
		user *model.User
		err error
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
				},
				dto: &dto.UserRegistration{
					Email: "user@mail.com",
				},
				err: nil,
			},
			expected: expected{
				user: &model.User{
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
				},
				dto: &dto.UserRegistration{
					Email: "user@mail.com",
				},
				err: errors.New("server internal error"),
			},
			expected: expected{
				user: nil,
				err: errors.New("server internal error"),
			},
		},
		{
			description: "should return error when registering same user 2 times",
			input: input{
				user: &model.User{
					Email: "user@mail.com",
				},
				dto: &dto.UserRegistration{
					Email: "user@mail.com",
				},
				err: errs.ErrUserAlreadyExist,
			},
			expected: expected{
				user: nil,
				err: errs.ErrUserAlreadyExist,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			mockRepo := new(mocks.UserRepository)
			service := service.NewUserService(&service.UserSConfig{
				Repository: mockRepo,
			})
			mockRepo.On("SignUp", tc.input.user).Return(tc.expected.user, tc.expected.err)

			result, err := service.SignUp(tc.dto)

			assert.Equal(t, tc.expected.user, result)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}
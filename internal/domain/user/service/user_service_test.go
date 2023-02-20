package service_test

import (
	"errors"
	"kedai/backend/be-kedai/internal/domain/user/model"

	"kedai/backend/be-kedai/internal/domain/user/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_GetByID(t *testing.T) {
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

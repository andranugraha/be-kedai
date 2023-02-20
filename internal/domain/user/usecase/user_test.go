package usecase_test

import (
	"errors"
	"kedai/backend/be-kedai/internal/domain/user/entity"
	"kedai/backend/be-kedai/internal/domain/user/usecase"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserUsecase_GetByID(t *testing.T) {
	type input struct {
		id   int
		data *entity.User
		err  error
	}
	type expected struct {
		user *entity.User
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
				data: &entity.User{
					Email:    "user@email.com",
					Username: "user_name",
					Profile: &entity.UserProfile{
						UserID: 1,
					},
				},
				err: nil,
			},
			expected: expected{
				user: &entity.User{
					Email:    "user@email.com",
					Username: "user_name",
					Profile: &entity.UserProfile{
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
			uc := usecase.NewUserUsecase(&usecase.UserUConfig{
				Repository: mockRepo,
			})

			actualUser, actualErr := uc.GetByID(tc.input.id)

			assert.Equal(t, tc.expected.user, actualUser)
			assert.Equal(t, actualErr, tc.expected.err)
		})
	}
}

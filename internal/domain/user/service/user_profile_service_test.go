package service_test

import (
	"errors"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/service"
	"kedai/backend/be-kedai/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateProfile(t *testing.T) {
	type input struct {
		userId     int
		request    *dto.UpdateProfileRequest
		mockReturn *model.UserProfile
		mockErr    error
	}
	type expected struct {
		res *dto.UpdateProfileResponse
		err error
	}

	var (
		dob, _      = time.Parse("2006-01-02", "2006-01-02")
		name        = "new name"
		phoneNumber = "0123456789"
		gender      = "others"
		photoUrl    = "http://photo.url/example.png"
	)

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when failed to update user profile",
			input: input{
				userId: 1,
				request: &dto.UpdateProfileRequest{
					Name:        "new name",
					PhoneNumber: "0123456789",
					DoB:         "2006-01-02",
					Gender:      "others",
					PhotoUrl:    "http://photo.url/example.png",
				},
				mockReturn: nil,
				mockErr:    errors.New("failed to update user"),
			},
			expected: expected{
				res: nil,
				err: errors.New("failed to update user"),
			},
		},
		{
			description: "should return updated profile when update profile successed",
			input: input{
				userId: 1,
				request: &dto.UpdateProfileRequest{
					Name:        "new name",
					PhoneNumber: "0123456789",
					DoB:         "2006-01-02",
					Gender:      "others",
					PhotoUrl:    "http://photo.url/example.png",
				},
				mockReturn: &model.UserProfile{
					ID:          1,
					Name:        &name,
					PhoneNumber: &phoneNumber,
					DoB:         &dob,
					Gender:      &gender,
					PhotoUrl:    &photoUrl,
				},
				mockErr: nil,
			},
			expected: expected{
				res: &dto.UpdateProfileResponse{
					ID:          1,
					Name:        name,
					PhoneNumber: phoneNumber,
					DoB:         dob,
					Gender:      gender,
					PhotoUrl:    photoUrl,
				},
				err: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			profileRepo := mocks.NewUserProfileRepository(t)
			profileRepo.On("Update", tc.input.userId, tc.input.request.ToUserProfile()).Return(tc.input.mockReturn, tc.input.mockErr)
			profileService := service.NewUserProfileService(&service.UserProfileSConfig{
				Repository: profileRepo,
			})

			actualRes, actualErr := profileService.UpdateProfile(tc.input.userId, tc.input.request)

			assert.Equal(t, tc.expected.res, actualRes)
			assert.Equal(t, tc.expected.err, actualErr)
		})
	}
}

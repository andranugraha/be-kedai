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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProfile(t *testing.T) {
	type input struct {
		userId     int
		request    *dto.UpdateProfileRequest
		beforeTest func(*mocks.UserProfileService)
	}
	type expected struct {
		statusCode int
		response   response.Response
	}

	var (
		request = dto.UpdateProfileRequest{
			Name:        "new name",
			PhoneNumber: "0123456789",
			DoB:         "2006-01-02",
			Gender:      "male",
			PhotoUrl:    "http://photo.url/example.png",
		}
		mockResponse = dto.UpdateProfileResponse{
			ID:          1,
			Name:        "new name",
			PhoneNumber: "0123456789",
			DoB:         time.Now(),
			Gender:      "male",
			PhotoUrl:    "http://photo.url/example.png",
		}
	)

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error with status code 400 when given invalid request body",
			input: input{
				userId: 1,
				request: &dto.UpdateProfileRequest{
					Name: "new name",
				},
				beforeTest: func(ups *mocks.UserProfileService) {},
			},
			expected: expected{
				statusCode: http.StatusBadRequest,
				response: response.Response{
					Code:    code.BAD_REQUEST,
					Message: "PhotoUrl is required",
				},
			},
		},
		{
			description: "should return error with status code 500 when failed to update user profile",
			input: input{
				userId:  1,
				request: &request,
				beforeTest: func(ups *mocks.UserProfileService) {
					ups.On("UpdateProfile", 1, &request).Return(nil, errors.New("failed to update profile"))
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
			description: "should return updated profile data with status code 200 when successed to update user profile",
			input: input{
				userId:  1,
				request: &request,
				beforeTest: func(ups *mocks.UserProfileService) {
					ups.On("UpdateProfile", 1, &request).Return(&mockResponse, nil)
				},
			},
			expected: expected{
				statusCode: http.StatusOK,
				response: response.Response{
					Code:    code.UPDATED,
					Message: "updated",
					Data:    &mockResponse,
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
			userProfileService := mocks.NewUserProfileService(t)
			tc.input.beforeTest(userProfileService)
			cfg := handler.HandlerConfig{
				UserProfileService: userProfileService,
			}
			h := handler.New(&cfg)
			c.Request, _ = http.NewRequest("PUT", "/v1/users/profiles", test.MakeRequestBody(tc.input.request))

			h.UpdateProfile(c)

			assert.Equal(t, tc.expected.statusCode, rec.Code)
			assert.Equal(t, string(expectedRes), rec.Body.String())
		})
	}
}

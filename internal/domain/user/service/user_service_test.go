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

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUp(t *testing.T) {
	type input struct {
		user *model.User
		dto  *dto.UserRegistrationRequest
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
					Email:    "user@mail.com",
					Password: "Password2",
				},
				dto: &dto.UserRegistrationRequest{
					Email:    "user@mail.com",
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
					Email:    "user@mail.com",
					Password: "Password1",
				},
				dto: &dto.UserRegistrationRequest{
					Email:    "user@mail.com",
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
					Email:    "user@mail.com",
					Password: "Password",
				},
				dto: &dto.UserRegistrationRequest{
					Email:    "user@mail.com",
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
					Email:    "user@mail.com",
					Password: "Password1user",
				},
				dto: &dto.UserRegistrationRequest{
					Email:    "user@mail.com",
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
					Email:    "user@mail.com",
					Password: "Password1",
				},
				dto: &dto.UserRegistrationRequest{
					Email:    "user@mail.com",
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

func TestRenewToken(t *testing.T) {
	type input struct {
		userId       int
		refreshToken string
		beforeTest   func(*mocks.UserCache)
	}
	type expected struct {
		token *dto.Token
		err   error
	}

	tests := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when refresh token is expired or does not exist",
			input: input{
				userId:       1,
				refreshToken: "token",
				beforeTest: func(uc *mocks.UserCache) {
					uc.On("FindToken", 1, "token").Return(redis.Nil)
				},
			},
			expected: expected{
				token: nil,
				err:   errs.ErrExpiredToken,
			},
		},
		{
			description: "should return error when failed to fetch token from redis",
			input: input{
				userId:       1,
				refreshToken: "token",
				beforeTest: func(uc *mocks.UserCache) {
					uc.On("FindToken", 1, "token").Return(errors.New("failed to check token"))
				},
			},
			expected: expected{
				token: nil,
				err:   errors.New("failed to check token"),
			},
		},
		{
			description: "should return error when failed to store renewed tokens",
			input: input{
				userId:       1,
				refreshToken: "token",
				beforeTest: func(uc *mocks.UserCache) {
					uc.On("FindToken", 1, "token").Return(nil)
					uc.On("DeleteToken", "user_1:token").Return(errors.New("failed to delete token"))
				},
			},
			expected: expected{
				token: nil,
				err:   errors.New("failed to delete token"),
			},
		},
		{
			description: "should return error when failed to store renewed tokens",
			input: input{
				userId:       1,
				refreshToken: "token",
				beforeTest: func(uc *mocks.UserCache) {
					uc.On("FindToken", 1, "token").Return(nil)
					uc.On("DeleteToken", "user_1:token").Return(nil)
					uc.On("StoreToken", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("failed to store token"))
				},
			},
			expected: expected{
				token: nil,
				err:   errors.New("failed to store token"),
			},
		},
	}

	for _, tc := range tests {
		userCache := mocks.NewUserCache(t)
		tc.beforeTest(userCache)
		userService := service.NewUserService(&service.UserSConfig{
			Redis: userCache,
		})

		actualToken, actualErr := userService.RenewToken(tc.input.userId, tc.input.refreshToken)

		assert.Equal(t, tc.expected.token, actualToken)
		assert.Equal(t, tc.expected.err, actualErr)
	}
}

func TestUpdateEmail(t *testing.T) {
	type input struct {
		userId     int
		request    *dto.UpdateEmailRequest
		beforeTest func(*mocks.UserRepository)
	}
	type expected struct {
		res *dto.UpdateEmailResponse
		err error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when email is used",
			input: input{
				userId: 1,
				request: &dto.UpdateEmailRequest{
					Email: "used.email@email.com",
				},
				beforeTest: func(ur *mocks.UserRepository) {
					ur.On("GetByEmail", "used.email@email.com").Return(&model.User{Email: "used.email@email.com"}, nil)
				},
			},
			expected: expected{
				res: nil,
				err: errs.ErrEmailUsed,
			},
		},
		{
			description: "should return error when failed to check if email is used or not",
			input: input{
				userId: 1,
				request: &dto.UpdateEmailRequest{
					Email: "used.email@email.com",
				},
				beforeTest: func(ur *mocks.UserRepository) {
					ur.On("GetByEmail", "used.email@email.com").Return(nil, errors.New("failed to check email"))
				},
			},
			expected: expected{
				res: nil,
				err: errors.New("failed to check email"),
			},
		},
		{
			description: "should return error when failed to update email",
			input: input{
				userId: 1,
				request: &dto.UpdateEmailRequest{
					Email: "new.email@email.com",
				},
				beforeTest: func(ur *mocks.UserRepository) {
					ur.On("GetByEmail", "new.email@email.com").Return(nil, errs.ErrUserDoesNotExist)
					ur.On("UpdateEmail", 1, "new.email@email.com").Return(nil, errors.New("failed to update email"))
				},
			},
			expected: expected{
				res: nil,
				err: errors.New("failed to update email"),
			},
		},
		{
			description: "should return updated email when update email successed",
			input: input{
				userId: 1,
				request: &dto.UpdateEmailRequest{
					Email: "new.email@email.com",
				},
				beforeTest: func(ur *mocks.UserRepository) {
					ur.On("GetByEmail", "new.email@email.com").Return(nil, errs.ErrUserDoesNotExist)
					ur.On("UpdateEmail", 1, "new.email@email.com").Return(&model.User{Email: "new.email@email.com"}, nil)
				},
			},
			expected: expected{
				res: &dto.UpdateEmailResponse{
					Email: "new.email@email.com",
				},
				err: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			userRepo := mocks.NewUserRepository(t)
			tc.beforeTest(userRepo)
			userService := service.NewUserService(&service.UserSConfig{
				Repository: userRepo,
			})

			res, err := userService.UpdateEmail(tc.input.userId, tc.input.request)

			assert.Equal(t, tc.expected.res, res)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func TestUpdateUsername(t *testing.T) {
	type input struct {
		userId     int
		request    *dto.UpdateUsernameRequest
		beforeTest func(*mocks.UserRepository)
	}
	type expected struct {
		res *dto.UpdateUsernameResponse
		err error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when username is not valid",
			input: input{
				userId: 1,
				request: &dto.UpdateUsernameRequest{
					Username: "new_u$ername",
				},
				beforeTest: func(ur *mocks.UserRepository) {},
			},
			expected: expected{
				res: nil,
				err: errs.ErrInvalidUsernamePattern,
			},
		},
		{
			description: "should return new username when update username successed",
			input: input{
				userId: 1,
				request: &dto.UpdateUsernameRequest{
					Username: "new_username",
				},
				beforeTest: func(ur *mocks.UserRepository) {
					ur.On("UpdateUsername", 1, "new_username").Return(&model.User{Username: "new_username"}, nil)
				},
			},
			expected: expected{
				res: &dto.UpdateUsernameResponse{Username: "new_username"},
				err: nil,
			},
		},
		{
			description: "should return error if failed to update username",
			input: input{
				userId: 1,
				request: &dto.UpdateUsernameRequest{
					Username: "new_username",
				},
				beforeTest: func(ur *mocks.UserRepository) {
					ur.On("UpdateUsername", 1, "new_username").Return(nil, errors.New("failed to update username"))
				},
			},
			expected: expected{
				res: nil,
				err: errors.New("failed to update username"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			userRepo := mocks.NewUserRepository(t)
			tc.beforeTest(userRepo)
			userService := service.NewUserService(&service.UserSConfig{
				Repository: userRepo,
			})

			actualRes, actualErr := userService.UpdateUsername(tc.input.userId, tc.input.request)

			assert.Equal(t, tc.expected.res, actualRes)
			assert.Equal(t, tc.expected.err, actualErr)
		})
	}
}

func TestSignOut(t *testing.T) {
	type input struct {
		data       dto.UserLogoutRequest
		beforeTest func(*mocks.UserCache)
	}
	type expected struct {
		err error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error if failed to sign out",
			input: input{
				data: dto.UserLogoutRequest{
					UserId:       1,
					RefreshToken: "refresh_token",
					AccessToken:  "access_token",
				},
				beforeTest: func(ur *mocks.UserCache) {
					ur.On("DeleteRefreshTokenAndAccessToken", 1, "refresh_token", "access_token").Return(errors.New("failed to sign out"))
				},
			},
			expected: expected{
				err: errors.New("failed to sign out"),
			},
		},
		{
			description: "should return nil if sign out successed",
			input: input{
				data: dto.UserLogoutRequest{
					UserId:       1,
					RefreshToken: "refresh_token",
					AccessToken:  "access_token",
				},
				beforeTest: func(ur *mocks.UserCache) {
					ur.On("DeleteRefreshTokenAndAccessToken", 1, "refresh_token", "access_token").Return(nil)
				},
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			userCache := mocks.NewUserCache(t)
			tc.beforeTest(userCache)
			userService := service.NewUserService(&service.UserSConfig{
				Redis: userCache,
			})

			actualErr := userService.SignOut(&tc.input.data)

			assert.Equal(t, tc.expected.err, actualErr)
		})
	}

}

func TestRequestPasswordChange(t *testing.T) {
	hashedPassword, _ := hash.HashAndSalt("Passwrod123")
	type input struct {
		request    *dto.RequestPasswordChangeRequest
		beforeTest func(*mocks.UserRepository, *mocks.MailUtils, *mocks.RandomUtils, *mocks.UserCache)
	}
	type expected struct {
		err error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when GetByID failed",
			input: input{
				request: &dto.RequestPasswordChangeRequest{
					UserId:          1,
					CurrentPassword: "Passwrod123",
					NewPassword:     "Passwrod1234",
				},
				beforeTest: func(ur *mocks.UserRepository, mu *mocks.MailUtils, ru *mocks.RandomUtils, uc *mocks.UserCache) {
					ur.On("GetByID", 1).Return(nil, errs.ErrUserDoesNotExist)

				},
			},
			expected: expected{
				err: errs.ErrUserDoesNotExist,
			},
		},
		{
			description: "should return error when StoreUserPasswordAndVerificationCode failed",
			input: input{
				request: &dto.RequestPasswordChangeRequest{
					UserId:          1,
					CurrentPassword: "Passwrod123",
					NewPassword:     "Passwrod1234",
				},
				beforeTest: func(ur *mocks.UserRepository, mu *mocks.MailUtils, ru *mocks.RandomUtils, uc *mocks.UserCache) {
					ur.On("GetByID", 1).Return(&model.User{ID: 1, Password: hashedPassword, Username: "test"}, nil)
					ru.On("GenerateAlphanumericString", mock.Anything).Return("code", nil)
					uc.On("StoreUserPasswordAndVerificationCode", 1, mock.Anything, mock.Anything).Return(errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				err: errs.ErrUserDoesNotExist,
			},
		},
		{
			description: "should return error when SendUpdatePasswordEmail failed",
			input: input{
				request: &dto.RequestPasswordChangeRequest{
					UserId:          1,
					CurrentPassword: "Passwrod123",
					NewPassword:     "Passwrod1234",
				},
				beforeTest: func(ur *mocks.UserRepository, mu *mocks.MailUtils, ru *mocks.RandomUtils, uc *mocks.UserCache) {
					ur.On("GetByID", 1).Return(&model.User{ID: 1, Password: hashedPassword, Username: "test"}, nil)
					ru.On("GenerateAlphanumericString", mock.Anything).Return("code", nil)
					uc.On("StoreUserPasswordAndVerificationCode", 1, mock.Anything, mock.Anything).Return(nil)
					mu.On("SendUpdatePasswordEmail", mock.Anything, mock.Anything, mock.Anything).Return(errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				err: errs.ErrUserDoesNotExist,
			},
		},
		{
			description: "should return nil when success",
			input: input{
				request: &dto.RequestPasswordChangeRequest{
					UserId:          1,
					CurrentPassword: "Passwrod123",
					NewPassword:     "Passwrod1234",
				},
				beforeTest: func(ur *mocks.UserRepository, mu *mocks.MailUtils, ru *mocks.RandomUtils, uc *mocks.UserCache) {
					ur.On("GetByID", 1).Return(&model.User{ID: 1, Password: hashedPassword, Username: "test"}, nil)
					ru.On("GenerateAlphanumericString", mock.Anything).Return("code", nil)
					uc.On("StoreUserPasswordAndVerificationCode", 1, mock.Anything, mock.Anything).Return(nil)
					mu.On("SendUpdatePasswordEmail", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				},
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			userRepo := mocks.NewUserRepository(t)
			mailUtils := mocks.NewMailUtils(t)
			randomUtils := mocks.NewRandomUtils(t)
			userCache := mocks.NewUserCache(t)
			tc.beforeTest(userRepo, mailUtils, randomUtils, userCache)
			userService := service.NewUserService(&service.UserSConfig{
				Repository:  userRepo,
				MailUtils:   mailUtils,
				RandomUtils: randomUtils,
				Redis:       userCache,
			})

			actualErr := userService.RequestPasswordChange(tc.input.request)

			assert.Equal(t, tc.expected.err, actualErr)
		})
	}

}

func TestCompletePasswordChange(t *testing.T) {
	var (
		verifcationCode      = "code"
		newPassword          = "Passwrod1234"
		wrongVerifcationCode = "wrongCode"
	)
	type input struct {
		request    *dto.CompletePasswordChangeRequest
		beforeTest func(*mocks.UserRepository, *mocks.UserCache)
	}
	type expected struct {
		err error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when FindUserPasswordAndVerificationCode failed",
			input: input{
				request: &dto.CompletePasswordChangeRequest{
					UserId:           1,
					VerificationCode: "code",
				},
				beforeTest: func(ur *mocks.UserRepository, uc *mocks.UserCache) {
					uc.On("FindUserPasswordAndVerificationCode", 1).Return(newPassword, verifcationCode, errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				err: errs.ErrUserDoesNotExist,
			},
		},
		{
			description: "should return error when verification code is wrong",
			input: input{
				request: &dto.CompletePasswordChangeRequest{
					UserId:           1,
					VerificationCode: wrongVerifcationCode,
				},
				beforeTest: func(ur *mocks.UserRepository, uc *mocks.UserCache) {
					uc.On("FindUserPasswordAndVerificationCode", 1).Return(newPassword, verifcationCode, nil)
				},
			},
			expected: expected{
				err: errs.ErrIncorrectVerificationCode,
			},
		},
		{
			description: "should return error when UpdatePassword failed",
			input: input{
				request: &dto.CompletePasswordChangeRequest{
					UserId:           1,
					VerificationCode: verifcationCode,
				},
				beforeTest: func(ur *mocks.UserRepository, uc *mocks.UserCache) {
					uc.On("FindUserPasswordAndVerificationCode", 1).Return(newPassword, verifcationCode, nil)
					ur.On("UpdatePassword", 1, mock.Anything).Return(nil, errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				err: errs.ErrUserDoesNotExist,
			},
		},

		{
			description: "should return error when DeleteUserPasswordAndVerificationCode failed",
			input: input{
				request: &dto.CompletePasswordChangeRequest{
					UserId:           1,
					VerificationCode: verifcationCode,
				},
				beforeTest: func(ur *mocks.UserRepository, uc *mocks.UserCache) {
					uc.On("FindUserPasswordAndVerificationCode", 1).Return(newPassword, verifcationCode, nil)
					ur.On("UpdatePassword", 1, mock.Anything).Return(nil, nil)
					uc.On("DeleteUserPasswordAndVerificationCode", 1).Return(errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				err: errs.ErrUserDoesNotExist,
			},
		},
		{
			description: "should return nil when success",
			input: input{
				request: &dto.CompletePasswordChangeRequest{
					UserId:           1,
					VerificationCode: verifcationCode,
				},
				beforeTest: func(ur *mocks.UserRepository, uc *mocks.UserCache) {
					uc.On("FindUserPasswordAndVerificationCode", 1).Return(newPassword, verifcationCode, nil)
					ur.On("UpdatePassword", 1, mock.Anything).Return(nil, nil)
					uc.On("DeleteUserPasswordAndVerificationCode", 1).Return(nil)
				},
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			userRepo := mocks.NewUserRepository(t)
			userCache := mocks.NewUserCache(t)
			tc.beforeTest(userRepo, userCache)
			userService := service.NewUserService(&service.UserSConfig{
				Repository: userRepo,
				Redis:      userCache,
			})

			actualErr := userService.CompletePasswordChange(tc.input.request)

			assert.Equal(t, tc.expected.err, actualErr)
		})
	}

}

func TestRequestPasswordReset(t *testing.T) {
	var (
		email = "email"
	)
	type input struct {
		request    *dto.RequestPasswordResetRequest
		beforeTest func(*mocks.UserRepository, *mocks.UserCache, *mocks.MailUtils, *mocks.RandomUtils)
	}
	type expected struct {
		err error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when FindUserByEmail failed",
			input: input{
				request: &dto.RequestPasswordResetRequest{
					Email: email,
				},
				beforeTest: func(ur *mocks.UserRepository, uc *mocks.UserCache, mu *mocks.MailUtils, ru *mocks.RandomUtils) {
					ur.On("GetByEmail", email).Return(nil, errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				err: errs.ErrUserDoesNotExist,
			},
		},
		{
			description: "should return error when StoreResetPasswordToken failed",
			input: input{
				request: &dto.RequestPasswordResetRequest{
					Email: email,
				},
				beforeTest: func(ur *mocks.UserRepository, uc *mocks.UserCache, mu *mocks.MailUtils, ru *mocks.RandomUtils) {
					ur.On("GetByEmail", email).Return(&model.User{}, nil)
					ru.On("GenerateSecureUniqueToken").Return("token")
					uc.On("StoreResetPasswordToken", mock.Anything, mock.Anything).Return(errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				err: errs.ErrUserDoesNotExist,
			},
		},
		{
			description: "should return error when SendResetPasswordEmail failed",
			input: input{
				request: &dto.RequestPasswordResetRequest{
					Email: email,
				},
				beforeTest: func(ur *mocks.UserRepository, uc *mocks.UserCache, mu *mocks.MailUtils, ru *mocks.RandomUtils) {
					ur.On("GetByEmail", email).Return(&model.User{}, nil)
					ru.On("GenerateSecureUniqueToken").Return("token")
					uc.On("StoreResetPasswordToken", mock.Anything, mock.Anything).Return(nil)
					mu.On("SendResetPasswordEmail", mock.Anything, mock.Anything).Return(errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				err: errs.ErrUserDoesNotExist,
			},
		},
		{
			description: "should return nil when success",
			input: input{
				request: &dto.RequestPasswordResetRequest{
					Email: email,
				},
				beforeTest: func(ur *mocks.UserRepository, uc *mocks.UserCache, mu *mocks.MailUtils, ru *mocks.RandomUtils) {
					ur.On("GetByEmail", email).Return(&model.User{}, nil)
					ru.On("GenerateSecureUniqueToken").Return("token")
					uc.On("StoreResetPasswordToken", mock.Anything, mock.Anything).Return(nil)
					mu.On("SendResetPasswordEmail", mock.Anything, mock.Anything).Return(nil)
				},
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			userRepo := mocks.NewUserRepository(t)
			userCache := mocks.NewUserCache(t)
			mailUtils := mocks.NewMailUtils(t)
			randomUtils := mocks.NewRandomUtils(t)
			tc.beforeTest(userRepo, userCache, mailUtils, randomUtils)
			userService := service.NewUserService(&service.UserSConfig{
				Repository:  userRepo,
				Redis:       userCache,
				MailUtils:   mailUtils,
				RandomUtils: randomUtils,
			})

			actualErr := userService.RequestPasswordReset(tc.input.request)

			assert.Equal(t, tc.expected.err, actualErr)
		})
	}

}

func TestCompletePasswordReset(t *testing.T) {
	var (
		token = "token"
	)

	type input struct {
		request    *dto.CompletePasswordResetRequest
		beforeTest func(*mocks.UserRepository, *mocks.UserCache)
	}
	type expected struct {
		err error
	}

	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when FindResetPasswordToken failed",
			input: input{
				request: &dto.CompletePasswordResetRequest{
					Token:       token,
					NewPassword: "newPassword",
				},
				beforeTest: func(ur *mocks.UserRepository, uc *mocks.UserCache) {
					uc.On("FindResetPasswordToken", token).Return(0, errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				err: errs.ErrUserDoesNotExist,
			},
		},
		{
			description: "should return error when GetByID failed",
			input: input{
				request: &dto.CompletePasswordResetRequest{
					Token:       token,
					NewPassword: "newPassword",
				},
				beforeTest: func(ur *mocks.UserRepository, uc *mocks.UserCache) {
					uc.On("FindResetPasswordToken", token).Return(1, nil)
					ur.On("GetByID", 1).Return(nil, errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				err: errs.ErrUserDoesNotExist,
			},
		},
		{
			description: "should return error when UpdatePassword failed",
			input: input{
				request: &dto.CompletePasswordResetRequest{
					Token:       token,
					NewPassword: "newPassword123",
				},
				beforeTest: func(ur *mocks.UserRepository, uc *mocks.UserCache) {
					uc.On("FindResetPasswordToken", token).Return(1, nil)
					ur.On("GetByID", 1).Return(&model.User{Username: "asd", ID: 1}, nil)
					ur.On("UpdatePassword", mock.Anything, mock.Anything).Return(nil, errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				err: errs.ErrUserDoesNotExist,
			},
		},
		{
			description: "should return error when DeleteResetPasswordToken failed",
			input: input{
				request: &dto.CompletePasswordResetRequest{
					Token:       token,
					NewPassword: "newPassword123",
				},
				beforeTest: func(ur *mocks.UserRepository, uc *mocks.UserCache) {
					uc.On("FindResetPasswordToken", token).Return(1, nil)
					ur.On("GetByID", 1).Return(&model.User{Username: "asd", ID: 1}, nil)
					ur.On("UpdatePassword", 1, "newPassword123").Return(nil, nil)
					uc.On("DeleteResetPasswordToken", token).Return(errs.ErrUserDoesNotExist)
				},
			},
			expected: expected{
				err: errs.ErrUserDoesNotExist,
			},
		},

		{
			description: "should return nil when success",
			input: input{
				request: &dto.CompletePasswordResetRequest{
					Token:       token,
					NewPassword: "newPassword123",
				},
				beforeTest: func(ur *mocks.UserRepository, uc *mocks.UserCache) {
					uc.On("FindResetPasswordToken", token).Return(1, nil)
					ur.On("GetByID", 1).Return(&model.User{Username: "asd", ID: 1}, nil)
					ur.On("UpdatePassword", 1, "newPassword123").Return(nil, nil)
					uc.On("DeleteResetPasswordToken", token).Return(nil)
				},
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			userRepo := mocks.NewUserRepository(t)
			userCache := mocks.NewUserCache(t)
			tc.beforeTest(userRepo, userCache)
			userService := service.NewUserService(&service.UserSConfig{
				Repository: userRepo,
				Redis:      userCache,
			})

			actualErr := userService.CompletePasswordReset(tc.input.request)

			assert.Equal(t, tc.expected.err, actualErr)
		})
	}

}

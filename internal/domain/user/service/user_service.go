package service

import (
	"errors"
	"fmt"
	"kedai/backend/be-kedai/config"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/cache"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
	"kedai/backend/be-kedai/internal/utils/credential"
	"kedai/backend/be-kedai/internal/utils/google"
	"kedai/backend/be-kedai/internal/utils/hash"
	jwttoken "kedai/backend/be-kedai/internal/utils/jwtToken"
	"kedai/backend/be-kedai/internal/utils/mail"
	"kedai/backend/be-kedai/internal/utils/random"
	"strings"

	"github.com/redis/go-redis/v9"
)

type UserService interface {
	GetByID(id int) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	SignUp(*dto.UserRegistrationRequest) (*dto.UserRegistrationResponse, error)
	SignIn(*dto.UserLogin, string) (*dto.Token, error)
	SignInWithGoogle(userLogin *dto.UserLoginWithGoogleRequest) (*dto.Token, error)
	SignUpWithGoogle(userReg *dto.UserRegistrationWithGoogleRequest) (*dto.Token, error)
	GetSession(userId int, token string) error
	RenewToken(userId int, refreshToken string) (*dto.Token, error)
	UpdateEmail(userId int, request *dto.UpdateEmailRequest) (*dto.UpdateEmailResponse, error)
	UpdateUsername(userId int, requst *dto.UpdateUsernameRequest) (*dto.UpdateUsernameResponse, error)
	RequestPasswordChange(request *dto.RequestPasswordChangeRequest) error
	CompletePasswordChange(request *dto.CompletePasswordChangeRequest) error
	RequestPasswordReset(request *dto.RequestPasswordResetRequest) error
	CompletePasswordReset(request *dto.CompletePasswordResetRequest) error
	ValidatePasswordChange(request *dto.RequestPasswordChangeRequest, user *model.User) error
	SignOut(*dto.UserLogoutRequest) error
	AdminSignIn(*dto.UserLogin) (*dto.Token, error)
}

type userServiceImpl struct {
	repository  repository.UserRepository
	redis       cache.UserCache
	randomUtils random.RandomUtils
	mailUtils   mail.MailUtils
}

type UserSConfig struct {
	Repository  repository.UserRepository
	Redis       cache.UserCache
	RandomUtils random.RandomUtils
	MailUtils   mail.MailUtils
}

func NewUserService(cfg *UserSConfig) UserService {
	return &userServiceImpl{
		repository:  cfg.Repository,
		redis:       cfg.Redis,
		mailUtils:   cfg.MailUtils,
		randomUtils: cfg.RandomUtils,
	}
}

func (s *userServiceImpl) GetByID(id int) (*model.User, error) {
	return s.repository.GetByID(id)
}

func (s *userServiceImpl) GetByUsername(username string) (*model.User, error) {
	return s.repository.GetByUsername(username)
}

func (s *userServiceImpl) SignUp(userReg *dto.UserRegistrationRequest) (*dto.UserRegistrationResponse, error) {
	isValidPassword := credential.VerifyPassword(userReg.Password)
	if !isValidPassword {
		return nil, errs.ErrInvalidPasswordPattern
	}

	emailSplit := strings.Split(userReg.Email, "@")
	isContainEmail := strings.Contains(strings.ToLower(userReg.Password), strings.ToLower(emailSplit[0]))
	if isContainEmail {
		return nil, errs.ErrContainEmail
	}

	user := userReg.ToUser()
	user.Email = strings.ToLower(user.Email)

	result, err := s.repository.SignUp(user)
	if err != nil {
		return nil, err
	}

	var response dto.UserRegistrationResponse
	response.FromUser(result)

	return &response, nil
}

func (s *userServiceImpl) SignUpWithGoogle(userReg *dto.UserRegistrationWithGoogleRequest) (*dto.Token, error) {
	isValidUsername := credential.VerifyUsername(userReg.Username)
	if !isValidUsername {
		return nil, errs.ErrInvalidUsernamePattern
	}

	isValidPassword := credential.VerifyPassword(userReg.Password)
	if !isValidPassword {
		return nil, errs.ErrInvalidPasswordPattern
	}

	if credential.ContainsUsername(userReg.Password, userReg.Username) {
		return nil, errs.ErrInvalidPasswordPattern
	}

	claim, err := google.ValidateGoogleToken(userReg.Credential)
	if err != nil {
		return nil, errs.ErrUnauthorized
	}

	result, err := s.repository.SignUp(&model.User{
		Email:    claim.Email,
		Username: userReg.Username,
		Password: userReg.Password,
	})

	if err != nil {
		return nil, err
	}

	defaultLevel := 0

	accessToken, _ := jwttoken.GenerateAccessToken(result, defaultLevel)
	refreshToken, _ := jwttoken.GenerateRefreshToken(result, defaultLevel)

	token := &dto.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	err = s.redis.StoreToken(result.ID, accessToken, refreshToken)
	if err != nil {
		return nil, err
	}

	return token, nil

}

func (s *userServiceImpl) SignIn(userLogin *dto.UserLogin, inputPw string) (*dto.Token, error) {
	user := userLogin.ToUser()
	user.Email = strings.ToLower(user.Email)

	result, err := s.repository.SignIn(user)
	if err != nil {
		return nil, err
	}

	isValid := hash.ComparePassword(result.Password, inputPw)
	if isValid {
		defaultLevel := 0
		accessToken, _ := jwttoken.GenerateAccessToken(result, defaultLevel)
		refreshToken, _ := jwttoken.GenerateRefreshToken(result, defaultLevel)

		token := &dto.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		err := s.redis.StoreToken(result.ID, accessToken, refreshToken)
		if err != nil {
			return nil, err
		}

		return token, nil
	}

	return nil, errs.ErrInvalidCredential
}

func (s *userServiceImpl) SignInWithGoogle(userLogin *dto.UserLoginWithGoogleRequest) (*dto.Token, error) {
	claim, err := google.ValidateGoogleToken(userLogin.Credential)
	if err != nil {
		return nil, errs.ErrUnauthorized
	}

	result, err := s.repository.SignIn(&model.User{Email: claim.Email})
	if err != nil {
		return nil, err
	}

	defaultLevel := 0
	accessToken, _ := jwttoken.GenerateAccessToken(result, defaultLevel)
	refreshToken, _ := jwttoken.GenerateRefreshToken(result, defaultLevel)

	token := &dto.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	err = s.redis.StoreToken(result.ID, accessToken, refreshToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *userServiceImpl) GetSession(userId int, accessToken string) error {
	return s.redis.FindToken(userId, accessToken)
}

func (s *userServiceImpl) RenewToken(userId int, refreshToken string) (*dto.Token, error) {
	err := s.redis.FindToken(userId, refreshToken)
	if errors.Is(err, redis.Nil) {
		return nil, errs.ErrExpiredToken
	}
	if err != nil {
		return nil, err
	}

	err = s.redis.DeleteToken(fmt.Sprintf("user_%d:%s", userId, refreshToken))
	if err != nil {
		return nil, err
	}

	defaultLevel := 0
	newAccessToken, _ := jwttoken.GenerateAccessToken(&model.User{ID: userId}, defaultLevel)
	newRefreshToken, _ := jwttoken.GenerateRefreshToken(&model.User{ID: userId}, defaultLevel)

	err = s.redis.StoreToken(userId, newAccessToken, newRefreshToken)
	if err != nil {
		return nil, err
	}

	token := dto.Token{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	return &token, nil
}

func (s *userServiceImpl) UpdateEmail(userId int, request *dto.UpdateEmailRequest) (*dto.UpdateEmailResponse, error) {
	email := strings.ToLower(request.Email)

	_, err := s.repository.GetByEmail(email)
	if err == nil {
		return nil, errs.ErrEmailUsed
	}

	if !errors.Is(err, errs.ErrUserDoesNotExist) {
		return nil, err
	}

	res, err := s.repository.UpdateEmail(userId, email)
	if err != nil {
		return nil, err
	}

	var response dto.UpdateEmailResponse
	response.FromUser(res)

	return &response, nil
}

func (s *userServiceImpl) UpdateUsername(userId int, request *dto.UpdateUsernameRequest) (*dto.UpdateUsernameResponse, error) {
	username := strings.ToLower(request.Username)

	if isUsernameValid := credential.VerifyUsername(username); !isUsernameValid {
		return nil, errs.ErrInvalidUsernamePattern
	}

	res, err := s.repository.UpdateUsername(userId, username)
	if err != nil {
		return nil, err
	}

	response := dto.UpdateUsernameResponse{
		Username: res.Username,
	}

	return &response, nil
}

func (s *userServiceImpl) SignOut(request *dto.UserLogoutRequest) error {
	return s.redis.DeleteRefreshTokenAndAccessToken(request.UserId, request.RefreshToken, request.AccessToken)
}

func (s *userServiceImpl) RequestPasswordChange(request *dto.RequestPasswordChangeRequest) error {
	user, err := s.repository.GetByID(request.UserId)
	if err != nil {
		return err
	}

	isValidPassword := hash.ComparePassword(user.Password, request.CurrentPassword)
	if !isValidPassword {
		return errs.ErrInvalidCredential
	}

	err = s.ValidatePasswordChange(request, user)
	if err != nil {
		return err
	}

	codeLength := 6
	verifCode := s.randomUtils.GenerateAlphanumericString(codeLength)

	err = s.redis.StoreUserPasswordAndVerificationCode(request.UserId, request.NewPassword, verifCode)
	if err != nil {
		return err
	}

	err = s.mailUtils.SendUpdatePasswordEmail(user.Email, verifCode)
	if err != nil {
		return err
	}

	return nil
}

func (s *userServiceImpl) ValidatePasswordChange(request *dto.RequestPasswordChangeRequest, user *model.User) error {
	isInvalidPassword := hash.ComparePassword(user.Password, request.NewPassword)
	if isInvalidPassword {
		return errs.ErrSamePassword
	}

	isValidPassword := credential.VerifyPassword(request.NewPassword)
	if !isValidPassword {
		return errs.ErrInvalidPasswordPattern
	}

	if credential.ContainsUsername(request.NewPassword, user.Username) {
		return errs.ErrContainUsername
	}

	return nil
}

func (s *userServiceImpl) CompletePasswordChange(request *dto.CompletePasswordChangeRequest) error {
	newPassword, verifcationCode, err := s.redis.FindUserPasswordAndVerificationCode(request.UserId)
	if err != nil {
		return err
	}

	if verifcationCode != request.VerificationCode {
		return errs.ErrIncorrectVerificationCode
	}

	_, err = s.repository.UpdatePassword(request.UserId, newPassword)
	if err != nil {
		return err
	}

	_ = s.redis.DeleteUserPasswordAndVerificationCode(request.UserId)

	return nil
}

func (s *userServiceImpl) RequestPasswordReset(request *dto.RequestPasswordResetRequest) error {
	user, err := s.repository.GetByEmail(request.Email)
	if err != nil {
		return err
	}

	secureToken := s.randomUtils.GenerateSecureUniqueToken()
	err = s.redis.StoreResetPasswordToken(user.ID, secureToken)
	if err != nil {
		return err
	}

	err = s.mailUtils.SendResetPasswordEmail(user.Email, secureToken)
	if err != nil {
		return err
	}

	return nil
}

func (s *userServiceImpl) CompletePasswordReset(request *dto.CompletePasswordResetRequest) error {
	userId, err := s.redis.FindResetPasswordToken(request.Token)
	if err != nil {
		return err
	}

	user, err := s.repository.GetByID(userId)
	if err != nil {
		return err
	}

	err = s.ValidatePasswordChange(&dto.RequestPasswordChangeRequest{
		UserId:      userId,
		NewPassword: request.NewPassword,
	}, user)
	if err != nil {
		return err
	}

	_, err = s.repository.UpdatePassword(userId, request.NewPassword)
	if err != nil {
		return err
	}

	_ = s.redis.DeleteResetPasswordToken(request.Token)

	return nil
}

func (s *userServiceImpl) AdminSignIn(userLogin *dto.UserLogin) (*dto.Token, error) {

	admin := userLogin.ToUser()

	isValid := hash.ComparePassword(config.AdminPassword, admin.Password)

	if admin.Email == config.AdminEmail && isValid {
		result := &model.User{
			ID:       0,
			Email:    admin.Email,
			Username: "admin",
		}

		defaultLevel := 0
		accessToken, _ := jwttoken.GenerateAccessToken(result, defaultLevel)
		refreshToken, _ := jwttoken.GenerateRefreshToken(result, defaultLevel)

		token := &dto.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		errToken := s.redis.StoreToken(result.ID, accessToken, refreshToken)
		if errToken != nil {
			return nil, errToken
		}

		return token, nil
	}

	return nil, errs.ErrInvalidCredential
}

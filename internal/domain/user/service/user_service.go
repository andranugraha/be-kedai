package service

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/cache"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
	"kedai/backend/be-kedai/internal/utils/hash"
	jwttoken "kedai/backend/be-kedai/internal/utils/jwtToken"
	pwValidator "kedai/backend/be-kedai/internal/utils/password"
	"strings"
)

type UserService interface {
	GetByID(id int) (*model.User, error)
	SignUp(*dto.UserRegistrationRequest) (*dto.UserRegistrationResponse, error)
	SignIn(*dto.UserLogin, string) (*dto.Token, error)
	GetSession(userId int, token string) error
}

type userServiceImpl struct {
	repository repository.UserRepository
	redis      cache.UserCache
}

type UserSConfig struct {
	Repository repository.UserRepository
	Redis      cache.UserCache
}

func NewUserService(cfg *UserSConfig) UserService {
	return &userServiceImpl{
		repository: cfg.Repository,
		redis:      cfg.Redis,
	}
}

func (s *userServiceImpl) GetByID(id int) (*model.User, error) {
	return s.repository.GetByID(id)
}

func (s *userServiceImpl) SignUp(userReg *dto.UserRegistrationRequest) (*dto.UserRegistrationResponse, error) {
	isValidPassword := pwValidator.VerifyPassword(userReg.Password)
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

func (s *userServiceImpl) SignIn(userLogin *dto.UserLogin, inputPw string) (*dto.Token, error) {
	user := userLogin.ToUser()
	user.Email = strings.ToLower(user.Email)

	result, err := s.repository.SignIn(user)
	if err != nil {
		return nil, err
	}

	isValid := hash.ComparePassword(result.Password, inputPw)
	if isValid {
		accessToken, _ := jwttoken.GenerateAccessToken(result)
		refreshToken, _ := jwttoken.GenerateRefreshToken(result)

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

func (s *userServiceImpl) GetSession(userId int, accessToken string) error {
	return s.redis.FindToken(userId, accessToken)
}

package service

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/cache"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
	"kedai/backend/be-kedai/internal/utils/hash"
	jwttoken "kedai/backend/be-kedai/internal/utils/jwtToken"
)

type UserService interface {
	GetByID(id int) (*model.User, error)
	SignUp(*dto.UserRegistration) (*dto.UserRegistration, error)
	SignIn(*dto.UserLogin, string) (*dto.Token, error)
}

type userServiceImpl struct {
	repository repository.UserRepository
	redis cache.UserCache
}

type UserSConfig struct {
	Repository repository.UserRepository
	Redis cache.UserCache
}

func NewUserService(cfg *UserSConfig) UserService {
	return &userServiceImpl{
		repository: cfg.Repository,
		redis: cfg.Redis,
	}
}

func (s *userServiceImpl) GetByID(id int) (*model.User, error) {
	return s.repository.GetByID(id)
}

func (s *userServiceImpl) SignUp(userReg *dto.UserRegistration) (*dto.UserRegistration, error) {
	user := userReg.ToUser()

	result, err := s.repository.SignUp(user)
	if err != nil {
		return nil, err
	}

	userReg.FromUser(result)

	return userReg, nil
}

func (s *userServiceImpl) SignIn(userLogin *dto.UserLogin, inputPw string) (*dto.Token, error) {
	user := userLogin.ToUser()

	result, err := s.repository.SignIn(user)
	if err != nil {
		return nil, err
	}

	isValid := hash.ComparePassword(result.Password, inputPw)
	if isValid {
		accessToken, _ := jwttoken.GenerateAccessToken(result)
		refreshToken, _ := jwttoken.GenerateRefreshToken(result)
		
		token := &dto.Token{
			AccessToken: accessToken,
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

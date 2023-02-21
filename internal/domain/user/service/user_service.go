package service

import (
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
	errs "kedai/backend/be-kedai/internal/common/error"
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
}

type UserSConfig struct {
	Repository repository.UserRepository
}

func NewUserService(cfg *UserSConfig) UserService {
	return &userServiceImpl{
		repository: cfg.Repository,
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
		token, _ := jwttoken.GenerateAccessToken(result)
		return token, nil
	}

	return nil, errs.ErrInvalidCredential
}

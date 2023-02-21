package service

import (
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/repository"
)

type UserService interface {
	SignUp(*dto.UserRegistration) (*dto.UserRegistration, error)
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

func (s *userServiceImpl) SignUp(userReg *dto.UserRegistration) (*dto.UserRegistration, error) {
	user := userReg.ToUser()

	result, err := s.repository.SignUp(user)
	if err != nil {
		return nil, err
	}

	userReg.FromUser(result)

	return userReg, nil
}

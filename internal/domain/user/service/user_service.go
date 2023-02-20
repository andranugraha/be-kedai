package service

import (
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
)

type UserService interface {
	SignUp(*dto.UserRegistration) (*model.User, error)
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

func (s *userServiceImpl) SignUp(userReg *dto.UserRegistration) (*model.User, error) {
	user := userReg.ToUser()

	return s.repository.SignUp(user)
}

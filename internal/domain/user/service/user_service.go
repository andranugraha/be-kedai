package service

import (
	model "kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
)

type UserService interface {
	GetByID(id int) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
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

func (s *userServiceImpl) GetByEmail(email string) (*model.User, error) {
	return s.repository.GetByEmail(email)
}

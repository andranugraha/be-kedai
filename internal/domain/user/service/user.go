package usecase

import (
	model "kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
)

type UserUsecase interface {
	GetByID(id int) (*model.User, error)
}

type userUsecaseImpl struct {
	repository repository.UserRepository
}

type UserUConfig struct {
	Repository repository.UserRepository
}

func NewUserUsecase(cfg *UserUConfig) UserUsecase {
	return &userUsecaseImpl{
		repository: cfg.Repository,
	}
}

func (u *userUsecaseImpl) GetByID(id int) (*model.User, error) {
	return u.repository.GetByID(id)
}

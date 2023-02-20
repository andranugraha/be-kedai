package service

import (
	"fmt"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
	"kedai/backend/be-kedai/internal/utils/hash"
	"math/rand"
	"strings"
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

	user.Password, _ = hash.HashAndSalt(user.Password)

	emailString := strings.Split(user.Email, "@")
	
	username := fmt.Sprintf("%s%d", emailString[0], rand.Intn(999))

	user.Username = username

	return s.repository.SignUp(user)
}
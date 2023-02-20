package repository

import (
	"fmt"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/utils/hash"
	"math/rand"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	SignUp(user *model.User) (*model.User, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

type UserRConfig struct {
	DB *gorm.DB
}

func NewUserRepository(cfg *UserRConfig) UserRepository {
	return &userRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *userRepositoryImpl) SignUp(user *model.User) (*model.User, error) {
	emailString := strings.Split(user.Email, "@")

	username := fmt.Sprintf("%s%d", emailString[0], rand.Intn(999))

	user.Username = username

	hashedPw, _ := hash.HashAndSalt(user.Password)

	user.Password = hashedPw

	err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
	if err.Error != nil {
		return nil, err.Error
	}

	if err.RowsAffected == 0 {
		return nil, errs.ErrUserAlreadyExist
	}

	return user, nil
}

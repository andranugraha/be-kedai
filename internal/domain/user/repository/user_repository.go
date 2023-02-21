package repository

import (
	"fmt"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/utils/hash"
	"math/rand"
	"strings"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	GetByID(ID int) (*model.User, error)
	SignUp(user *model.User) (*model.User, error)
	SignIn(user *model.User) (*model.User, error)
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

func (r *userRepositoryImpl) GetByID(ID int) (*model.User, error) {
	var user model.User

	err := r.db.Where("user_id = ?", ID).Preload("Profile").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserDoesNotExist
		}

		return nil, err
	}

	return &user, nil
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

func (r *userRepositoryImpl) SignIn(user *model.User) (*model.User, error) {
	err := r.db.Where("email = ?", user.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserDoesNotExist
		}
		return nil, err
	}

	return user, nil
}
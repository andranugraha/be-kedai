package repository

import (
	"errors"
	"fmt"
	errs "kedai/backend/be-kedai/internal/common/error"

	"kedai/backend/be-kedai/internal/domain/user/cache"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/utils/hash"
	"math/rand"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	GetByID(ID int) (*model.User, error)
	SignUp(user *model.User) (*model.User, error)
	SignIn(user *model.User) (*model.User, error)
	UpdateEmail(id int, payload *model.User) (*model.User, error)
	UpdateUsername(id int, username string) (*model.User, error)
}

type userRepositoryImpl struct {
	db        *gorm.DB
	userCache cache.UserCache
}

type UserRConfig struct {
	DB        *gorm.DB
	UserCache cache.UserCache
}

func NewUserRepository(cfg *UserRConfig) UserRepository {
	return &userRepositoryImpl{
		db:        cfg.DB,
		userCache: cfg.UserCache,
	}
}

func (r *userRepositoryImpl) GetByID(ID int) (*model.User, error) {
	var user model.User

	err := r.db.Where("id = ?", ID).Preload("Shop").Preload("Profile").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserDoesNotExist
		}

		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) GetByEmail(email string) (*model.User, error) {
	var user model.User

	err := r.db.Where("email = ?", email).First(&user).Error
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
			return nil, errs.ErrInvalidCredential
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) UpdateEmail(id int, payload *model.User) (*model.User, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("id = ?", id).Clauses(clause.OnConflict{DoNothing: true}).Updates(payload)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errs.ErrEmailUsed
		}

		if err := r.userCache.DeleteAllByID(id); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (r *userRepositoryImpl) UpdateUsername(userId int, username string) (*model.User, error) {
	res := r.db.Model(&model.User{}).Where("id = ?", userId).Clauses(clause.OnConflict{DoNothing: true}).Update("username", username)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errs.ErrUsernameUsed
	}

	return &model.User{ID: 1, Username: username}, nil
}

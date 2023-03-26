package repository

import (
	"errors"
	"fmt"
	"kedai/backend/be-kedai/config"
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
	GetByEmail(email string) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	SignUp(user *model.User) (*model.User, error)
	SignIn(user *model.User) (*model.User, error)
	UpdateEmail(userId int, email string) (*model.User, error)
	UpdateUsername(id int, username string) (*model.User, error)
	UpdatePassword(id int, password string) (*model.User, error)
	AdminSignIn(email string, password string) error
}

type userRepositoryImpl struct {
	db              *gorm.DB
	userCache       cache.UserCache
	userProfileRepo UserProfileRepository
}

type UserRConfig struct {
	DB              *gorm.DB
	UserCache       cache.UserCache
	UserProfileRepo UserProfileRepository
}

func NewUserRepository(cfg *UserRConfig) UserRepository {
	return &userRepositoryImpl{
		db:              cfg.DB,
		userCache:       cfg.UserCache,
		userProfileRepo: cfg.UserProfileRepo,
	}
}

func (r *userRepositoryImpl) GetByID(ID int) (*model.User, error) {
	var user model.User

	err := r.db.Where("id = ?", ID).Preload("Profile").First(&user).Error
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

	if user.Username == "" {
		username := fmt.Sprintf("%s%d", emailString[0], rand.Intn(999))
		user.Username = username
	}

	_, err := r.GetByUsername(user.Username)
	if err == nil {
		return nil, errs.ErrUsernameUsed
	}

	hashedPw, _ := hash.HashAndSalt(user.Password)
	user.Password = hashedPw

	err = r.db.Where("email = ?", user.Email).First(&model.UsedEmail{}).Error
	if err == nil {
		return nil, errs.ErrEmailUsed
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	tx := r.db.Begin()
	defer tx.Commit()

	res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(user)
	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		tx.Rollback()
		return nil, errs.ErrUserAlreadyExist
	}

	_, err = r.userProfileRepo.Create(tx, &model.UserProfile{UserID: user.ID})
	if err != nil {
		tx.Rollback()
		return nil, err
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

func (r *userRepositoryImpl) UpdateEmail(userId int, email string) (*model.User, error) {
	err := r.db.Where("email = ?", email).First(&model.UsedEmail{}).Error
	if err == nil {
		return nil, errs.ErrEmailUsed
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var user model.User
	err = r.db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.User{}).Where("id = ?", userId).Update("email", email).Error; err != nil {
			return err
		}

		if err := tx.Create(
			&model.UsedEmail{
				Email: user.Email,
			}).Error; err != nil {
			return err
		}

		if err := r.userCache.DeleteAllByID(userId); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &model.User{ID: userId, Email: email}, nil
}

func (r *userRepositoryImpl) UpdateUsername(userId int, username string) (*model.User, error) {
	res := r.db.Model(&model.User{}).Where("id = ?", userId).Clauses(clause.OnConflict{DoNothing: true}).Update("username", username)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errs.ErrUsernameUsed
	}

	return &model.User{ID: userId, Username: username}, nil
}

func (r *userRepositoryImpl) GetByUsername(username string) (*model.User, error) {
	var user model.User

	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserDoesNotExist
		}

		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) UpdatePassword(userId int, password string) (*model.User, error) {
	hashedPw, _ := hash.HashAndSalt(password)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.User{}).Where("id = ?", userId).Update("password", hashedPw).Error; err != nil {
			return err
		}

		if err := r.userCache.DeleteAllByID(userId); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &model.User{ID: userId}, nil
}

func (r *userRepositoryImpl) AdminSignIn(email string, password string) error {
	if email == config.AdminEmail && password == config.AdminPassword {
		return nil
	}
	return errs.ErrInvalidCredential
}

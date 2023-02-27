package repository

import (
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserProfileRepository interface {
	Create(tx *gorm.DB, payload *model.UserProfile) (*model.UserProfile, error)
	Update(userId int, payload *model.UserProfile) (*model.UserProfile, error)
}

type userProfileRepositoryImpl struct {
	db *gorm.DB
}

type UserProfileRConfig struct {
	DB *gorm.DB
}

func NewUserProfileRepository(cfg *UserProfileRConfig) UserProfileRepository {
	return &userProfileRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *userProfileRepositoryImpl) Update(userId int, payload *model.UserProfile) (*model.UserProfile, error) {
	err := r.db.Where("user_id = ?", userId).Clauses(clause.Returning{}).Updates(payload).Error
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (r *userProfileRepositoryImpl) Create(tx *gorm.DB, payload *model.UserProfile) (*model.UserProfile, error) {
	err := tx.Create(payload).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return payload, nil
}

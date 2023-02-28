package repository

import (
	spErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SealabsPayRepository interface {
	GetByUserID(userID int) ([]*model.SealabsPay, error)
	Create(sealabsPay *model.SealabsPay) error
}

type sealabsPayRepositoryImpl struct {
	db *gorm.DB
}

type SealabsPayRConfig struct {
	DB *gorm.DB
}

func NewSealabsPayRepository(config *SealabsPayRConfig) SealabsPayRepository {
	return &sealabsPayRepositoryImpl{
		db: config.DB,
	}
}

func (r *sealabsPayRepositoryImpl) GetByUserID(userID int) ([]*model.SealabsPay, error) {
	var sealabsPays []*model.SealabsPay

	err := r.db.Where("user_id = ?", userID).Find(sealabsPays).Error
	if err != nil {
		return nil, err
	}

	return sealabsPays, nil
}

func (r *sealabsPayRepositoryImpl) Create(sealabsPay *model.SealabsPay) error {
	err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&sealabsPay)
	if err.Error != nil {
		return err.Error
	}

	if err.RowsAffected == 0 {
		return spErr.ErrSealabsPayAlreadyRegistered
	}

	return nil
}

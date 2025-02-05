package service

import (
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
	"time"
)

type SealabsPayService interface {
	GetSealabsPaysByUserID(userID int) ([]*model.SealabsPay, error)
	RegisterSealabsPay(*dto.CreateSealabsPayRequest) (*model.SealabsPay, error)
	GetValidSealabsPayByCardNumberAndUserID(cardNumber string, userID int) (*model.SealabsPay, error)
}

type sealabsPayServiceImpl struct {
	sealabsPayRepo repository.SealabsPayRepository
}

type SealabsPaySConfig struct {
	SealabsPayRepo repository.SealabsPayRepository
}

func NewSealabsPayService(config *SealabsPaySConfig) SealabsPayService {
	return &sealabsPayServiceImpl{
		sealabsPayRepo: config.SealabsPayRepo,
	}
}

func (s *sealabsPayServiceImpl) GetSealabsPaysByUserID(userID int) ([]*model.SealabsPay, error) {
	return s.sealabsPayRepo.GetByUserID(userID)
}

func (s *sealabsPayServiceImpl) RegisterSealabsPay(req *dto.CreateSealabsPayRequest) (*model.SealabsPay, error) {
	expiryDate, _ := time.Parse("01/06", req.ExpiryDate)
	expiryDate = time.Date(expiryDate.Year(), expiryDate.Month()+1, 0, 0, 0, 0, 0, time.UTC)

	sealabsPay := &model.SealabsPay{
		CardNumber: req.CardNumber,
		CardName:   req.CardName,
		ExpiryDate: expiryDate,
		UserID:     req.UserID,
	}

	err := s.sealabsPayRepo.Create(sealabsPay)
	if err != nil {
		return nil, err
	}

	return sealabsPay, nil
}

func (s *sealabsPayServiceImpl) GetValidSealabsPayByCardNumberAndUserID(cardNumber string, userID int) (*model.SealabsPay, error) {
	return s.sealabsPayRepo.GetValidByCardNumberAndUserID(cardNumber, userID)
}

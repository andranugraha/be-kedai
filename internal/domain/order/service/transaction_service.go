package service

import (
	"kedai/backend/be-kedai/internal/domain/order/model"
	"kedai/backend/be-kedai/internal/domain/order/repository"
)

type TransactionService interface {
	GetByID(id int) (*model.Transaction, error)
}

type transactionServiceImpl struct {
	transactionRepo repository.TransactionRepository
}

type TransactionSConfig struct {
	TransactionRepo repository.TransactionRepository
}

func NewTransactionService(config *TransactionSConfig) TransactionService {
	return &transactionServiceImpl{
		transactionRepo: config.TransactionRepo,
	}
}

func (s *transactionServiceImpl) GetByID(id int) (*model.Transaction, error) {
	return s.transactionRepo.GetByID(id)
}

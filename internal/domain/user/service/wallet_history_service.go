package service

import (
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
)

type WalletHistoryService interface {
	GetWalletHistoryById(userId int) ([]*model.WalletHistory, error)	
}

type walletHistoryImpl struct {
	walletHistoryRepository repository.WalletHistoryRepository
	walletService	WalletService
}

type WalletHistorySConfig struct {
	WalletHistoryRepository repository.WalletHistoryRepository
	WalletService WalletService
}

func NewWalletHistoryService(cfg *WalletHistorySConfig) WalletHistoryService {
	return &walletHistoryImpl{
		walletHistoryRepository: cfg.WalletHistoryRepository,
		walletService: cfg.WalletService,
	}
}

func (s *walletHistoryImpl) GetWalletHistoryById(userId int) ([]*model.WalletHistory, error) {
	wallet, err := s.walletService.GetWalletByUserID(userId)
	if err != nil {
		return nil, err
	}

	return s.walletHistoryRepository.GetWalletHistoryById(wallet.ID)
}

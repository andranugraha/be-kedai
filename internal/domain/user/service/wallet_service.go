package service

import (
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
	"kedai/backend/be-kedai/internal/utils/hash"
)

const(
	TopUp = "Top-up"
)

type WalletService interface {
	RegisterWallet(userID int, pin string) (*model.Wallet, error)
	GetWalletByUserID(userID int) (*model.Wallet, error)
	TopUp(userId int, req dto.TopUpRequest) (*model.WalletHistory, error)
}

type walletServiceImpl struct {
	walletRepo repository.WalletRepository
}

type WalletSConfig struct {
	WalletRepo repository.WalletRepository
}

func NewWalletService(cfg *WalletSConfig) WalletService {
	return &walletServiceImpl{
		walletRepo: cfg.WalletRepo,
	}
}

func (s *walletServiceImpl) RegisterWallet(userID int, pin string) (*model.Wallet, error) {
	hashedPin, _ := hash.HashAndSalt(pin)

	wallet := &model.Wallet{
		UserID:  userID,
		Pin:     hashedPin,
		Balance: 0,
	}

	return s.walletRepo.Create(wallet)
}

func (s *walletServiceImpl) GetWalletByUserID(userID int) (*model.Wallet, error) {
	return s.walletRepo.GetByUserID(userID)
}

func (s *walletServiceImpl) TopUp(userId int, req dto.TopUpRequest) (*model.WalletHistory, error) {
	var history model.WalletHistory
	
	wallet, err := s.walletRepo.GetByUserID(userId)
	if err != nil {
		return nil, err
	}
	
	history.WalletId = wallet.ID
	history.Type = TopUp
	history.Amount = req.Amount
	history.Reference = req.TxnId

	return s.walletRepo.TopUp(&history, wallet)
}

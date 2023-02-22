package service

import (
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
	"kedai/backend/be-kedai/internal/utils/hash"
)

type WalletService interface {
	RegisterWallet(userID int, pin string) (*model.Wallet, error)
	GetByUserID(userID int) (*model.Wallet, error)
}

type walletServiceImpl struct {
	walletRepo  repository.WalletRepository
	userService UserService
}

type WalletSConfig struct {
	WalletRepo  repository.WalletRepository
	UserService UserService
}

func NewWalletService(cfg *WalletSConfig) WalletService {
	return &walletServiceImpl{
		walletRepo:  cfg.WalletRepo,
		userService: cfg.UserService,
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

func (s *walletServiceImpl) GetByUserID(userID int) (*model.Wallet, error) {
	_, err := s.userService.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return s.walletRepo.GetByUserID(userID)
}

package service

import (
	"kedai/backend/be-kedai/internal/domain/user/cache"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
	"kedai/backend/be-kedai/internal/utils/hash"
	jwttoken "kedai/backend/be-kedai/internal/utils/jwtToken"

	errRes "kedai/backend/be-kedai/internal/common/error"
)

const (
	TopUp = "Top-up"
)

type WalletService interface {
	RegisterWallet(userID int, pin string) (*model.Wallet, error)
	GetWalletByUserID(userID int) (*model.Wallet, error)
	TopUp(userId int, req dto.TopUpRequest) (*model.WalletHistory, error)
	StepUp(userId int, req dto.StepUpRequest) (*dto.Token, error)
}

type walletServiceImpl struct {
	walletRepo repository.WalletRepository
	redis      cache.UserCache
}

type WalletSConfig struct {
	WalletRepo repository.WalletRepository
	Redis      cache.UserCache
}

func NewWalletService(cfg *WalletSConfig) WalletService {
	return &walletServiceImpl{
		walletRepo: cfg.WalletRepo,
		redis:      cfg.Redis,
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

func (s *walletServiceImpl) StepUp(userId int, req dto.StepUpRequest) (*dto.Token, error) {
	wallet, err := s.walletRepo.GetByUserID(userId)
	if err != nil {
		return nil, err
	}

	if !hash.ComparePassword(wallet.Pin, req.Pin) {
		return nil, errRes.ErrWrongPin
	}

	var (
		user = &model.User{
			ID: userId,
		}
		stepUpLevel = 1
	)
	accessToken, _ := jwttoken.GenerateAccessToken(user, stepUpLevel)
	refreshToken, _ := jwttoken.GenerateRefreshToken(user, stepUpLevel)

	token := &dto.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	err = s.redis.StoreToken(userId, accessToken, refreshToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

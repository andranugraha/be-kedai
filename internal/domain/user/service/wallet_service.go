package service

import (
	"fmt"
	"kedai/backend/be-kedai/config"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/cache"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
	"kedai/backend/be-kedai/internal/utils/hash"
	jwttoken "kedai/backend/be-kedai/internal/utils/jwtToken"
)

const (
	TopUp = "Top-up"
)

type WalletService interface {
	RegisterWallet(userID int, pin string) (*model.Wallet, error)
	GetWalletByUserID(userID int) (*model.Wallet, error)
	TopUp(userId int, req dto.TopUpRequest) (*model.WalletHistory, error)
	StepUp(userId int, req dto.StepUpRequest) (*dto.Token, error)
	CheckIsWalletBlocked(userID int) error
}

type walletServiceImpl struct {
	walletRepo  repository.WalletRepository
	userCache   cache.UserCache
	walletCache cache.WalletCache
}

type WalletSConfig struct {
	WalletRepo  repository.WalletRepository
	UserCache   cache.UserCache
	WalletCache cache.WalletCache
}

func NewWalletService(cfg *WalletSConfig) WalletService {
	return &walletServiceImpl{
		walletRepo:  cfg.WalletRepo,
		userCache:   cfg.UserCache,
		walletCache: cfg.WalletCache,
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

	signString := fmt.Sprintf("%s:%d:%s", req.CardNumber, int(req.Amount), config.MerchantCode)
	hashedString := hash.HashSHA256(signString)

	isValid := hash.CompareSignature(req.Signature, hashedString)
	if !isValid {
		return nil, errs.ErrInvalidSignature
	}

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

	err = s.walletCache.CheckIsWalletBlocked(wallet.ID)
	if err != nil {
		return nil, err
	}

	if !hash.ComparePassword(wallet.Pin, req.Pin) {
		errorCount, err := s.walletCache.FindWalletStepUpErrorCount(wallet.ID)
		if err != nil {
			return nil, err
		}

		const maxErrorCount = 3
		if errorCount != nil && *errorCount+1 >= maxErrorCount {
			err = s.walletCache.DeleteErrorCount(wallet.ID)
			if err != nil {
				return nil, err
			}

			return nil, s.walletCache.BlockWallet(wallet.ID)
		}

		err = s.walletCache.StoreOrIncrementWalletStepUpErrorCount(wallet.ID)
		if err != nil {
			return nil, err
		}

		return nil, errs.ErrWrongPin
	}

	err = s.walletCache.DeleteErrorCount(wallet.ID)
	if err != nil {
		return nil, err
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

	err = s.userCache.StoreToken(userId, accessToken, refreshToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *walletServiceImpl) CheckIsWalletBlocked(userID int) error {
	wallet, err := s.walletRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	return s.walletCache.CheckIsWalletBlocked(wallet.ID)
}

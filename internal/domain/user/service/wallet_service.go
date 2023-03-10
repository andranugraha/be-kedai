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
	"kedai/backend/be-kedai/internal/utils/mail"
	"kedai/backend/be-kedai/internal/utils/random"
)

const (
	TopUp = model.WalletHistoryTypeTopup
)

type WalletService interface {
	RegisterWallet(userID int, pin string) (*model.Wallet, error)
	GetWalletByUserID(userID int) (*model.Wallet, error)
	TopUp(userId int, req dto.TopUpRequest) (*model.WalletHistory, error)
	StepUp(userId int, req dto.StepUpRequest) (*dto.Token, error)
	CheckIsWalletBlocked(userID int) error
	RequestPinChange(userID int, request *dto.ChangePinRequest) error
	CompletePinChange(userID int, request *dto.CompleteChangePinRequest) error
	RequestPinReset(userID int) error
	CompletePinReset(userID int, request *dto.CompleteResetPinRequest) error
	GetWalletDetailByUserID(userID int) (*dto.GetWalletResponse, error)
}

type walletServiceImpl struct {
	userService UserService
	walletRepo  repository.WalletRepository
	userCache   cache.UserCache
	walletCache cache.WalletCache
	randomUtils random.RandomUtils
	mailUtils   mail.MailUtils
}

type WalletSConfig struct {
	UserService UserService
	WalletRepo  repository.WalletRepository
	UserCache   cache.UserCache
	WalletCache cache.WalletCache
	RandomUtils random.RandomUtils
	MailUtils   mail.MailUtils
}

func NewWalletService(cfg *WalletSConfig) WalletService {
	return &walletServiceImpl{
		walletRepo:  cfg.WalletRepo,
		userCache:   cfg.UserCache,
		walletCache: cfg.WalletCache,
		userService: cfg.UserService,
		randomUtils: cfg.RandomUtils,
		mailUtils:   cfg.MailUtils,
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

func (s *walletServiceImpl) GetWalletDetailByUserID(userID int) (*dto.GetWalletResponse, error) {
	wallet, err := s.walletRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.GetWalletResponse{
		ID:        wallet.ID,
		Balance:   wallet.Balance,
		Number:    wallet.Number,
		IsBlocked: s.walletCache.CheckIsWalletBlocked(wallet.ID) != nil,
	}, nil
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

func (s *walletServiceImpl) RequestPinChange(userID int, request *dto.ChangePinRequest) error {
	wallet, err := s.walletRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	if isPinValid := hash.ComparePassword(wallet.Pin, request.CurrentPin); !isPinValid {
		return errs.ErrPinMismatch
	}

	codeLength := 6
	verificationCode := s.randomUtils.GenerateAlphanumericString(codeLength)
	newPin, _ := hash.HashAndSalt(request.NewPin)

	err = s.walletCache.StorePinAndVerificationCode(userID, newPin, verificationCode)
	if err != nil {
		return err
	}

	user, err := s.userService.GetByID(userID)
	if err != nil {
		return err
	}

	return s.mailUtils.SendUpdatePinEmail(user.Email, verificationCode)
}

func (s *walletServiceImpl) CompletePinChange(userID int, request *dto.CompleteChangePinRequest) error {
	newPin, verificationCode, err := s.walletCache.FindPinAndVerificationCode(userID)
	if err != nil {
		return err
	}

	if verificationCode != request.VerificationCode {
		return errs.ErrIncorrectVerificationCode
	}

	err = s.walletRepo.ChangePin(userID, newPin)
	if err != nil {
		return err
	}

	_ = s.walletCache.DeletePinAndVerificationCode(userID)

	return nil
}

func (s *walletServiceImpl) RequestPinReset(userID int) error {
	user, err := s.userService.GetByID(userID)
	if err != nil {
		return err
	}

	codeLength := 6
	verificationCode := s.randomUtils.GenerateAlphanumericString(codeLength)
	err = s.walletCache.StoreResetPinToken(user.ID, verificationCode)
	if err != nil {
		return err
	}

	return s.mailUtils.SendResetPinEmail(user.Email, verificationCode)
}

func (s *walletServiceImpl) CompletePinReset(userID int, request *dto.CompleteResetPinRequest) error {
	err := s.walletCache.FindResetPinToken(request.Token)
	if err != nil {
		return err
	}

	newPin, _ := hash.HashAndSalt(request.NewPin)

	err = s.walletRepo.ChangePin(userID, newPin)
	if err != nil {
		return err
	}

	_ = s.walletCache.DeleteResetPinToken(request.Token)

	return nil
}

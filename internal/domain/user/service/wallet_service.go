package service

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/cache"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
	"kedai/backend/be-kedai/internal/utils/hash"
	"kedai/backend/be-kedai/internal/utils/mail"
	"kedai/backend/be-kedai/internal/utils/random"
)

const (
	TopUp = "Top-up"
)

type WalletService interface {
	RegisterWallet(userID int, pin string) (*model.Wallet, error)
	GetWalletByUserID(userID int) (*model.Wallet, error)
	TopUp(userId int, req dto.TopUpRequest) (*model.WalletHistory, error)
	RequestPinChange(userID int, request *dto.ChangePinRequest) error
	CompletePinChange(userID int, request *dto.CompleteChangePinRequest) error
	RequestPinReset(userID int) error
	CompletePinReset(userID int, request *dto.CompleteResetPinRequest) error
}

type walletServiceImpl struct {
	userService UserService
	walletRepo  repository.WalletRepository
	walletCache cache.WalletCache
	randomUtils random.RandomUtils
	mailUtils   mail.MailUtils
}

type WalletSConfig struct {
	UserService UserService
	WalletRepo  repository.WalletRepository
	WalletCache cache.WalletCache
	RandomUtils random.RandomUtils
	MailUtils   mail.MailUtils
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
	newPin, verifcationCode, err := s.walletCache.FindPinAndVerificationCode(userID)
	if err != nil {
		return err
	}

	if verifcationCode != request.VerificationCode {
		return errs.ErrIncorrectVerificationCode
	}

	err = s.walletRepo.ChangePin(userID, newPin)
	if err != nil {
		return err
	}

	return s.walletCache.DeletePinAndVerificationCode(userID)
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

	err = s.mailUtils.SendResetPasswordEmail(user.Email, verificationCode)
	if err != nil {
		return err
	}

	return nil
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

	return s.walletCache.DeleteResetPinToken(request.Token)
}

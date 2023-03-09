package service

import (
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
)

type WalletHistoryService interface {
	GetHistoryDetailById(userId int, ref string) (*model.WalletHistory, error)
	GetWalletHistoryById(req dto.WalletHistoryRequest, userId int) (*commonDto.PaginationResponse, error)
}

type walletHistoryImpl struct {
	walletHistoryRepository repository.WalletHistoryRepository
	walletService           WalletService
}

type WalletHistorySConfig struct {
	WalletHistoryRepository repository.WalletHistoryRepository
	WalletService           WalletService
}

func NewWalletHistoryService(cfg *WalletHistorySConfig) WalletHistoryService {
	return &walletHistoryImpl{
		walletHistoryRepository: cfg.WalletHistoryRepository,
		walletService:           cfg.WalletService,
	}
}

func (s *walletHistoryImpl) GetHistoryDetailById(userId int, ref string) (*model.WalletHistory, error) {
	wallet, err := s.walletService.GetWalletByUserID(userId)
	if err != nil {
		return nil, err
	}

	result, err := s.walletHistoryRepository.GetHistoryDetailById(ref, wallet)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *walletHistoryImpl) GetWalletHistoryById(req dto.WalletHistoryRequest, userId int) (*commonDto.PaginationResponse, error) {
	wallet, err := s.walletService.GetWalletByUserID(userId)
	if err != nil {
		return nil, err
	}

	result, rows, pages, err := s.walletHistoryRepository.GetWalletHistoryById(req, wallet.ID)
	if err != nil {
		return nil, err
	}

	return &commonDto.PaginationResponse{
		Data:       result,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalRows:  rows,
		TotalPages: pages,
	}, nil
}

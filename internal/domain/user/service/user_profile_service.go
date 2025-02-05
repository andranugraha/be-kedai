package service

import (
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
)

type UserProfileService interface {
	UpdateProfile(userId int, request *dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error)
	GetProfile(userId int) (*model.UserProfile, error)
}

type userProfileServiceImpl struct {
	repository repository.UserProfileRepository
}

type UserProfileSConfig struct {
	Repository repository.UserProfileRepository
}

func NewUserProfileService(cfg *UserProfileSConfig) UserProfileService {
	return &userProfileServiceImpl{
		repository: cfg.Repository,
	}
}

func (s *userProfileServiceImpl) UpdateProfile(userId int, request *dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error) {
	payload := request.ToUserProfile()

	res, err := s.repository.Update(userId, payload)
	if err != nil {
		return nil, err
	}

	var response dto.UpdateProfileResponse
	response.FromUserProfile(res)

	return &response, nil
}

func (s *userProfileServiceImpl) GetProfile(userId int) (*model.UserProfile, error) {
	return s.repository.GetByID(userId)
}

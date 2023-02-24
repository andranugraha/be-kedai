package service

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/user/cache"
	"kedai/backend/be-kedai/internal/domain/user/dto"
	"kedai/backend/be-kedai/internal/domain/user/model"
	"kedai/backend/be-kedai/internal/domain/user/repository"
	"kedai/backend/be-kedai/internal/utils/credential"
	"kedai/backend/be-kedai/internal/utils/google"
	"kedai/backend/be-kedai/internal/utils/hash"
	jwttoken "kedai/backend/be-kedai/internal/utils/jwtToken"
	"strings"
)

type UserService interface {
	GetByID(id int) (*model.User, error)
	SignUp(*dto.UserRegistrationRequest) (*dto.UserRegistrationResponse, error)
	SignIn(*dto.UserLogin, string) (*dto.Token, error)
	SignInWithGoogle(userLogin *dto.UserLoginWithGoogleRequest) (*dto.Token, error)
	GetSession(userId int, token string) error
	UpdateEmail(userId int, request *dto.UpdateEmailRequest) (*dto.UpdateEmailResponse, error)
	UpdateUsername(userId int, requst *dto.UpdateUsernameRequest) (*dto.UpdateUsernameResponse, error)
}

type userServiceImpl struct {
	repository repository.UserRepository
	redis      cache.UserCache
}

type UserSConfig struct {
	Repository repository.UserRepository
	Redis      cache.UserCache
}

func NewUserService(cfg *UserSConfig) UserService {
	return &userServiceImpl{
		repository: cfg.Repository,
		redis:      cfg.Redis,
	}
}

func (s *userServiceImpl) GetByID(id int) (*model.User, error) {
	return s.repository.GetByID(id)
}

func (s *userServiceImpl) SignUp(userReg *dto.UserRegistrationRequest) (*dto.UserRegistrationResponse, error) {
	isValidPassword := credential.VerifyPassword(userReg.Password)
	if !isValidPassword {
		return nil, errs.ErrInvalidPasswordPattern
	}

	emailSplit := strings.Split(userReg.Email, "@")
	isContainEmail := strings.Contains(strings.ToLower(userReg.Password), strings.ToLower(emailSplit[0]))
	if isContainEmail {
		return nil, errs.ErrContainEmail
	}

	user := userReg.ToUser()
	user.Email = strings.ToLower(user.Email)

	result, err := s.repository.SignUp(user)
	if err != nil {
		return nil, err
	}

	var response dto.UserRegistrationResponse
	response.FromUser(result)

	return &response, nil
}

func (s *userServiceImpl) SignIn(userLogin *dto.UserLogin, inputPw string) (*dto.Token, error) {
	user := userLogin.ToUser()
	user.Email = strings.ToLower(user.Email)

	result, err := s.repository.SignIn(user)
	if err != nil {
		return nil, err
	}

	isValid := hash.ComparePassword(result.Password, inputPw)
	if isValid {
		accessToken, _ := jwttoken.GenerateAccessToken(result)
		refreshToken, _ := jwttoken.GenerateRefreshToken(result)

		token := &dto.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		err := s.redis.StoreToken(result.ID, accessToken, refreshToken)
		if err != nil {
			return nil, err
		}

		return token, nil
	}

	return nil, errs.ErrInvalidCredential
}

func (s *userServiceImpl) SignInWithGoogle(userLogin *dto.UserLoginWithGoogleRequest) (*dto.Token, error) {
	claim, err := google.ValidateGoogleToken(userLogin.Credential)
	if err != nil {
		return nil, errs.ErrUnauthorized
	}

	result, err := s.repository.SignIn(&model.User{Email: claim.Email})
	if err != nil {
		return nil, err
	}

	accessToken, _ := jwttoken.GenerateAccessToken(result)
	refreshToken, _ := jwttoken.GenerateRefreshToken(result)

	token := &dto.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	err = s.redis.StoreToken(result.ID, accessToken, refreshToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *userServiceImpl) GetSession(userId int, accessToken string) error {
	return s.redis.FindToken(userId, accessToken)
}

func (s *userServiceImpl) UpdateEmail(userId int, request *dto.UpdateEmailRequest) (*dto.UpdateEmailResponse, error) {
	email := strings.ToLower(request.Email)

	res, err := s.repository.UpdateEmail(userId, email)
	if err != nil {
		return nil, err
	}

	var response dto.UpdateEmailResponse
	response.FromUser(res)

	return &response, nil
}

func (s *userServiceImpl) UpdateUsername(userId int, request *dto.UpdateUsernameRequest) (*dto.UpdateUsernameResponse, error) {
	username := strings.ToLower(request.Username)

	if isUsernameValid := credential.VerifyUsername(username); !isUsernameValid {
		return nil, errs.ErrInvalidUsernamePattern
	}

	res, err := s.repository.UpdateUsername(userId, username)
	if err != nil {
		return nil, err
	}

	response := dto.UpdateUsernameResponse{
		Username: res.Username,
	}

	return &response, nil
}

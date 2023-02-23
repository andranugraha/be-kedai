package dto

import "kedai/backend/be-kedai/internal/domain/user/model"

type UserRegistrationRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password,omitempty" binding:"required,min=8,max=16"`
}

type UserRegistrationResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginWithGoogleRequest struct {
	Credential string `json:"credential" binding:"required"`
}

type UserLoginWithGoogle struct {
	Email string
}

func (d *UserRegistrationRequest) ToUser() *model.User {
	return &model.User{
		Email:    d.Email,
		Password: d.Password,
	}
}

func (d *UserRegistrationResponse) FromUser(user *model.User) {
	d.Email = user.Email
	d.Username = user.Username
}

func (d *UserLogin) ToUser() *model.User {
	return &model.User{
		Email:    d.Email,
		Password: d.Password,
	}
}

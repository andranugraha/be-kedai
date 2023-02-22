package dto

import "kedai/backend/be-kedai/internal/domain/user/model"

type UserRegistration struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password,omitempty" binding:"required,min=8,max=16"`
	Username string `json:"username"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (d *UserRegistration) ToUser() *model.User {
	return &model.User{
		Email:    d.Email,
		Password: d.Password,
	}
}

func (d *UserRegistration) FromUser(user *model.User) {
	d.Email = user.Email
	d.Username = user.Username
}

func (d *UserLogin) ToUser() *model.User {
	return &model.User{
		Email:    d.Email,
		Password: d.Password,
	}
}

package dto

import "kedai/backend/be-kedai/internal/domain/user/model"

type UserRegistration struct {
	Email    string       `json:"email" binding:"required,email"`
	Password string       `json:"password" binding:"required,min=6"`
}

func (d *UserRegistration) ToUser() *model.User {
	return &model.User{
		Email: d.Email,
		Password: d.Password,
	}
}
package dto

import (
	"kedai/backend/be-kedai/internal/domain/user/model"
	"time"
)

type UpdateProfileRequest struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required,numeric,min=10,max=15"`
	DoB         string `json:"dob" binding:"required,datetime=2006-01-02"`
	Gender      string `json:"gender" binding:"required,oneof=male female others"`
	PhotoUrl    string `json:"photoUrl" binding:"required,url"`
}

type UpdateProfileResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phoneNumber"`
	DoB         time.Time `json:"dob"`
	Gender      string    `json:"gender"`
	PhotoUrl    string    `json:"photoUrl"`
}

func (d *UpdateProfileRequest) ToUserProfile() *model.UserProfile {
	dob, _ := time.Parse("2006-01-02", d.DoB)

	return &model.UserProfile{
		Name:        &d.Name,
		PhoneNumber: &d.PhoneNumber,
		DoB:         &dob,
		Gender:      &d.Gender,
		PhotoUrl:    &d.PhotoUrl,
	}
}

func (d *UpdateProfileResponse) FromUserProfile(profile *model.UserProfile) {
	d.ID = profile.ID
	d.Name = *profile.Name
	d.PhoneNumber = *profile.PhoneNumber
	d.DoB = *profile.DoB
	d.Gender = *profile.Gender
	d.PhotoUrl = *profile.PhotoUrl
}

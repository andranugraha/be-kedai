package model

import "gorm.io/gorm"

type ReviewMedia struct {
	ID       int    `json:"id"`
	Url      string `json:"url"`
	ReviewId int    `json:"reviewId"`

	gorm.Model `json:"-"`
}

func (ReviewMedia) TableName() string {
	return "review_medias"
}

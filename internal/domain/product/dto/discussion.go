package dto

import (
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"
	"time"

	"gorm.io/gorm"
)

type Discussion struct {
	ID         int              `json:"id"`
	UserID     int              `json:"-"`
	Username   string           `json:"username"`
	UserUrl    string           `json:"userUrl"`
	User       *userModel.User  `json:"-"`
	ShopId     int              `json:"-"`
	Shop       *shopModel.Shop  `json:"-" gorm:"foreignKey:ShopId"`
	ShopName   string           `json:"shopName,omitempty"`
	ShopUrl    string           `json:"shopUrl,omitempty"`
	ProductID  int              `json:"productId"`
	Message    string           `json:"message"`
	Date       time.Time        `json:"date"`
	Reply      *DiscussionReply `json:"reply" gorm:"foreignKey:ID"`
	ReplyCount int              `json:"replyCount"`

	gorm.Model `json:"-"`
}

type DiscussionReply struct {
	ID        int             `json:"id"`
	UserID    int             `json:"-"`
	Username  string          `json:"username"`
	UserUrl   string          `json:"userUrl"`
	User      *userModel.User `json:"-"`
	ShopId    int             `json:"-"`
	Shop      *shopModel.Shop `json:"-" gorm:"foreignKey:ShopId"`
	ShopName  string          `json:"shopName,omitempty"`
	ShopUrl   string          `json:"shopUrl,omitempty"`
	ProductID int             `json:"productId"`
	ParentID  int             `json:"parentId"`
	Message   string          `json:"message"`
	Date      string          `json:"date"`

	gorm.Model `json:"-"`
}

type DiscussionReq struct {
	ID        int       `json:"id"`
	UserID    int       `json:"-"`
	ProductID int       `json:"productId" binding:"required"`
	Message   string    `json:"message" binding:"required"`
	Date      time.Time `json:"date"`
	ParentID  *int      `json:"parentId"`
	ShopID    int       `json:"-"`
	IsSeller  *bool     `json:"isSeller"`
}

type GetDiscussionReq struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

func (req *GetDiscussionReq) Validate() {
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 1 {
		req.Limit = 5
	}
}

func (req *GetDiscussionReq) GetOffset() int {
	return (req.Page - 1) * req.Limit
}

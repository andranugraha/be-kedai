package dto

import (
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	shopModel "kedai/backend/be-kedai/internal/domain/shop/model"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"
	"time"

	"kedai/backend/be-kedai/internal/common/constant"

	"gorm.io/gorm"
)

type Discussion struct {
	ID         int                   `json:"id"`
	UserID     int                   `json:"-"`
	Username   string                `json:"username" gorm:"-"`
	UserUrl    string                `json:"userUrl" gorm:"-"`
	User       *userModel.User       `json:"-"`
	ShopId     int                   `json:"-"`
	Shop       *shopModel.Shop       `json:"-" gorm:"foreignKey:ShopId"`
	ShopName   string                `json:"shopName,omitempty" gorm:"-"`
	ShopUrl    string                `json:"shopUrl,omitempty" gorm:"-"`
	ProductID  int                   `json:"productId"`
	Product    *productModel.Product `json:"product,omitempty"`
	Message    string                `json:"message"`
	Date       time.Time             `json:"date"`
	Reply      *DiscussionReply      `json:"reply" gorm:"foreignKey:ID"`
	ReplyCount int                   `json:"replyCount" gorm:"-"`

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
	ID        int        `json:"id"`
	UserID    int        `json:"-"`
	ProductID int        `json:"productId" binding:"required"`
	Message   string     `json:"message" binding:"required"`
	Date      *time.Time `json:"date" gorm:"default:CURRENT_TIMESTAMP"`
	ParentID  *int       `json:"parentId"`
	ShopID    *int       `json:"-"`
	IsSeller  bool       `json:"isSeller"`
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
		req.Limit = constant.DefaultDiscussionLimit
	}

	if req.Limit > 10 {
		req.Limit = constant.MaxDiscussionLimit
	}
}

func (req *GetDiscussionReq) Offset() int {
	return (req.Page - 1) * req.Limit
}

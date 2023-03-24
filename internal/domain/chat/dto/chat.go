package dto

import (
	chatModel "kedai/backend/be-kedai/internal/domain/chat/model"
	"time"
)

type ListOfChatsParamRequest struct {
	Search string `form:"search"`
	Status string `form:"status"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

type ChatParamRequest struct {
	Page       int `form:"page"`
	LimitByDay int `form:"limitByDay"`
}

func (cpr *ChatParamRequest) Validate() {
	if cpr.LimitByDay < 1 {
		cpr.LimitByDay = 366
	}
	if cpr.Page < 1 {
		cpr.Page = 1
	}
}

func (cpr *ChatParamRequest) Offset() int {
	return int((cpr.Page - 1) * cpr.LimitByDay)
}

type SendChatBodyRequest struct {
	Message string `json:"message" binding:"required"`
	Type    string `json:"type"`
}

type UserChatProfile struct {
	Username string `json:"username"`
	ImageUrl string `json:"imageUrl"`
}

type ShopChatProfile struct {
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
	ShopSlug string `json:"shopSlug"`
}

type UserListOfChatResponse struct {
	Shop          *ShopChatProfile `json:"shop"`
	RecentMessage string           `json:"recentMessage"`
	UnreadCount   int              `json:"unreadCount"`
}

type SellerListOfChatResponse struct {
	User          *UserChatProfile `json:"user"`
	RecentMessage string           `json:"recentMessage"`
	UnreadCount   int              `json:"unreadCount"`
}

type ChatResponse struct {
	ID         int       `json:"id"`
	Message    string    `json:"message"`
	Time       time.Time `json:"time"`
	Type       string    `json:"type"`
	IsIncoming bool      `json:"isIncoming"`
}

func ConvertChatToOutput(c *chatModel.Chat, role string) *ChatResponse {
	// Role: user | seller
	if c == nil {
		return nil
	}
	return &ChatResponse{
		ID:         c.ID,
		Message:    c.Message,
		Time:       c.CreatedAt,
		Type:       c.Type,
		IsIncoming: incomingSelector(c.Issuer, role),
	}
}

func incomingSelector(issuer string, role string) bool {
	if role == "user" {
		return issuer != "user"
	}

	if role == "seller" {
		return issuer != "seller"
	}

	return true
}

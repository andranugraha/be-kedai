package dto

import (
	chatModel "kedai/backend/be-kedai/internal/domain/chat/model"
	"time"
)

type ChatRequest struct {
	RoomId  string `json:"roomId" binding:"required"`
	Message string `json:"message" binding:"required"`
	Type    string `json:"type" binding:"required"`
}

type UserChatProfile struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	PhotoUrl string `json:"photoUrl"`
}

type ChatListResponse struct {
	RoomId        string           `json:"roomId"`
	LatestMessage string           `json:"latestMessage"`
	UnreadCount   int              `json:"unreadMessageCount"`
	OppositeUser  *UserChatProfile `json:"oppositeUser"`
	CreatedAt     time.Time        `json:"createdAt"`
}

type ChatDetailResponse struct {
	ID        int              `json:"id"`
	RoomId    string           `json:"roomId"`
	Message   string           `json:"message"`
	Type      string           `json:"type"`
	User      *UserChatProfile `json:"user"`
	CreatedAt time.Time        `json:"createdAt"`
}

func ConvertChatDetailToOutput(c *chatModel.Chat) *ChatDetailResponse {
	if c == nil {
		return nil
	}
	return &ChatDetailResponse{
		ID:      c.ID,
		RoomId:  c.RoomId,
		Message: c.Message,
		Type:    c.Type,
		User: &UserChatProfile{
			ID:       c.User.ID,
			Name:     *c.User.Profile.Name,
			PhotoUrl: *c.User.Profile.PhotoUrl,
		},
		CreatedAt: c.CreatedAt,
	}
}

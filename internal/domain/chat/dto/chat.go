package dto

import (
	chatModel "kedai/backend/be-kedai/internal/domain/chat/model"
	"time"
)

type ListOfChatsParamRequest struct {
	Search string `form:"search"`
	Status string `form:"status"`
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
	ID       int     `json:"-"`
	Username string  `json:"username"`
	ImageUrl *string `json:"imageUrl"`
}

type ShopChatProfile struct {
	ID       int     `json:"-"`
	Name     string  `json:"name"`
	ImageUrl *string `json:"imageUrl"`
	ShopSlug string  `json:"shopSlug"`
}

type UserListOfChatResponse struct {
	Shop              *ShopChatProfile `json:"shop"`
	RecentMessage     string           `json:"recentMessage"`
	RecentMessageType string           `json:"recentMessageType"`
	UnreadCount       int              `json:"unreadCount"`
}

type SellerListOfChatResponse struct {
	User              *UserChatProfile `json:"user"`
	RecentMessage     string           `json:"recentMessage"`
	RecentMessageType string           `json:"recentMessageType"`
	UnreadCount       int              `json:"unreadCount"`
}

type ChatResponse struct {
	ID                  int       `json:"id"`
	Message             string    `json:"message"`
	Time                time.Time `json:"time"`
	Type                string    `json:"type"`
	IsIncoming          bool      `json:"isIncoming"`
	IsFirstMessageOfDay *bool     `json:"isFirstMessageOfDay,omitempty"`
}

func ConvertChatToOutput(c *chatModel.Chat, role string) *ChatResponse {
	// Role: user | seller
	if c == nil {
		return nil
	}
	return &ChatResponse{
		ID:      c.ID,
		Message: c.Message,
		Time: time.Date(
			c.CreatedAt.Year(),
			c.CreatedAt.Month(),
			c.CreatedAt.Day(),
			c.CreatedAt.Hour(),
			c.CreatedAt.Minute(),
			c.CreatedAt.Second(),
			c.CreatedAt.Nanosecond(),
			time.FixedZone("WIB", 7*60*60)).Add(7 * time.Hour),
		Type:       c.Type,
		IsIncoming: incomingSelector(c.Issuer, role),
	}
}

func AssignFirstMessagesOfDay(chatResponses []*ChatResponse) { // pointer mutation direcly
	// Assumption: chat responses ordered by newest to oldest
	// Reverse slice first
	for i := 0; i < len(chatResponses)/2; i++ {
		j := len(chatResponses) - i - 1
		chatResponses[i], chatResponses[j] = chatResponses[j], chatResponses[i]
	}

	// Assign true if the message is the first of the day
	var lastDate time.Time
	trueVal := true
	falseVal := false
	for i, chatResponse := range chatResponses {
		if i == 0 {
			chatResponse.IsFirstMessageOfDay = &trueVal
			lastDate = time.Date(
				chatResponse.Time.Year(),
				chatResponse.Time.Month(),
				chatResponse.Time.Day(),
				0, 0, 0, 0,
				time.FixedZone("WIB", 7*60*60)).Add(7 * time.Hour)
			continue
		}
		currentDate := time.Date(
			chatResponse.Time.Year(),
			chatResponse.Time.Month(),
			chatResponse.Time.Day(),
			0, 0, 0, 0,
			time.FixedZone("WIB", 7*60*60)).Add(7 * time.Hour)
		// time comparation in WIB (+07:00)
		if currentDate.After(lastDate) {
			chatResponse.IsFirstMessageOfDay = &trueVal
			lastDate = currentDate
		} else {
			chatResponse.IsFirstMessageOfDay = &falseVal
		}
	}

	// Reverse slice again
	for i := 0; i < len(chatResponses)/2; i++ {
		j := len(chatResponses) - i - 1
		chatResponses[i], chatResponses[j] = chatResponses[j], chatResponses[i]
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

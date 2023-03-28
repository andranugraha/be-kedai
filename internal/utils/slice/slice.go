package slice

import "kedai/backend/be-kedai/internal/domain/chat/dto"

func Contains(s []int, number int) bool {
	for _, v := range s {
		if v == number {
			return true
		}
	}
	return false
}

func SellerRemoveElement(slice []*dto.SellerListOfChatResponse, index int) []*dto.SellerListOfChatResponse {
	return append(slice[:index], slice[index+1:]...)
}

func UserRemoveElement(slice []*dto.UserListOfChatResponse, index int) []*dto.UserListOfChatResponse {
	return append(slice[:index], slice[index+1:]...)
}

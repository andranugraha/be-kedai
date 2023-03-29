package slice_test

import (
	"kedai/backend/be-kedai/internal/domain/chat/dto"
	sliceUtil "kedai/backend/be-kedai/internal/utils/slice"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	number1 := 3
	number2 := 6

	result1 := sliceUtil.Contains(s, number1)
	result2 := sliceUtil.Contains(s, number2)

	assert.True(t, result1)
	assert.False(t, result2)
}

func TestSellerRemoveElement(t *testing.T) {
	slice := []*dto.SellerListOfChatResponse{{}, {}, {}}
	index := 1
	expectedResult := []*dto.SellerListOfChatResponse{{}, {}}

	result := sliceUtil.SellerRemoveElement(slice, index)

	assert.Equal(t, expectedResult, result)
}

func TestUserRemoveElement(t *testing.T) {
	slice := []*dto.UserListOfChatResponse{{}, {}, {}}
	index := 2
	expectedResult := []*dto.UserListOfChatResponse{{}, {}}

	result := sliceUtil.UserRemoveElement(slice, index)

	assert.Equal(t, expectedResult, result)
}

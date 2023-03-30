package slice_test

import (
	"fmt"
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

func BenchmarkContains(b *testing.B) {
	var table = []struct {
		name   string
		s      []int
		number int
	}{
		{name: "small_slice", s: []int{1, 2, 3, 4, 5}, number: 3},
		{name: "medium_slice", s: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, number: 9},
		{name: "large_slice", s: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, number: 19},
	}

	for _, v := range table {
		b.Run(fmt.Sprintf("input_type_%s", v.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sliceUtil.Contains(v.s, v.number)
			}
		})
	}
}

func BenchmarkSellerRemoveElement(b *testing.B) {
	// create slice of length n to test with
	n := 100
	slice := make([]*dto.SellerListOfChatResponse, n)
	for i := 0; i < n; i++ {
		slice[i] = &dto.SellerListOfChatResponse{
			User: &dto.UserChatProfile{
				ID:       i,
				Username: fmt.Sprintf("user%d", i),
				ImageUrl: nil,
			},
			RecentMessage:     "",
			RecentMessageType: "",
			UnreadCount:       0,
		}
	}

	// perform benchmark on removing element at index 0
	b.Run("remove_first", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sliceUtil.SellerRemoveElement(slice, 0)
		}
	})

	// perform benchmark on removing element at index n-1
	b.Run("remove_last", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sliceUtil.SellerRemoveElement(slice, n-1)
		}
	})

	// perform benchmark on removing element at index n/2
	b.Run("remove_middle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sliceUtil.SellerRemoveElement(slice, n/2)
		}
	})
}

func BenchmarkUserRemoveElement(b *testing.B) {
	// create slice of length n to test with
	n := 100
	slice := make([]*dto.UserListOfChatResponse, n)
	for i := 0; i < n; i++ {
		slice[i] = &dto.UserListOfChatResponse{
			Shop: &dto.ShopChatProfile{
				ID:       i,
				Name:     fmt.Sprintf("shop%d", i),
				ShopSlug: fmt.Sprintf("shop-%d", i),
				ImageUrl: nil,
			},
			RecentMessage:     "",
			RecentMessageType: "",
			UnreadCount:       0,
		}
	}

	// perform benchmark on removing element at index 0
	b.Run("remove_first", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sliceUtil.UserRemoveElement(slice, 0)
		}
	})

	// perform benchmark on removing element at index n-1
	b.Run("remove_last", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sliceUtil.UserRemoveElement(slice, n-1)
		}
	})

	// perform benchmark on removing element at index n/2
	b.Run("remove_middle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sliceUtil.UserRemoveElement(slice, n/2)
		}
	})
}

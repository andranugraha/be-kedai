package service_test

import (
	"kedai/backend/be-kedai/internal/domain/shop/model"
	"kedai/backend/be-kedai/internal/domain/shop/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateShopGuest(t *testing.T) {
	type input struct {
		shopId int
	}

	type expected struct {
		shopGuest *model.ShopGuest
		err       error
	}

	cases := []struct {
		description string
		input       input
		expected    expected
	}{
		{
			description: "return error and shop guest",
			input: input{
				shopId: 1,
			},
			expected: expected{
				shopGuest: &model.ShopGuest{
					UUID:   "1",
					ShopId: 1,
				},
				err: nil,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			shopGuestRepository := new(mocks.ShopGuestRepository)
			shopGuestRepository.On("CreateShopGuest", &model.ShopGuest{ShopId: c.input.shopId}).Return(c.expected.shopGuest, c.expected.err)

			shopGuestService := service.NewShopGuestService(&service.ShopGuestSConfig{
				ShopGuestRepository: shopGuestRepository,
			})

			shopGuest, err := shopGuestService.CreateShopGuest(c.input.shopId)

			assert.Equal(t, c.expected.shopGuest, shopGuest)
			assert.Equal(t, c.expected.err, err)
		})
	}

}

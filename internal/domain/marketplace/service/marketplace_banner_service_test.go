package service_test

import (
	"errors"
	"kedai/backend/be-kedai/internal/domain/marketplace/model"
	"kedai/backend/be-kedai/internal/domain/marketplace/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMarketplaceBanner(t *testing.T) {
	banners := []*model.MarketplaceBanner{}
	type input struct {
		beforeTest func(*mocks.MarketplaceBannerRepository)
	}
	type expected struct {
		result []*model.MarketplaceBanner
		err    error
	}
	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when fails to fetch marketplace banner",
			input: input{
				beforeTest: func(m *mocks.MarketplaceBannerRepository) {
					m.On("GetMarketplaceBanner").Return(nil, errors.New("internal server error"))
				},
			},
			expected: expected{
				result: nil,
				err:    errors.New("internal server error"),
			},
		},
		{
			description: "should return list of active banners on success",
			input: input{
				beforeTest: func(m *mocks.MarketplaceBannerRepository) {
					m.On("GetMarketplaceBanner").Return(banners, nil)
				},
			},
			expected: expected{
				result: banners,
				err:    nil,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			m := mocks.NewMarketplaceBannerRepository(t)
			tc.beforeTest(m)

			s := service.NewMarketplaceBannerService(&service.MarketplaceBannerSConfig{
				MarketplaceBannerRepository: m,
			})

			result, err := s.GetMarketplaceBanner()

			assert.Equal(t, tc.expected.err, err)
			assert.Equal(t, tc.expected.result, result)
		})
	}
}

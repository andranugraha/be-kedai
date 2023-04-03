package service_test

import (
	"errors"
	commonErr "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/marketplace/dto"
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

func TestAddMarketplaceBanner(t *testing.T) {
	banner := &model.MarketplaceBanner{}
	type input struct {
		body       *dto.MarketplaceBannerRequest
		beforeTest func(m *mocks.MarketplaceBannerRepository)
	}
	type expected struct {
		result *model.MarketplaceBanner
		err    error
	}
	cases := []struct {
		description string
		input
		expected
	}{
		{
			description: "should return error when fails validate datetime format",
			input: input{
				body:       &dto.MarketplaceBannerRequest{StartDate: "", EndDate: ""},
				beforeTest: func(m *mocks.MarketplaceBannerRepository) {},
			},
			expected: expected{
				result: nil,
				err:    commonErr.ErrInvalidRFC3999Nano,
			},
		},
		{
			description: "should return error when back date",
			input: input{
				body:       &dto.MarketplaceBannerRequest{StartDate: "2023-03-28T12:00:00.13Z", EndDate: "2023-03-28T12:00:00.12Z"},
				beforeTest: func(m *mocks.MarketplaceBannerRepository) {},
			},
			expected: expected{
				result: nil,
				err:    commonErr.ErrBackDate,
			},
		},
		{
			description: "should return success when successfully add banner",
			input: input{
				body: &dto.MarketplaceBannerRequest{StartDate: "2023-03-28T12:00:00.12Z", EndDate: "2023-03-28T12:00:00.13Z"},
				beforeTest: func(m *mocks.MarketplaceBannerRepository) {
					m.On("AddMarketplaceBanner", &dto.MarketplaceBannerRequest{StartDate: "2023-03-28T12:00:00.12Z", EndDate: "2023-03-28T12:00:00.13Z"}).Return(banner, nil)
				},
			},
			expected: expected{
				result: banner,
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

			result, err := s.AddMarketplaceBanner(tc.body)

			assert.Equal(t, tc.expected.err, err)
			assert.Equal(t, tc.expected.result, result)
		})
	}
}

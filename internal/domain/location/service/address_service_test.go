package service_test

import (
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/service"
	"kedai/backend/be-kedai/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchAddress(t *testing.T) {
	tests := []struct {
		name    string
		req     *dto.SearchAddressRequest
		want    []*dto.SearchAddressResponse
		wantErr error
	}{
		{
			name: "should return list of address when search address success",
			req: &dto.SearchAddressRequest{
				Keyword: "Jalan Puncak Pesanggrahan VI No. 5",
			},
			want:    []*dto.SearchAddressResponse{},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			addressRepo := mocks.NewAddressRepository(t)
			addressRepo.On("SearchAddress", test.req).Return(test.want, nil)
			addressService := service.NewAddressService(&service.AddressSConfig{
				AddressRepo: addressRepo,
			})

			got, err := addressService.SearchAddress(test.req)

			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantErr)
		})
	}
}

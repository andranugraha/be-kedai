package repository

import (
	"context"
	"kedai/backend/be-kedai/internal/domain/location/dto"

	"googlemaps.github.io/maps"
	"gorm.io/gorm"
)

type AddressRepository interface {
	SearchAddress(req *dto.SearchAddressRequest) ([]*dto.SearchAddressResponse, error)
}

type addressRepositoryImpl struct {
	db         *gorm.DB
	googleMaps *maps.Client
}

type AddressRConfig struct {
	DB         *gorm.DB
	GoogleMaps *maps.Client
}

func NewAddressRepository(cfg *AddressRConfig) AddressRepository {
	return &addressRepositoryImpl{
		db:         cfg.DB,
		googleMaps: cfg.GoogleMaps,
	}
}

func (c *addressRepositoryImpl) SearchAddress(req *dto.SearchAddressRequest) (addresses []*dto.SearchAddressResponse, err error) {
	ctx := context.Background()
	defer ctx.Done()

	autoCompleteRequest := &maps.PlaceAutocompleteRequest{
		Input: req.Keyword,
		Components: map[maps.Component][]string{
			maps.ComponentCountry: {"id"},
		},
		Language: "id",
	}

	autocomplete, err := c.googleMaps.PlaceAutocomplete(ctx, autoCompleteRequest)
	if err != nil {
		return
	}

	for _, place := range autocomplete.Predictions {
		addresses = append(addresses, &dto.SearchAddressResponse{
			PlaceID:     place.PlaceID,
			Description: place.Description,
		})
	}

	return
}

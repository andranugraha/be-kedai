package cache

import (
	"context"
	"encoding/json"
	"fmt"
	commonDto "kedai/backend/be-kedai/internal/common/dto"
	"kedai/backend/be-kedai/internal/domain/location/dto"
	"kedai/backend/be-kedai/internal/domain/location/model"

	"github.com/redis/go-redis/v9"
)

type LocationCache interface {
	GetSubdistricts(req dto.GetSubdistrictsRequest) []*model.Subdistrict
	StoreSubdistricts(req dto.GetSubdistrictsRequest, subdistricts []*model.Subdistrict)
	GetDistricts(req dto.GetDistrictsRequest) []*model.District
	StoreDistricts(req dto.GetDistrictsRequest, districts []*model.District)
	GetCities(req dto.GetCitiesRequest) *commonDto.PaginationResponse
	StoreCities(req dto.GetCitiesRequest, cities *commonDto.PaginationResponse)
	GetProvinces() []*model.Province
	StoreProvinces(provinces []*model.Province)
}

type locationCacheImpl struct {
	rdc *redis.Client
}

type LocationCConfig struct {
	RDC *redis.Client
}

func NewLocationCache(cfg *LocationCConfig) LocationCache {
	return &locationCacheImpl{
		rdc: cfg.RDC,
	}
}

func (c *locationCacheImpl) GetSubdistricts(req dto.GetSubdistrictsRequest) []*model.Subdistrict {
	key := fmt.Sprintf("subdistricts:district_id_%d", req.DistrictID)
	subdistricts, err := c.rdc.Get(context.Background(), key).Result()
	if err != nil {
		return nil
	}

	subdistrictsObject := []*model.Subdistrict{}
	err = json.Unmarshal([]byte(subdistricts), &subdistrictsObject)
	if err != nil {
		return nil
	}

	return subdistrictsObject
}

func (c *locationCacheImpl) StoreSubdistricts(req dto.GetSubdistrictsRequest, subdistricts []*model.Subdistrict) {
	key := fmt.Sprintf("subdistricts:district_id_%d", req.DistrictID)
	subdistrictsJSON, _ := json.Marshal(subdistricts)

	_ = c.rdc.Set(context.Background(), key, subdistrictsJSON, 0)
}

func (c *locationCacheImpl) GetDistricts(req dto.GetDistrictsRequest) []*model.District {
	key := fmt.Sprintf("districts:city_id_%d", req.CityID)
	districts, err := c.rdc.Get(context.Background(), key).Result()
	if err != nil {
		return nil
	}

	districtsObject := []*model.District{}
	err = json.Unmarshal([]byte(districts), &districtsObject)
	if err != nil {
		return nil
	}

	return districtsObject
}

func (c *locationCacheImpl) StoreDistricts(req dto.GetDistrictsRequest, districts []*model.District) {
	key := fmt.Sprintf("districts:city_id_%d", req.CityID)
	districtsJSON, _ := json.Marshal(districts)

	_ = c.rdc.Set(context.Background(), key, districtsJSON, 0)
}

func (c *locationCacheImpl) GetCities(req dto.GetCitiesRequest) *commonDto.PaginationResponse {
	key := fmt.Sprintf("cities:province_id_%d:limit_%d:page_%d:sort_%s", req.ProvinceID, req.Limit, req.Page, req.Sort)
	cities, err := c.rdc.Get(context.Background(), key).Result()
	if err != nil {
		return nil
	}

	citiesObject := &commonDto.PaginationResponse{}
	err = json.Unmarshal([]byte(cities), &citiesObject)
	if err != nil {
		return nil
	}

	return citiesObject
}

func (c *locationCacheImpl) StoreCities(req dto.GetCitiesRequest, cities *commonDto.PaginationResponse) {
	key := fmt.Sprintf("cities:province_id_%d:limit_%d:page_%d:sort_%s", req.ProvinceID, req.Limit, req.Page, req.Sort)
	citiesJSON, _ := json.Marshal(cities)

	_ = c.rdc.Set(context.Background(), key, citiesJSON, 0)
}

func (c *locationCacheImpl) GetProvinces() []*model.Province {
	provinces, err := c.rdc.Get(context.Background(), "provinces").Result()
	if err != nil {
		return nil
	}

	provincesObject := []*model.Province{}
	err = json.Unmarshal([]byte(provinces), &provincesObject)
	if err != nil {
		return nil
	}

	return provincesObject
}

func (c *locationCacheImpl) StoreProvinces(provinces []*model.Province) {
	provincesJSON, _ := json.Marshal(provinces)

	_ = c.rdc.Set(context.Background(), "provinces", provincesJSON, 0)
}

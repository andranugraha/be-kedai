package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"kedai/backend/be-kedai/internal/common/dto"
	categoryDto "kedai/backend/be-kedai/internal/domain/product/dto"
	"time"

	"github.com/redis/go-redis/v9"
)

type CategoryCache interface {
	GetAll(query categoryDto.GetCategoriesRequest) *dto.PaginationResponse
	StoreCategories(query categoryDto.GetCategoriesRequest, categories *dto.PaginationResponse)
	ResetCategories()
}

type categoryCacheImpl struct {
	rdc *redis.Client
}

type CategoryCConfig struct {
	RDC *redis.Client
}

func NewCategoryCache(cfg *CategoryCConfig) CategoryCache {
	return &categoryCacheImpl{
		rdc: cfg.RDC,
	}
}

func (c *categoryCacheImpl) GetAll(query categoryDto.GetCategoriesRequest) *dto.PaginationResponse {
	key := fmt.Sprintf("categories:page_%d:limit_%d:depth_%d:with_price_%t:parent_id_%d", query.Page, query.Limit, query.Depth, query.WithPrice, query.ParentID)
	categories, err := c.rdc.Get(context.Background(), key).Result()
	if err != nil {
		return nil
	}

	paginationObject := &dto.PaginationResponse{}
	err = json.Unmarshal([]byte(categories), paginationObject)
	if err != nil {
		return nil
	}

	return paginationObject
}

func (c *categoryCacheImpl) StoreCategories(query categoryDto.GetCategoriesRequest, categories *dto.PaginationResponse) {
	key := fmt.Sprintf("categories:page_%d:limit_%d:depth_%d:with_price_%t:parent_id_%d", query.Page, query.Limit, query.Depth, query.WithPrice, query.ParentID)
	categoriesJSON, _ := json.Marshal(categories)

	expireTime := time.Hour * 3
	_ = c.rdc.Set(context.Background(), key, categoriesJSON, expireTime)
}

func (c *categoryCacheImpl) ResetCategories() {
	keys, _ := c.rdc.Keys(context.Background(), "categories:*").Result()

	if len(keys) > 0 {
		_ = c.rdc.Del(context.Background(), keys...)
	}
}

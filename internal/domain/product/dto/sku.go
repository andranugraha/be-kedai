package dto

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/model"
	"kedai/backend/be-kedai/internal/utils/random"
	"strconv"
	"strings"
)

type GetSKURequest struct {
	VariantID string `form:"variantId" binding:"required"`
}

func (d *GetSKURequest) ToIntList() ([]int, error) {
	variantList := strings.Split(d.VariantID, ",")

	variant1, err := strconv.Atoi(variantList[0])
	if err != nil || variant1 < 1 {
		return nil, errs.ErrInvalidVariantID
	}

	if len(variantList) == 1 {
		return []int{variant1}, nil
	}

	variant2, err := strconv.Atoi(variantList[1])
	if err != nil || variant2 < 1 {
		return nil, errs.ErrInvalidVariantID
	}

	return []int{variant1, variant2}, nil
}

type CreateSKURequest struct {
	Sku        string  `json:"sku" binding:"omitempty,max=16"`
	VariantIDs []int   `json:"variantIds" binding:"required,max=2,dive,gte=0" gorm:"-"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	Stock      int     `json:"stock" binding:"required,gte=0"`
}

func (d *CreateProductRequest) GenerateSKU(variantGroup []*model.VariantGroup) []*model.Sku {
	randomUtils := random.NewRandomUtils(&random.RandomUtilsConfig{})
	skuLength := 16

	if d.VariantGroups == nil {
		return []*model.Sku{{
			Price: d.Price,
			Stock: d.Stock,
			Sku:   randomUtils.GenerateAlphanumericString(skuLength),
		}}
	}

	skuList := []*model.Sku{}
	for _, s := range d.SKU {
		var sku string
		if s.Sku == "" {
			sku = randomUtils.GenerateAlphanumericString(skuLength)
		} else {
			sku = s.Sku
		}

		variants := []model.Variant{}
		for groupIdx, variantIdx := range s.VariantIDs {
			variants = append(variants, *variantGroup[groupIdx].Variant[variantIdx])
		}

		skuList = append(skuList, &model.Sku{
			Price:    s.Price,
			Stock:    s.Stock,
			Sku:      sku,
			Variants: variants,
		})
	}

	return skuList
}

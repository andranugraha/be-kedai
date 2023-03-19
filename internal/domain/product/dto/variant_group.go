package dto

import (
	"kedai/backend/be-kedai/internal/domain/product/model"
)

type CreateVariantGroupRequest struct {
	Name    string               `json:"name" binding:"required,max=14"`
	Variant []*AddVariantRequest `json:"variants" binding:"required,min=1,max=50,dive"`
}

func (d *CreateProductRequest) GenerateVariantGroups() []*model.VariantGroup {
	if d.VariantGroups == nil {
		return nil
	}

	variantGroups := []*model.VariantGroup{}
	for _, req := range d.VariantGroups {
		variants := []*model.Variant{}

		for _, v := range req.Variant {
			variants = append(variants, &model.Variant{
				Value:    v.Name,
				MediaUrl: v.MediaUrl,
			})
		}

		variantGroups = append(variantGroups, &model.VariantGroup{
			Name:    req.Name,
			Variant: variants,
		})
	}

	return variantGroups
}

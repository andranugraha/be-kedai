package dto

import "kedai/backend/be-kedai/internal/domain/product/model"

type CreateVariantGroupRequest struct {
	Name     string   `json:"name" binding:"required,max=20"`
	Options  []string `json:"options" binding:"required,min=1,max=50,dive,max=14"`
	MediaUrl string   `json:"mediaUrl" binding:"omitempty,url"`
}

func (d *CreateProductRequest) GenerateVariantGroups() []*model.VariantGroup {
	if d.VariantGroups == nil {
		return nil
	}

	variantGroups := []*model.VariantGroup{}
	for _, req := range d.VariantGroups {
		variants := []*model.Variant{}
		for _, v := range req.Options {
			variants = append(variants, &model.Variant{
				Value: v,
			})
		}

		variants[0].MediaUrl = req.MediaUrl

		variantGroups = append(variantGroups, &model.VariantGroup{
			Name:    req.Name,
			Variant: variants,
		})
	}

	return variantGroups
}

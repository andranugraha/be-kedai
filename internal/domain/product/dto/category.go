package dto

import (
	"kedai/backend/be-kedai/internal/domain/product/model"
)

type GetCategoriesRequest struct {
	Depth     int  `form:"depth"`
	ParentID  int  `form:"parentId"`
	WithPrice bool `form:"withPrice"`
	Limit     int  `form:"limit"`
	Page      int  `form:"page"`
}

func (r *GetCategoriesRequest) Validate() {
	if r.Limit < 0 {
		r.Limit = 0
	}
	if r.Page < 1 {
		r.Page = 1
	}
}

func (r *GetCategoriesRequest) Offset() int {
	return int((r.Page - 1) * r.Limit)
}

type CategoryDTO struct {
	Name     string        `json:"name" binding:"required"`
	ImageURL string        `json:"image_url" binding:"required"`
	ParentID *int          `json:"parent_id,omitempty"`
	Children []CategoryDTO `json:"children,omitempty"`
}

func (cdto CategoryDTO) ToModel() *model.Category {
	categoryModel := &model.Category{
		Name:     cdto.Name,
		ImageURL: cdto.ImageURL,
	}
	for _, childDTO := range cdto.Children {
		childModel := childDTO.ToModel()

		categoryModel.Children = append(categoryModel.Children, childModel)
	}

	return categoryModel
}

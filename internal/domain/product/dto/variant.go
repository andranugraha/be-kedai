package dto

type AddVariantRequest struct {
	Name     string `json:"name" binding:"required,max=20"`
	MediaUrl string `json:"mediaUrl" binding:"omitempty,url"`
}

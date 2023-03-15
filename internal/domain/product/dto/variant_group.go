package dto

type CreateVariantGroupRequest struct {
	Name     string   `json:"name" binding:"required,max=20"`
	Options  []string `json:"options" binding:"required,min=1,max=50,dive,max=14"`
	MediaUrl string   `json:"mediaUrl" binding:"omitempty,url"`
}

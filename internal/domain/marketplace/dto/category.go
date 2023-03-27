package dto

type CategoryDTO struct {
	Name     string        `json:"name"`
	ImageURL string        `json:"image_url"`
	ParentID *int          `json:"parent_id,omitempty"`
	Children []CategoryDTO `json:"children,omitempty"`
}

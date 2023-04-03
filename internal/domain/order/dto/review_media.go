package dto

import "kedai/backend/be-kedai/internal/domain/order/model"

type ReviewMediaRequest struct {
	Url string `json:"url" binding:"required"`
}

func (req *ReviewMediaRequest) ToModel() *model.ReviewMedia {
	return &model.ReviewMedia{
		Url: req.Url,
	}
}

func (req *TransactionReviewRequest) ToReviewMediaModels() []*model.ReviewMedia {
	var reviewMedias []*model.ReviewMedia
	for _, reviewMedia := range req.ReviewMedias {
		reviewMedias = append(reviewMedias, reviewMedia.ToModel())
	}

	return reviewMedias
}

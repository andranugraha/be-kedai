package dto

import "kedai/backend/be-kedai/internal/domain/order/model"

type TransactionReviewRequest struct {
	Description   *string `json:"description" binding:"max=500"`
	Rating        int     `json:"rating" binding:"required,min=1,max=5,numeric"`
	TransactionId int     `json:"transactionId" binding:"required,min=1,numeric"`

	ReviewMedias []*ReviewMediaRequest `json:"reviewMedias" binding:"max=5"`

	UserId int
}

func (req *TransactionReviewRequest) ToModel() *model.TransactionReview {
	return &model.TransactionReview{
		Description:   req.Description,
		Rating:        req.Rating,
		TransactionId: req.TransactionId,
		ReviewMedias:  req.ToReviewMediaModels(),
	}
}

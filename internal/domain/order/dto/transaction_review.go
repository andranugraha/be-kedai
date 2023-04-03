package dto

import (
	"kedai/backend/be-kedai/internal/domain/order/model"
	productModel "kedai/backend/be-kedai/internal/domain/product/model"
	userModel "kedai/backend/be-kedai/internal/domain/user/model"
	"time"
)

type TransactionReviewRequest struct {
	Description   *string `json:"description" binding:"max=500"`
	Rating        int     `json:"rating" binding:"required,min=1,max=5,numeric"`
	TransactionId int     `json:"transactionId" binding:"required,min=1,numeric"`

	ReviewMedias []ReviewMediaRequest `json:"reviewMedias" binding:"max=5,dive"`

	UserId int
}

type TransactionReviewResponse struct {
	Description  *string              `json:"description"`
	Rating       int                  `json:"rating"`
	ReviewDate   time.Time            `json:"reviewDate"`
	ReviewMedias []*model.ReviewMedia `json:"reviewMedias"`
}

type ReviewResponse struct {
	User              UserReviewResponse        `json:"user"`
	TransactionReview TransactionReviewResponse `json:"transactionReview"`
	Variant           []*productModel.Variant   `json:"variant"`
}

type UserReviewResponse struct {
	Username string `json:"username"`
	PhotoUrl string `json:"photoUrl"`
}

func (res *UserReviewResponse) ToResponse(user *userModel.User) {
	res.Username = user.Username

	if user.Profile != nil {
		res.PhotoUrl = *user.Profile.PhotoUrl
	}

}

func (res *TransactionReviewResponse) ToResponse(review *model.TransactionReview) {
	res.Description = review.Description
	res.Rating = review.Rating
	res.ReviewDate = review.ReviewDate
	res.ReviewMedias = review.ReviewMedias
}

func (res *ReviewResponse) ToResponse(review *model.TransactionReview) {
	res.TransactionReview.ToResponse(review)
	res.User.ToResponse(review.Transaction.User)
	for _, variant := range review.Transaction.Sku.Variants {
		res.Variant = append(res.Variant, &variant)
	}
}

func ConvertReviewsToResponses(reviews []*model.TransactionReview) (responses []*ReviewResponse) {
	if len(reviews) == 0 {
		return []*ReviewResponse{}
	}

	for _, review := range reviews {
		reviewResponse := &ReviewResponse{}
		reviewResponse.ToResponse(review)
		responses = append(responses, reviewResponse)
	}

	return responses
}

func (req *TransactionReviewRequest) ToModel() *model.TransactionReview {
	if *req.Description == "" {
		req.Description = nil
	}

	return &model.TransactionReview{
		Description:   req.Description,
		Rating:        req.Rating,
		TransactionId: req.TransactionId,
		ReviewMedias:  req.ToReviewMediaModels(),
	}
}

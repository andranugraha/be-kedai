package service

import (
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"kedai/backend/be-kedai/internal/domain/product/repository"
	"kedai/backend/be-kedai/internal/domain/shop/service"
)

type DiscussionService interface {
	GetDiscussionByProductID(productID int) ([]*dto.Discussion, error)
	GetChildDiscussionByParentID(parentID int) ([]*dto.DiscussionReply, error)
	PostDiscussion(discussion *dto.DiscussionReq) error
}

type discussionServiceImpl struct {
	discussionRepository repository.DiscussionRepository
	shopService          service.ShopService
}

type DiscussionSConfig struct {
	DiscussionRepository repository.DiscussionRepository
	ShopService          service.ShopService
}

func NewDiscussionService(cfg *DiscussionSConfig) DiscussionService {
	return &discussionServiceImpl{
		discussionRepository: cfg.DiscussionRepository,
		shopService:          cfg.ShopService,
	}
}

func (d *discussionServiceImpl) GetDiscussionByProductID(productID int) ([]*dto.Discussion, error) {
	return d.discussionRepository.GetDiscussionByProductID(productID)
}

func (d *discussionServiceImpl) GetChildDiscussionByParentID(parentID int) ([]*dto.DiscussionReply, error) {
	return d.discussionRepository.GetChildDiscussionByParentID(parentID)
}

func (d *discussionServiceImpl) PostDiscussion(discussion *dto.DiscussionReq) error {

	if discussion.IsSeller {
		shop, err := d.shopService.FindShopByUserId(discussion.UserID)
		if err != nil {
			return err
		}
		discussion.ShopID = shop.ID
	}

	return d.discussionRepository.PostDiscussion(discussion)
}

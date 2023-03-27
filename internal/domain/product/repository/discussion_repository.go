package repository

import (
	"kedai/backend/be-kedai/internal/domain/product/dto"
	"log"
	"math"

	"gorm.io/gorm"
)

type DiscussionRepository interface {
	GetDiscussionByProductID(productID int, req dto.GetDiscussionReq) (data []*dto.Discussion, limit int, page int, totalRows int, totalPages int, err error)
	GetChildDiscussionByParentID(parentID int) ([]*dto.DiscussionReply, error)
	PostDiscussion(discussion *dto.DiscussionReq) error
}

type discussionRepositoryImpl struct {
	db *gorm.DB
}

type DiscussionRConfig struct {
	DB *gorm.DB
}

func NewDiscussionRepository(cfg *DiscussionRConfig) DiscussionRepository {
	return &discussionRepositoryImpl{
		db: cfg.DB,
	}
}

func (d *discussionRepositoryImpl) GetDiscussionByProductID(productID int, req dto.GetDiscussionReq) (data []*dto.Discussion, limit int, page int, totalRows int, totalPages int, err error) {

	var discussions []*dto.Discussion
	var count int64
	err = d.db.Model(&dto.Discussion{}).Where("product_id = ? AND parent_id IS NULL", productID).Count(&count).Error
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	err = d.db.Where("product_id = ? AND parent_id IS NULL", productID).Preload("User").Preload("User.Profile").Preload("Shop").Limit(req.Limit).Offset((req.Page - 1) * req.Limit).Find(&discussions).Error
	if err != nil {
		return nil, 0, 0, 0, 0, err
	}

	if len(discussions) == 0 {
		return nil, 0, 0, 0, 0, nil
	}

	for i, discussion := range discussions {
		var replies []*dto.DiscussionReply
		var repliesCount int
		err := d.db.Table("discussions").Where("parent_id = ?", discussion.ID).Preload("User").Preload("User.Profile").Preload("Shop").Find(&replies).Error
		if err != nil {
			return nil, 0, 0, 0, 0, err
		}

		discussions[i].Username = discussion.User.Username
		discussions[i].UserUrl = *discussion.User.Profile.PhotoUrl
		if discussions[i].ShopId != 0 {
			discussions[i].ShopName = discussion.Shop.Name
			discussions[i].ShopUrl = *discussion.Shop.PhotoUrl
		}
		repliesCount = len(replies)

		if repliesCount >= 1 {
			discussions[i].Reply = replies[0]
			discussions[i].Reply.Username = replies[0].User.Username
			discussions[i].Reply.UserUrl = *replies[0].User.Profile.PhotoUrl
			if discussions[i].Reply.ShopId != 0 {
				discussions[i].Reply.ShopName = replies[0].Shop.Name
				discussions[i].Reply.ShopUrl = *replies[0].Shop.PhotoUrl
			}
		}
		discussions[i].ReplyCount = repliesCount
	}
	return discussions, req.Limit, req.Page, int(count), int(math.Ceil(float64(count) / float64(req.Limit))), nil
}

func (d *discussionRepositoryImpl) GetChildDiscussionByParentID(parentID int) ([]*dto.DiscussionReply, error) {
	var replies []*dto.DiscussionReply
	err := d.db.Table("discussions").Where("parent_id = ?", parentID).Preload("User").Preload("User.Profile").Preload("Shop").Find(&replies).Error
	if err != nil {
		return nil, err
	}

	for i, reply := range replies {
		replies[i].Username = reply.User.Username
		replies[i].UserUrl = *reply.User.Profile.PhotoUrl
		log.Println(reply.ShopId)
		if replies[i].ShopId != 0 {
			replies[i].ShopName = reply.Shop.Name
			replies[i].ShopUrl = *reply.Shop.PhotoUrl
		}
	}

	return replies, nil
}

func (d *discussionRepositoryImpl) PostDiscussion(discussion *dto.DiscussionReq) error {

	err := d.db.Table("discussions").Create(discussion).Error
	if err != nil {
		return err
	}

	return nil
}

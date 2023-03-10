package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	orderRepo "kedai/backend/be-kedai/internal/domain/order/repository"
	"kedai/backend/be-kedai/internal/domain/shop/dto"
	"kedai/backend/be-kedai/internal/domain/shop/model"
	userRepo "kedai/backend/be-kedai/internal/domain/user/repository"
	"math"

	"gorm.io/gorm"
)

type ShopRepository interface {
	FindShopById(id int) (*model.Shop, error)
	FindShopByUserId(userId int) (*model.Shop, error)
	FindShopBySlug(slug string) (*model.Shop, error)
	FindShopByKeyword(req dto.FindShopRequest) ([]*dto.FindShopResponse, int64, int, error)
	UpdateShopAddressIdByUserId(tx *gorm.DB, userId int, addressId int) error
	GetShopFinanceOverview(shopId int) (*dto.ShopFinanceOverviewResponse, error)
}

type shopRepositoryImpl struct {
	db                 *gorm.DB
	invoicePerShopRepo orderRepo.InvoicePerShopRepository
	walletHistoryRepo  userRepo.WalletHistoryRepository
}

type ShopRConfig struct {
	DB                 *gorm.DB
	InvoicePerShopRepo orderRepo.InvoicePerShopRepository
	WalletHistoryRepo  userRepo.WalletHistoryRepository
}

func NewShopRepository(cfg *ShopRConfig) ShopRepository {
	return &shopRepositoryImpl{
		db:                 cfg.DB,
		invoicePerShopRepo: cfg.InvoicePerShopRepo,
		walletHistoryRepo:  cfg.WalletHistoryRepo,
	}
}

func (r *shopRepositoryImpl) FindShopById(id int) (*model.Shop, error) {
	var shop model.Shop

	err := r.db.Where("id = ?", id).First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrShopNotFound
		}

		return nil, err
	}

	return &shop, err
}

func (r *shopRepositoryImpl) FindShopByUserId(userId int) (*model.Shop, error) {
	var shop model.Shop

	err := r.db.Where("user_id = ?", userId).First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrShopNotFound
		}

		return nil, err
	}

	return &shop, err
}

func (r *shopRepositoryImpl) FindShopBySlug(slug string) (*model.Shop, error) {
	var shop model.Shop

	err := r.db.Where("slug = ?", slug).Preload("ShopCategory").First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrShopNotFound
		}
	}

	return &shop, nil
}

func (r *shopRepositoryImpl) FindShopByKeyword(req dto.FindShopRequest) ([]*dto.FindShopResponse, int64, int, error) {
	var (
		shopList  []*dto.FindShopResponse
		totalRows int64
		totalPage int
		isActive  = true
	)

	db := r.db.Select(`shops.*, count(p.id) as product_count`).
		Joins("left join products p on shops.id = p.shop_id and p.is_active = ?", isActive).
		Group("shops.id").Where("shops.name ILIKE ?", "%"+req.Keyword+"%")

	countQuery := db.Session(&gorm.Session{})
	countQuery.Model(&model.Shop{}).Distinct("shops.id").Count(&totalRows)
	totalPage = int(math.Ceil(float64(totalRows) / float64(req.Limit)))

	err := db.Model(&model.Shop{}).Order("rating desc").Limit(req.Limit).Offset(req.Offset()).Find(&shopList).Error
	if err != nil {
		return nil, 0, 0, err
	}

	return shopList, totalRows, totalPage, nil
}

func (r *shopRepositoryImpl) UpdateShopAddressIdByUserId(tx *gorm.DB, userId int, addressId int) error {
	res := tx.Model(&model.Shop{}).Where("user_id = ?", userId).Update("address_id", addressId)
	if err := res.Error; err != nil {
		tx.Rollback()
		return err
	}

	if res.RowsAffected == 0 {
		tx.Rollback()
		return errs.ErrShopNotFound
	}

	return nil
}

func (r *shopRepositoryImpl) GetShopFinanceOverview(shopId int) (*dto.ShopFinanceOverviewResponse, error) {
	toRelease, err := r.invoicePerShopRepo.GetShopFinanceToRelease(shopId)
	if err != nil {
		return nil, err
	}

	released, err := r.walletHistoryRepo.GetShopFinanceReleased(shopId)
	if err != nil {
		return nil, err
	}

	return &dto.ShopFinanceOverviewResponse{
		ToRelease: toRelease,
		Released:  *released,
	}, nil
}

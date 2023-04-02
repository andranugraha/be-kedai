package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SkuRepository interface {
	GetByID(ID int) (*model.Sku, error)
	GetByVariantIDs(variantIDs []int) (*model.Sku, error)
	GetByProductId(productID int) ([]*model.Sku, error)
	ReduceStock(tx *gorm.DB, skuID int, quantity int) error
	IncreaseStock(tx *gorm.DB, skuID int, quantity int) error
	Create(tx *gorm.DB, skus []*model.Sku) error
	Update(tx *gorm.DB, productId int, skus []*model.Sku) error
}

type skuRepositoryImpl struct {
	db                         *gorm.DB
	productPromotionRepository ProductPromotionRepository
}

type SkuRConfig struct {
	DB                         *gorm.DB
	ProductPromotionRepository ProductPromotionRepository
}

func NewSkuRepository(cfg *SkuRConfig) SkuRepository {
	return &skuRepositoryImpl{
		db:                         cfg.DB,
		productPromotionRepository: cfg.ProductPromotionRepository,
	}
}

func (r *skuRepositoryImpl) GetByID(ID int) (*model.Sku, error) {
	var sku model.Sku

	err := r.db.Where("id = ?", ID).Preload("Product").First(&sku).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrProductDoesNotExist
		}

		return nil, err
	}

	return &sku, err
}

func (r *skuRepositoryImpl) GetByVariantIDs(variantIDs []int) (*model.Sku, error) {
	var sku model.Sku
	var err error

	if len(variantIDs) < 2 {
		err = r.db.
			Select("skus.id, skus.price, skus.stock, skus.product_id").
			Joins("JOIN product_variants pv1 ON skus.id = pv1.sku_id").
			Where("pv1.variant_id = ?", variantIDs[0]).
			First(&sku).Error
	} else {
		err = r.db.
			Select("skus.id, skus.price, skus.stock, skus.product_id").
			Joins("JOIN product_variants pv1 ON skus.id = pv1.sku_id").
			Joins("JOIN product_variants pv2 ON pv1.sku_id = pv2.sku_id AND pv1.id != pv2.id").
			Where("pv1.variant_id = ? AND pv2.variant_id = ?", variantIDs[0], variantIDs[1]).
			First(&sku).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrSKUDoesNotExist
		}

		return nil, err
	}

	return &sku, nil
}

func (r *skuRepositoryImpl) ReduceStock(tx *gorm.DB, skuID int, quantity int) error {
	err := tx.Model(&model.Sku{}).
		Where("id = ?", skuID).
		Where("stock >= ?", quantity).
		Update("stock", gorm.Expr("stock - ?", quantity))
	if err.Error != nil {
		tx.Rollback()
		return err.Error
	}

	if err.RowsAffected == 0 {
		tx.Rollback()
		return errs.ErrProductQuantityNotEnough
	}

	return nil
}

func (r *skuRepositoryImpl) IncreaseStock(tx *gorm.DB, skuID int, quantity int) error {
	err := tx.Model(&model.Sku{}).
		Where("id = ?", skuID).
		Update("stock", gorm.Expr("stock + ?", quantity))
	if err.Error != nil {
		tx.Rollback()
		return err.Error
	}

	return nil
}

func (r *skuRepositoryImpl) Create(tx *gorm.DB, skus []*model.Sku) error {
	res := tx.Omit("Variants").Clauses(clause.OnConflict{DoNothing: true}).Create(&skus)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errs.ErrSKUUsed
	}

	return nil
}

func (r *skuRepositoryImpl) GetByProductId(productID int) ([]*model.Sku, error) {
	var skus []*model.Sku

	err := r.db.Where("product_id = ?", productID).
		Preload("Variants").
		Find(&skus).Error
	if err != nil {
		return nil, err
	}

	return skus, nil
}

func (r *skuRepositoryImpl) Update(tx *gorm.DB, productId int, skus []*model.Sku) error {
	retrievedSkus, err := r.GetByProductId(productId)
	if err != nil {
		return err
	}

	var union []*model.Sku

	for _, sku := range retrievedSkus {
		found := false
		for _, newSku := range skus {
			if sku.Sku == newSku.Sku {
				sku.Stock = newSku.Stock
				sku.Price = newSku.Price
				sku.Variants = newSku.Variants
				union = append(union, sku)
				found = true
				break
			}
		}

		if !found {
			err = tx.Delete(sku).Error
			if err != nil {
				tx.Rollback()
				return err
			}

			if errDeletePromotion := r.productPromotionRepository.Delete(tx, sku.ID); errDeletePromotion != nil {
				tx.Rollback()
				return errDeletePromotion
			}

			if errDelete := tx.Model(sku).Unscoped().Association("Variants").Clear(); errDelete != nil {
				tx.Rollback()
				return errDelete
			}
		}
	}

	for _, sku := range skus {
		found := false
		for _, retrievedSku := range retrievedSkus {
			if sku.Sku == retrievedSku.Sku {
				found = true
				break
			}
		}

		if !found {
			union = append(union, sku)
		}
	}

	if len(union) == 0 {
		return nil
	}

	err = tx.Clauses(clause.OnConflict{
		OnConstraint: ("skus_sku_key"),
		UpdateAll:    true,
	}).Clauses(clause.Returning{
		Columns: []clause.Column{{Name: "id"}},
	}).Save(&union).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, sku := range union {
		for _, variant := range sku.Variants {
			if err := tx.Model(&model.ProductVariant{}).Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&model.ProductVariant{
				SkuId:     sku.ID,
				VariantId: variant.ID,
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return nil
}

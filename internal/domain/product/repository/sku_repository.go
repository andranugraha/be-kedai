package repository

import (
	"errors"
	errs "kedai/backend/be-kedai/internal/common/error"
	"kedai/backend/be-kedai/internal/domain/product/model"

	"gorm.io/gorm"
)

type SkuRepository interface {
	GetByID(ID int) (*model.Sku, error)
	GetByVariantIDs(variantIDs []int) (*model.Sku, error)
	ReduceStock(tx *gorm.DB, skuID int, quantity int) error
	IncreaseStock(tx *gorm.DB, skuID int, quantity int) error
}

type skuRepositoryImpl struct {
	db *gorm.DB
}

type SkuRConfig struct {
	DB *gorm.DB
}

func NewSkuRepository(cfg *SkuRConfig) SkuRepository {
	return &skuRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *skuRepositoryImpl) GetByID(ID int) (*model.Sku, error) {
	var sku model.Sku

	err := r.db.Where("id = ?", ID).First(&sku).Error
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
			Joins("JOIN product_variants pv1 ON skus.id = pv1.sku_id").
			Where("pv1.variant_id = ?", variantIDs[0]).
			First(&sku).Error
	} else {
		err = r.db.
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

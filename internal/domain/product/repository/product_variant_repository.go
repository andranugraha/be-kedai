package repository

import (
	"kedai/backend/be-kedai/internal/domain/product/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductVariantRepository interface {
	Create(tx *gorm.DB, payload []*model.ProductVariant) error
	GetBySkuId(id int) ([]*model.ProductVariant, error)
	Update(tx *gorm.DB, groupId int, payload []*model.ProductVariant) error
}

type productVariantRepositoryImpl struct {
	db *gorm.DB
}

type ProductVariantRConfig struct {
	DB *gorm.DB
}

func NewProductVariantRepository(cfg *ProductVariantRConfig) ProductVariantRepository {
	return &productVariantRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *productVariantRepositoryImpl) Create(tx *gorm.DB, payload []*model.ProductVariant) error {
	if err := tx.Create(&payload).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *productVariantRepositoryImpl) GetBySkuId(id int) ([]*model.ProductVariant, error) {
	var variants []*model.ProductVariant

	err := r.db.Where("sku_id = ?", id).Find(&variants).Error
	if err != nil {
		return nil, err
	}
	return variants, nil
}

func (r *productVariantRepositoryImpl) Update(tx *gorm.DB, groupId int, payload []*model.ProductVariant) error {

	retrievedVariants, err := r.GetBySkuId(groupId)
	if err != nil {
		return err
	}

	var variants []*model.ProductVariant

	for _, variant := range payload {
		found := false
		for _, retrievedVariant := range retrievedVariants {
			if variant.SkuId == retrievedVariant.SkuId {
				variant.ID = retrievedVariant.ID
				variants = append(variants, variant)
				found = true
				break
			}
		}
		if !found {
			if err := tx.Delete(&variant).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "sku_id"}, {Name: "variant_id"}},
			DoNothing: true,
		}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"sku_id", "variant_id"}),
		}).Save(&variants).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

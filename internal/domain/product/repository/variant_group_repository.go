package repository

import (
	"kedai/backend/be-kedai/internal/domain/product/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VariantGroupRepository interface {
	Create(tx *gorm.DB, variantGroups []*model.VariantGroup) error
	GetByProductId(id int) ([]*model.VariantGroup, error)
	Update(tx *gorm.DB, productId int, variantGroups []*model.VariantGroup) ([]*model.VariantGroup, error)
}

type variantGroupRepositoryImpl struct {
	db *gorm.DB
}

type VariantGroupRConfig struct {
	DB *gorm.DB
}

func NewVariantGroupRepository(cfg *VariantGroupRConfig) VariantGroupRepository {
	return &variantGroupRepositoryImpl{
		db: cfg.DB,
	}
}

func (r *variantGroupRepositoryImpl) Create(tx *gorm.DB, variantGroups []*model.VariantGroup) error {
	return tx.Create(&variantGroups).Error
}

func (r *variantGroupRepositoryImpl) GetByProductId(id int) ([]*model.VariantGroup, error) {
	var variantGroups []*model.VariantGroup

	err := r.db.Where("product_id = ?", id).
		Preload("Variant").
		Find(&variantGroups).Error
	if err != nil {
		return nil, err
	}
	return variantGroups, nil
}

func (r *variantGroupRepositoryImpl) Update(tx *gorm.DB, productId int, variantGroups []*model.VariantGroup) ([]*model.VariantGroup, error) {

	retrievedVarGroups, err := r.GetByProductId(productId)
	if err != nil {
		return nil, err
	}

	var (
		union        []model.VariantGroup
		defaultUnion []model.VariantGroup
	)

	for _, variantGroup := range retrievedVarGroups {
		found := false
		for _, vg := range variantGroups {
			if vg.Name == variantGroup.Name {
				vg.ID = variantGroup.ID
				found = true
				union = append(union, *vg)

				for _, retrievedVariant := range variantGroup.Variant {
					found := false
					for _, variant := range vg.Variant {
						if variant.Value == retrievedVariant.Value {
							variant.ID = retrievedVariant.ID
							found = true
						}
					}
					if !found {
						if err := tx.Delete(retrievedVariant).Error; err != nil {
							tx.Rollback()
							return nil, err
						}

						if errClear := tx.Unscoped().Where("variant_id=?", retrievedVariant.ID).Delete(&model.ProductVariant{}).Error; errClear != nil {
							tx.Rollback()
							return nil, errClear
						}
					}
				}
			}
		}

		if !found {
			if err := tx.Delete(variantGroup).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

	}

	for _, vg := range variantGroups {
		if vg.ID == 0 {
			union = append(union, *vg)
		}
	}

	defaultUnion = append(defaultUnion, union...)

	if err := tx.
		Clauses(clause.Returning{}).
		Save(&union).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var res []*model.VariantGroup
	for _, variantGroup := range union {
		for _, vg := range defaultUnion {
			if vg.Name == variantGroup.Name {
				for _, variant := range vg.Variant {
					variant.GroupId = variantGroup.ID
				}

				if err := tx.
					Clauses(clause.Returning{}).
					Clauses(clause.OnConflict{
						Columns: []clause.Column{{Name: "value"}, {Name: "group_id"}},
						DoUpdates: clause.AssignmentColumns([]string{
							"media_url",
						}),
					}).
					Save(&vg.Variant).Error; err != nil {
					tx.Rollback()

					return nil, err
				}

				res = append(res, &vg)
			}
		}
	}

	return res, nil
}

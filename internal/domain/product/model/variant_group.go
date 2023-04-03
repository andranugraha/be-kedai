package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VariantGroup struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	ProductID int        `json:"productId"`
	Variant   []*Variant `json:"variants,omitempty" gorm:"foreignKey:GroupId"`

	gorm.Model `json:"-"`
}

func (v *VariantGroup) BeforeDelete(tx *gorm.DB) (err error) {
	var deletedVariants []*Variant
	err = tx.Clauses(clause.Returning{}).Delete(&deletedVariants, "group_id = ?", v.ID).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	var deletedIds []int

	for _, deletedVariant := range deletedVariants {
		deletedIds = append(deletedIds, deletedVariant.ID)
	}

	return tx.Unscoped().Delete(&ProductVariant{}, "variant_id IN ?", deletedIds).Error
}

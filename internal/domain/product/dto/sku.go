package dto

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"strconv"
	"strings"
)

type GetSKURequest struct {
	VariantID string `form:"variantId" binding:"required"`
}

func (d *GetSKURequest) ToIntList() ([]int, error) {
	variantList := strings.Split(d.VariantID, ",")

	variant1, err := strconv.Atoi(variantList[0])
	if err != nil || variant1 < 1 {
		return nil, errs.ErrInvalidVariantID
	}

	if len(variantList) == 1 {
		return []int{variant1}, nil
	}

	variant2, err := strconv.Atoi(variantList[1])
	if err != nil || variant2 < 1 {
		return nil, errs.ErrInvalidVariantID
	}

	return []int{variant1, variant2}, nil
}

package dto

import (
	errs "kedai/backend/be-kedai/internal/common/error"
	"strconv"
	"strings"
)

type GetSKURequest struct {
	VariantID string `form:"variantId"`
}

func (d *GetSKURequest) ToIntList() ([]int, error) {
	variantList := strings.Split(d.VariantID, ",")
	if len(variantList) < 2 {
		return nil, errs.ErrIncompleteVariantIDArguments
	}

	variant1, err := strconv.Atoi(variantList[0])
	if err != nil || variant1 < 1 {
		return nil, errs.ErrInvalidVariantID
	}

	variant2, err := strconv.Atoi(variantList[1])
	if err != nil || variant2 < 1 {
		return nil, errs.ErrInvalidVariantID
	}

	return []int{variant1, variant2}, nil
}

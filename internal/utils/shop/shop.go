package shop

import (
	"unicode"

	"github.com/forPelevin/gomoji"
)

func ValidateShopName(shopName string) bool {
	if containEmoji := gomoji.ContainsEmoji(shopName); containEmoji {
		return false
	}

	containLetter := false
	for _, c := range shopName {
		if unicode.IsLetter(c) {
			containLetter = true
			continue
		}

		if unicode.IsNumber(c) || c == ' ' {
			continue
		}

		return false
	}

	return containLetter
}

package product

import (
	"unicode"

	"github.com/forPelevin/gomoji"
)

func ValidateProductName(productName string) bool {
	if productName == "" {
		return true
	}

	if containEmoji := gomoji.ContainsEmoji(productName); containEmoji {
		return false
	}

	containLetter := false
	for _, c := range productName {
		if unicode.IsLetter(c) {
			containLetter = true
		}
	}

	return containLetter
}

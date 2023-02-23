package password

import (
	"unicode"

	"github.com/forPelevin/gomoji"
)

func VerifyPassword(pw string) bool {
	var containUpper bool
	var containLower bool
	var containNumeric bool

	containEmoji := gomoji.ContainsEmoji(pw)
	if containEmoji {
		return false
	}

	for _, c := range pw {
		switch {
		case unicode.IsUpper(c):
			containUpper = true
		case unicode.IsLower(c):
			containLower = true
		case unicode.IsNumber(c):
			containNumeric = true
		}
	}

	if containNumeric && containUpper && containLower {
		return true
	}

	return false
}

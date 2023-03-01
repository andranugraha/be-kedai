package credential

import (
	"kedai/backend/be-kedai/internal/utils/hash"
	"strings"
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

func VerifyChangePassword(oldPw string, newPw string, username string) bool {
	if hash.ComparePassword(oldPw, newPw) {
		return false
	}

	return !ContainsUsername(newPw, username)
}

func ContainsUsername(pw string, username string) bool {
	return strings.Contains(strings.ToLower(pw), strings.ToLower(username))

}

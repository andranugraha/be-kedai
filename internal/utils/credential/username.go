package credential

import (
	"unicode"

	"github.com/forPelevin/gomoji"
)

func VerifyUsername(username string) bool {
	if containEmoji := gomoji.ContainsEmoji(username); containEmoji {
		return false
	}

	containLetter := false
	for _, c := range username {
		if unicode.IsLetter(c) {
			containLetter = true
			continue
		}

		if unicode.IsNumber(c) || c == '_' || c == '.' {
			continue
		}

		return false
	}

	return containLetter
}

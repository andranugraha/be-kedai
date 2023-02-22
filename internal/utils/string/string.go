package string

import "unicode"

func VerifyPassword(pw string) bool {
	var containUpper bool
	var containLower bool
	var containNumeric bool
	var containLetter bool

	for _, c := range pw {
		switch {
		case unicode.IsUpper(c):
			containUpper = true
		case unicode.IsLower(c):
			containLower = true
		case unicode.IsNumber(c):
			containNumeric = true
		case unicode.IsLetter(c):
			containLetter = true
		}
	}

	if containNumeric && containUpper && containLower && containLetter {
		return true
	}

	return false
}

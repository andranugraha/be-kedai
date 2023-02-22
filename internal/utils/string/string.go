package string

import "unicode"

func VerifyPassword(pw string) bool {
	var containUpper bool
	var containNumeric bool

	for _, c := range pw {
		switch {
		case unicode.IsUpper(c):
			containUpper = true
		case unicode.IsNumber(c):
			containNumeric = true
		}
	}

	if containNumeric && containUpper {
		return true
	}

	return false
}

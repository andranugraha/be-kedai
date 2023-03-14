package strings

import "strings"

func GenerateSlug(str string) string {
	str = strings.ToLower(str)
	str = strings.ReplaceAll(str, " ", "-")
	str = strings.ReplaceAll(str, ":", "")
	str = strings.ReplaceAll(str, "!", "")
	str = strings.ReplaceAll(str, "?", "")
	str = strings.ReplaceAll(str, "(", "")
	str = strings.ReplaceAll(str, ")", "")
	str = strings.ReplaceAll(str, ",", "")
	str = strings.ReplaceAll(str, ".", "")
	str = strings.ReplaceAll(str, "/", "")
	str = strings.ReplaceAll(str, "\"", "")
	str = strings.ReplaceAll(str, "'", "")
	str = strings.ReplaceAll(str, "’", "")
	str = strings.ReplaceAll(str, "‘", "")
	str = strings.ReplaceAll(str, "”", "")
	str = strings.ReplaceAll(str, "“", "")
	str = strings.ReplaceAll(str, "–", "")
	str = strings.ReplaceAll(str, "—", "")
	return str
}

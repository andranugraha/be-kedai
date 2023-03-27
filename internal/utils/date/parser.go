package date

import (
	"regexp"
	"time"
)

func IsValidRFC3999NanoDate(date string) bool {
	pattern := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,9}0*|\,\d{1,3})(Z|[+-]\d{2}:\d{2})$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(date)
}

func ParseRFC3999NanoTime(str string, defaultValue time.Time) time.Time {
	t, err := time.Parse(time.RFC3339Nano, str)
	if err != nil {
		return defaultValue
	}
	return t
}

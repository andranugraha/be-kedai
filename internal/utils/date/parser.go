package date

import (
	"regexp"
	"time"
)

func IsValidRFC3999NanoDate(date string) bool {
	/*
		Breaking down the pattern:
			^ indicates the start of the string.
			\d{4} matches exactly four digits, representing the year.
			- matches the hyphen character.
			\d{2} matches exactly two digits, representing the month.
			- matches the hyphen character.
			\d{2} matches exactly two digits, representing the day.
			T matches the literal character "T", separating the date from the time.
			\d{2} matches exactly two digits, representing the hour.
			: matches the colon character.
			\d{2} matches exactly two digits, representing the minute.
			: matches the colon character.
			\d{2} matches exactly two digits, representing the second.
			(\.\d{1,9}0*|\,\d{1,3}) matches either a dot (".") followed by one to nine digits optionally followed by any number of zeros, OR a comma (",") followed by one to three digits. This is used to match the fractional seconds portion of the timestamp.
			(Z|[+-]\d{2}:\d{2}) matches either the letter "Z" (representing UTC time), OR a plus or minus sign, followed by exactly two digits (representing the time zone offset hours), followed by a colon, followed by exactly two digits (representing the time zone offset minutes). This is used to match the time zone offset portion of the timestamp.
			$ indicates the end of the string.

		So in summary, this regex pattern matches a string representing a date-time in the format of "YYYY-MM-DDTHH:MM:SS.sssssssssZ"
			or "YYYY-MM-DDTHH:MM:SS.sssssssss±HH:MM" (where "ssssssssss" represents the fractional seconds and "±HH:MM" represents the time zone
			offset from UTC).
	*/
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

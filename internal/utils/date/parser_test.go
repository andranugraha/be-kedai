package date_test

import (
	. "kedai/backend/be-kedai/internal/utils/date"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsValidRFC3999NanoDate(t *testing.T) {
	// Test valid dates
	validDates := []string{
		"2023-03-27T10:11:12.123456789Z",
		"2023-03-27T10:11:12.123456789+02:00",
	}
	for _, date := range validDates {
		assert.True(t, IsValidRFC3999NanoDate(date))
	}

	// Test invalid dates
	invalidDates := []string{
		"2023-03-27",
		"2023-03-27T10:11:12",
		"2023-03-27T10:11:12+02",
		"2023-03-27T10:11:12.12345678Z",
		"2023-03-27T10:11:12.1234567890Z",
		"2023-03-27T10:11:12.123456789+02:0",
		"2023-03-27T10:11:12.123456789+2:00",
		"2023-03-27T10:11:12.123456789+02:000",
		"2023-03-27T10:11:12.123456789Z+02:00",
		"2023-03-27T10:11:12z",
		"2023-03-27T10:11:12z+02:00",
		"2023-03-27T10:11:12.123456789Z+02",
		"2023-03-27T10:11:12.123456789Z+2:00",
		"2023-03-27T10:11:12.123456789Z+02:000",
		"2023-03-27T10:11:12.123456789+02:00Z",
		"2023-03-27T10:11:12Z02:00",
		"2023-03-27T10:11:12.123Z",
	}
	for _, date := range invalidDates {
		assert.False(t, IsValidRFC3999NanoDate(date))
	}
}

func TestParseRFC3999NanoTime(t *testing.T) {
	defaultValue := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	// Test valid input
	expected := time.Date(2023, 3, 27, 13, 45, 0, 0, time.UTC)
	result := ParseRFC3999NanoTime("2023-03-27T13:45:00.000Z", defaultValue)
	assert.Equal(t, expected, result)

	// Test invalid input
	expected = defaultValue
	result = ParseRFC3999NanoTime("invalid datetime string", defaultValue)
	assert.Equal(t, expected, result)
}

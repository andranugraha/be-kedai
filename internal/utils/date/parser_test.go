package date_test

import (
	"fmt"
	. "kedai/backend/be-kedai/internal/utils/date"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsValidRFC3999NanoDate(t *testing.T) {
	// Test valid dates
	validDates := []string{
		"2023-03-27T10:11:12.123456789Z",
		"2023-03-27T10:11:12.12345678Z",
		"2023-03-27T10:11:12.123456789+02:00",
		"2023-03-27T10:11:12.1234567890Z",
		"2023-03-27T10:11:12.123Z",
	}
	for _, date := range validDates {
		assert.True(t, IsValidRFC3999NanoDate(date))
	}

	// Test invalid dates
	invalidDates := []string{
		"2023-03-27",
		"2023-03-27T10:11:12",
		"2023-03-27T10:11:12+02",
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

func BenchmarkIsValidRFC3999NanoDate(b *testing.B) {
	var table = []struct {
		name string
		date string
	}{
		{name: "valid", date: "2023-03-30T08:00:00Z"},
		{name: "invalid", date: "2023-03-30T08:00:00.123456789Z1"},
	}
	for _, v := range table {
		b.Run(fmt.Sprintf("date_"+v.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				IsValidRFC3999NanoDate(v.date)
			}
		})
	}
}

func BenchmarkParseRFC3999NanoTime(b *testing.B) {
	var table = []struct {
		name         string
		str          string
		defaultValue time.Time
	}{
		{name: "valid", str: "2023-03-30T10:45:30.123456789Z", defaultValue: time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC)},
		{name: "invalid", str: "invalid_time", defaultValue: time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC)},
	}
	for _, v := range table {
		b.Run(fmt.Sprintf("input_type_"+v.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParseRFC3999NanoTime(v.str, v.defaultValue)
			}
		})
	}
}

package date_test

import (
	. "kedai/backend/be-kedai/internal/utils/date"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDaysBetween(t *testing.T) {
	start := time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, time.March, 1, 0, 0, 0, 1, time.UTC)
	expected := 1

	result := DaysBetween(start, end)

	assert.Equal(t, expected, result, "Expected the number of days between %v and %v to be %d, but got %d", start, end, expected, result)
}

package date_test

import (
	"fmt"
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

func BenchmarkDaysBetween(b *testing.B) {
	var table = []struct {
		name  string
		start time.Time
		end   time.Time
	}{
		{name: "normal", start: time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC), end: time.Date(2023, time.March, 15, 0, 0, 0, 0, time.UTC)},
		{name: "short", start: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), end: time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)},
		{name: "long", start: time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC), end: time.Date(4023, time.March, 1, 0, 0, 0, 0, time.UTC)},
	}
	for _, v := range table {
		b.Run(fmt.Sprintf("time_interval_size_"+v.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				DaysBetween(v.start, v.end)
			}
		})
	}
}

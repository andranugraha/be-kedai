package date

import (
	"math"
	"time"
)

func DaysBetween(start time.Time, end time.Time) int {
	duration := end.Sub(start)
	days := int(math.Ceil(duration.Hours() / 24))
	return days
}

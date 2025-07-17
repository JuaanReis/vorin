package print

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	seconds := int(d.Seconds()) % 60
	minutes := int(d.Minutes()) % 60
	hours := int(d.Hours())
	ms := d.Milliseconds() % 1000

	if hours > 0 {
		return fmt.Sprintf("%02dh %02dm %02ds %03dms", hours, minutes, seconds, ms)
	} else if minutes > 0 {
		return fmt.Sprintf("%02dm %02ds %03dms", minutes, seconds, ms)
	} else if seconds > 0 {
		return fmt.Sprintf("%02ds %03dms", seconds, ms)
	}
	return fmt.Sprintf("%dms", ms)
}

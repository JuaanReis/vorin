package internal

import (
	"fmt"
	"sync/atomic"
	"time"
)

func UpdateProgressBar(total int, current *int32, errors *int32, reqPerSec *int32, startTime time.Time, doneBar chan struct{}) {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		curr := atomic.LoadInt32(current)
		errs := atomic.LoadInt32(errors)
		duration := time.Since(startTime).Truncate(time.Second)
		media := float64(curr) / time.Since(startTime).Seconds()
		green := "\033[32m"
		red := "\033[31m"
		reset := "\033[0m"

		fmt.Printf("\r%s[Vorin]%s Progress: [%d/%d] | %s%.1f req/s%s | Time: %s | %s%d errors%s",
			green, reset, curr, total, green, media, reset, duration, red, errs, reset)


		if int(curr) >= total {
			break
		}
	}
	
	fmt.Println()
	close(doneBar)
}

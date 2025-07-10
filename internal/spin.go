package internal

import (
	"fmt"
	"time"
)

func Spinner(prefix string, done <-chan bool) {
	symbols := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Print("\r")
			return
		default:
			fmt.Printf("\r%s %s", prefix, symbols[i%len(symbols)])
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}
}

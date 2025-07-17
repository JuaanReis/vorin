package collector

import (
	"strings"
)

func CountLinesAndSize(body string) (int, int) {
    lines := strings.Count(body, "\n")
    return lines, len(body)
}
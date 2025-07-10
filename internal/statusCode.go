package internal

import (
	"strconv"
	"strings"
)

func ParseStatusCodes(input string) map[int]bool {
	result := make(map[int]bool)
	codes := strings.Split(input, ",")
	for _, codeStr := range codes {
		codeStr = strings.TrimSpace(codeStr)
		if code, err := strconv.Atoi(codeStr); err == nil {
			result[code] = true
		}
	}
	return result
}

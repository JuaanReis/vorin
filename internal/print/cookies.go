package print

import (
	"fmt"
	"strings"
)

type CookiesFlags []string

func (h *CookiesFlags) String() string {
	return strings.Join(*h, ", ")
}

func (h *CookiesFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func ParseCookiesFlags(cookies CookiesFlags) map[string]string {
	result := make(map[string]string)
	for _, h := range cookies {
		parts := strings.SplitN(h, "=", 2)
		if len(parts) != 2 {
			fmt.Printf("[WARNING] Invalid cookie format: %s\n", h)
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		result[key] = val
	}
	return result
}
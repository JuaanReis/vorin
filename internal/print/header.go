package print

import (
	"fmt"
	"strings"
)

type HeaderFlags []string

func (h *HeaderFlags) String() string {
	return strings.Join(*h, ", ")
}

func (h *HeaderFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func ParseHeaderFlags(headers HeaderFlags) map[string]string {
	result := make(map[string]string)
	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			fmt.Printf("[WARNING] Invalid header format: %s\n", h)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		result[key] = val
	}
	return result
}
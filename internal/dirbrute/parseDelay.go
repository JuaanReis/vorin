package dirbrute

import (
	"strings"
	"fmt"
	"strconv"
)

func ParseDelay(input string) (int, int, error) {
	if strings.Contains(input, "-") {
		parts := strings.Split(input, "-")
		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("invalid delay format")
		}
		min, err1 := strconv.Atoi(parts[0])
		max, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil || min < 0 || max < min {
			return 0, 0, fmt.Errorf("invalid delay values")
		}
		return min, max, nil
	}

	val, err := strconv.Atoi(input)
	if err != nil || val < 0 {
		return 0, 0, fmt.Errorf("invalid delay value")
	}
	return val, val, nil
}

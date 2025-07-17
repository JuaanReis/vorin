package modules

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseDelay(input string) (float64, float64, error) {

	if input == "0" {
		return 0, 0, nil
	}

	if strings.Contains(input, "-") {
		parts := strings.Split(input, "-")
		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("invalid delay format")
		}
		min, err1 := strconv.ParseFloat(parts[0], 64)
		max, err2 := strconv.ParseFloat(parts[1], 64)
		if err1 != nil || err2 != nil || min < 0 || max < min {
			return 0, 0, fmt.Errorf("invalid delay values")
		}
		return min, max, nil
	}

	val, err := strconv.ParseFloat(input, 64)
	if err != nil || val < 0 {
		return 0, 0, fmt.Errorf("invalid delay value")
	}
	return val, val, nil
}

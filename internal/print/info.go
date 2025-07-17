package print

import (
	"fmt"
)

func PrintInfo(title string, value string, width int) {
	fmt.Printf(" $  %-*s : %s\n", width, title, value)
}

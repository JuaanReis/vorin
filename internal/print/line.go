package print

import (
	"fmt"
)

func PrintLine(char string, length int, text ...string) {
	if len(text) == 0 {
		for i := 0; i < length; i++ {
			fmt.Print(char)
		}
		fmt.Println()
		return
	}

	label := " " + text[0] + " "
	side := (length - len(label)) / 2

	if side < 0 {
		side = 0
	}

	for i := 0; i < side; i++ {
		fmt.Print(char)
	}

	fmt.Print(label)

	for i := 0; i < length-side-len(label); i++ {
		fmt.Print(char)
	}

	fmt.Println()
}
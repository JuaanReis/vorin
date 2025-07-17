package print

import (
	"fmt"
)

func PrintError(msg string) {
	fmt.Printf("\033[31m[ERROR]\033[0m %s\n", msg)
}
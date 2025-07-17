package print

import (
	"os"
	"fmt"
)

func FatalIfErr(err error) {
	if err != nil {
		fmt.Printf("[ERROR] %v:\n", err)
		os.Exit(1)
	}
}
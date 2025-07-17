package print

import (
	"fmt"
	"os"
)

func Version() string {
	version, err := os.ReadFile("../assets/version.txt")
	FatalIfErr(err)
	pack := fmt.Sprintf("Vorin v%s", version)
	return pack
} 
package flags

import (
	"fmt"
	"os"
	"strings"
)

func ValidateFlags(cfg *CLIConfig) {
	if !strings.HasPrefix(cfg.URL, "http://") && !strings.HasPrefix(cfg.URL, "https://") {
		fmt.Println("\033[31m[ERROR]\033[0m URL must start with http:// or https://")
		os.Exit(1)
	}

	if cfg.Rate < 0 {
		fmt.Println("\033[31m[ERROR]\033[0m -rate must be >= 0")
		os.Exit(1)
	}

	if cfg.Silence && (cfg.Live || cfg.Verbose) {
		fmt.Println("\033[31m[ERROR]\033[0m -silence cannot be used with -live or -verbose")
		os.Exit(1)
	}

	if cfg.Method == "GET" && cfg.Wordlist == "" {
		cfg.Wordlist = "./assets/wordlist/common.txt"
	}

	if cfg.Method == "GET" && !strings.Contains(cfg.URL, "FUZZ") {
		fmt.Println("\033[31m[ERROR]\033[0m URL must contain 'FUZZ' placeholder")
		os.Exit(1)
	}

	if cfg.Method == "POST" && cfg.Payload == "" {
		fmt.Println("\033[31m[ERROR]\033[0m The -data flag (payload) is required for POST requests")
		os.Exit(1)
	}
}

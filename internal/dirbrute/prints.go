package dirbrute

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
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

func printInfo(title string, value string, width int) {
	fmt.Printf(" $  %-*s : %s\n", width, title, value)
}

func PrintHeader(url, wordlist, threads string, delay string, timeout string, header map[string]string, valid map[int]bool, stealth bool, proxy string, silence bool) {
	if !silence {
		fmt.Println()
		PrintLine("_", 80, "Vorin v1.0.1")
		PrintLine(" ", 80)
		if stealth {
			printInfo("\033[31mStealth\033[0m", "Activate", 19)
		} else {
			printInfo("Stealth", "Disabled", 10)
		}
		printInfo("URL", url, 10)
		printInfo("Wordlist", wordlist, 10)
		printInfo("Threads", threads, 10)
		if delay == "0" {
			printInfo("Delay", "Disable", 10)
		} else {
			printInfo("Delay", delay, 10)
		}
		printInfo("Timeout", timeout, 10)
		codes := make([]int, 0, len(valid))
		for code := range valid {
			codes = append(codes, code)
		}
		sort.Ints(codes)

		statusStr := []string{}
		for _, code := range codes {
			statusStr = append(statusStr, strconv.Itoa(code))
		}
		printInfo("Code HTTP", strings.Join(statusStr, ", "), 10)
		if len(header) > 0 {
			if stealth {
				sel := []string{"User-Agent", "Accept", "Accept-Language"}
				var preview []string
				for _, k := range sel {
					if v, ok := header[k]; ok {
						preview = append(preview, fmt.Sprintf("%s: %s", k, v))
					}
				}
				printInfo("Header", "(stealth) "+strings.Join(preview, " | "), 10)
			} else {
				headers := HeadersToString(header)
				printInfo("Header", headers, 10)
			}
		}
		if proxy != "" {
			printInfo("\033[31mProxy\033[0m", proxy, 19)
		} else {
			printInfo("Proxy", "Disabled", 10)
		}
		PrintLine("_", 80)
		fmt.Println()
	}
}
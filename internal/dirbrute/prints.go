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

func PrintHeader(url, wordlist, threads string, delay string, timeout string, header map[string]string, valid map[int]bool, stealth bool, proxy string, silence bool, bypass bool, extension string, rate string, filterT string, filterB string, filterL int, filterS int, shuffle bool, randomAgent bool, live bool, contentB string, contentT string, regexB string, regexT string, statusOnly bool, retries int, compare string, randomIp bool) {
	if !silence {
		fmt.Println()
		PrintLine("_", 80, "Vorin v1.3.0")
		PrintLine(" ", 80)
		if live {
			printInfo("Live mode", "Activate", 18)
		}
		if statusOnly {
			printInfo("Status Only", "Activate", 18)
		}
		if stealth {
			printInfo("\033[31mStealth\033[0m", "Activate", 27)
		}
		if bypass {
			printInfo("\033[31mBypass\033[0m", "Activate", 27)
		}
		printInfo("URL", url, 18)
		printInfo("Wordlist", wordlist, 18)
		if rate != "" {
			printInfo("Rate", rate, 18)
		}
		printInfo("Threads", threads, 18)
		if delay == "0" {
			printInfo("Delay", "Disable", 18)
		} else {
			printInfo("Delay", delay, 18)
		}
		if compare != "" {
			printInfo("compared to", compare, 18)
		}
		printInfo("Timeout", timeout, 18)
		codes := make([]int, 0, len(valid))
		for code := range valid {
			codes = append(codes, code)
		}
		sort.Ints(codes)

		statusStr := []string{}
		for _, code := range codes {
			statusStr = append(statusStr, strconv.Itoa(code))
		}
		if retries > 0 {
			printInfo("Retries", strconv.Itoa(retries), 18)
		}
		printInfo("Code HTTP", strings.Join(statusStr, ", "), 18)
		if extension != "" {
			printInfo("Extensions", extension, 18)
		}
		if len(header) > 0 {
    	if stealth {
        printInfo("Header", "(stealth) randomized headers per request", 18)
    	} else {
        headers := HeadersToString(header)
        printInfo("Header", headers, 18)
    	}
		}
		if proxy != "" {
			printInfo("\033[31mProxy\033[0m", proxy, 18)
		}
		if contentT != "" {
			printInfo("Content title", contentT, 18)
		}
		if contentB != "" {
			printInfo("Content body", contentB, 18)
		}
		if filterB != "" {
			printInfo("Filter body", filterB, 18)
		}
		if filterL > 0 {
			printInfo("Filter line", strconv.Itoa(filterL), 18)
		}
		if filterS > 0 {
			printInfo("Filter size", strconv.Itoa(filterS), 18)
		}
		if filterT != "" {
			printInfo("Filter title", filterT, 18)
		}
		if shuffle {
			printInfo("Shuffle", "Activate", 18)
		}
		if randomIp {
			printInfo("Random Ip", "Activate", 18)
		}
		if randomAgent {
			printInfo("Random-Agent", "Activate", 18)
		}
		if regexB != "" {
			printInfo("Regex body", regexB, 18)
		}
		if regexT != "" {
			printInfo("Regex title", regexT, 18)
		}
		fmt.Println()
		PrintLine("_", 80)
	}
}

func PrintError(msg string) {
	fmt.Printf("\033[31m[ERROR]\033[0m %s\n", msg)
}
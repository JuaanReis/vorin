package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"vorin/internal/dirbrute"
	"strings"
)

func headersToString(headers map[string]string) string {
	var sb strings.Builder
	for k, v := range headers {
		sb.WriteString(fmt.Sprintf("%s: %s, ", k, v))
	}
	str := sb.String()
	return strings.TrimSuffix(str, ", ")
}

func printInfo(title string, value string, width int) {
	fmt.Printf(" $  %-*s : %s\n", width, title, value)
}

func printHeader(url, wordlist, threads string, delay string, timeout string, header map[string]string) {
	fmt.Println()
	printLine("_", 60, "Vorin v1.0")
	printLine(" ", 60)
	printInfo("URL", url, 10)
	printInfo("Wordlist", wordlist, 10)
	printInfo("Threads", threads, 10)
	printInfo("Delay", delay, 10)
	printInfo("Timeout", timeout, 10)
	if len(header) > 0 {
		headers := headersToString(header)
		printInfo("Header", headers, 10)
	}
	printLine("_", 60)
	fmt.Println()
}

func printLine(char string, length int, text ...string) {
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

type headerFlags []string

func (h *headerFlags) String() string {
	return strings.Join(*h, ", ")
}

func (h *headerFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func ParseHeaderFlags(headers headerFlags) map[string]string {
	result := make(map[string]string)
	for _, h := range headers {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			fmt.Printf("[WARNING] Invalid header format: %s\n", h)
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		result[key] = val
	}
	return result
}

func main() {
	var headers headerFlags
	url := flag.String("u", "", "Target URL")
	threads := flag.Int("t", 50, "Number of concurrent threads")
	wordlist := flag.String("w", "assets/wordlist/common.txt", "Path to wordlist")
	delayFlag := flag.String("d", "1-3", "Delay entre requisições, ex: -d 1-5")
	timeout := flag.Int("timeout", 5, "Request time")
	flag.Var(&headers, "H", "Custom headers. Ex: -H 'Authorization: Bearer x' -H 'X-Test: true'")
	flag.Parse()

	if *url == "" {
		fmt.Println("\033[31m[ERROR]\033[0m The -url flag cannot be empty")
		os.Exit(1)
	}

	minDelay := 0
	maxDelay := 0

	minDelay, maxDelay, err := dirbrute.ParseDelay(*delayFlag)
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}

	delayStr := ""
	if minDelay == maxDelay {
		delayStr = fmt.Sprintf("%ds", minDelay)
	} else {
		delayStr = fmt.Sprintf("%ds-%ds", minDelay, maxDelay)
	}

	customHeader := ParseHeaderFlags(headers)

	printHeader(*url, *wordlist, strconv.Itoa(*threads), delayStr, strconv.Itoa(*timeout), customHeader)

	fmt.Println()
	printLine("_", 60, "Results")
	fmt.Println()

	resultado := dirbrute.Parser(*url, *threads, *wordlist, minDelay, maxDelay, *timeout, customHeader)

	printLine("_", 60)

	if len(resultado) == 0 {
		fmt.Println(dirbrute.Red + "\n[!!] No path found" + dirbrute.Reset)
	}
}

package print

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func HeadersToString(headers map[string]string) string {
	var sb strings.Builder
	for k, v := range headers {
		sb.WriteString(fmt.Sprintf("%s: %s, ", k, v))
	}
	str := sb.String()
	return strings.TrimSuffix(str, ", ")
}

func CookiesToString(cookies map[string]string) string {
	var parts []string
	for k, v := range cookies {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, "; ")
}


func PrintHeader(
	banner string, url, wordlist, threads, delay, timeout string,
	header map[string]string, valid map[int]bool,
	stealth bool, proxy string, silence bool,
	extension, rate, filterT, filterB string,
	filterL, filterS int, shuffle, randomAgent, live bool, regexB, regexT string,
	statusOnly bool, retries int, compare string,
	randomIp bool, method, payload, userlist, passlist string,
	redirect, logo bool, filterCode string, verbose bool, cookies map[string]string,
	calibrate bool) {
	pack := Version()

	if silence {
		return
	}

	fmt.Println()
	PrintLine("_", 80, pack)
	PrintLine(" ", 80)

	if !logo {
		fmt.Println(banner)
		PrintLine(" ", 80)
	}

	PrintInfo("URL", url, 18)

	// Method Info
	PrintInfo("Method", method, 18)
	if method == "POST" {
		PrintInfo("Payload", payload, 18)
		PrintInfo("Userlist", userlist, 18)
		PrintInfo("Passlist", passlist, 18)
	}

	// Modes
	if live {
		PrintInfo("Live mode", "Activate", 18)
	}
	if verbose {
		PrintInfo("Verbose", "Activate", 18)
	}
	if statusOnly {
		PrintInfo("Status Only", "Activate", 18)
	}
	if redirect {
		PrintInfo("Follow redirects", "Activate", 18)
	}
	if stealth {
		PrintInfo("\033[31mStealth\033[0m", "Activate", 27)
	}
	if calibrate {
		PrintInfo("Calibrate", "Activate", 18)
	}

	// Target & Config
	if wordlist != "" {
		PrintInfo("Wordlist", wordlist, 18)
	}
	if rate != "0" {
		PrintInfo("Rate", rate, 18)
	}
	PrintInfo("Threads", threads, 18)
	if delay != "0.0s" {
		PrintInfo("Delay", delay, 18)
	}
	PrintInfo("Timeout", timeout, 18)
	if retries > 0 {
		PrintInfo("Retries", strconv.Itoa(retries), 18)
	}
	if compare != "" {
		PrintInfo("Compared to", compare, 18)
	}

	// HTTP Status Codes
	codes := make([]int, 0, len(valid))
	for code := range valid {
		codes = append(codes, code)
	}
	sort.Ints(codes)
	if len(codes) > 0 {
		statusStr := make([]string, len(codes))
		for i, code := range codes {
			statusStr[i] = strconv.Itoa(code)
		}
		PrintInfo("Code HTTP", strings.Join(statusStr, ", "), 18)
	}

	if extension != "" {
		PrintInfo("Extensions", extension, 18)
	}
	if filterCode != "" {
		PrintInfo("Filter code", filterCode, 18)
	}

	// Headers
	if len(header) > 0 {
		switch {
		case stealth:
			PrintInfo("Header", "(stealth) randomized headers per request", 18)
		default:
			PrintInfo("Header", HeadersToString(header), 18)
		}
	}

	if len(cookies) > 0 {
		PrintInfo("Cookies", CookiesToString(cookies), 18)
	}

	// Proxy
	if proxy != "" {
		PrintInfo("\033[31mProxy\033[0m", proxy, 27)
	}

	// Filters
	if filterB != "" {
		PrintInfo("Filter body", filterB, 18)
	}
	if filterT != "" {
		PrintInfo("Filter title", filterT, 18)
	}
	if filterL > 0 {
		PrintInfo("Filter line", strconv.Itoa(filterL), 18)
	}
	if filterS > 0 {
		PrintInfo("Filter size", strconv.Itoa(filterS), 18)
	}
	if regexB != "" {
		PrintInfo("Regex body", regexB, 18)
	}
	if regexT != "" {
		PrintInfo("Regex title", regexT, 18)
	}

	// Extras
	if shuffle {
		PrintInfo("Shuffle", "Activate", 18)
	}
	if randomIp {
		PrintInfo("Random IP", "Activate", 18)
	}
	if randomAgent {
		PrintInfo("Random-Agent", "Activate", 18)
	}

	PrintLine("_", 80)
}

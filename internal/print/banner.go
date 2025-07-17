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

func PrintHeader(banner string, url, wordlist, threads string, delay string, timeout string, header map[string]string, valid map[int]bool, stealth bool, proxy string, silence bool, bypass bool, extension string, rate string, filterT string, filterB string, filterL int, filterS int, shuffle bool, randomAgent bool, live bool, contentB string, contentT string, regexB string, regexT string, statusOnly bool, retries int, compare string, randomIp bool, method string, payload string, userlist string, passlist string, redirect bool, logo bool, filterCode string, verbose bool) {
	pack := Version()
	if !silence {
		fmt.Println()
		PrintLine("_", 80, pack)
		PrintLine(" ", 80)
		if !logo {
			fmt.Println(banner)
			PrintLine(" ", 80)
		}
		if method == "POST" {
			PrintInfo("Method", method, 18)
			PrintInfo("Payload", payload, 18)
			PrintInfo("Userlist", userlist, 18)
			PrintInfo("Passlist", passlist, 18)
		} else {
			PrintInfo("Method", method, 18)
		}
		if live {
			PrintInfo("Live mode", "Activate", 18)
		}
		if verbose && !silence{
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
		if bypass {
			PrintInfo("\033[31mBypass\033[0m", "Activate", 27)
		}
		PrintInfo("URL", url, 18)
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

		if compare != "" {
			PrintInfo("compared to", compare, 18)
		}

		if retries > 0 {
			PrintInfo("Retries", strconv.Itoa(retries), 18)
		}

		PrintInfo("Timeout", timeout, 18)

		codes := make([]int, 0, len(valid))
		for code := range valid {
			codes = append(codes, code)
		}

		sort.Ints(codes)
		statusStr := []string{}
		for _, code := range codes {
			statusStr = append(statusStr, strconv.Itoa(code))
		}

		if len(statusStr) != 0 {
			PrintInfo("Code HTTP", strings.Join(statusStr, ", "), 18)
		}

		if extension != "" {
			PrintInfo("Extensions", extension, 18)
		}

		if filterCode != "" {
			PrintInfo("Filter code", filterCode, 18)
		}

		if len(header) > 0 {
			if stealth && bypass {
				PrintInfo("Header", "randomized headers per request", 18)
			} else if stealth {
				PrintInfo("Header", "(stealth) randomized headers per request", 18)
			} else if bypass {
				PrintInfo("Header", "(bypass) randomized headers per request", 18)
			} else {
				headers := HeadersToString(header)
				PrintInfo("Header", headers, 18)
			}
		}
		if proxy != "" {
			PrintInfo("\033[31mProxy\033[0m", proxy, 18)
		}
		if contentT != "" {
			PrintInfo("Content title", contentT, 18)
		}
		if contentB != "" {
			PrintInfo("Content body", contentB, 18)
		}
		if filterB != "" {
			PrintInfo("Filter body", filterB, 18)
		}
		if filterL > 0 {
			PrintInfo("Filter line", strconv.Itoa(filterL), 18)
		}
		if filterS > 0 {
			PrintInfo("Filter size", strconv.Itoa(filterS), 18)
		}
		if filterT != "" {
			PrintInfo("Filter title", filterT, 18)
		}
		if shuffle {
			PrintInfo("Shuffle", "Activate", 18)
		}
		if randomIp {
			PrintInfo("Random Ip", "Activate", 18)
		}
		if randomAgent {
			PrintInfo("Random-Agent", "Activate", 18)
		}
		if regexB != "" {
			PrintInfo("Regex body", regexB, 18)
		}
		if regexT != "" {
			PrintInfo("Regex title", regexT, 18)
		}
		PrintLine("_", 80)
	}
}

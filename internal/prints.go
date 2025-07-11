package internal

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"os"
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

func PrintHeader(banner string, url, wordlist, threads string, delay string, timeout string, header map[string]string, valid map[int]bool, stealth bool, proxy string, silence bool, bypass bool, extension string, rate string, filterT string, filterB string, filterL int, filterS int, shuffle bool, randomAgent bool, live bool, contentB string, contentT string, regexB string, regexT string, statusOnly bool, retries int, compare string, randomIp bool, method string, payload string, userlist string, passlist string, redirect bool) {
	pack := Version()
	if !silence {
		fmt.Println()
		PrintLine("_", 80, pack)
		PrintLine(" ", 80)
		fmt.Println(banner)
		PrintLine(" ", 80)
		if method == "POST" {
			printInfo("Method", method, 18)
			printInfo("Payload", payload, 18)
			printInfo("Userlist", userlist, 18)
			printInfo("Passlist", passlist, 18)
		} else {
			printInfo("Method", method, 18)
		}
		if live {
			printInfo("Live mode", "Activate", 18)
		}
		if statusOnly {
			printInfo("Status Only", "Activate", 18)
		}
		if redirect {
			printInfo("Follow redirects", "Activate", 18)
		}
		if stealth {
			printInfo("\033[31mStealth\033[0m", "Activate", 27)
		}
		if bypass {
			printInfo("\033[31mBypass\033[0m", "Activate", 27)
		}
		printInfo("URL", url, 18)
		if wordlist != "" {
			printInfo("Wordlist", wordlist, 18)
		}
		if rate != "0" {
			printInfo("Rate", rate, 18)
		}
		printInfo("Threads", threads, 18)
		if delay != "0.0-0.0" {
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

func FatalIfErr(err error) {
	if err != nil {
		fmt.Printf("[ERROR] %v:\n", err)
		os.Exit(1)
	}
}

func PrintHelp() {
    fmt.Println(`
 Vorin - Advanced Fuzzer
 
 Usage:
   vorin [OPTIONS]
 
 Options:
   -u, --url           Target URL (required)
   -t, --threads       Number of concurrent threads (default: 50)
   -w, --wordlist      Path to wordlist for GET fuzzing (default: assets/wordlist/common.txt)
   -userlist           User wordlist file for POST (default: assets/username/top-usernames-shortlist.txt)
   -passlist           Password wordlist file for POST (default: assets/password/rockyou-20.txt)
   -P                  Data sent to the server (payload template, ex: "user=USERFUZZ&password=PASSFUZZ")
   -d, --delay         Delay between requests, e.g. -d 1-5 (default: 0.1-0.2)
   -timeout            Request timeout in seconds (default: 5)
   -H                  Custom headers. Ex: -H 'Authorization: Bearer x' -H 'X-Test: true'
   -s                  Status codes to be considered valid (ex: -s 200,301,302)
   --stealth           Stealth mode, slower, less chance of getting caught
   --proxy             Proxy URL (ex: http://127.0.0.1:8080 or socks5://...)
   --silence           Disables any UI
   --live              Print when finding a result (slower)
   --save-json         Output file path to save results as JSON
   --bypass            Enable WAF bypass techniques
   --ext               Additional extensions, separated by commas (e.g. .php,.bak)
   --rate              Maximum number of requests per second (default: 20)
   --filter-size       Filter pages by size (ex: --filter-size 5)
   --filter-line       Filter pages by number of lines (ex: --filter-line 2)
   --filter-body       Filter pages by words in body (ex: --filter-body "Not Found")
   --filter-title      Filter pages by title (ex: --filter-title "404")
   --random-agent      Use a random user agent per request
   --shuffle           Shuffle the wordlist
   --title-contains    Returns the path containing the title content
   --body-contains     Returns the path containing the body content
   --regex-body        Apply regex to the body (ex: --regex-body "admin|login|dashboard")
   --regex-title       Apply regex to the title (ex: --regex-title "admin|login")
   --redirect          Follow status code 3xx redirection
   --status-only       Output only shows the status code and path
   --retries           Maximum number of attempts in a request
   --compare           Path to be compared to wordlist
   --random-ip         Use a random IP per request
   --method            HTTP method to use (GET, POST) (default: GET)
   -h, --help          Show this help message and exit
 
 Examples:
   # Simple GET fuzzing
   vorin -u "https://target/Fuzz" -w wordlist.txt
 
   # POST brute-force with user and password lists
   vorin -method post -u "https://target/login" -userlist users.txt -passlist passwords.txt -P "user=USERFUZZ&password=PASSFUZZ"
 
   # Save results as JSON and use custom headers
   vorin -u "https://target/Fuzz" -w wordlist.txt --save-json out.json -H "Authorization: Bearer token"
 
   # Use a proxy and random user agent
   vorin -u "https://target/Fuzz" --proxy "http://127.0.0.1:8080" --random-agent
 
   # Filter by status code and title
   vorin -u "https://target/Fuzz" -s "200,403" --filter-title "forbidden"
 
 `)
}

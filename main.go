package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/JuaanReis/vorin/internal"
)

func main() {
	banner, err := os.ReadFile("assets/banner/banner.txt")
	internal.FatalIfErr(err)
	bannerString := string(banner)
    var headers internal.HeaderFlags
    var statusCodeFlags string
    help := flag.Bool("help", false, "Show help message")
    flag.BoolVar(help, "h", false, "Show help message (shorthand)")
    url := flag.String("u", "", "Target URL")
    threads := flag.Int("t", 35, "Number of concurrent threads")
    wordlist := flag.String("w", "", "Path to wordlist")
    payload := flag.String("P", "", "data sent to the server")
    delayFlag := flag.String("d", "", "Delay between requests, e.g. -d 1-5")
    timeout := flag.Int("timeout", 5, "Request time")
    flag.Var(&headers, "H", "Custom headers. Ex: -H 'Authorization: Bearer x' -H 'X-Test: true'")
    flag.StringVar(&statusCodeFlags, "s", "200,301,302,401,403", "status codes to be considered valid (ex: -s 200,301,302)")
    stealth := flag.Bool("stealth", false, "stealth mode, slower less chance of getting caught")
    proxy := flag.String("proxy", "", "Proxy URL (ex: http://127.0.0.1:8080 or socks5://...)")
    silence := flag.Bool("silence", false, "Disables any UI")
    live := flag.Bool("live", false, "print when finding a result (slower)")
    outputFile := flag.String("save-json", "", "Output file path to save results as JSON")
    bypass := flag.Bool("bypass", false, "Enable WAF bypass techniques")
    extension := flag.String("ext", "", "Additional extensions, separated by commas (e.g. .php, .bak)")
    rate := flag.Int("rate", 0, "Maximum number of requests per second (RPS). Set 0 to disable rate limiting")
    filterSize := flag.Int("filter-size", 0, "filter pages by size (ex: -filter-size 5)")
    filterLine := flag.Int("filter-line", 0, "filters pages by number of lines (ex: -filter-size 2)")
    filterBody := flag.String("filter-body", "", "filters pages by words in body page (ex: -filter-body Not Found)")
    filterTitle := flag.String("filter-title", "", "filters pages by title (ex: -filter-title 404|forbiden)")
    randomAgent := flag.Bool("random-agent", false, "uses a random user agent per request")
    shuffle := flag.Bool("shuffle", false, "shuffle the wordlist")
    titleContains := flag.String("title-contains", "", "returns the path containing the title content")
    bodyContains := flag.String("body-contains", "", "returns the path containing the body content")
    regexBody := flag.String("regex-body", "", "Apply regex to the body (ex: -regex-body admin|login|dashboard)")
    regexTitle := flag.String("regex-title", "", "Applies regex to the title (ex: regex-title admin|login)")
    redirect := flag.Bool("redirect", false, "follow status code 3xx redirection")
    statusOnly := flag.Bool("status-only", false, "the output only shows the status code and path")
    retries := flag.Int("retries", 0, "Maximum number of attempts in a request")
    compare := flag.String("compare", "", "Path to be compared to wordlist")
    randomIp := flag.Bool("random-ip", false, "uses a random user agent per request")
    method := flag.String("method", "GET", "HTTP method to use (GET, POST)")
    userlist := flag.String("userlist", "assets/username/top-usernames-shortlist.txt", "User wordlist file")
    passlist := flag.String("passlist", "assets/password/rockyou-20.txt", "Password wordlist file")
	if *delayFlag == "" {
		*delayFlag = "0.0-0.0"
	}
    flag.Parse()
 
	 if *help {
        internal.PrintHelp()
        os.Exit(0)
    }

	if *url == "" {
		fmt.Println("\033[31m[ERROR]\033[0m The -url flag cannot be empty")
		os.Exit(1)
	}

	if *rate < 0 {
		fmt.Println("[ERROR]: -rate must be >= 0 (0 means no limit)")
		os.Exit(1)
	}

	if *silence && *live {
		fmt.Println("[ERROR] You cannot use --live and --silence at the same time.")
		os.Exit(1)
	}

	if *wordlist == "" && *method == "GET" {
		*wordlist = "assets/wordlist/common.txt"
	}

	chosenMethod := strings.ToUpper(*method)

	if *wordlist == "" && chosenMethod != "POST" {
		fmt.Printf("[ERROR] the flag -w (wordlist) is required for GET requests")
		os.Exit(1)
	}

	statusCodeFlags = strings.ReplaceAll(statusCodeFlags, " ", "")

	minDelay := float64(0)
	maxDelay := float64(0)

	minDelay, maxDelay, err = internal.ParseDelay(*delayFlag)
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}

	customHeader := internal.ParseHeaderFlags(headers)

	if *stealth {
		if *rate == 0 {
			*rate = 15
		}
		if *threads == 35 {
			*threads = 30
		}
		if *timeout == 5 {
			*timeout = 7
		}
		if minDelay == 0.1 && maxDelay == 0.2 {
			minDelay = 0.2
			maxDelay = 0.2
		}
		customHeader = internal.GetRandomHeaders()
	}

	if *bypass {
		if *rate == 0 {
			*rate = 15
		}
		if *threads == 35 {
			*threads = 30
		}
		if *timeout == 5 {
			*timeout = 8
		}
		if minDelay == 0.0 && maxDelay == 0.0 {
			minDelay = 0.2
			maxDelay = 0.3
		}
	}

	valid := internal.ParseStatusCodes(statusCodeFlags)

	if *threads <= 0 || *threads >= 250 {
		internal.PrintError("Thread count must be between 1 and 249.")
		os.Exit(1)
	}

	delayStr := ""
	if minDelay == maxDelay {
		delayStr = fmt.Sprintf("%.1fs", minDelay)
	} else {
		delayStr = fmt.Sprintf("%.1fs-%.1fs", minDelay, maxDelay)
	}

	var rateStr string 
	if *rate > 0 {
		rateStr = fmt.Sprintf("%-3dreq/s", *rate)
	} else {
		rateStr = "0"
	}

	if *outputFile == "" {
		internal.PrintHeader(bannerString, *url, *wordlist, strconv.Itoa(*threads), delayStr, fmt.Sprintf("%ds", *timeout), customHeader, valid, *stealth, *proxy, *silence, *bypass, *extension, rateStr, *filterBody, *filterTitle, *filterLine, *filterSize, *shuffle, *randomAgent, *live, *bodyContains, *titleContains, *regexBody, *regexTitle, *statusOnly, *retries, *compare, *randomIp, chosenMethod, *payload, *userlist, *passlist, *redirect)
	}

	if !*silence {
		fmt.Println()
		internal.PrintLine("_", 80, "Results")
		fmt.Println()
	}

	var listExtension []string
	if *extension != "" {
		listExtension = strings.Split(*extension, ",")
	}

	if *bypass && len(listExtension) > 0 && listExtension[0] != "" && !*stealth {
		if *rate == 0 {
			*rate = 20
		}
		if *threads == 30 {
			*threads = 35
		}
		if *timeout == 8 {
			*timeout = 6
		}
		minDelay = 0.4
		maxDelay = 0.4
	}

	var resultado []internal.Resultado
	var temp time.Duration

	configGet := internal.ParserConfigGet{
		Endereco:           *url,
		Threads:            *threads,
		Wordlist:           *wordlist,
		MinDelay:           minDelay,
		MaxDelay:           maxDelay,
		Timeout:            *timeout,
		CustomHeaders:      customHeader,
		Code:               valid,
		Stealth:            *stealth,
		Proxy:              *proxy,
		Silence:            *silence,
		Live:               *live,
		Bypass:             *bypass,
		Extension:          listExtension,
		RateLimit:          *rate,
		FilterSize:         *filterSize,
		FilterLine:         *filterLine,
		FilterTitle:        *filterTitle,
		RandomAgent:        *randomAgent,
		Shuffle:            *shuffle,
		FilterTitleContent: *titleContains,
		FilterBodyContent:  *bodyContains,
		FilterBody:         *filterBody,
		RegexBody:          *regexBody,
		RegexTitle:         *regexTitle,
		Redirect:           *redirect,
		StatusOnly:         *statusOnly,
		Retries:            *retries,
		Compare:            *compare,
		RandomIp:           *randomIp,
	}

	configPost := internal.ParserConfigPost {
		Endereco: *url,
		Threads: *threads,
		Userlist: *userlist,
		Passlist: *passlist,
		PayloadTemplate: *payload,
		MinDelay: minDelay,
		MaxDelay: maxDelay,
		Timeout: *timeout,
		CustomHeaders: customHeader,
		RandomAgent: *randomAgent,
		Shuffle: *shuffle,
		Live: *live,
		StatusOnly: *statusOnly,
		RegexBody: *regexBody,
		RegexTitle: *regexTitle,
		Silence: *silence,
	}

	switch chosenMethod {
	case "GET":
		if !strings.Contains(*url, "Fuzz") {
			fmt.Println("\033[31m[ERROR]\033[0m URL must contain 'Fuzz' placeholder")
			os.Exit(1)
		}
		resultado, temp = internal.ParserGET(configGet)
	case "POST":
		if *payload == "" {
			fmt.Println("\033[31m[ERRO]\033[0m The payload flag cannot be empty")
			os.Exit(1)
		}
		resultado, temp = internal.ParserPost(configPost)
	}

	resultadoJson := internal.PrepareResultsForJSON(resultado)

	if *statusOnly && *live {
		if *outputFile != "" {
			err := internal.SaveJson(resultadoJson, *outputFile)
			if err != nil {
				fmt.Printf("Error saving JSON to %s: %v\n", *outputFile, err)
				os.Exit(1)
			}
			fmt.Printf("Results saved to %s\n", *outputFile)
		}
	} else if *statusOnly {
		if chosenMethod != "POST" {
			for _, v := range resultado {
				fmt.Printf("%s[%3d]%s %-26s\n",
					v.Color, v.Status, internal.Reset,
					v.URL,
				)
			}
		} else {
			for _, v := range resultado {
        		fmt.Printf("%s[%3d]%s user=%s pass=%s\n", v.Color, v.Status, internal.Reset, v.User, v.Pass)
    		}
		}
		if *outputFile != "" {
			err := internal.SaveJson(resultadoJson, *outputFile)
			if err != nil {
				fmt.Printf("Error saving JSON to %s: %v\n", *outputFile, err)
				os.Exit(1)
			}
			fmt.Printf("Results saved to %s\n", *outputFile)
		}
	} else if !*live {
		if *outputFile != "" {
			err := internal.SaveJson(resultadoJson, *outputFile)
			if err != nil {
				fmt.Printf("Error saving JSON to %s: %v\n", *outputFile, err)
				os.Exit(1)
			}
			fmt.Printf("Results saved to %s\n", *outputFile)
		} else {
			if chosenMethod != "POST" { 
				for _, v := range resultado {
					fmt.Printf("%s[%3d]%s  %-26s Words: %-6d Size: %-6dB Lines: %-5d %-6s %-11s\n",
						v.Color, v.Status, internal.Reset,
						v.URL,
						v.Text,
						v.Size,
						v.Lines,
						v.Time,
						v.Label,
					)
				}
			} else {
				for _, v := range resultado {
					fmt.Printf("%s[%3d]%s user=%s pass=%s Size: %-6dB Lines: %-5d %-6s %-11s\n",
                        v.Color, v.Status, internal.Reset,
                        v.User, v.Pass,
                        v.Size,
                        v.Lines,
                        v.Time,
                        v.Label,
                	)
				}
			}
		}
	}
	if !*silence {
		internal.PrintLine("_", 80)
		fmt.Printf("\n%s[✓]%s Scan completed in %s%s%s\n\n", internal.Green, internal.Reset, internal.Blue, internal.FormatDuration(temp), internal.Reset)
	}
	if len(resultado) == 0 {
		fmt.Println(internal.Red + "\n[!!] No path found" + internal.Reset)
	}
}
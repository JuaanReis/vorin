package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"github.com/JuaanReis/vorin/internal/dirbrute"
	"strings"
)

func main() {
	var headers dirbrute.HeaderFlags
	var statusCodeFlags string
	url := flag.String("u", "", "Target URL")
	threads := flag.Int("t", 50, "Number of concurrent threads")
	wordlist := flag.String("w", "assets/wordlist/common.txt", "Path to wordlist")
	delayFlag := flag.String("d", "0.1-0.2", "Delay between requests, e.g. -d 1-5")
	timeout := flag.Int("timeout", 5, "Request time")
	flag.Var(&headers, "H", "Custom headers. Ex: -H 'Authorization: Bearer x' -H 'X-Test: true'")
	flag.StringVar(&statusCodeFlags, "s", "200,301,302,401,403", "status codes to be considered valid (ex: -s 200,301,302)")
	stealth := flag.Bool("stealth", false, "stealth mode, slower less chance of getting caught")
	proxy := flag.String("proxy", "", "Proxy URL (ex: http://127.0.0.1:8080 or socks5://...)")
	silence := flag.Bool("silence", false, "Disables any UI")
	live := flag.Bool("live", false, "print when finding a result (slower)")
	outputFile := flag.String("o", "", "Output file path to save results as JSON")
	bypass := flag.Bool("bypass", false, "Enable WAF bypass techniques")
	extension := flag.String("ext", "", "Additional extensions, separated by commas (e.g. .php, .bak)")
	rate := flag.Int("rate", 20, "Maximum number of requests per second (RPS). Set 0 to disable rate limiting")
	filterSize := flag.Int("filter-size", 0, "filter pages by size (ex: -filter-size 5)")
	filterLine := flag.Int("filter-line", 0, "filters pages by number of lines (ex: -filter-size 2)")
	filterTitle := flag.String("filter-title", "", "filters pages by title (ex: -filter-title 404|forbiden)")
	randomAgent := flag.Bool("random-agent", false, "uses a random user agent per request")
	shuffle := flag.Bool("shuffle", false, "shuffle the wordlist")
	flag.Parse()

	if *url == "" {
		fmt.Println("\033[31m[ERROR]\033[0m The -url flag cannot be empty")
		os.Exit(1)
	}

	if *rate < 0 {
		fmt.Println("[ERROR]: -rate must be >= 0 (0 means no limit)")
		os.Exit(1)
	}


	if !strings.Contains(*url, "Fuzz") {
		fmt.Println("\033[31m[ERROR]\033[0m URL must contain 'Fuzz' placeholder")
		os.Exit(1)
	}

	if *silence && *live {
		fmt.Println("[ERROR] You cannot use --live and --silence at the same time.")
		os.Exit(1)
	}

	statusCodeFlags = strings.ReplaceAll(statusCodeFlags, " ", "")

	minDelay := float64(0)
	maxDelay := float64(0)

	minDelay, maxDelay, err := dirbrute.ParseDelay(*delayFlag)
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}

	customHeader := dirbrute.ParseHeaderFlags(headers)

	if *stealth {
		if *rate == 20 { *rate = 15 }
		if *threads == 50 { *threads = 30 }
		if *timeout == 5 { *timeout = 7 }
		if minDelay == 0.1 && maxDelay == 0.2 {
			minDelay = 0.2
			maxDelay = 0.2
		}
		customHeader = dirbrute.GetRandomHeaders()
	}

	if *bypass {
		if *rate == 20 { *rate = 15 }
		if *threads == 50 { *threads = 30 }
		if *timeout == 5 { *timeout = 8 }
		if minDelay == 0.1 && maxDelay == 0.2 {
			minDelay = 0.2
			maxDelay = 0.3
		}
	}

	valid := dirbrute.ParseStatusCodes(statusCodeFlags)

	if *threads <= 0|| *threads >= 250 {
		dirbrute.PrintError("Thread count must be between 1 and 249.")
		os.Exit(1)
	}

	delayStr := ""
	if minDelay == maxDelay {
		delayStr = fmt.Sprintf("%.1fs", minDelay)
	} else {
		delayStr = fmt.Sprintf("%.1fs-%.1fs", minDelay, maxDelay)
	}

	rateStr := fmt.Sprintf("%dr/s", *rate)
	if *outputFile == "" {
		dirbrute.PrintHeader(*url, *wordlist, strconv.Itoa(*threads), delayStr, fmt.Sprintf("%ds", *timeout), customHeader, valid, *stealth, *proxy, *silence, *bypass, *extension, rateStr)
	}

	if !*silence {
		fmt.Println()
		dirbrute.PrintLine("_", 80, "Results")
		fmt.Println()
	}

	var listExtension []string
	if *extension != "" {
		listExtension = strings.Split(*extension, ",")
	}

	if *bypass && len(listExtension) > 0 && listExtension[0] != "" && !*stealth {
		if *rate == 20 { *rate = 20 }
		if *threads == 30 { *threads = 35 }
		if *timeout == 8 { *timeout = 6 }
		minDelay = 0.4
		maxDelay = 0.4
	}

	resultado, temp := dirbrute.Parser(*url, *threads, *wordlist, minDelay, maxDelay, *timeout, customHeader, valid, *stealth, *proxy, *silence, *live, *bypass, listExtension, *rate, *filterSize, *filterLine, *filterTitle, *randomAgent, *shuffle)

	resultadoJson := dirbrute.PrepareResultsForJSON(resultado)

	if !*live {
		if *outputFile != "" {
    	err := dirbrute.SaveJson(resultadoJson, *outputFile)
    	if err != nil {
        fmt.Printf("Error saving JSON to %s: %v\n", *outputFile, err)
        os.Exit(1)
    	}
    	fmt.Printf("Results saved to %s\n", *outputFile)
		} else {
    	for _, v := range resultado {
    	    fmt.Printf("%s[%3d]%s  %-26s Size: %-6dB Lines: %-5d %-6s %-11s\n",
    	        v.Color, v.Status, dirbrute.Reset,
    	        v.URL,
    	        v.Size,
    	        v.Lines,
    	        v.Time,
    	        v.Label,
    	    )
    	}
		}
	}

	if !*silence {
		dirbrute.PrintLine("_", 80)
		fmt.Printf("\n%s[✓]%s Scan completed in %s%s%s\n\n", dirbrute.Green, dirbrute.Reset, dirbrute.Blue, dirbrute.FormatDuration(temp), dirbrute.Reset)
	}

	if len(resultado) == 0 {
		fmt.Println(dirbrute.Red + "\n[!!] No path found" + dirbrute.Reset)
	}
}

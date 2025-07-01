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
	delayFlag := flag.String("d", "0", "Delay entre requisições, ex: -d 1-5")
	timeout := flag.Int("timeout", 5, "Request time")
	flag.Var(&headers, "H", "Custom headers. Ex: -H 'Authorization: Bearer x' -H 'X-Test: true'")
	flag.StringVar(&statusCodeFlags, "s", "200,301,302", "status codes to be considered valid (ex: -s 200,301,302)")
	stealth := flag.Bool("stealth", false, "stealth mode, slower less chance of getting caught")
	proxy := flag.String("proxy", "", "Proxy URL (ex: http://127.0.0.1:8080 or socks5://...)")
	silence := flag.Bool("silence", false, "Disables any UI")
	live := flag.Bool("live", false, "print when finding a result (slower)")
	outputFile := flag.String("o", "", "Output file path to save results as JSON")
	flag.Parse()

	if *url == "" {
		fmt.Println("\033[31m[ERROR]\033[0m The -url flag cannot be empty")
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

	minDelay := 0
	maxDelay := 0

	minDelay, maxDelay, err := dirbrute.ParseDelay(*delayFlag)
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}

	customHeader := dirbrute.ParseHeaderFlags(headers)

	if *stealth {
		minDelay = 3
		maxDelay = 5
		*threads = 40
		*timeout = 9
		statusCodeFlags = "200,301,302"
		customHeader = dirbrute.GetRandomHeaders()
	}

	valid := dirbrute.ParseStatusCodes(statusCodeFlags)

	if *threads <= 0 || *threads >= 250 {
		fmt.Println("[ERROR]: you can't put threads too high (> 250)")
		os.Exit(1)
	}

	delayStr := ""
	if minDelay == maxDelay {
		delayStr = fmt.Sprintf("%ds", minDelay)
	} else {
		delayStr = fmt.Sprintf("%ds-%ds", minDelay, maxDelay)
	}

	fmt.Print("\033[H\033[2J")


	if *outputFile == "" {
		dirbrute.PrintHeader(*url, *wordlist, strconv.Itoa(*threads), delayStr, fmt.Sprintf("%ds", *timeout), customHeader, valid, *stealth, *proxy, *silence)
	}

	fmt.Println()
	dirbrute.PrintLine("_", 80, "Results")
	fmt.Println()

	resultado, temp := dirbrute.Parser(*url, *threads, *wordlist, minDelay, maxDelay, *timeout, customHeader, valid, *stealth, *proxy, *silence, *live)

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

	dirbrute.PrintLine("_", 80)

	fmt.Printf("\n%s[✓]%s Scan completed in %s%s%s\n\n", dirbrute.Green, dirbrute.Reset, dirbrute.Blue, dirbrute.FormatDuration(temp), dirbrute.Reset)

	if len(resultado) == 0 {
		fmt.Println(dirbrute.Red + "\n[!!] No path found" + dirbrute.Reset)
	}
}

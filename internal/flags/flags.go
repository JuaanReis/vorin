package flags

import (
	"flag"
	"io"
	"os"
)

type CLIConfig struct {
	URL              string
	Threads          int
	Wordlist         string
	Method           string
	Payload          string
	Timeout          int
	Delay            string
	Proxy            string
	Rate             int
	Stealth          bool
	Silence          bool
	Live             bool
	RandomAgent      bool
	RandomIp         bool
	Shuffle          bool
	StatusOnly       bool
	Verbose          bool
	Redirect         bool
	OutputFile       string
	NoBanner         bool
	Compare          string
	Extension        string
	StatusCodeFlags  string
	FilterCodeFlags  string
	FilterSize       int
	FilterLine       int
	FilterBody       string
	FilterTitle      string
	RegexBody        string
	RegexTitle       string
	Retries          int
	Userlist         string
	Passlist         string
	Help             bool
	HeaderFlags      []string
	Cookies          []string
	Calibrate        bool
}

func ParseFlags() *CLIConfig {
	cfg := &CLIConfig{}

	flag.StringVar(&cfg.URL, "u", "", "Target URL")
	flag.StringVar(&cfg.URL, "url", "", "Target URL")
	flag.StringVar(&cfg.Wordlist, "wordlist", "internal/wordlist/common.txt", "Path to wordlist")
	flag.StringVar(&cfg.Wordlist, "w", "internal/wordlist/common.txt", "Path to wordlist")
	flag.IntVar(&cfg.Threads, "t", 35, "Number of concurrent threads")
	flag.IntVar(&cfg.Threads, "thread", 35, "Number of concurrent threads")
	flag.StringVar(&cfg.Method, "method", "GET", "HTTP method to use (GET, POST)")
	flag.StringVar(&cfg.Method, "X", "GET", "HTTP method to use (GET, POST)")
	flag.StringVar(&cfg.Payload, "data", "", "Data sent to the server")
	flag.StringVar(&cfg.Payload, "D", "", "Data sent to the server")
	flag.IntVar(&cfg.Timeout, "timeout", 5, "Request timeout in seconds")
	flag.IntVar(&cfg.Timeout, "T", 5, "Request timeout in seconds")
	flag.StringVar(&cfg.Delay, "d", "0.0-0.0", "Delay between requests (e.g., 0.1-0.3)")
	flag.StringVar(&cfg.Delay, "delay", "0.0-0.0", "Delay between requests (e.g., 0.1-0.3)")
	flag.StringVar(&cfg.Proxy, "proxy", "", "Proxy URL")
	flag.IntVar(&cfg.Rate, "rate", 0, "Max requests per second")
	flag.BoolVar(&cfg.Stealth, "stealth", false, "Enable stealth mode")
	flag.BoolVar(&cfg.Silence, "silence", false, "Disable UI")
	flag.BoolVar(&cfg.Live, "live", false, "Print live results")
	flag.BoolVar(&cfg.RandomAgent, "random-agent", false, "Random user-agent per request")
	flag.BoolVar(&cfg.RandomIp, "spoof-ip", false, "Random IP per request")
	flag.BoolVar(&cfg.Shuffle, "shuffle", false, "Shuffle wordlist")
	flag.BoolVar(&cfg.StatusOnly, "status-only", false, "Only print status codes")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&cfg.Verbose, "v", false, "Verbose output")
	flag.BoolVar(&cfg.Redirect, "redirect", false, "Follow 3xx redirects")
	flag.StringVar(&cfg.OutputFile, "save-json", "", "Save results to JSON")
	flag.StringVar(&cfg.OutputFile, "o", "", "Save results to JSON")
	flag.BoolVar(&cfg.NoBanner, "no-banner", false, "Disable banner")
	flag.StringVar(&cfg.Compare, "compare", "", "Compare response with known false positive")
	flag.StringVar(&cfg.Compare, "c", "", "Compare response with known false positive")
	flag.BoolVar(&cfg.Calibrate, "calibrate", false, "Detect false path")
	flag.BoolVar(&cfg.Calibrate, "C", false, "Detect false path")
	flag.StringVar(&cfg.StatusCodeFlags, "status-code", "", "Valid status codes")
	flag.StringVar(&cfg.StatusCodeFlags, "sc", "", "Valid status codes")
	flag.StringVar(&cfg.FilterCodeFlags, "filter-code", "", "Filter by status codes")
	flag.StringVar(&cfg.FilterCodeFlags, "fc", "", "Filter by status codes")
	flag.IntVar(&cfg.FilterSize, "filter-size", 0, "Filter by response size")
	flag.IntVar(&cfg.FilterSize, "fs", 0, "Filter by response size")
	flag.IntVar(&cfg.FilterLine, "filter-line", 0, "Filter by response line count")
	flag.IntVar(&cfg.FilterLine, "fl", 0, "Filter by response line count")
	flag.StringVar(&cfg.FilterBody, "filter-body", "", "Filter by body content")
	flag.StringVar(&cfg.FilterBody, "fb", "", "Filter by body content")
	flag.StringVar(&cfg.FilterTitle, "filter-title", "", "Filter by title content")
	flag.StringVar(&cfg.FilterTitle, "ft", "", "Filter by title content")
	flag.StringVar(&cfg.RegexBody, "regex-body", "", "Regex on body")
	flag.StringVar(&cfg.RegexBody, "rb", "", "Regex on body")
	flag.StringVar(&cfg.RegexTitle, "regex-title", "", "Regex on title")
	flag.StringVar(&cfg.RegexTitle, "rt", "", "Regex on title")
	flag.IntVar(&cfg.Retries, "retries", 0, "Retry count")
	flag.StringVar(&cfg.Userlist, "userlist", "internal/username/top-usernames-shortlist.txt", "User list")
	flag.StringVar(&cfg.Userlist, "ul", "internal/username/top-usernames-shortlist.txt", "User list")
	flag.StringVar(&cfg.Passlist, "passlist", "internal/password/rockyou-20.txt", "Pass list")
	flag.StringVar(&cfg.Passlist, "pl", "internal/password/rockyou-20.txt", "Pass list")
	flag.BoolVar(&cfg.Help, "help", false, "Show help message")
	flag.BoolVar(&cfg.Help, "h", false, "Show help (shorthand)")
	flag.Func("H", "Custom headers (e.g., -H 'Key: Value')", func(h string) error {
		cfg.HeaderFlags = append(cfg.HeaderFlags, h)
		return nil
	})
	flag.Func("b", "Custom cookie (e.g., -b 'name=value')", func(c string) error {
		cfg.Cookies = append(cfg.Cookies, c)
		return nil
	})
	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard) 

	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil || cfg.Help || len(os.Args[1:]) == 0 {
		PrintHelp()
		os.Exit(0)
	}
	flag.Parse()

	return cfg
}

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"vorin/internal/dirbrute"
)

func printInfo(title string, value string, width int) {
	fmt.Printf(" :: %-*s : %s\n", width, title, value)
}

func printHeader(method, url, wordlist, threads string) {
	fmt.Println()
	printLine("_", 60)
	fmt.Println("         VORIN :: Directory Brute Forcer")
	printLine("_", 60)
	printInfo("Method", method, 15)
	printInfo("URL", url, 15)
	printInfo("Wordlist", wordlist, 15)
	printInfo("Threads", threads, 15)
	printLine("_", 60)
	fmt.Println()
}

func printLine(char string, length int) {
	for i := 0; i < length; i++ {
		fmt.Print(char)
	}
	fmt.Println()
}

func main() {
	url := flag.String("u", "", "Target URL")
	threads := flag.Int("t", 50, "Number of concurrent threads")
	wordlist := flag.String("w", "assets/wordlist/common.txt", "Path to wordlist")
	method := flag.String("X", "GET", "HTTP Method (default: GET)")
	flag.Parse()

	if *url == "" {
		fmt.Println("\033[31m[ERROR]\033[0m The -url flag cannot be empty")
		os.Exit(1)
	}
	
	printHeader(*method, *url, *wordlist, strconv.Itoa(*threads))

	resultado := dirbrute.Parser(*url, *threads, *wordlist)

	if len(resultado) == 0 {
		fmt.Println("\033[33m[-]\033[0m No path found.")
	}
}

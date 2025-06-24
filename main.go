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

func printHeader(url, wordlist, threads string, delay string) {
	fmt.Println()
	printLine("_", 60, "Vorin")
  printLine(" ", 60)
	printInfo("URL", url, 15)
	printInfo("Wordlist", wordlist, 15)
	printInfo("Threads", threads, 15)
	printInfo("Delay", delay, 15)
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

func main() {
	url := flag.String("u", "", "Target URL")
	threads := flag.Int("t", 50, "Number of concurrent threads")
	wordlist := flag.String("w", "assets/wordlist/common.txt", "Path to wordlist")
	delay := flag.Int("d", 5, "Timeout between requests")
	flag.Parse()

	if *url == "" {
		fmt.Println("\033[31m[ERROR]\033[0m The -url flag cannot be empty")
		os.Exit(1)
	}

	printHeader(*url, *wordlist, strconv.Itoa(*threads), strconv.Itoa(*delay))

	fmt.Println()
	printLine("_", 60, "Results")
	fmt.Println()

	resultado := dirbrute.Parser(*url, *threads, *wordlist, *delay)

	printLine("_", 60)


	if len(resultado) == 0 {
		fmt.Println(dirbrute.Red + "[!!] No path found" + dirbrute.Reset)
	}
}

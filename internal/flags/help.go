package flags

import (
	"fmt"
)

func PrintHelp() {
    fmt.Println(`
 Vorin - Advanced Fuzzer
 
 Usage:
   vorin [OPTIONS]
 
 Options:
   -u                 Target URL (required)
   -t                 Number of concurrent threads (default: 50)
   -w                 Path to wordlist for GET fuzzing (default: assets/wordlist/common.txt)
   -userlist          User wordlist file for POST (default: assets/username/top-usernames-shortlist.txt)
   -passlist          Password wordlist file for POST (default: assets/password/rockyou-20.txt)
   -data              Data sent to the server (payload template, ex: "user=USERFUZZ&password=PASSFUZZ")
   -d                 Delay between requests, e.g. -d 1-5 (default: 0.1-0.2)
   -timeout           Request timeout in seconds (default: 5)
   -H                 Custom headers. Ex: -H 'Authorization: Bearer x' -H 'X-Test: true'
   -status-code       Status codes to be considered valid (ex: -s 200,301,302)
   -stealth           Stealth mode, slower, less chance of getting caught
   -proxy             Proxy URL (ex: http://127.0.0.1:8080 or socks5://...)
   -silence           Disables any UI
   -live              Print when finding a result (slower)
   -no-banner         Disable a ascii art
   -verbose           Shows more details of the path such as the entire path and the path used in the fuzz
   -save-json         Output file path to save results as JSON
   -bypass            Enable WAF bypass techniques
   -ext               Additional extensions, separated by commas (e.g. .php,.bak)
   -rate              Maximum number of requests per second (default: 20)
   -filter-size       Filter pages by size (ex: --filter-size 5)
   -filter-line       Filter pages by number of lines (ex: --filter-line 2)
   -filter-body       Filter pages by words in body (ex: --filter-body "Invalid request")
   -filter-title      Filter pages by title (ex: --filter-title "404 Not Found")
   -filter-code       Filter pages by status code (ex: --filter-code "404, 500, 502")
   -random-agent      Use a random user agent per request
   -shuffle           Shuffle the wordlist
   -title-contains    Returns the path containing the title content
   -body-contains     Returns the path containing the body content
   -regex-body        Apply regex to the body (ex: --regex-body "admin|login|dashboard")
   -regex-title       Apply regex to the title (ex: --regex-title "admin|login|mysql|root")
   -redirect          Follow status code 3xx redirection
   -status-only       Output only shows the status code and path
   -retries           Maximum number of attempts in a request
   -compare           Path to be compared to wordlist
   -spoof-ip          Use a random IP per request
   -method            HTTP method to use (GET, POST) (default: GET)
   -help/-h           Show this help message and exit
 
 Examples:
   # Simple GET fuzzing
   vorin -u "https://target/FUZZ" -w wordlist.txt
 
   # POST brute-force with user and password lists
   vorin -method post -u "https://target/login" -userlist users.txt -passlist passwords.txt -data "user=USERFUZZ&password=PASSFUZZ"
 
   # Save results as JSON and use custom headers
   vorin -u "https://target/FUZZ" -w wordlist.txt -save-json out.json -H "Authorization: Bearer token"
 
   # Use a proxy and random user agent
   vorin -u "https://target/FUZZ" -proxy "http://127.0.0.1:8080" -random-agent
 
   # Filter by status code and title
   vorin -u "https://target/FUZZ" -s "200,403" -filter-title "forbidden"
 
 `)
}

package model

type ParserConfigGet struct {
	Endereco           string
	Threads            int
	Wordlist           string
	MinDelay           float64
	MaxDelay           float64
	Timeout            int
	CustomHeaders      map[string]string
	Code               map[int]bool
	Stealth            bool
	Proxy              string
	Silence            bool
	Live               bool
	Bypass             bool
	Extension          []string
	RateLimit          int
	FilterSize         int
	FilterLine         int
	FilterTitle        string
	RandomAgent        bool
	Shuffle            bool
	FilterTitleContent string
	FilterBodyContent  string
	FilterBody         string
	RegexBody          string
	RegexTitle         string
	Redirect           bool
	StatusOnly         bool
	Retries            int
	Compare            string
	RandomIp           bool
	FilterCode         map[int]bool
	Verbose            bool
	Title              string
	RedirectDepth              int
}

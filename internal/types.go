package internal

import (
	"time"
)

type Resultado struct {
	Status int
	Resposta string
	URL    string
	Title  string
	Text int
	Size   int
	Lines  int
	Time   time.Duration
	Label  string
	Color  string
	User string
	Pass string
}

type ResultadoJSON struct {
	Status int    `json:"status"`
	URL    string `json:"path"`
	Title  string `json:"title"`
	Size   int    `json:"size"`
	Lines  int    `json:"lines"`
	TimeMs int64  `json:"time_ms"`
	Label  string `json:"label"`
}

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
}

type ParserConfigPost struct{
	Endereco string
	Threads int
	Userlist string
	Passlist string
	PayloadTemplate string
	MinDelay float64
	MaxDelay float64
	Timeout int
	CustomHeaders map[string]string
	RandomAgent bool
	Shuffle bool
	Live bool
	StatusOnly bool
	RegexBody string
	RegexTitle string
	Silence bool
}

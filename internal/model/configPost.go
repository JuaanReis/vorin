package model

type ParserConfigPost struct {
	Endereco        string
	Threads         int
	Userlist        string
	Passlist        string
	PayloadTemplate string
	MinDelay        float64
	MaxDelay        float64
	Timeout         int
	CustomHeaders   map[string]string
	RandomAgent     bool
	Shuffle         bool
	Live            bool
	StatusOnly      bool
	RegexBody       string
	RegexTitle      string
	Silence         bool
	FilterCode      map[int]bool
	Verbose         bool
	RandomIp        bool
	Retries         int
	Proxy           string
	RateLimit       int
	Cookies         map[string]string
	Calibrate       bool
}

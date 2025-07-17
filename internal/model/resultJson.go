package model

type ResultadoJSON struct {
	Status int    `json:"status"`
	URL    string `json:"path"`
	Title  string `json:"title"`
	Size   int    `json:"size"`
	Lines  int    `json:"lines"`
	TimeMs int64  `json:"time_ms"`
	Label  string `json:"label"`
}

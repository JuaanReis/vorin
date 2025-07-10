package internal

import (
	"encoding/json"
	"os"
	"regexp"
)

func SaveJson(resultados []ResultadoJSON, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(resultados)
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

var ansi = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSICodes(input string) string {
	return ansi.ReplaceAllString(input, "")
}

func PrepareResultsForJSON(results []Resultado) []ResultadoJSON {
	cleaned := make([]ResultadoJSON, 0, len(results))
	for _, r := range results {
		cleaned = append(cleaned, ResultadoJSON{
			Status: r.Status,
			URL:    r.URL,
			Title:  r.Title,
			Size:   r.Size,
			Lines:  r.Lines,
			TimeMs: r.Time.Milliseconds(),
			Label:  stripANSICodes(r.Label),
		})
	}
	return cleaned
}

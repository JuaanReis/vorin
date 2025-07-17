package collector

import (
	"github.com/JuaanReis/vorin/internal/analyzer"
	"strings"
)

func DataTaget(body []byte) (int, int, string, int, string, int) {
	stringBody := string(body)
	structureOnly := analyzer.CleanStructure(stringBody)
	structureSize := len(structureOnly)
	htmlSize := len(body)
	title := strings.TrimSpace(strings.ToLower(analyzer.GetTitle(stringBody)))
	lines := len(strings.Split(stringBody, "\n"))
	text := ExtractTextFromHTML(stringBody)
	words := CountWords(text)
	return structureSize, htmlSize, title, lines, stringBody, words
}
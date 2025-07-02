package dirbrute

import (
	"strings"
)

func DataTaget(body []byte) (int, int, string, int, string){
	stringBody := string(body)
	structureOnly := cleanStructure(stringBody)
	structureSize := len(structureOnly)
	htmlSize := len(body)
	title := strings.TrimSpace(strings.ToLower(getTitle(stringBody)))
	lines := len(strings.Split(stringBody, "\n"))
	return structureSize, htmlSize, title, lines, stringBody
}
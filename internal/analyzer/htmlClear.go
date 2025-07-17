package analyzer

import (
	"regexp"
	"strings"
)

func CleanHTML(input string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	clean := re.ReplaceAllString(input, " ")
	clean = strings.Join(strings.Fields(clean), " ")
	clean = RemoveAccents(clean)
	return strings.ToLower(clean)
}
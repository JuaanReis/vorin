package analyzer

import (
	"regexp"
)

func CleanStructure(html string) string {
	re := regexp.MustCompile(`>([^<]*)<`)
	return re.ReplaceAllString(html, `><`)
}
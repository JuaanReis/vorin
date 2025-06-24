package dirbrute

import (
	"regexp"
	"strings"
)

func getTitle(html string) string {
	re := regexp.MustCompile("(?i)<title>(.*?)</title>")
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func content404(html string) bool {
	err := []string {
		"404", "not found", "pagina não encontrada", "no content", "erro", "não existe", "not found", "does not exist",
	}
	htmlLower := strings.ToLower(html)

	for _, e := range err {
		if strings.Contains(htmlLower, e) {
			return true
		}
	}
	return false
}

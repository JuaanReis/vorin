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
		"404",
		"erro 404",
		"error 404",
		"Error",
		"not found",
		"page not found",
		"404 not found",
		"does not exist",
		"doesn't exist",
		"no such file",
		"no such page",
		"not available",
		"n/a",
		"resource not found",
		"file not found",
		"invalid url",
		"the requested url was not found",
		"page could not be found",
		"no results found",
		"empty result",
		"response empty",
		"document not found",
		"webpage not found",
		"unable to locate",
		"link is broken",
		"invalid link",
		"this page does not exist",
		"this page isn't available",
		"no content",
		"status: 404",

		// Português
		"página não encontrada",
		"pagina não encontrada",
		"página inexistente",
		"página inválida",
		"não existe",
		"não foi encontrado",
		"conteúdo não encontrado",
		"sem resultados",
		"sem conteúdo",
		"link quebrado",
		"endereço inválido",
		"erro",
		"erro ao carregar página",
		"página removida",
	}

	htmlLower := strings.ToLower(html)

	for _, e := range err {
		if strings.Contains(htmlLower, e) {
			return true
		}
	}
	return false
}

package analyzer

import (
	"strings"
)

func Content404(html string) bool {
	errKeywords := []string{
		"404", "error 404", "not found", "page not found", "does not exist",
		"no such page", "invalid url", "status: 404", "empty result",
		"document not found", "webpage not found", "broken link",
		"this page does not exist", "link is broken", "no content", "Invalid",

		"pagina nao encontrada", "pagina inexistente", "nao existe", "conteudo nao encontrado",
		"sem resultados", "endereco invalido", "erro 404", "pagina removida",
	}
	cleaned := CleanHTML(html)

	for _, keyword := range errKeywords {
		if strings.Contains(cleaned, keyword) {
			return true
		}
	}

	return false
}
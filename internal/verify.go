package internal

import (
	"regexp"
	"strings"
)

func cleanHTML(input string) string {

	re := regexp.MustCompile(`<[^>]*>`)
	clean := re.ReplaceAllString(input, " ")

	clean = strings.Join(strings.Fields(clean), " ")

	clean = removeAccents(clean)

	return strings.ToLower(clean)
}

func removeAccents(s string) string {
	replacer := strings.NewReplacer(
		"á", "a", "à", "a", "ã", "a", "â", "a",
		"Á", "A", "À", "A", "Ã", "A", "Â", "A",
		"é", "e", "ê", "e",
		"É", "E", "Ê", "E",
		"í", "i", "Í", "I",
		"ó", "o", "ô", "o", "õ", "o",
		"Ó", "O", "Ô", "O", "Õ", "O",
		"ú", "u", "Ú", "U",
		"ç", "c", "Ç", "C",
	)
	return replacer.Replace(s)
}

func content404(html string) bool {
	errKeywords := []string{
		"404", "error 404", "not found", "page not found", "does not exist",
		"no such page", "invalid url", "status: 404", "empty result",
		"document not found", "webpage not found", "broken link",
		"this page does not exist", "link is broken", "no content", "Invalid",

		"pagina nao encontrada", "pagina inexistente", "nao existe", "conteudo nao encontrado",
		"sem resultados", "endereco invalido", "erro 404", "pagina removida",
	}
	cleaned := cleanHTML(html)

	for _, keyword := range errKeywords {
		if strings.Contains(cleaned, keyword) {
			return true
		}
	}

	return false
}

func getTitle(html string) string {
	re := regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func cleanStructure(html string) string {
	re := regexp.MustCompile(`>([^<]*)<`)
	return re.ReplaceAllString(html, `><`)
}

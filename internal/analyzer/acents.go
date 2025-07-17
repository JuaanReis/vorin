package analyzer

import (
	"strings"
)

func RemoveAccents(s string) string {
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
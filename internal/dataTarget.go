package internal

import (
	"strings"
	"golang.org/x/net/html"
)

func CountWords(text string) int {
	words := strings.Fields(text)
	return len(words)
}

func DataTaget(body []byte) (int, int, string, int, string, int) {
	stringBody := string(body)
	structureOnly := cleanStructure(stringBody)
	structureSize := len(structureOnly)
	htmlSize := len(body)
	title := strings.TrimSpace(strings.ToLower(getTitle(stringBody)))
	lines := len(strings.Split(stringBody, "\n"))
	text := ExtractTextFromHTML(stringBody)
	words := CountWords(text)
	return structureSize, htmlSize, title, lines, stringBody, words
}

func CountLinesAndSize(body string) (int, int) {
    lines := strings.Count(body, "\n")
    return lines, len(body)
}

func DataTargetFake(bodyAle []byte) (int, string) {
	stringBodyAle := string(bodyAle)
	structureOnly := cleanStructure(stringBodyAle)
	fakeStructureSize := len(structureOnly)
	titleAle := strings.TrimSpace(strings.ToLower(getTitle(stringBodyAle)))
	return fakeStructureSize, titleAle
}


func ExtractTextFromHTML(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return ""
	}
	var b strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data + " ")
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return b.String()
}